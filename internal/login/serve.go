package login

import (
	"context"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/database"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/networking"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/persistence"
	shinelog "github.com/shine-o/shine.engine.emulator/pkg/log"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var log = shinelog.NewLogger("login", "./output", logrus.DebugLevel)

func metrics() {
	enabled := viper.GetBool("metrics.enabled")
	if enabled {
		port := viper.GetString("metrics.prometheus.port")
		log.Infof("metrics enabled at :%v/metrics", port)
		http.Handle("/metrics", promhttp.Handler())
		log.Info(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
	}
}

// Start the login service
// that is, use networking library to handle TCP connection
// configure networking library to use handlers implemented in this package for packets
func Start(cmd *cobra.Command, args []string) {
	go metrics()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	initRedis()

	go newRPCServer("login")

	defer cancel()

	persistence.InitDB(database.ConnectionParams{
		User:     viper.GetString("database.postgres.db_user"),
		Password: viper.GetString("database.postgres.db_password"),
		Host:     viper.GetString("database.postgres.host"),
		Port:     viper.GetString("database.postgres.port"),
		Database: viper.GetString("database.postgres.db_name"),
		Schema:   viper.GetString("database.postgres.schema"),
	})

	defer persistence.CloseDB()

	s := networking.Settings{}

	if xk, err := hex.DecodeString(viper.GetString("crypt.xorKey")); err != nil {
		log.Error(err)
		os.Exit(1)
	} else {
		s.XorKey = xk
	}

	s.XorLimit = uint16(viper.GetInt("crypt.xorLimit"))

	if path, err := filepath.Abs(viper.GetString("protocol.commands")); err != nil {
		log.Error(err)
	} else {
		s.CommandsFilePath = path
	}

	l := login{}
	l.load()

	// note: use factory

	ss := networking.ShineService{
		Name:     "login",
		Settings: s,
		ShinePacketRegistry: networking.ShinePacketRegistry{
			networking.NC_MISC_SEED_ACK: {
				Handler: ncMiscSeedAck,
			},
			networking.NC_USER_CLIENT_VERSION_CHECK_REQ: {
				Handler: ncUserClientVersionCheckReq,
			},
			networking.NC_USER_US_LOGIN_REQ: {
				Handler: ncUserUsLoginReq,
			},
			networking.NC_USER_WORLD_STATUS_REQ: {
				Handler: ncUserWorldStatusReq,
			},
			networking.NC_USER_WORLDSELECT_REQ: {
				// networking.NC_USER_WORLDSELECT_ACK: {
				Handler: ncUserWorldSelectReq,
			},
			networking.NC_USER_NORMALLOGOUT_CMD: {
				Handler: ncUserNormalLogoutCmd,
			},
			// this no longer works
			networking.NC_USER_LOGIN_WITH_OTP_REQ: {
				Handler: ncUserLoginWithOtpReq,
			},
		},
		SessionFactory: sessionFactory{},
	}

	ss.Listen(ctx, viper.GetString("serve.port"))
}
