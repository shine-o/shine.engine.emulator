package service

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/go-pg/pg/v9"
	"github.com/google/logger"
	"github.com/shine-o/shine.engine.networking"
	lw "github.com/shine-o/shine.engine.protocol-buffers/login-world"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"io/ioutil"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
)

type world struct {
	id    string
	name  string
	port  string
	extIP string
}

type activeWorlds struct {
	activeWorlds map[string]*world
	mu           sync.Mutex
}

var (
	log *logger.Logger
	aw  *activeWorlds
	grpcc *RPCClients
	worldDB  * pg.DB
)

func init() {
	log = logger.Init("world service default logger", true, false, ioutil.Discard)
}

// Start the world service
// that is, use networking library to handle TCP connection
// configure networking library to use handlers implemented in this package for packets
func Start(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)


	initRedis()
	gRPCClients(ctx)
	selfRPC(ctx)
	worldDB = dbConn(ctx, "world")

	defer cancel()
	defer cleanupRPC()
	defer worldDB.Close()
	// for each world in serve.worlds
	aw = &activeWorlds{
		activeWorlds: make(map[string]*world),
	}
	startWorlds(ctx)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	select {
	case <-c:
		cancel()
	}
	<-c
}

// reminder: remove this as the service will be responsible only for one server
func startWorlds(ctx context.Context) {
	if viper.IsSet("serve.worlds") {
		// snippet for loading yaml array
		services := make([]map[string]string, 0)
		var m map[string]string
		servicesI := viper.Get("serve.worlds")
		servicesS := servicesI.([]interface{})
		for _, s := range servicesS {
			serviceMap := s.(map[interface{}]interface{})
			m = make(map[string]string)
			for k, v := range serviceMap {
				m[k.(string)] = v.(string)
			}
			services = append(services, m)
		}

		for _, s := range services {
			w := world{
				id:    s["id"],
				name:  s["name"],
				port:  s["port"],
				extIP: s["external_ip"],
			}
			go startWorld(ctx, w)
		}
	}
}

func startWorld(ctx context.Context, w world) {
	select {
	case <-ctx.Done():
		return
	default:
		log = logger.Init(fmt.Sprintf("%v logger", w.name), true, false, ioutil.Discard)
		log.Infof(" [%v] starting the world on port: %v", w.name, w.port)

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
		ch[3087] = userLoginWorldReq
		ch[2061] = miscGametimeReq
		ch[3123] = userWillWorldSelectReq
		ch[5121] = avatarCreateReq

		hw := networking.NewHandlerWarden(ch)

		ss := networking.NewShineService(s, hw)

		wsf := &sessionFactory{
			worldID: w.id,
		}
		ss.UseSessionFactory(wsf)

		aw.mu.Lock()
		aw.activeWorlds[w.id] = &w
		aw.mu.Unlock()

		ss.Listen(ctx, w.port)
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

		lw.RegisterWorldServer(s, &server{})

		log.Infof("Loading gRPC server connections %v@::%v", service["name"], service["port"])

		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}
}
