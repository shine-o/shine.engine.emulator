package zone

import (
	"context"
	"errors"
	zm "github.com/shine-o/shine.engine.core/grpc/zone-master"
	"github.com/spf13/viper"
)

type runningMaps map[int]*zoneMap

type zone struct {
	rm runningMaps
	// static events, i should add a wrapper type for these
	//send    sendEvents
	//recv    recvEvents
	events
	dynamic
}


// instead of accessing global variables for data
// fire a query event struct, which will be populated with the requested data by a worker (event receiver)
var zoneEvents sendEvents

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

func (z *zone) load() {
	var registerMaps []int32
	rm := make(runningMaps)
	zoneMaps := loadMaps()
	for i, m := range zoneMaps {
		registerMaps = append(registerMaps, int32(m.data.ID))
		events := []eventIndex{
			playerHandle,
			playerHandleMaintenance,
			queryPlayer, queryMonster,
			playerAppeared, playerDisappeared, playerJumped, playerWalks, playerRuns, playerStopped,
		}

		for _, index := range events {
			c := make(chan event, 5)
			zoneMaps[i].recv[index] = c
			zoneMaps[i].send[index] = c
		}

		rm[m.data.ID] = &zoneMaps[i]

		go zoneMaps[i].run()
	}

	zEvents := []eventIndex{
		playerSHN,
		playerData,
		queryMap,
		playerLogoutStart, playerLogoutCancel, playerLogoutConclude,
	}

	z.recv = make(recvEvents)
	z.send = make(sendEvents)

	for _, index := range zEvents {
		c := make(chan event, 5)
		z.recv[index] = c
		z.send[index] = c
	}

	zoneEvents = z.send

	z.dynamic= dynamic{
		events:  make(map[string]events),
	}

	err := registerZone(registerMaps)
	if err != nil {
		// close all event channels
		for i, _ := range zoneMaps {
			for j, _ := range zoneMaps[i].send {
				close(zoneMaps[i].send[j])
			}
		}
		for _, e := range zEvents {
			close(z.send[e])
		}
		log.Fatal(err)
	}
	z.rm = rm
	go z.run()
}

func (z *zone) run() {
	// run query workers
	go z.mapQueries()
	go z.security()
	go z.playerSession()
}
