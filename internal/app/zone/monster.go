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

func (m *monster) roam(zm *zoneMap) {
	log.Infof("[monster_ticks] roam for handle %v", m.handle)

	tick := time.NewTicker(time.Duration(int64(networking.RandomIntBetween(15, 20))) * time.Second)

	m.Lock()
	m.tickers = append(m.tickers, tick)
	m.Unlock()
	defer tick.Stop()

	// first, walks to a position
	// 		while it walks, notify the player
	// after it reaches destination, stop there for a few seconds
	//
	var state = monsterRoaming

	for {
		select {
		case <-tick.C:

			m.RLock()
			rd := m.regenData.RangeDegree

			if rd == 0 {
				continue
			}

			if rd < 0 {
				rd = 40
			}
			m.RUnlock()

			if state == monsterRoaming {
				var (
					x, y  int
					tries = 100
					done  = false
				)

				for tries != 0 {

					if done {
						break
					}

					m.RLock()

					x = networking.RandomIntBetween(int(m.location.x), int(m.location.x)+rd)
					y = networking.RandomIntBetween(int(m.location.y), int(m.location.y)+rd)
					m.RUnlock()

					rX, rY := igCoordToBitmap(x, y)

					if canWalk(zm.walkableX, zm.walkableY, uint32(rX), uint32(rY)) {
						done = true
					}

					tries--

				}

				if tries == 0 {
					break
				}

				m.RLock()
				nc := structs.NcActSomeoneMoveWalkCmd{
					Handle: m.handle,
					From: structs.ShineXYType{
						X: m.x,
						Y: m.y,
					},
					To: structs.ShineXYType{
						X: uint32(x),
						Y: uint32(y),
					},
					Speed:    uint16(m.mobInfo.WalkSpeed),
					MoveAttr: structs.NcActSomeoneMoveWalkCmdAttr{},
				}
				m.RUnlock()
				state = monsterGoBack

				m.Lock()
				m.location.x = uint32(x)
				m.location.y = uint32(y)
				m.Unlock()

				go func() {

					e := monsterWalksEvent{
						nc: &nc,
					}
					zm.send[monsterWalks] <- &e

				}()

			}

			if state == monsterGoBack {

				m.RLock()
				nc := structs.NcActSomeoneMoveRunCmd{
					Handle: m.handle,
					From: structs.ShineXYType{
						X: m.location.x,
						Y: m.location.y,
					},
					To: structs.ShineXYType{
						X: m.fallback.x,
						Y: m.fallback.y,
					},
					Speed:    uint16(m.mobInfo.RunSpeed),
				}
				m.RUnlock()

				state = monsterRoaming

				m.Lock()
				m.location.x = m.fallback.x
				m.location.y = m.fallback.y
				m.Unlock()

				go func() {
					e := monsterRunsEvent{
						nc: &nc,
					}
					zm.send[monsterRuns] <- &e
				}()

			}


			//if state == monsterIdling {
			//
			//
			//}
			//
			//if state == monsterWalking {
			//	state = monsterRunning
			//}
			//
			//if state == monsterWalking {
			//	state = monsterIdling
			//}
			//
			//if state == monsterFighting {
			//	continue
			//}

			//nc

			//zoneEvents[monsterRuns]    <- &emptyEvent{}
			//zoneEvents[monsterStopped] <- &emptyEvent{}

		}
	}
}
