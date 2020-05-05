package service

import (
	"context"
	"errors"
	zm "github.com/shine-o/shine.engine.core/grpc/zone-master"
	"github.com/spf13/viper"
)

type runningMaps map[int]*zoneMap

func loadZone() map[int]*zoneMap {
	var registerMaps []int32
	rm := make(runningMaps)
	zoneMaps := loadMaps()
	for i, m := range zoneMaps {
		registerMaps = append(registerMaps, int32(m.data.ID))

		events := make(map[uint32]chan event)

		events[playerAppeared] = make(chan event, 5)
		events[playerDisappeared] = make(chan event, 5)
		events[playerJumped] = make(chan event, 5)
		events[playerMoved] = make(chan event, 5)
		events[playerStopped] = make(chan event, 5)

		zoneMaps[i].recv = make(map[uint32]<-chan event)
		zoneMaps[i].send = make(map[uint32]chan<- event)

		for k, v := range events {
			zoneMaps[i].recv[k] = v
			zoneMaps[i].send[k] = v
		}

		go zoneMaps[i].run()
		rm[m.data.ID] = &zoneMaps[i]
	}

	err := registerZone(registerMaps)
	if err != nil {
		// close all event channels
		for _, m := range zoneMaps {
			for _, e := range m.send {
				close(e)
			}
		}
		log.Fatal(err)
	}
	return rm
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
