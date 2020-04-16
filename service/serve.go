package service

import (
	"context"
	"encoding/hex"
	"errors"
	"github.com/google/logger"
	zm "github.com/shine-o/shine.engine.core/grpc/zone-master"
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

func init() {
	log = logger.Init("zone master logger", true, false, ioutil.Discard)
}

func Start(cmd *cobra.Command, args []string) {
	ctx := context.Background()

	log = logger.Init("world logger", true, false, ioutil.Discard)

	zonePort := viper.GetString("serve.port")

	log.Infof("starting the service on port: %v", zonePort)

	// register against the zone master
	err := registerToZone()
	if err != nil {
		log.Fatal(err)
	}

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

	ss.Listen(ctx, zonePort)
}

func registerToZone() error {
	conn, err := newRPCClient("master")

	if err != nil {
		return err
	}
	c := zm.NewMasterClient(conn)
	rpcCtx, _ := context.WithTimeout(context.Background(), gRPCTimeout)

	zr, err := c.RegisterZone(rpcCtx, &zm.ZoneDetails{
		Maps: viper.GetStringSlice("maps"),
		Conn: &zm.ConnectionInfo{
			IP:   viper.GetString("serve.external_ip"),
			Port: viper.GetInt32("serve.port"),
		},
	})

	if err != nil {
		return err
	}
	if !zr.Success {
		return errors.New("failed to register against the zone master")
	}

	viper.SetDefault("world.ip", zr.World.IP)
	viper.SetDefault("world.port", zr.World.Port)

	return nil
}
