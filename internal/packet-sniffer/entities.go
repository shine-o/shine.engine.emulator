package packet_sniffer

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/structs"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var em EntitiesMovements

type EntitiesMovements struct {
	Entities map[uint16][]Movement
	sync.Mutex
}

type Movement struct {
	Timestamp time.Time
	X, Y      uint32
}

type PacketData struct {
	// opcode => hex string
	Packets map[uint16][]string
	sync.Mutex
}

var pd PacketData

func persistPacketData(dp decodedPacket) {
	ds := hex.EncodeToString(dp.packet.Base.Data)
	pd.Lock()
	pd.Packets[dp.packet.Base.OperationCode] = append(pd.Packets[dp.packet.Base.OperationCode], ds)
	pd.Unlock()
}

// store info of packets that contain coordinates
func persistMovement(dp decodedPacket) {
	switch dp.packet.Base.OperationCode {
	// server
	// has handle identifier.
	// TODO: fetch the packets that have info about surrounding entities
	case 8211:
		// NC_ACT_SOMEONESTOP_CMD
		nc := structs.NcActSomeoneStopCmd{}
		err := structs.Unpack(dp.packet.Base.Data, &nc)
		if err != nil {
			log.Error(err)
		}
		em.Lock()
		em.Entities[nc.Handle] = append(em.Entities[nc.Handle], Movement{
			Timestamp: time.Now(),
			X:         nc.Location.X,
			Y:         nc.Location.Y,
		})
		em.Unlock()
	case 8216:
		// NC_ACT_SOMEONEMOVEWALK_CMD
	case 8218:
		// NC_ACT_SOMEONEMOVERUN_CMD
		nc := structs.NcActSomeoneMoveRunCmd{}
		err := structs.Unpack(dp.packet.Base.Data, &nc)
		if err != nil {
			log.Error(err)
		}
		em.Lock()
		em.Entities[nc.Handle] = append(em.Entities[nc.Handle], Movement{
			Timestamp: time.Now(),
			X:         nc.To.X,
			Y:         nc.To.Y,
		})
		em.Unlock()

	// client
	// only register the player moving
	case 8215:
		// NC_ACT_MOVEWALK_CMD
		nc := structs.NcActMoveWalkCmd{}
		err := structs.Unpack(dp.packet.Base.Data, &nc)
		if err != nil {
			log.Error(err)
		}
		em.Lock()
		em.Entities[1] = append(em.Entities[1], Movement{
			Timestamp: time.Now(),
			X:         nc.To.X,
			Y:         nc.To.Y,
		})
		em.Unlock()
	case 8217:
		// NC_ACT_MOVERUN_CMD
		nc := structs.NcActMoveRunCmd{}
		err := structs.Unpack(dp.packet.Base.Data, &nc)
		if err != nil {
			log.Error(err)
		}
		em.Lock()
		em.Entities[1] = append(em.Entities[1], Movement{
			Timestamp: time.Now(),
			X:         nc.To.X,
			Y:         nc.To.Y,
		})
		em.Unlock()

	// map enter
	case 7175:
		// NC_BRIEFINFO_CHARACTER_CMD
		nc := structs.NcBriefInfoCharacterCmd{}
		err := structs.Unpack(dp.packet.Base.Data, &nc)
		if err != nil {
			log.Error(err)
		}
		for _, c := range nc.Characters {
			em.Lock()
			em.Entities[c.Handle] = append(em.Entities[c.Handle], Movement{
				Timestamp: time.Now(),
				X:         c.Coordinates.XY.X,
				Y:         c.Coordinates.XY.Y,
			})
			em.Unlock()
		}
	case 7177:
		// NC_BRIEFINFO_MOB_CMD
		nc := structs.NcBriefInfoMobCmd{}
		err := structs.Unpack(dp.packet.Base.Data, &nc)
		if err != nil {
			log.Error(err)
		}
		for _, m := range nc.Mobs {
			em.Lock()
			em.Entities[m.Handle] = append(em.Entities[m.Handle], Movement{
				Timestamp: time.Now(),
				X:         m.Coord.XY.X,
				Y:         m.Coord.XY.Y,
			})
			em.Unlock()
		}
	case 7194:
		// NC_BRIEFINFO_REGENMOVER_CMD
		nc := structs.NcBriefInfoRegenMoverCmd{}
		err := structs.Unpack(dp.packet.Base.Data, &nc)
		if err != nil {
			log.Error(err)
		}
		em.Lock()
		em.Entities[nc.Handle] = append(em.Entities[nc.Handle], Movement{
			Timestamp: time.Now(),
			X:         nc.Coordinates.XY.X,
			Y:         nc.Coordinates.XY.Y,
		})
		em.Unlock()
	}
}

func exportPackets() {
	pathName, err := filepath.Abs(fmt.Sprintf("output/packets-%v.json", time.Now().Unix()))
	if err != nil {
		log.Fatal(err)
	}
	f, err := os.OpenFile(pathName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	//_,_ = f.Write([]byte("{"))

	pd.Lock()
	b, err := json.Marshal(pd.Packets)
	if err != nil {
		log.Error(err)
	}
	_, _ = f.Write(b)
	pd.Unlock()

	f.Close()
}

func exportEntitiesMovements() {
	log.Info("printing entity movements")
	pathName, err := filepath.Abs("output/movements.json")
	if err != nil {
		log.Fatal(err)
	}
	f, err := os.OpenFile(pathName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}

	//_,_ = f.Write([]byte("{"))

	em.Lock()
	b, err := json.Marshal(em.Entities)
	if err != nil {
		log.Error(err)
	}
	_, _ = f.Write(b)
	em.Unlock()

	f.Close()
}
