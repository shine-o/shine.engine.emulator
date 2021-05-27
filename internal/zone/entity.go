package zone

import (
	"reflect"
	"sync"

	"github.com/shine-o/shine.engine.emulator/internal/pkg/errors"
	path "github.com/shine-o/shine.engine.emulator/internal/pkg/pathfinding"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/structs"
)

const (
	// lengthX = 512
	// lengthY = 512
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

var _ entity = (*monster)(nil)

type entity interface {
	getHandle() uint16
	movement
	proximity
	target
}

type movement interface {
	getLocation() location
	move(*zoneMap, int, int) error
}

type proximity interface {
	newNearbyEntitiesTicker(*zoneMap)
	oldNearbyEntitiesTicker()
	getNearbyEntities() <-chan entity
	addNearbyEntity(entity)
	removeNearbyEntity(entity)
	alreadyNearbyEntity(entity) bool
	notifyAboutNewEntity(entity)
	notifyAboutRemovedEntity(entity)
	// if data needs to be broadcast to player
	// data can be for types: npc, monster, player
	getNewEntityNearbyPacketData() interface{}
}

type target interface {
	selects(entity) (order byte)
	currentlySelected() entity
	selectedBy(entity)
	// packet data for current entity
	// will return packet with current order
	getTargetPacketData() *structs.NcBatTargetInfoCmd
	// packet data for selected entity by current entity
	// will return packet with current order +1
	getNextTargetPacketData() *structs.NcBatTargetInfoCmd
}

type location struct {
	mapID   int
	mapName string
	x, y, d int
}

type entityType int

type baseEntity struct {
	handle    uint16
	eType     entityType
	fallback  location
	previous  location
	current   location
	next      location
	proximity *entityProximity
	events    events
	sync.RWMutex
}

type entityProximity struct {
	entities map[uint16]entity
	sync.RWMutex
}

func getNearbyEntities(ep *entityProximity) <-chan entity {
	ep.RLock()
	ch := make(chan entity, len(ep.entities))
	ep.RUnlock()

	go func(ep *entityProximity, send chan<- entity) {
		ep.RLock()
		for _, e := range ep.entities {
			send <- e
		}
		ep.RUnlock()
		close(send)
	}(ep, ch)

	return ch
}

type targeting struct {
	selectionOrder    byte
	currentlySelected entity
	players           map[uint16]*player
	monsters          map[uint16]*monster
	npc               map[uint16]*npc
	sync.RWMutex
}

func (t *targeting) selectedBy(e entity) {
	t.Lock()
	switch e.(type) {
	case *player:
		t.players[e.getHandle()] = e.(*player)
		break
	case *monster:
		t.monsters[e.getHandle()] = e.(*monster)
		break
	case *npc:
		t.npc[e.getHandle()] = e.(*npc)
		break
	default:
		log.Error(errors.Err{
			Code: errors.ZoneBadEntityType,
			Details: errors.ErrDetails{
				"got_type": reflect.TypeOf(e).String(),
			},
		})
	}
	t.Unlock()
}

func (t *targeting) selectedByPlayers() <-chan *player {
	t.RLock()
	ch := make(chan *player, len(t.players))
	t.RUnlock()

	go func(t *targeting, send chan<- *player) {
		t.RLock()
		for _, p := range t.players {
			send <- p
		}
		t.RUnlock()
		close(send)
	}(t, ch)

	return ch
}

func (t *targeting) selectedByMonsters() <-chan *monster {
	t.RLock()
	ch := make(chan *monster, len(t.monsters))
	t.RUnlock()

	go func(t *targeting, send chan<- *monster) {
		t.RLock()
		for _, m := range t.monsters {
			send <- m
		}
		t.RUnlock()
		close(send)
	}(t, ch)

	return ch
}

func (t *targeting) selectedByNPCs() <-chan *npc {
	t.RLock()
	ch := make(chan *npc, len(t.npc))
	t.RUnlock()

	go func(t *targeting, send chan<- *npc) {
		t.RLock()
		for _, n := range t.npc {
			send <- n
		}
		t.RUnlock()
		close(send)
	}(t, ch)

	return ch
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

func (b *baseEntity) getHandle() uint16 {
	b.RLock()
	h := b.handle
	b.RUnlock()
	return h
}

func (b *baseEntity) getLocation() location {
	b.RLock()
	defer b.RUnlock()
	return b.current
}

func (b *baseEntity) move(m *zoneMap, igX, igY int) error {
	rX, rY := bitmapCoordinates(igX, igY)

	if !path.CanWalk(m.rawNodes, rX, rY) {
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

func withinRange(e1, e2 entity) bool {
	l1 := e1.getLocation()
	l2 := e2.getLocation()
	if l1.mapID != l2.mapID {
		return false
	}

	if entityInRange(l1, l2) {
		return true
	}

	return false
}

// check if entities are within each other's interaction range
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
