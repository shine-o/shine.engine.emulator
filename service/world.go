package service

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/google/logger"
	networking "github.com/shine-o/shine.engine.networking"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	lw "shine.engine.protocol-buffers/login-world"
	"sync"
)

var log *logger.Logger

func init() {
	log = logger.Init("WorldLogger", true, true, ioutil.Discard)
	log.Info("WorldLogger init()")
}

func Start(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	s := &networking.Settings{}

	if xk, err := hex.DecodeString(viper.GetString("crypt.xorKey")); err != nil {
		log.Error(err)
		os.Exit(1)
	} else {
		s.XorKey = xk
	}

	s.XorLimit = uint16(viper.GetInt("crypt.xorLimit"))

	if path, err := filepath.Abs(viper.GetString("protocol.nc-data")); err != nil {
		log.Error(err)
	} else {
		s.CommandsFilePath = path
	}

	ch := make(map[uint16]func(ctx context.Context, pc *networking.Command))

	// for now, nada

	hw := networking.NewHandlerWarden(ch)

	ss := networking.NewShineService(s, hw)

	go selfRPC(ctx)

	ss.Listen()
}

type InRPC struct {
	services map[string]*grpc.ClientConn
	mu       sync.Mutex
}

type OutRPC struct {
	services map[string]*grpc.ClientConn
	mu       sync.Mutex
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
				go gRpcServers(ctx, v)
			}
		}
	}
}

func gRpcServers(ctx context.Context, service map[string]string) {
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
		//switch service["name"] { in case more than one rpc service
		//case "world":
		//}
		lw.RegisterWorldServer(s, &server{})

		log.Infof("Loading gRPC server connections %v@::%v", service["name"], service["port"])

		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}
}
