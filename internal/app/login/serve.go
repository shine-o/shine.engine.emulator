package login

import (
	"context"
	"encoding/hex"
	"github.com/google/logger"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/networking"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
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

	lf, err := os.OpenFile("./output/login.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0660)
	if err != nil {
		logger.Fatalf("Failed to create output file: %v", err)
	}

	log = logger.Init("login", true, false, lf)
}

// Start the login service
// that is, use networking library to handle TCP connection
// configure networking library to use handlers implemented in this package for packets
func Start(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	db = dbConn(ctx, "accounts")
	initRedis()

	go newRPCServer("login")

	defer db.Close()
	defer cancel()

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
		Settings: s,
		ShineHandler: networking.ShineHandler{
			2055: ncMiscSeedAck,
			3173: ncUserClientVersionCheckReq,
			3162: ncUserUsLoginReq,
			3076: ncUserXtrapReq,
			3099: ncUserWorldStatusReq,
			3083: ncUserWorldSelectReq,
			3096: ncUserNormalLogoutCmd,
			3127: ncUserLoginWithOtpReq,
		},
		SessionFactory: sessionFactory{},
	}

	ss.Listen(ctx, viper.GetString("serve.port"))
}
