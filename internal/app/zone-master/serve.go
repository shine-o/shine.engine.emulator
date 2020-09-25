package zonemaster

import (
	"fmt"
	"github.com/google/logger"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/grpc/zone-master"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"os"
)

var (
	log *logger.Logger
)

func init()  {
	if _, err := os.Stat("./output"); os.IsNotExist(err) {
		err := os.Mkdir("./output", 0660)
		if err != nil {
			logger.Fatalf("Failed to create output folder: %v", err)
		}
	}

	lf, err := os.OpenFile("./output/zone-master.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0660)
	if err != nil {
		logger.Fatalf("Failed to create output file: %v", err)
	}

	log = logger.Init("zone-master", true, false, lf)
}


// Start initializes an intermediary service for the diverse zone services to connect to and acknowledge their status
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

	zm.RegisterMasterServer(s, &server{})

	log.Infof("Loading gRPC server connection master@::%v", port)
	if err := s.Serve(lis); err != nil {
		log.Errorf("failed to serve: %v", err)
	}
}
