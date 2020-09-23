package zone

import (
	mobs "github.com/shine-o/shine.engine.emulator/internal/pkg/game-data/monsters"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/game-data/shn"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/networking"
	"github.com/shine-o/shine.engine.emulator/pkg/structs"
	"sync"
	"time"
)

type monster struct {
	baseEntity
	hp, sp        uint32
	mobInfo       *shn.MobInfo
	mobInfoServer *shn.MobInfoServer
	regenData     *mobs.RegenEntry
	tickers       []*time.Ticker
	status
	sync.RWMutex
}

func (m *monster) getHandle() uint16 {
	m.RLock()
	h := m.handle
	m.RUnlock()
	return h
}

func (m *monster) alive() {

	// tick for roaming around, idle action, return to base

}

func (m *monster) dead() {

	// initial removal from monsters...
	// trigger monsterDied event
	// create ticker so it can respawn again
}

func (m *monster) monsterWalk(c chan<- location, zm *zoneMap) {
	defer func() {
		n := recover()
		if n != nil {
			log.Infof("recovered from panic %v", n)
		}
	}()

	tick := time.NewTicker(time.Duration(int64(networking.RandomIntBetween(5256, 21234))) * time.Millisecond)
loop:
	for {
		select {
		case <-tick.C:
			m.RLock()
			var (
				x         = m.current.x
				y         = m.current.y
				walkSpeed = int(m.mobInfo.WalkSpeed / 10)
			)
			m.RUnlock()

			var (
				lx int
				ly int
			)

			rX, rY := igCoordToBitmap(x, y)

			switch networking.RandomIntBetween(0, 8) {
			case 1:
				if rX < rX+walkSpeed {
					lx = networking.RandomIntBetween(rX, rX+walkSpeed)
				}
				if rY < rY+walkSpeed {
					ly = networking.RandomIntBetween(rY, rY+walkSpeed)
				}
			case 2:
				if rX-(walkSpeed*4) < rX+walkSpeed {
					lx = networking.RandomIntBetween(rX-(walkSpeed*4), rX+walkSpeed)
				}
				if rY-(walkSpeed*4) < rY+walkSpeed {
					ly = networking.RandomIntBetween(rY-(walkSpeed*4), rY+walkSpeed)
				}
			case 3:
				if rY-(walkSpeed*4) < rY {
					lx = rX
					ly = networking.RandomIntBetween(rY-(walkSpeed*4), rY)
				}
			case 4:
				if rY+(walkSpeed*4) < rY {
					lx = rX
					ly = networking.RandomIntBetween(rY+(walkSpeed*4), rY)
				}
			case 5:
				if rX-(walkSpeed*4) < rX {
					lx = networking.RandomIntBetween(rX-(walkSpeed*4), rX)
					ly = rY
				}
			case 6:
				if rX+(walkSpeed*4) < rX {
					lx = networking.RandomIntBetween(rX+(walkSpeed*4), rX)
					ly = rY
				}
			case 7:
				m.RLock()
				x = m.fallback.x
				y = m.fallback.y
				m.RUnlock()
				c <- location{
					x: x,
					y: y,
				}
				continue loop
			}

			igX, igY := bitmapCoordToIg(lx, ly)

			if canWalk(zm.walkableX, zm.walkableY, lx, ly) {
				c <- location{
					x: igX,
					y: igY,
				}
			}
		}
	}
}

func (m *monster) roam(zm *zoneMap) {
	ch := make(chan location)
	go m.monsterWalk(ch, zm)
	for {
		select {
		case l := <-ch:
			m.RLock()
			nc := structs.NcActSomeoneMoveWalkCmd{
				Handle: m.handle,
				From: structs.ShineXYType{
					X: uint32(m.current.x),
					Y: uint32(m.current.y),
				},
				To: structs.ShineXYType{
					X: uint32(l.x),
					Y: uint32(l.y),
				},
				Speed:    uint16(m.mobInfo.WalkSpeed),
				MoveAttr: structs.NcActSomeoneMoveWalkCmdAttr{},
			}
			m.RUnlock()

			m.Lock()
			m.current.x = l.x
			m.current.y = l.y
			m.Unlock()

			e := monsterWalksEvent{
				nc: &nc,
				m:  m,
			}

			zm.send[monsterWalks] <- &e
		}
	}
}
