package zone

import (
	"context"
	"errors"
	"github.com/go-pg/pg/v9"
	zm "github.com/shine-o/shine.engine.emulator/internal/pkg/grpc/zone-master"
	"github.com/spf13/viper"
)

type runningMaps map[int]*zoneMap

type zone struct {
	rm runningMaps
	events
	*dynamicEvents
	worldDB *pg.DB
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
			unknownHandle,
		}

		for _, index := range events {
			c := make(chan event, 500)
			zoneMaps[i].recv[index] = c
			zoneMaps[i].send[index] = c
		}

		rm[m.data.ID] = &zoneMaps[i]

		go zoneMaps[i].run()
	}

	zEvents := []eventIndex{
		playerMapLogin,
		playerSHN,
		playerData,
		heartbeatUpdate,
		queryMap,
		playerLogoutStart, playerLogoutCancel, playerLogoutConclude,
		persistPlayerPosition,
	}

	z.recv = make(recvEvents)
	z.send = make(sendEvents)

	for _, index := range zEvents {
		c := make(chan event, 5)
		z.recv[index] = c
		z.send[index] = c
	}

	zoneEvents = z.send

	z.dynamicEvents = &dynamicEvents{
		events: make(map[string]events),
	}

	err := registerZone(registerMaps)
	if err != nil {
		// close all event channels
		for i := range zoneMaps {
			for j := range zoneMaps[i].send {
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
	num := viper.GetInt("workers.num_zone_workers")
	for i := 0; i <= num; i++ {
		go z.mapQueries()
		go z.security()
		go z.playerSession()
		go z.playerGameData()
	}
}
