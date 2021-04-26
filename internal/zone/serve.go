package zone

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/database"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/networking"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/persistence"
	shinelog "github.com/shine-o/shine.engine.emulator/pkg/log"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net/http"
	"path/filepath"
)

var (
	log = shinelog.NewLogger("zone", "../../output", logrus.DebugLevel)
)

// Start initializes the TCP server and all the needed services and configuration for the zone
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

	zonePort := viper.GetString("serve.port")

	log.Infof("starting the service on port: %v", zonePort)

	z := zone{}

	z.load()
	go z.run()

	persistence.InitDB(database.ConnectionParams{
		User:     viper.GetString("world_database.db_user"),
		Password: viper.GetString("world_database.db_password"),
		Host:     viper.GetString("world_database.host"),
		Port:     viper.GetString("world_database.port"),
		Database: viper.GetString("world_database.db_name"),
		Schema:   viper.GetString("world_database.schema"),
	})

	defer persistence.CloseDB()

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
		Name:     "zone",
		Settings: s,
		ShinePacketRegistry: networking.ShinePacketRegistry{
			networking.NC_MISC_SEED_ACK: networking.ShinePacket{
				Handler: ncMiscSeedAck,
			},
			networking.NC_MISC_HEARTBEAT_ACK: networking.ShinePacket{
				Handler: ncMiscHeartBeatAck,
				//NcStruct: nil,
			},
			networking.NC_MAP_LOGIN_REQ: networking.ShinePacket{
				Handler: ncMapLoginReq,
			},
			networking.NC_MAP_LOGINCOMPLETE_CMD: networking.ShinePacket{
				Handler: ncMapLoginCompleteCmd,
			},
			networking.NC_CHAR_LOGOUTREADY_CMD: networking.ShinePacket{
				Handler: ncCharLogoutReadyCmd,
			},
			networking.NC_CHAR_LOGOUTCANCEL_CMD: networking.ShinePacket{
				Handler: ncCharLogoutCancelCmd,
			},
			networking.NC_ACT_MOVEWALK_CMD: networking.ShinePacket{
				Handler: ncActMoveWalkCmd,
			},
			networking.NC_ACT_MOVERUN_CMD: networking.ShinePacket{
				Handler: ncActMoveRunCmd,
			},
			networking.NC_ACT_JUMP_CMD: networking.ShinePacket{
				Handler: ncActJumpCmd,
			},
			networking.NC_ACT_STOP_REQ: networking.ShinePacket{
				Handler: ncActStopReq,
			},
			networking.NC_BRIEFINFO_INFORM_CMD: networking.ShinePacket{
				Handler: ncBriefInfoInformCmd,
			},
			networking.NC_BAT_TARGETTING_REQ: networking.ShinePacket{
				Handler: ncBatTargetingReq,
			},
			networking.NC_BAT_UNTARGET_REQ: networking.ShinePacket{
				Handler: ncBatUntargetReq,
			},
			networking.NC_USER_NORMALLOGOUT_CMD: networking.ShinePacket{
				Handler: ncUserNormalLogoutCmd,
			},
			networking.NC_ACT_NPCCLICK_CMD: networking.ShinePacket{
				Handler: ncActNpcClickCmd,
			},
			networking.NC_MENU_SERVERMENU_ACK: networking.ShinePacket{
				Handler: ncMenuServerMenuAck,
			},
			networking.NC_ITEM_RELOC_REQ: networking.ShinePacket{
				Handler: ncItemRelocReq,
			},
			networking.NC_ITEM_EQUIP_REQ: networking.ShinePacket{
				Handler: ncItemEquipReq,
			},
			networking.NC_ITEM_UNEQUIP_REQ: networking.ShinePacket{
				Handler: ncItemUnEquipReq,
			},
		},
		SessionFactory: sessionFactory{},
	}

	ss.Listen(ctx, zonePort)
}
