package zone

import (
	"fmt"
	"sync"
)

const playerHandleMin uint16 = 8000
const playerHandleMax uint16 = 12000
const playerAttemptsMax uint16 = 50

type players struct {
	handleIndex uint16
	active      map[uint16]*player
	sync.RWMutex
}

func (p *players) newHandle() (uint16, error) {
	var attempts uint16 = 0
	min := playerHandleMin
	max := playerHandleMax
	maxAttempts := playerAttemptsMax

	p.RLock()
	index := p.handleIndex
	p.RUnlock()

	for {
		if attempts == maxAttempts {
			return 0, fmt.Errorf("\nmaximum number of attempts reached, no handle is available")
		}

		index++

		if index == max {
			index = min
		}

		p.Lock()
		p.handleIndex = index
		p.Unlock()

		p.RLock()
		if _, used := p.active[index]; used {
			attempts++
			continue
		}
		p.RUnlock()

		return index, nil
	}
}

func playerInRange(viewer, target *player) bool {

	target.RLock()
	targetX := (target.x * 8) / 50
	targetY := (target.y * 8) / 50
	targetHandle := target.handle
	target.RUnlock()

	viewer.RLock()
	viewerX := (viewer.x * 8) / 50
	viewerY := (viewer.y * 8) / 50
	viewerHandle := viewer.handle
	viewer.RUnlock()

	vertical := targetY <= viewerY+lengthY && targetY >= viewerY || targetY >= (viewerY-lengthY) && targetY <= viewerY
	horizontal := targetX <= (viewerX+lengthX) && targetX >= viewerX || targetX >= (viewerX-lengthX) && targetX <= viewerX

	if vertical && horizontal {

		viewer.Lock()
		viewer.knownNearbyPlayers[target.handle] = target
		viewer.Unlock()

		log.Infof("%v is in range of %v", targetHandle, viewerHandle)
		return true
	}

	log.Infof("%v is in not in range of %v", targetHandle, viewerHandle)

	return false
}

