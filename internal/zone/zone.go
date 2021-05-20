package zone

import (
	"context"
	"errors"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/data"
	zm "github.com/shine-o/shine.engine.emulator/internal/pkg/grpc/zone-master"
	"github.com/spf13/viper"
	"sync"
)

type zone struct {
	rm            *runningMaps
	events        *events
	dynamicEvents *dynamicEvents
	sync.RWMutex
}

type runningMaps struct {
	list map[int]*zoneMap
	sync.RWMutex
}

var (
	monsterData    *data.MonsterData
	mapData        *data.MapData
	npcData        *data.NpcData
	itemsData      *data.ItemData
	zoneEvents     sendEvents
	maps           *runningMaps
	handlerManager = &handler{
		index: 0,
		inUse: make(map[uint16]bool),
	}
	newHandler    = make(chan *handlerPetition, 1500)
	removeHandler = make(chan *handlerPetition, 1500)
	queryHandler  = make(chan *handlerPetition, 1500)
)

func init() {
	go handlerManager.handleWorker()
}

func (z *zone) load() error {
	loadGameData()

	z.rm = &runningMaps{
		list: make(map[int]*zoneMap),
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

	var wg sync.WaitGroup
	var sem = make(chan int, 10)
	for _, id := range normalMaps {
		wg.Add(1)
		sem <- 1

		go func(id int) {
			defer wg.Done()
			z.addMap(id)
			<-sem
		}(id)
	}

	wg.Wait()

	maps = z.rm

	return nil
}

func (z *zone) addMap(mapID int) {
	m, err := loadMap(mapID)
	if err != nil {
		log.Error(err)
		return
	}
	z.rm.Lock()
	z.rm.list[m.data.ID] = m
	z.rm.Unlock()

	go m.run()
}

func (z *zone) run() {
	num := viper.GetInt("workers.num_zone_workers")
	for i := 0; i <= num; i++ {
		go z.security()
		go z.playerSession()
		go z.playerGameData()
	}
}

func (z *zone) allMaps() <-chan *zoneMap {
	z.rm.RLock()
	ch := make(chan *zoneMap, len(z.rm.list))
	z.rm.RUnlock()

	go func(rm *runningMaps, send chan<- *zoneMap) {
		rm.RLock()
		for _, rm := range rm.list {
			send <- rm
		}
		rm.RUnlock()
		close(send)
	}(z.rm, ch)

	return ch
}

func loadGameData() {
	filesPath := viper.GetString("shine_folder")

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
