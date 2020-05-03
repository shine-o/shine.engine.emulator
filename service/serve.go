package service

import (
	"context"
	"encoding/hex"
	"github.com/go-pg/pg/v9"
	"github.com/google/logger"
	"github.com/shine-o/shine.engine.core/database"
	"github.com/shine-o/shine.engine.core/networking"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"path/filepath"
)

var (
	log *logger.Logger
	db  *pg.DB
)

func init() {
	log = logger.Init("zone master logger", true, false, ioutil.Discard)
}

type zoneParameters struct {
	rm runningMaps
}

func Start(cmd *cobra.Command, args []string) {
	ctx := context.Background()

	log = logger.Init("world logger", true, false, ioutil.Discard)

	zonePort := viper.GetString("serve.port")

	log.Infof("starting the service on port: %v", zonePort)

	var rm runningMaps
	rm = loadZone()

	db = database.Connection(ctx, database.ConnectionParams{
		User:     viper.GetString("world_database.db_user"),
		Password: viper.GetString("world_database.db_password"),
		Host:     viper.GetString("world_database.host"),
		Port:     viper.GetString("world_database.port"),
		Database: viper.GetString("world_database.db_name"),
		Schema:   viper.GetString("world_database.schema"),
	})

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

	ss := networking.ShineService{
		Settings:        s,
		ShineHandler:    networking.ShineHandler{
			2055: NcMiscSeedAck,
			6145: NcMapLoginReq,
		},
		SessionFactory:  sessionFactory{},
		// here I should add maps loaded in this zone
		// that way when a command comes in, i can send events to the map the player is situated in
		// the map logic routines will handle the received event
		ExtraParameters: zoneParameters{
			rm: rm,
		},
	}

	ss.Listen(ctx, zonePort)
}

