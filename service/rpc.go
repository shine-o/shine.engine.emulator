package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/shine-o/shine.engine.core/grpc/zone-master"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"net"
)

type server struct {
	zm.UnimplementedMasterServer
}

var errBadRPCClient = errors.New("gRPC client is not present in the config file")

func (s *server) RegisterZone(ctx context.Context, zd *zm.ZoneDetails) (*zm.ZoneRegistered, error) {
	// keep in redis:
	// for each map, add new entry to registeredMaps with the map name as key and the zone connection info as value
	var rm registeredMaps

	res, err := redisClient.Get("zone-master").Result()
	if err != nil {
		rm = make(registeredMaps)
	} else {
		err = json.Unmarshal([]byte(res), &rm)
	}

	for _, m := range zd.Maps {
		rm[m] = ZoneInfo{
			IP:   zd.Conn.IP,
			Port: zd.Conn.Port,
		}
		log.Infof("registering map %v for zone %v:%v", m, zd.Conn.IP, zd.Conn.Port)
	}
	err = persist(&rm)
	if err != nil {
		return &zm.ZoneRegistered{}, fmt.Errorf("failed to persist data to redis: %v", err)
	}
	return &zm.ZoneRegistered{
		Success: true,
		ZoneID:  "",
	}, nil
}

func (s *server) WhereIsMap(ctx context.Context, zd *zm.MapQuery) (*zm.ConnectionInfo, error) {
	var rm registeredMaps

	res, err := redisClient.Get("zone-master").Result()
	if err != nil {
		return &zm.ConnectionInfo{}, fmt.Errorf("failed to fetch data from redis: %v", err)
	}

	err = json.Unmarshal([]byte(res), &rm)
	if err != nil {
		return &zm.ConnectionInfo{}, fmt.Errorf("failed to unmarshal data from redis: %v", err)
	}

	ci, ok := rm[zd.ID]
	if !ok {
		return &zm.ConnectionInfo{}, fmt.Errorf("no ConnectionInfofor map %v: %v",zd.ID, err)
	}

	return &zm.ConnectionInfo{
		IP:                   ci.IP,
		Port:                 ci.Port,
	}, nil
}

func newRPCClient(name string) (*grpc.ClientConn, error) {
	clientKey := fmt.Sprintf("gRPC.clients.%v", name)
	if viper.IsSet(clientKey) {
		host := viper.GetString(fmt.Sprintf("%v.%v", clientKey, "host"))
		port := viper.GetString(fmt.Sprintf("%v.%v", clientKey, "port"))
		address := fmt.Sprintf("%v:%v", host, port)

		conn, err := grpc.Dial(address, grpc.WithInsecure())
		if err != nil {
			log.Errorf("could not connect service %v : %v", name, err)
			return &grpc.ClientConn{}, err
		}
		log.Infof("[gRPC] client connection: %v@%v:%v", name, host, port)
		return conn, nil
	}
	return &grpc.ClientConn{}, errBadRPCClient
}

func newRPCServer(name string) {
	clientKey := fmt.Sprintf("gRPC.servers.%v", name)
	if viper.IsSet(clientKey) {
		port := viper.GetString(fmt.Sprintf("%v.%v", clientKey, "port"))
		address := fmt.Sprintf(":%v", port)
		lis, err := net.Listen("tcp", address)
		if err != nil {
			log.Errorf("could listen on port %v for service %v : %v", name, port, err)
		}
		s := grpc.NewServer()
		zm.RegisterMasterServer(s, &server{})

		log.Infof("Loading gRPC server connection %v@::%v", name, port)
		if err := s.Serve(lis); err != nil {
			log.Errorf("failed to serve: %v", err)
		}
	}
}
