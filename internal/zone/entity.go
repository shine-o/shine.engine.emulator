package zone

import (
	"github.com/shine-o/shine.engine.emulator/internal/pkg/errors"
	"sync"
)

type entity interface {
	basicActions
	getHandle() uint16
	getLocation() (int, int)
}

type basicActions interface {
	move(m *zoneMap, x, y int) error
}

type location struct {
	mapID     int
	mapName   string
	x, y, d   int
	movements []movement
}

type movement struct {
	x, y uint32
}

type entityType int

type baseEntity struct {
	handle   uint16
	eType    entityType
	fallback location
	previous location
	current  location
	next     location
	events   events
	// dangerZone: only to be used when loading or other situation!!
	sync.RWMutex
}

type targeting struct {
	selectionOrder byte
	selectingP     *player
	selectingN     *npc
	selectedByP    []*player
	selectedByN    []*npc
	sync.RWMutex
}

type entityState struct {
	idling   chan bool
	fighting chan bool
	chasing  chan bool
	fleeing  chan bool
	sync.RWMutex
}

type mover struct {
	baseEntity
}

const (
	//lengthX = 512
	//lengthY = 512
	lengthX = 256
	lengthY = 256
)

const (
	isMonster entityType = iota
	isPlayer
	isNPC
)

var _ entity = (*player)(nil)

var _ entity = (*npc)(nil)

func (b *baseEntity) getHandle() uint16 {
	b.RLock()
	h := b.handle
	b.RUnlock()
	return h
}

func (b *baseEntity) getLocation() (int, int) {
	return b.current.x, b.current.y
}

func (b *baseEntity) move(m *zoneMap, igX, igY int) error {
	rX, rY := bitmapCoordinates(igX, igY)

	if !canWalk(m.walkableX, m.walkableY, rX, rY) {
		return errors.Err{
			Code: errors.ZoneMapCollisionDetected,
			Details: errors.ErrDetails{
				"entity": b.getHandle(),
				"igX":    igX,
				"igY":    igY,
			},
		}
	}

	b.Lock()
	b.previous.x = b.current.x
	b.previous.y = b.current.y
	b.current.x = igX
	b.current.y = igY
	b.Unlock()

	return nil
}

func entityInRange(e1, e2 location) bool {
	viewerX, viewerY := bitmapCoordinates(e1.x, e1.y)
	targetX, targetY := bitmapCoordinates(e2.x, e2.y)

	maxY := viewerY + lengthY
	minY := viewerY - lengthY

	maxX := viewerX + lengthX
	minX := viewerX - lengthX

	vertical := (targetY <= maxY && targetY >= viewerY) || (targetY >= minY && targetY <= viewerY)
	horizontal := (targetX <= maxX && targetX >= viewerX) || (targetX >= minX && targetX <= viewerX)

	if vertical && horizontal {
		return true
	}

	return false
}
