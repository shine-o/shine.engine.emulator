package zone

import (
	mobs "github.com/shine-o/shine.engine.emulator/internal/pkg/game-data/monsters"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/game-data/shn"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/networking"
	"github.com/shine-o/shine.engine.emulator/pkg/structs"
	"sync"
	"time"
)

const (
	monsterRoaming = iota
	monsterGoBack
	monsterIdling
	monsterFighting
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

func (m * monster) monsterWalk(c chan <- * location, zm * zoneMap)  {
	tick := time.NewTicker(time.Duration(int64(networking.RandomIntBetween(4200, 8150))) * time.Millisecond)
	//tick := time.NewTicker(3 * time.Second)
	for {
		select {
		case <- tick.C:
			m.RLock()
			var (
				x = m.x
				y = m.y
				walkSpeed = (m.mobInfo.WalkSpeed/10) * 2
			)
			m.RUnlock()

			var (
				lx, ly uint32
			)

			// problem is, I need to generate the same kind of coordinates as the game does, I cannot invent them!
			//rstX := (rX * 50.0) / 8.0
			rX, rY := igCoordToBitmap(x, y)

			switch networking.RandomIntBetween(0,5) {
				case 1:
					lx = rX
					ly = rY - walkSpeed
				case 2:
					lx = rX
					ly = rY + walkSpeed
				case 3:
					lx = rX - walkSpeed
					ly = rY
				case 4:
					lx = rX + walkSpeed
					ly = rY
			}

			if ly > m.fallback.y + 400 && lx > m.fallback.x + 400 {
				c <- &location{
					x:         m.fallback.x,
					y:        m.fallback.y,
				}
				break
			}


			if canWalk(zm.walkableX, zm.walkableY, lx, ly) {
				igX, igY := bitmapCoordToIg(lx, ly)

				c <- &location{
					x:         igX,
					y:         igY,
				}
			}
		}
	}
}

func (m *monster) roam(zm *zoneMap) {
	ch := make(chan * location)
	go m.monsterWalk(ch, zm)
	for {
		select {
			case l := <- ch:
				m.RLock()
				nc := structs.NcActSomeoneMoveWalkCmd{
					Handle: m.handle,
					From: structs.ShineXYType{
						X: m.x,
						Y: m.y,
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
				m.location.x = l.x
				m.location.y = l.y
				m.Unlock()

				e := monsterWalksEvent{
					nc: &nc,
					m:  m,
				}

				zm.send[monsterWalks] <- &e
		}
	}
}
