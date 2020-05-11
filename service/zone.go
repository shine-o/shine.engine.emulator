package service

import (
	"context"
	"errors"
	zm "github.com/shine-o/shine.engine.core/grpc/zone-master"
	"github.com/spf13/viper"
)

type runningMaps map[int]*zoneMap

type zone struct {
	rm      runningMaps
	queries recvEvents
	send    sendEvents
	recv    recvEvents
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
		events := []eventIndex{registerPlayerHandle, handleCleanUp,queryPlayer, queryMonster, playerAppeared, playerDisappeared, playerJumped, playerMoved, playerStopped}

		for _, index := range events {
			c := make(chan event, 5)
			zoneMaps[i].recv[index] = c
			zoneMaps[i].send[index] = c
		}

		rm[m.data.ID] = &zoneMaps[i]

		go zoneMaps[i].run()
	}

	events := []eventIndex{clientSHN, loadPlayerData, queryMap}
	z.recv = make(recvEvents)
	z.send = make(sendEvents)

	for _, index := range events {
		c := make(chan event, 5)
		z.recv[index] = c
		z.send[index] = c
	}

	zoneEvents = z.send

	err := registerZone(registerMaps)
	if err != nil {
		// close all event channels
		for _, m := range zoneMaps {
			for _, e := range m.send {
				close(e)
			}
		}
		for _, e := range events {
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

func (z *zone) loadQueries() sendEvents {
	queries := make(sendEvents)
	z.loadMapQueries(queries)
	return queries
}

func (z *zone) loadMapQueries(queries sendEvents) {
	loadMapQueries := []eventIndex{queryMap}
	for _, index := range loadMapQueries {
		c := make(chan event, 5)
		queries[index] = c
		z.queries[index] = c
	}
}
