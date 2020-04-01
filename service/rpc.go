package service

import (
	"context"
	"fmt"
	lw "github.com/shine-o/shine.engine.protocol-buffers/login-world"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net"
	"os"
	"sync"
	"time"
)

type server struct {
	lw.UnimplementedLoginServer
}

// RPCClients that this service will use to communicate with other services
type RPCClients struct {
	services map[string]*grpc.ClientConn
	mu       sync.Mutex
}

const gRPCTimeout = time.Second * 2

func (s *server) AccountInfo(ctx context.Context, req *lw.User) (*lw.UserInfo, error) {

	// query user by user_name
	var userID uint64
	err := db.Model((*User)(nil)).Column("id").Where("user_name = ?", req.UserName).Limit(1).Select(&userID)
	if err != nil {
		return &lw.UserInfo{}, status.Errorf(codes.FailedPrecondition, "failed to fetch user record %v", err)
	}
	return &lw.UserInfo{
		UserID: userID,
	}, nil
}

// dial gRPC services that are needed.
func gRPCClients(ctx context.Context) {
	select {
	case <-ctx.Done():
		return
	default:
		inRPC := &RPCClients{
			services: make(map[string]*grpc.ClientConn),
		}

		if viper.IsSet("gRPC.services.external") {
			// snippet for loading yaml array
			services := make([]map[string]string, 0)
			var m map[string]string
			servicesI := viper.Get("gRPC.services.external")
			servicesS := servicesI.([]interface{})
			for _, s := range servicesS {
				serviceMap := s.(map[interface{}]interface{})
				m = make(map[string]string)
				for k, v := range serviceMap {
					m[k.(string)] = v.(string)
				}
				services = append(services, m)
			}

			for _, v := range services {
				address := fmt.Sprintf("%v:%v", v["host"], v["port"])
				conn, err := grpc.Dial(address, grpc.WithInsecure())
				if err != nil {
					log.Errorf("could not connect service %v : %v", v["name"], err)
					os.Exit(1)
				}
				log.Infof("[gRPC] client connection: %v@%v:%v", v["name"], v["host"], v["port"])
				inRPC.services[v["name"]] = conn
				go statusConn(ctx, v, conn)
			}
			grpcc = inRPC
		}
	}
}

// listen on gRPC TCP connections related to this project
// not needed for now, as login is not expecting to act as server
func selfRPC(ctx context.Context) {
	select {
	case <-ctx.Done():
		return
	default:
		if viper.IsSet("gRPC.services.self") {
			// snippet for loading yaml array
			services := make([]map[string]string, 0)
			var m map[string]string
			servicesI := viper.Get("gRPC.services.self")
			servicesS := servicesI.([]interface{})
			for _, s := range servicesS {
				serviceMap := s.(map[interface{}]interface{})
				m = make(map[string]string)
				for k, v := range serviceMap {
					m[k.(string)] = v.(string)
				}
				services = append(services, m)
			}
			for _, v := range services {
				go gRPCServers(ctx, v)
			}
		}
	}
}

func gRPCServers(ctx context.Context, service map[string]string) {
	select {
	case <-ctx.Done():
		return
	default:
		address := fmt.Sprintf(":%v", service["port"])
		lis, err := net.Listen("tcp", address)
		if err != nil {
			log.Errorf("could listen on port %v for service %v : %v", service["port"], service["name"], err)
		}
		s := grpc.NewServer()

		lw.RegisterLoginServer(s, &server{})

		log.Infof("Loading gRPC server connections %v@::%v", service["name"], service["port"])

		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}
}

func statusConn(ctx context.Context, service map[string]string, conn *grpc.ClientConn) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			time.Sleep(120 * time.Second)
			log.Infof("[gRPC] [client] [%v]: %v@%v:%v ", conn.GetState(), service["name"], service["host"], service["port"])
		}
	}
}

func cleanupRPC() {
	grpcc.mu.Lock()
	for _, s := range grpcc.services {
		if err := s.Close(); err != nil {
			log.Error(err)
		} else {
			log.Info("Closing down external gRPC connection")
		}
	}
	grpcc.mu.Unlock()
}
