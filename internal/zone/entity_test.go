package zone

import (
	"testing"

	"github.com/shine-o/shine.engine.emulator/internal/pkg/persistence"
)

func TestMoveEntityAB(t *testing.T) {
	zm, err := loadMap(1)
	if err != nil {
		t.Fatal(err)
	}

	persistence.CleanDB()

	char := persistence.NewDummyCharacter("mage", false, "dummy")

	p := &player{
		baseEntity: &baseEntity{},
	}

	err = p.load(char.Name)

	zm.addEntity(p)

	x := 4089
	y := 3214

	err = p.move(zm, x, y)

	if err != nil {
		t.Fatal(err)
	}

	if p.baseEntity.current.x != x || p.baseEntity.current.y != y {
		t.Fatalf("mismatched coordinates %v %v", p.baseEntity.current.x, p.baseEntity.current.y)
	}
}

func TestMoveEntityCollision(t *testing.T) {
	t.Fail()
}

// Entities A, B
// A, B are spawned apart from one another
// A enters B range
// A and B know about each other's existence
func TestEntityWithinRange(t *testing.T) {
	// Roumen
	zm, err := loadMap(1)
	if err != nil {
		t.Fatal(err)
	}

	// load players A and B in two distinct places
	pA := &player{
		baseEntity: &baseEntity{
			handle: 1,
			proximity: &entityProximity{
				entities: make(map[uint16]entity),
			},
		},
	}

	pB := &player{
		baseEntity: &baseEntity{
			handle: 2,
			proximity: &entityProximity{
				entities: make(map[uint16]entity),
			},
		},
	}

	zm.addEntity(pA)
	zm.addEntity(pB)

	// Move entities to designated
	x1, y1 := gameCoordinates(1060, 767)
	x2, y2 := gameCoordinates(700, 763)

	_ = pA.move(zm, x1, y1)
	_ = pB.move(zm, x2, y2)

	// entities A and B should not be in range
	if withinRange(pA, pB) {
		t.Fatal("entity A should not be in range of entity B")
	}

	// move entity A closer to entity B
	_ = pA.move(zm, x2+10, y2+10)

	// manually find entities
	addWithinRangeEntities(pA, zm)
	addWithinRangeEntities(pB, zm)

	// entities A and B should now be in range
	if !withinRange(pA, pB) {
		t.Fatal("entity A should be in range of entity B")
	}

	// assert entity A is stored in entity B's proximity list
	_, ok := pB.baseEntity.proximity.entities[pA.getHandle()]
	if !ok {
		t.Fatal("entity A is not stored in entity B's proximity list")
	}

	_, ok = pA.baseEntity.proximity.entities[pB.getHandle()]
	if !ok {
		t.Fatal("entity B is not stored in entity A's proximity list")
	}
}

// A, B know about each other's existence
// A moves out of range
// A and B don't know about each other's existence anymore
func TestEntityOutOfRange(t *testing.T) {
	// Roumen
	zm, err := loadMap(1)
	if err != nil {
		t.Fatal(err)
	}

	// load players A and B in two distinct places
	pA := &player{
		baseEntity: &baseEntity{
			handle: 1,
			proximity: &entityProximity{
				entities: make(map[uint16]entity),
			},
		},
	}

	pB := &player{
		baseEntity: &baseEntity{
			handle: 2,
			proximity: &entityProximity{
				entities: make(map[uint16]entity),
			},
		},
	}

	zm.addEntity(pA)
	zm.addEntity(pB)

	// Move entities to designated
	x1, y1 := gameCoordinates(1060, 767)
	x2, y2 := gameCoordinates(700, 763)

	_ = pA.move(zm, x1, y1)
	_ = pB.move(zm, x2, y2)

	// move entity A closer to entity B
	_ = pA.move(zm, x2+10, y2+10)

	// manually find entities
	addWithinRangeEntities(pA, zm)
	addWithinRangeEntities(pB, zm)

	// move entity back to its original position
	_ = pA.move(zm, x1, y1)
	removeOutOfRangeEntities(pA)
	removeOutOfRangeEntities(pB)

	// entities A and B should not be in range
	if withinRange(pA, pB) {
		t.Fatal("entity A should not be in range of entity B")
	}

	// assert entity A is NOT stored in entity B's proximity list
	_, ok := pB.baseEntity.proximity.entities[pA.getHandle()]
	if ok {
		t.Fatal("entity A should not be stored in entity B's proximity list")
	}

	_, ok = pA.baseEntity.proximity.entities[pB.getHandle()]
	if ok {
		t.Fatal("entity B should not be stored in entity A's proximity list")
	}
}

