package service

import (
	"context"
	"encoding/hex"
	"github.com/google/logger"
	"github.com/shine-o/shine.engine.core/networking"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"path/filepath"
)

var (
	log *logger.Logger
)

// Start the login service
// that is, use networking library to handle TCP connection
// configure networking library to use handlers implemented in this package for packets
func Start(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	log = logger.Init("LoginLogger", true, false, ioutil.Discard)
	log.Info("login logger init()")

	db = dbConn(ctx, "accounts")
	initRedis()

	go newRPCServer("login")

	defer db.Close()
	defer cancel()

	s := &networking.Settings{}

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

	// note: use factory
	ch := &networking.CommandHandlers{
		3173: userClientVersionCheckReq,
		3162: userUsLoginReq,
		3076: userXtrapReq,
		3099: userWorldStatusReq,
		3083: userWorldSelectReq,
		3096: userNormalLogoutCmd,
		3127: userLoginWithOtpReq,
	}

	hw := networking.NewHandlerWarden(ch)

	ss := networking.NewShineService(s, hw)
	wsf := &sessionFactory{}
	ss.UseSessionFactory(wsf)
	ss.Listen(ctx, viper.GetString("serve.port"))
}
