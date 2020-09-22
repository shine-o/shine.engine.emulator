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

	tick := time.NewTicker(time.Duration(int64(networking.RandomIntBetween(8256, 21234))) * time.Millisecond)
loop:
	for {
		select {
		case <-tick.C:
			m.RLock()
			var (
				x         = m.current.x
				y         = m.current.y
				walkSpeed = m.mobInfo.WalkSpeed / 10
			)
			m.RUnlock()

			var (
				lx uint32
				ly uint32
			)

			rX, rY := igCoordToBitmap(x, y)

			switch networking.RandomIntBetween(0, 8) {
			case 1:
				if rX < rX+walkSpeed {
					lx = uint32(networking.RandomIntBetween(int(rX), int(rX+walkSpeed)))
				}
				if rY < rY+walkSpeed {
					ly = uint32(networking.RandomIntBetween(int(rY), int(rY+walkSpeed)))
				}
			case 2:
				if int(rX-(walkSpeed*4)) < int(rX+walkSpeed) {
					lx = uint32(networking.RandomIntBetween(int(rX-(walkSpeed*4)), int(rX+walkSpeed)))
				}
				if int(rY-(walkSpeed*4)) < int(rY+walkSpeed) {
					ly = uint32(networking.RandomIntBetween(int(rY-(walkSpeed*4)), int(rY+walkSpeed)))
				}
			case 3:
				if int(rY-(walkSpeed*4)) < int(rY) {
					lx = rX
					ly = uint32(networking.RandomIntBetween(int(rY-(walkSpeed*4)), int(rY)))
				}
			case 4:
				if int(rY+(walkSpeed*4)) < int(rY) {
					lx = rX
					ly = uint32(networking.RandomIntBetween(int(rY+(walkSpeed*4)), int(rY)))
				}
			case 5:
				if int(rX-(walkSpeed*4)) < int(rX) {
					lx = uint32(networking.RandomIntBetween(int(rX-(walkSpeed*4)), int(rX)))
					ly = rY
				}
			case 6:
				if int(rX+(walkSpeed*4)) < int(rX) {
					lx = uint32(networking.RandomIntBetween(int(rX+(walkSpeed*4)), int(rX)))
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
					X: m.current.x,
					Y: m.current.y,
				},
				To: structs.ShineXYType{
					X: l.x,
					Y: l.y,
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
