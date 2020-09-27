package zone

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/database"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/networking"
	shinelog "github.com/shine-o/shine.engine.emulator/pkg/log"
	"github.com/sirupsen/logrus"
	//"github.com/google/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net/http"
	"path/filepath"
)


var log = shinelog.NewLogger("zone", "./output", logrus.DebugLevel)

// Start initializes the TCP server and all the needed services and configuration for the zone
func Start(cmd *cobra.Command, args []string) {
	//defer func() {
	//	if r := recover(); r != nil {
	//		log.Error(r)
	//	}
	//}()
	go func() {
		enabled := viper.GetBool("metrics.enabled")
		if enabled {
			port := viper.GetString("metrics.prometheus.port")
			log.Infof("metrics enabled at :%v/metrics", port)
			http.Handle("/metrics", promhttp.Handler())
			log.Info(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
		}
	}()
	ctx := context.Background()

	zonePort := viper.GetString("serve.port")

	log.Infof("starting the service on port: %v", zonePort)

	z := zone{}

	z.load()

	db := database.Connection(ctx, database.ConnectionParams{
		User:     viper.GetString("world_database.db_user"),
		Password: viper.GetString("world_database.db_password"),
		Host:     viper.GetString("world_database.host"),
		Port:     viper.GetString("world_database.port"),
		Database: viper.GetString("world_database.db_name"),
		Schema:   viper.GetString("world_database.schema"),
	})

	defer db.Close()

	z.worldDB = db

	s := networking.Settings{}

	xk, err := hex.DecodeString(viper.GetString("crypt.xorKey"))

	if err != nil {
		log.Fatal(err)
	}

	s.XorKey = xk

	s.XorLimit = uint16(viper.GetInt("crypt.xorLimit"))

	path, err := filepath.Abs(viper.GetString("protocol.commands"))

	if err != nil {
		log.Fatal(err)
	}

	s.CommandsFilePath = path

	ss := networking.ShineService{
		Name: "zone",
		Settings: s,
		ShineHandler: networking.ShineHandler{
			2055: ncMiscSeedAck,
			2053: ncMiscHeartBeatAck,
			6145: ncMapLoginReq,
			6147: ncMapLoginCompleteCmd,
			4209: ncCharLogoutReadyCmd,
			4210: ncCharLogoutCancelCmd,
			8215: ncActMoveWalkCmd,
			8217: ncActMoveRunCmd,
			8228: ncActJumpCmd,
			8210: ncActStopReq,
			7169: ncBriefInfoInformCmd,
			9217: ncBatTargetingReq,
		},
		SessionFactory: sessionFactory{},
	}

	ss.Listen(ctx, zonePort)
}
