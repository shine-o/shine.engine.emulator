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
	log = logger.Init("service service default logger", true, false, ioutil.Discard)
}

// Start the service service
// that is, use networking library to handle TCP connection
// configure networking library to use handlers implemented in this package for packets
func Start(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	initRedis()

	go newRPCServer("world")

	db = database.Connection(ctx, database.ConnectionParams{
		User:     viper.GetString("database.postgres.db_user"),
		Password: viper.GetString("database.postgres.db_password"),
		Host:     viper.GetString("database.postgres.host"),
		Port:     viper.GetString("database.postgres.port"),
		Database: viper.GetString("database.postgres.db_name"),
		Schema:   viper.GetString("database.postgres.schema"),
	})

	defer db.Close()
	log = logger.Init("world logger", true, false, ioutil.Discard)

	worldName := viper.GetString("world.name")
	worldPort := viper.GetString("world.port")

	log.Infof(" [%v] starting the service on port: %v", worldName, worldPort)

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

	sh := networking.ShineHandler{
		2055: NcMiscSeedAck,
		3087:  NcUserLoginWorldReq,
		2061:  NcMiscGameTimeReq,
		3123:  NcUserWillWorldSelectReq,
		5121:  NcAvatarCreateReq,
		5127:  NcAvatarEraseReq,
		4097:  NcCharLoginReq,
		28684: NcCharOptionGetWindowPosReq,
		28676: NcCharOptionGetShortcutSizeReq,
		31750: NcPrisonGetReq,
	}

	ss := networking.ShineService{
		Settings:        s,
		ShineHandler:    sh,
		SessionFactory:  sessionFactory{
			worldID: viper.GetInt("world.id"),
		},
	}
	
	ss.Listen(ctx, worldPort)
}
