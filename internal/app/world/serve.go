package world

import (
	"context"
	"encoding/hex"
	"github.com/google/logger"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/database"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/networking"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"path/filepath"
)

var (
	log *logger.Logger
)

func init() {
	log = logger.Init("world service default logger", true, false, ioutil.Discard)
}

// Start the service service
// that is, use networking library to handle TCP connection
// configure networking library to use handlers implemented in this package for packets
func Start(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	defer cancel()

	initRedis()

	w := world{}
	w.load()
	w.run()

	go newRPCServer("world")

	db := database.Connection(ctx, database.ConnectionParams{
		User:     viper.GetString("database.postgres.db_user"),
		Password: viper.GetString("database.postgres.db_password"),
		Host:     viper.GetString("database.postgres.host"),
		Port:     viper.GetString("database.postgres.port"),
		Database: viper.GetString("database.postgres.db_name"),
		Schema:   viper.GetString("database.postgres.schema"),
	})

	defer db.Close()

	w.db = db

	worldName := viper.GetString("world.name")
	worldPort := viper.GetString("world.port")
	xorKey := viper.GetString("crypt.xorKey")

	log.Infof(" [%v] starting the service on port: %v", worldName, worldPort)

	s := networking.Settings{}

	xk, err := hex.DecodeString(xorKey)

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

	sh := networking.ShineHandler{
		2055:  ncMiscSeedAck,
		3087:  ncUserLoginWorldReq,
		2061:  ncMiscGameTimeReq,
		3123:  ncUserWillWorldSelectReq,
		5121:  ncAvatarCreateReq,
		5127:  ncAvatarEraseReq,
		4097:  ncCharLoginReq,
		28684: ncCharOptionGetWindowPosReq,
		28676: ncCharOptionGetShortcutSizeReq,
		31750: ncPrisonGetReq,
		28727: ncCharOptionImproveSetShortcutDataReq,
	}

	ss := networking.ShineService{
		Settings:     s,
		ShineHandler: sh,
		SessionFactory: sessionFactory{
			worldID: viper.GetInt("world.id"),
		},
	}

	ss.Listen(ctx, worldPort)
}
