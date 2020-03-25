package service

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"os"
	"sync"
	"time"
)

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

