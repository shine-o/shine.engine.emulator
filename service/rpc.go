package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/shine-o/shine.engine.core/grpc/login"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net"
)

type server struct {
	login.UnimplementedLoginServer
}

var errBadRPCClient = errors.New("gRPC client is not present in the config file")

func (s *server) AccountInfo(ctx context.Context, req *login.User) (*login.UserInfo, error) {
	var userID uint64
	err := db.Model((*User)(nil)).Column("id").Where("user_name = ?", req.UserName).Limit(1).Select(&userID)
	if err != nil {
		return &login.UserInfo{}, status.Errorf(codes.FailedPrecondition, "failed to fetch user record %v", err)
	}
	return &login.UserInfo{
		UserID: userID,
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
		login.RegisterLoginServer(s, &server{})

		log.Infof("loading gRPC server connection %v@::%v", name, port)
		if err := s.Serve(lis); err != nil {
			log.Errorf("failed to serve: %v", err)
		}
	}
}
