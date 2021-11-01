package worldmaster

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"

	wm "github.com/shine-o/shine.engine.emulator/internal/pkg/grpc/world-master"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

type server struct {
	wm.UnimplementedMasterServer
}

var errBadRPCClient = errors.New("gRPC client is not present in the config file")

func (s *server) RegisterWorld(ctx context.Context, wd *wm.WorldDetails) (*wm.WorldRegistered, error) {
	// keep in redis:
	// for each map, add new entry to registeredMaps with the map name as key and the zone connection info as value
	var rw registeredWorlds

	res, err := redisClient.Get("world-master").Result()
	if err != nil {
		rw = make(registeredWorlds)
	} else {
		err = json.Unmarshal([]byte(res), &rw)
	}

	rw[wd.ID] = WorldInfo{
		ID:   int(wd.ID),
		Name: wd.Name,
		IP:   wd.Conn.IP,
		Port: wd.Conn.Port,
	}

	log.Infof("registering world %v  %v:%v", wd.Name, wd.Conn.IP, wd.Conn.Port)

	err = persist(&rw)
	if err != nil {
		return &wm.WorldRegistered{}, fmt.Errorf("failed to persist data to redis: %v", err)
	}
	return &wm.WorldRegistered{
		Success: true,
	}, nil
}

func (s *server) GetWorlds(ctx context.Context, wd *wm.Empty) (*wm.Worlds, error) {
	rw := make(registeredWorlds)
	worlds := &wm.Worlds{}

	res, err := redisClient.Get("world-master").Result()
	if err != nil {
		return worlds, err
	}

	err = json.Unmarshal([]byte(res), &rw)

	if err != nil {
		return worlds, err
	}

	for _, w := range rw {
		worlds.List = append(worlds.List, &wm.WorldDetails{
			ID:   int32(w.ID),
			Name: w.Name,
			Conn: &wm.ConnectionInfo{
				IP:   w.IP,
				Port: w.Port,
			},
		})
	}
	return worlds, nil
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
	// clientKey := fmt.Sprintf("gRPC.servers.%v", name)
	// TODO: check if viper.IsSet(clientKey)?
	port := viper.GetString("serve.port")
	address := fmt.Sprintf(":%v", port)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Errorf("could listen on port %v for service %v : %v", name, port, err)
	}
	s := grpc.NewServer()
	wm.RegisterMasterServer(s, &server{})

	log.Infof("Loading gRPC server connection %v@::%v", name, port)
	if err := s.Serve(lis); err != nil {
		log.Errorf("failed to serve: %v", err)
	}
}
