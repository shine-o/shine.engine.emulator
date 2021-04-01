package zone

import (
	"github.com/shine-o/shine.engine.emulator/internal/pkg/data"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/networking"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/structs"
	"sync"
	"time"
)

type monsterLocationActivity struct {
	location
	run bool
}

type monster struct {
	baseEntity
	hp, sp        uint32
	mobInfo       *data.MobInfo
	mobInfoServer *data.MobInfoServer
	regenData     *data.RegenEntry
	tickers       []*time.Ticker
	status
	sync.RWMutex
}

func (m *monster) alive() {

	// tick for roaming around, idle action, return to base

}

func (m *monster) dead() {
	// initial removal from monsters...
	// trigger monsterDied event
	// create ticker so it can respawn again
}

func (m *monster) monsterActivity(c chan<- monsterLocationActivity, zm *zoneMap) {
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
				if rX-(walkSpeed) < rX+walkSpeed {
					lx = networking.RandomIntBetween(rX-(walkSpeed), rX+walkSpeed)
				}
				if rY-(walkSpeed) < rY+walkSpeed {
					ly = networking.RandomIntBetween(rY-(walkSpeed), rY+walkSpeed)
				}
			case 3:
				if rY-(walkSpeed) < rY {
					lx = rX
					ly = networking.RandomIntBetween(rY-(walkSpeed), rY)
				}
			case 4:
				if rY+(walkSpeed) < rY {
					lx = rX
					ly = networking.RandomIntBetween(rY+(walkSpeed), rY)
				}
			case 5:
				if rX-(walkSpeed) < rX {
					lx = networking.RandomIntBetween(rX-(walkSpeed), rX)
					ly = rY
				}
			case 6:
				if rX+(walkSpeed) < rX {
					lx = networking.RandomIntBetween(rX+(walkSpeed), rX)
					ly = rY
				}
			case 7:
				m.RLock()
				x = m.fallback.x
				y = m.fallback.y
				m.RUnlock()
				c <- monsterLocationActivity{
					location: location{
						x: x,
						y: y,
					},
					run: true,
				}
				continue loop
			}

			igX, igY := bitmapCoordToIg(lx, ly)

			if canWalk(zm.walkableX, zm.walkableY, lx, ly) {
				c <- monsterLocationActivity{
					location: location{
						x: igX,
						y: igY,
					},
				}
			}
		}
	}
}

func (m *monster) roam(zm *zoneMap) {
	ch := make(chan monsterLocationActivity)
	go m.monsterActivity(ch, zm)
	for {
		select {
		case l := <-ch:
			m.RLock()
			if l.run {
				nc := structs.NcActSomeoneMoveRunCmd{
					Handle: m.handle,
					From: structs.ShineXYType{
						X: uint32(m.current.x),
						Y: uint32(m.current.y),
					},
					To: structs.ShineXYType{
						X: uint32(l.x),
						Y: uint32(l.y),
					},
					Speed: uint16(m.mobInfo.RunSpeed),
				}
				m.RUnlock()

				m.Lock()
				m.current.x = l.x
				m.current.y = l.y
				m.Unlock()

				e := monsterRunsEvent{
					nc: &nc,
					m:  m,
				}

				zm.send[monsterRuns] <- &e
			} else {
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
}

func (m *monster) getHandle() uint16 {
	m.RLock()
	h := m.handle
	m.RUnlock()
	return h
}

func (m *monster) ncBatTargetInfoCmd() *structs.NcBatTargetInfoCmd {
	var nc structs.NcBatTargetInfoCmd
	m.RLock()
	nc = structs.NcBatTargetInfoCmd{
		Order:         0,
		Handle:        m.handle,
		TargetHP:      m.hp,
		TargetMaxHP:   m.mobInfo.MaxHP, //todo: use the same player stat system for mobs and NPCs
		TargetSP:      m.sp,
		TargetMaxSP:   uint32(m.mobInfoServer.MaxSP), //todo: use the same player stat system for mobs and NPCs
		TargetLP:      0,
		TargetMaxLP:   0,
		TargetLevel:   byte(m.mobInfo.Level),
		HpChangeOrder: 0,
	}
	m.RUnlock()
	return &nc
}
