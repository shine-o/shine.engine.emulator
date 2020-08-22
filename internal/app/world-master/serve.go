package world_master

import (
	"fmt"
	"github.com/google/logger"
	wm "github.com/shine-o/shine.engine.emulator/internal/pkg/grpc/world-master"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"io/ioutil"
	"net"
)

var (
	log *logger.Logger
)

func init() {
	log = logger.Init("zone master logger", true, false, ioutil.Discard)
}

func Start(cmd *cobra.Command, args []string) {
	initRedis()
	port := viper.GetString("serve.port")
	address := fmt.Sprintf(":%v", port)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Errorf("could listen on port %v: %v", port, err)
	}
	s := grpc.NewServer()
	wm.RegisterMasterServer(s, &server{})

	log.Infof("Loading gRPC server connection master@::%v", port)
	if err := s.Serve(lis); err != nil {
		log.Errorf("failed to serve: %v", err)
	}
}
