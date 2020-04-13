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

var log *logger.Logger

func Start(cmd *cobra.Command, args []string) {
	ctx := context.Background()

	log = logger.Init("world logger", true, false, ioutil.Discard)

	zonePort := viper.GetString("serve.port")

	log.Infof("starting the service on port: %v", zonePort)

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

	ch := &networking.CommandHandlers{}

	hw := networking.NewHandlerWarden(ch)

	ss := networking.NewShineService(s, hw)

	//wsf := &sessionFactory{
	//	worldID: viper.GetInt("world.id"),
	//}
	//ss.UseSessionFactory(wsf)

	ss.Listen(ctx, zonePort)
}
