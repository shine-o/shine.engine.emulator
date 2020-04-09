package service

import (
	"context"
	"errors"
	"fmt"
	lw "github.com/shine-o/shine.engine.protocol-buffers/login-world"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"net"
	"time"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	lw.UnimplementedWorldServer
}

const gRPCTimeout = time.Second * 2
var errBadRpcClient = errors.New("gRPC client is not present in the config file")

func (s *server) GetWorldData(ctx context.Context, req *lw.WorldQuery) (*lw.WorldData, error) {
	worldID := viper.GetInt("world.id")
	worldName := viper.GetString("world.name")
	worldIP := viper.GetString("world.external_ip")
	worldPort := int32(viper.GetInt("world.port"))
	return &lw.WorldData{
		WorldNumber:          int32( worldID),
		WorldName:            worldName,
		WorldStatus:          6,
		IP:   worldIP,
		Port: worldPort,
	}, nil
}

func newRpcConn(name string) (*grpc.ClientConn, error){
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

	return &grpc.ClientConn{}, errBadRpcClient
}

func newRpcServer(name string) {
	clientKey := fmt.Sprintf("gRPC.servers.%v", name)
	if viper.IsSet(clientKey) {
		port := viper.GetString(fmt.Sprintf("%v.%v", clientKey, "port"))
		address := fmt.Sprintf(":%v",  port)
		lis, err := net.Listen("tcp", address)
		if err != nil {
			log.Errorf("could listen on port %v for service %v : %v", name, port, err)
		}
		s := grpc.NewServer()
		lw.RegisterWorldServer(s, &server{})

		log.Infof("Loading gRPC server connections %v@::%v", name, port)
		if err := s.Serve(lis); err != nil {
			log.Errorf("failed to serve: %v", err)
		}
	}
}