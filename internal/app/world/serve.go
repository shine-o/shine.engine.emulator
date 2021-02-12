package world

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/database"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/networking"
	shinelog "github.com/shine-o/shine.engine.emulator/pkg/log"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net/http"
	"path/filepath"
)

var log = shinelog.NewLogger("world", "./output", logrus.DebugLevel)

// Start the service service
// that is, use networking library to handle TCP connection
// configure networking library to use handlers implemented in this package for packets
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

	//sh := networking.ShinePacketRegistry{
	//	2055:  ncMiscSeedAck,
	//	3087:  ncUserLoginWorldReq,
	//	2061:  ncMiscGameTimeReq,
	//	3123:  ncUserWillWorldSelectReq,
	//	5121:  ncAvatarCreateReq,
	//	5127:  ncAvatarEraseReq,
	//	4097:  ncCharLoginReq,
	//	28684: ncCharOptionGetWindowPosReq,
	//	28676: ncCharOptionGetShortcutSizeReq,
	//	31750: ncPrisonGetReq,
	//	28727: ncCharOptionImproveSetShortcutDataReq,
	//	3103:  ncUserAvatarListReq,
	//}
	sh := networking.ShinePacketRegistry{
		networking.NC_MISC_SEED_ACK: networking.ShinePacket{
			Handler:  ncMiscSeedAck,
		},
		//networking.NC_USER_LOGINWORLD_REQ: networking.ShinePacket{
		networking.NC_USER_LOGIN_DB: networking.ShinePacket{
			Handler: ncUserLoginWorldReq,
		},
		networking.NC_MISC_GAMETIME_REQ: networking.ShinePacket{
			Handler: ncMiscGameTimeReq,
		},
		networking.NC_USER_WILL_WORLD_SELECT_REQ: networking.ShinePacket{
			Handler: ncUserWillWorldSelectReq,
		},
		networking.NC_AVATAR_CREATE_REQ: networking.ShinePacket{
			Handler: ncAvatarCreateReq,
		},
		networking.NC_AVATAR_ERASE_REQ: networking.ShinePacket{
			Handler: ncAvatarEraseReq,
		},
		networking.NC_CHAR_LOGIN_REQ: networking.ShinePacket{
			Handler: ncCharLoginReq,
		},
		networking.NC_CHAR_OPTION_GET_WINDOWPOS_REQ: networking.ShinePacket{
			Handler: ncCharOptionGetWindowPosReq,
		},
		networking.NC_CHAR_OPTION_GET_SHORTCUTSIZE_REQ: networking.ShinePacket{
			Handler: ncCharOptionGetShortcutSizeReq,
		},
		networking.NC_PRISON_GET_REQ: networking.ShinePacket{
			Handler: ncPrisonGetReq,
		},
		networking.NC_CHAR_OPTION_IMPROVE_SET_SHORTCUTDATA_REQ: networking.ShinePacket{
			Handler: ncCharOptionImproveSetShortcutDataReq,
		},
		networking.NC_USER_AVATAR_LIST_REQ: networking.ShinePacket{
		//networking.NC_USER_WORLD_STATUS_ACK: networking.ShinePacket{
			Handler: ncUserAvatarListReq,
		},
	}
	ss := networking.ShineService{
		Name:         "world",
		Settings:     s,
		ShinePacketRegistry: sh,
		SessionFactory: sessionFactory{
			worldID: viper.GetInt("world.id"),
		},
	}

	ss.Listen(ctx, worldPort)
}
