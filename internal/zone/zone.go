package zone

import (
	"context"
	"errors"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/data"
	zm "github.com/shine-o/shine.engine.emulator/internal/pkg/grpc/zone-master"
	"github.com/spf13/viper"
	"sync"
)

type runningMaps map[int]*zoneMap

type zone struct {
	rm runningMaps
	*events
	*dynamicEvents
	//worldDB *pg.DB
	sync.RWMutex
	*handler
}

var (
	zoneEvents  sendEvents
	monsterData *data.MonsterData
	mapData     *data.MapData
	npcData     *data.NpcData
	itemsData   *data.ItemData
)

func (z *zone) load() {

	shinePath := viper.GetString("shine_folder")

	loadGameData(shinePath)

	z.rm = make(runningMaps)

	zEvents := []eventIndex{
		playerMapLogin, playerSHN, playerData,
		heartbeatUpdate, queryMap, playerLogoutStart,
		playerLogoutCancel, playerLogoutConclude, persistPlayerPosition,
		changeMap,
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

	h := &handler{
		handleIndex: 0,
		usedHandles: make(map[uint16]bool),
	}

	z.handler = h

	normalMaps := viper.GetIntSlice("normal_maps")

	var registerMaps []int32

	var wg sync.WaitGroup
	var sem = make(chan int, 10)
	for _, id := range normalMaps {
		wg.Add(1)
		sem <- 1
		registerMaps = append(registerMaps, int32(id))

		go func(id int) {
			defer wg.Done()
			z.addMap(id)
		}(id)

		<-sem
	}

	wg.Wait()

	err := registerZone(registerMaps)

	if err != nil {
		log.Fatal(err)
	}

	go z.run()
}

func (z *zone) addMap(mapId int) {
	md, ok := mapData.Maps[mapId]

	if !ok {
		log.Fatalf("no map data for map with id %v", mapId)
	}

	walkableX, walkableY, err := walkingPositions(md.SHBD)

	if err != nil {
		log.Fatal(err)
	}

	for _, m := range z.rm {
		if m.data.MapInfoIndex == md.Info.MapName.Name {
			log.Errorf("duplicate shn map index id %v %v, skipping", mapId, m.data.MapInfoIndex)
			return
		}
	}

	m := &zoneMap{
		data:      md,
		walkableX: walkableX,
		walkableY: walkableY,
		entities: &entities{
			players: &players{
				handler: z.handler,
				active:  make(map[uint16]*player),
			},
			monsters: &monsters{
				handler: z.handler,
				active:  make(map[uint16]*monster),
			},
			npcs: &npcs{
				handler: z.handler,
				active:  make(map[uint16]*npc),
			},
		},
		events: events{
			send: make(sendEvents),
			recv: make(recvEvents),
		},
		metrics: metrics{
			players: promauto.NewGauge(prometheus.GaugeOpts{
				Name: fmt.Sprintf("players_in_%v", md.Info.MapName.Name),
				Help: "Total number of active players.",
			}),
			monsters: promauto.NewGauge(prometheus.GaugeOpts{
				Name: fmt.Sprintf("monsters_in_%v", md.Info.MapName.Name),
				Help: "Total number of active monsters.",
			}),
			npcs: promauto.NewGauge(prometheus.GaugeOpts{
				Name: fmt.Sprintf("npcs_in_%v", md.Info.MapName.Name),
				Help: "Total number of active non player characters.",
			}),
		},
	}

	m.metrics.monsters.Set(0)
	m.metrics.players.Set(0)
	m.metrics.npcs.Set(0)

	events := []eventIndex{
		playerHandle,
		playerHandleMaintenance,
		queryPlayer, queryMonster,
		playerAppeared, playerDisappeared, playerJumped, playerWalks, playerRuns, playerStopped,
		unknownHandle, monsterAppeared, monsterDisappeared, monsterWalks, monsterRuns,
		playerSelectsEntity, playerUnselectsEntity, playerClicksOnNpc, playerPromptReply,
	}

	for _, index := range events {
		c := make(chan event, 500)
		m.recv[index] = c
		m.send[index] = c
	}

	z.Lock()
	z.rm[m.data.ID] = m
	z.Unlock()

	go m.run()

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

func loadGameData(filesPath string) {

	var wg sync.WaitGroup

	wg.Add(4)
	go func(path string) {
		defer wg.Done()
		md, err := data.LoadMonsterData(path)
		if err != nil {
			log.Fatal(err)
		}
		monsterData = md
	}(filesPath)

	go func(path string) {
		defer wg.Done()
		md, err := data.LoadMapData(path)
		if err != nil {
			log.Fatal(err)
		}
		mapData = md
	}(filesPath)

	go func(path string) {
		defer wg.Done()
		nd, err := data.LoadNPCData(path)
		if err != nil {
			log.Fatal(err)
		}
		npcData = nd
	}(filesPath)

	go func(path string) {
		defer wg.Done()
		id, err := data.LoadItemData(path)
		if err != nil {
			log.Fatal(err)
		}
		itemsData = id
	}(filesPath)

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