// Entities A, B, C
// A selects C
// B selects A
// A is aware of being selected by B
// C is aware of being selected by A
// B is notified about A selected C
// A is notified about selected C
func TestEntitySelectingEntityAwareness(t *testing.T) {
	eA := &player{
		baseEntity: &baseEntity{
			handle: 1,
		},
		targeting: &targeting{
			players:  make(map[uint16]*player),
			monsters: make(map[uint16]*monster),
			npc:      make(map[uint16]*npc),
		},
	}

	eB := &player{
		baseEntity: &baseEntity{
			handle: 2,
		},
		targeting: &targeting{
			players:  make(map[uint16]*player),
			monsters: make(map[uint16]*monster),
			npc:      make(map[uint16]*npc),
		},
	}

	eC := &monster{
		baseEntity: &baseEntity{
			handle: 3,
		},
		targeting: &targeting{
			players:  make(map[uint16]*player),
			monsters: make(map[uint16]*monster),
			npc:      make(map[uint16]*npc),
		},
	}

	eA.selects(eC)
	eC.selectedBy(eA)

	eB.selects(eA)
	eA.selectedBy(eB)

	// C is aware of being selected by A
	_, ok := eC.targeting.players[eA.getHandle()]
	if !ok {
		t.Fatal("A must be aware that its selected by B")
	}

	// A is aware of being selected by B
	_, ok = eA.targeting.players[eB.getHandle()]
	if !ok {
		t.Fatal("A must be aware that its selected by B")
	}

	e := eB.selected().selected()

	if e == nil {
		t.Fatalf("B must be aware that A is selecting C")
	}
}

func TestEntityUnSelectsEntity(t *testing.T) {
	// A unselects something
	// INFO : 2021/05/23 14:30:52.071356 handlers.go:271: 2021-05-23 14:30:52.058449 +0200 CEST 2388->9120 outbound NC_BAT_UNTARGET_REQ {"packetType":"small","length":2,"department":9,"command":"8","opCode":9224,"data":"","rawData":"020824","friendlyName":""}

	// B gets this info about A unselecting something
	// INFO : 2021/05/23 14:30:52.211170 handlers.go:271: 2021-05-23 14:30:52.202183 +0200 CEST 9120->2377 inbound NC_BAT_TARGETINFO_CMD {"packetType":"small","length":32,"department":9,"command":"2","opCode":9218,"data":"01ffff000000000000000000000000000000000000000000000000000000","rawData":"20022401ffff000000000000000000000000000000000000000000000000000000","friendlyName":""}
	eA := &player{
		baseEntity: &baseEntity{
			handle: 1,
		},
		targeting: &targeting{
			players:  make(map[uint16]*player),
			monsters: make(map[uint16]*monster),
			npc:      make(map[uint16]*npc),
		},
	}

	eB := &player{
		baseEntity: &baseEntity{
			handle: 2,
		},
		targeting: &targeting{
			players:  make(map[uint16]*player),
			monsters: make(map[uint16]*monster),
			npc:      make(map[uint16]*npc),
		},
	}

	eC := &monster{
		baseEntity: &baseEntity{
			handle: 3,
		},
		targeting: &targeting{
			players:  make(map[uint16]*player),
			monsters: make(map[uint16]*monster),
			npc:      make(map[uint16]*npc),
		},
	}

	eA.selects(eC)
	eC.selectedBy(eA)

	eB.selects(eA)
	eA.selectedBy(eB)

	eC.unselectedBy(eA)
	eA.unselectedBy(eB)

	_, ok := eC.targeting.players[eA.getHandle()]
	if ok {
		t.Fatal("A must NOT be aware that its selected by B")
	}

	_, ok = eA.targeting.players[eB.getHandle()]
	if ok {
		t.Fatal("A must NOT be aware that its selected by B")
	}

}
