package worldmaster

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	wm "github.com/shine-o/shine.engine.emulator/internal/pkg/grpc/world-master"
	shinelog "github.com/shine-o/shine.engine.emulator/pkg/log"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"net"
	"net/http"
)

var log = shinelog.NewLogger("world-master", "./output", logrus.DebugLevel)

// Start initializes an intermediary service for the diverse world services to connect to and acknowledge their status
func Start(cmd *cobra.Command, args []string) {
	go func() {
		enabled := viper.GetBool("metrics.enabled")
		if enabled {
			port := viper.GetString("metrics.prometheus.port")
			log.Infof("metrics enabled at :%v/metrics", port)
			http.Handle("/metrics", promhttp.Handler())
			log.Info(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
		}
	}()
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
