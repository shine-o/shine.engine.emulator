package service

import (
	"context"
	"encoding/hex"
	"github.com/google/logger"
	"github.com/shine-o/shine.engine.networking"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"path/filepath"
)

var (
	log   *logger.Logger
	grpcc *RPCClients
)

// Start the login service
// that is, use networking library to handle TCP connection
// configure networking library to use handlers implemented in this package for packets
func Start(cmd *cobra.Command, args []string) {
	log = logger.Init("LoginLogger", true, false, ioutil.Discard)
	log.Info("LoginLogger init()")
	initDatabase()
	initRedis()
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cleanupRPC()
	defer cancel()

	s := &networking.Settings{}

	if xk, err := hex.DecodeString(viper.GetString("crypt.xorKey")); err != nil {
		log.Error(err)
		os.Exit(1)
	} else {
		s.XorKey = xk
	}

	s.XorLimit = uint16(viper.GetInt("crypt.xorLimit"))

	if path, err := filepath.Abs(viper.GetString("protocol.nc-data")); err != nil {
		log.Error(err)
	} else {
		s.CommandsFilePath = path
	}

	// note: use factory
	ch := make(map[uint16]func(ctx context.Context, pc *networking.Command))
	ch[3173] = userClientVersionCheckReq
	ch[3162] = userUsLoginReq
	ch[3076] = userXtrapReq
	ch[3099] = userWorldStatusReq
	ch[3083] = userWorldSelectReq
	ch[3096] = userNormalLogoutCmd
	ch[3127] = userLoginWithOtpReq

	hw := networking.NewHandlerWarden(ch)

	ss := networking.NewShineService(s, hw)
	wsf := &sessionFactory{}
	ss.UseSessionFactory(wsf)
	gRPCClients(ctx)
	ss.Listen(ctx, viper.GetString("serve.port"))
}
