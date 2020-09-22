package zone

import (
	"context"
	"errors"
	"github.com/go-pg/pg/v9"
	mobs "github.com/shine-o/shine.engine.emulator/internal/pkg/game-data/monsters"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/game-data/world"
	zm "github.com/shine-o/shine.engine.emulator/internal/pkg/grpc/zone-master"
	"github.com/spf13/viper"
	"sync"
)

type runningMaps map[int]*zoneMap

type zone struct {
	rm runningMaps
	*events
	*dynamicEvents
	worldDB *pg.DB
	sync.RWMutex
}

// instead of accessing global variables for data
// fire a query event struct, which will be populated with the requested data by a worker (event receiver)
var (
	zoneEvents  sendEvents
	monsterData mobs.MonsterData
	mapData     map[int]*world.Map
)

func (z *zone) load() {
	loadGameData()

	z.rm = make(runningMaps)

	zEvents := []eventIndex{
		playerMapLogin,
		playerSHN,
		playerData,
		heartbeatUpdate,
		queryMap,
		playerLogoutStart, playerLogoutCancel, playerLogoutConclude,
		persistPlayerPosition,
	}

	z.events = &events{
		send: make(sendEvents),
		recv: make(recvEvents),
	}

	for _, index := range zEvents {
		c := make(chan event, 500)
		z.events.recv[index] = c
		z.events.send[index] = c
	}

	zoneEvents = z.events.send

	z.dynamicEvents = &dynamicEvents{
		events: make(map[string]events),
	}

	normalMaps := viper.GetIntSlice("normal_maps")

	var registerMaps []int32

	var wg sync.WaitGroup
	for _, id := range normalMaps {
		wg.Add(1)
		registerMaps = append(registerMaps, int32(id))
		go z.addMap(id, &wg)
	}

	wg.Wait()

	err := registerZone(registerMaps)

	if err != nil {
		log.Fatal(err)
	}

	go z.run()
}

func loadGameData() {
	shinePath := viper.GetString("shine_folder")
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		md, err := mobs.LoadMonsterData(shinePath)
		if err != nil {
			log.Fatal(err)
		}
		monsterData = md
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		md, err := world.LoadMapData(shinePath)
		if err != nil {
			log.Fatal(err)
		}
		mapData = md
	}()

	wg.Wait()
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
