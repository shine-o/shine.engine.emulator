package service

import (
	"context"
	"fmt"
	lw "github.com/shine-o/shine.engine.protocol-buffers/login-world"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"os"
	"strconv"
	"sync"
	"time"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	lw.UnimplementedWorldServer
}

const gRPCTimeout = time.Second * 2

// There can be many Worlds, each World with its own Zones
// for simplicity we use hardcoded one for now :)
func (s *server) AvailableWorlds(ctx context.Context, in *lw.ClientMetadata) (*lw.WorldsInfo, error) {
	select {
	case <-ctx.Done():
		return &lw.WorldsInfo{}, status.Errorf(codes.Canceled, "context was canceled")
	default:
		var law activeWorlds
		aw.mu.Lock()
		law = *aw
		aw.mu.Unlock()

		worlds := make([]*lw.WorldInfo, 0)
		for _, w := range law.activeWorlds {
			id, err := strconv.Atoi(w.id)
			if err != nil {
				log.Error(err)
				continue
			}

			worlds = append(worlds, &lw.WorldInfo{
				WorldNumber: int32(id),
				WorldName:   w.name,
				WorldStatus: 6,
			})
		}

		return &lw.WorldsInfo{
			Worlds: worlds,
		}, nil
	}
}

func (s *server) ConnectionInfo(ctx context.Context, req *lw.SelectedWorld) (*lw.WorldConnectionInfo, error) {
	select {
	case <-ctx.Done():
		return nil, status.Errorf(codes.Canceled, "context was canceled")
	default:
		aw.mu.Lock()
		inx := fmt.Sprintf("%v", req.Num)
		w, available := aw.activeWorlds[inx]
		aw.mu.Unlock()

		if !available {
			log.Error(aw)
			err := status.Errorf(codes.Unavailable, "requested world with id %v is not available", req.Num)
			log.Error(err)
			return nil, err
		}

		port, err := strconv.Atoi(w.port)
		if err != nil {
			log.Error(err)
			return nil, status.Errorf(codes.FailedPrecondition, "incorrect world port %v", port)
		}
		return &lw.WorldConnectionInfo{
			Name: w.name,
			IP:   w.extIP,
			Port: int32(port),
		}, nil
	}
}

// RPCClients that this service will use to communicate with other services
type RPCClients struct {
	services map[string]*grpc.ClientConn
	mu       sync.Mutex
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
				//conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
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
