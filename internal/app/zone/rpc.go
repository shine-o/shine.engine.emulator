package zone

import (
	"errors"
	"fmt"
	z "github.com/shine-o/shine.engine.emulator/internal/pkg/grpc/zone"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"net"
	"time"
)

const gRPCTimeout = time.Second * 2

type server struct {
	z.UnimplementedZoneServer
}

var errBadRPCClient = errors.New("gRPC client is not present in the config file")

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
		z.RegisterZoneServer(s, &server{})

		log.Infof("Loading gRPC server connection %v@::%v", name, port)
		if err := s.Serve(lis); err != nil {
			log.Errorf("failed to serve: %v", err)
		}
	}
}
