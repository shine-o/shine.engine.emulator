package service

import (
	"context"
	"errors"
	zm "github.com/shine-o/shine.engine.core/grpc/zone-master"
	"github.com/spf13/viper"
)

type runningMaps map[uint32]zoneMap

func loadZone() {
	var registerMaps []int32
	zoneMaps := loadMaps()
	for _, m := range zoneMaps {
		registerMaps = append(registerMaps, int32(m.data.ID))

		events := make(map[uint32]chan<- event)

		events[entityAppeared] = make(chan event)
		events[entityDisappeared] = make(chan event)
		events[entityJumped] = make(chan event)
		events[entityMoved] = make(chan event)
		events[entityStopped] = make(chan event)

		go m.run()
	}

	err := registerZone(registerMaps)

	if err != nil {
		log.Fatal(err)
	}
}

func registerZone(mapIDs []int32) error {
	zoneIP := viper.GetString("serve.external_ip")
	zonePort := viper.GetInt32("serve.port")

	conn, err := newRPCClient("zone_master")

	if err != nil {
		return err
	}
	c := zm.NewMasterClient(conn)
	rpcCtx, _ := context.WithTimeout(context.Background(), gRPCTimeout)

	zr, err := c.RegisterZone(rpcCtx, &zm.ZoneDetails{
		Maps: mapIDs,
		Conn: &zm.ConnectionInfo{
			IP:   zoneIP,
			Port: zonePort,
		},
	})

	if err != nil {
		return err
	}

	if !zr.Success {
		return errors.New("failed to register against the zone master")
	}
	return nil
}
