package service

import (
	"context"
	"errors"
	zm "github.com/shine-o/shine.engine.core/grpc/zone-master"
	"github.com/spf13/viper"
)

func registerZone() error {
	conn, err := newRPCClient("zone_master")

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
