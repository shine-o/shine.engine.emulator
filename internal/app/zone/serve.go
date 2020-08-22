package zone

import (
	"context"
	"encoding/hex"
	"github.com/go-pg/pg/v9"
	"github.com/google/logger"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/database"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/networking"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"path/filepath"
)

var (
	// grr todo: move these to the zone object, as all logic operations will have access to that object in some form
	log *logger.Logger
	db  *pg.DB
)

func init() {
	log = logger.Init("zone master logger", true, false, ioutil.Discard)
}

func Start(cmd *cobra.Command, args []string) {
	ctx := context.Background()

	log = logger.Init("world logger", true, false, ioutil.Discard)

	zonePort := viper.GetString("serve.port")

	log.Infof("starting the service on port: %v", zonePort)

	z := zone{}
	z.load()

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
		},
		SessionFactory: sessionFactory{},
	}

	ss.Listen(ctx, zonePort)
}
