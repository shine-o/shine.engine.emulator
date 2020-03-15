package service

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/google/logger"
	"github.com/shine-o/shine.engine.networking"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type RpcClients struct {
	services map[string]*grpc.ClientConn
	mu       sync.Mutex
}

var (
	log   *logger.Logger
	grpcc *RpcClients
)

func init() {
	log = logger.Init("LoginLogger", true, true, ioutil.Discard)
	log.Info("LoginLogger init()")
}

func Start(cmd *cobra.Command, args []string) {
	initDatabase()
	initRedis()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cleanupRPC()
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
	ch[3173] = userClientVersionCheckReq
	ch[3162] = userUsLoginReq
	ch[3076] = userXtrapReq
	ch[3099] = userWorldStatusReq
	ch[3083] = userWorldSelectReq
	ch[3096] = userNormalLogoutCmd

	hw := networking.NewHandlerWarden(ch)

	ss := networking.NewShineService(s, hw)
	wsf := &sessionFactory{}
	ss.UseSessionFactory(wsf)
	grpcc = gRpcClients(ctx)

	ss.Listen(ctx)
}

// dial gRPC services that are needed.
func gRpcClients(ctx context.Context) *RpcClients {
	select {
	case <-ctx.Done():
		return &RpcClients{}
	default:
		inRPC := &RpcClients{
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
				conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
				if err != nil {
					log.Error("could not connect service %v : %v", v["name"], err)
					os.Exit(1)
				}
				log.Infof("Loading gRPC client connections %v@%v:%v", v["name"], v["host"], v["port"])
				inRPC.services[v["name"]] = conn
				go statusConn(ctx, v, conn)
			}
		}
		return inRPC
	}
}

func statusConn(ctx context.Context, service map[string]string, conn *grpc.ClientConn) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			time.Sleep(15 * time.Second)
			log.Infof("[%v] gRPC client connection: %v@%v:%v ", conn.GetState(), service["name"], service["host"], service["port"])
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
