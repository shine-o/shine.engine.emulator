package zone

import (
	"testing"
)

func TestLoadMap(t *testing.T) {
	// Roumen
	zm, err := loadMap(1)

	if err != nil {
		t.Fatal(err)
	}

	if zm.events == nil {
		t.Fatal("value should not be nil")
	}

	if zm.data == nil {
		t.Fatal("value should not be nil")
	}

	if zm.rawNodes == nil {
		t.Fatal("value should not be nil")
	}

	if zm.presetNodes == nil {
		t.Fatal("value should not be nil")
	}

	if zm.presetNodesWithMargin == nil {
		t.Fatal("value should not be nil")
	}

	if zm.entities == nil {
		t.Fatal("value should not be nil")
	}

	expectedMapEvents := []eventIndex{
		playerHandle,
		playerHandleMaintenance,
		queryPlayer, queryMonster,
		playerAppeared, playerDisappeared, playerJumped, playerWalks, playerRuns, playerStopped,
		unknownHandle, monsterAppeared, monsterDisappeared, monsterWalks, monsterRuns,
		playerSelectsEntity, playerUnselectsEntity, playerClicksOnNpc, playerPromptReply, itemIsMoved, itemEquip, itemUnEquip,
	}

	if len(zm.events.send) != len(expectedMapEvents) || len(zm.events.recv) != len(expectedMapEvents) {
		t.Fatalf("mismatched amount of events %v %v %v ", len(zm.events.send), len(zm.events.recv), len(expectedMapEvents))
	}

	for _, e := range expectedMapEvents {
		_, ok := zm.events.send[e]
		if !ok {
			t.Errorf("missing zone event %v", e)
		}
		_, ok = zm.events.recv[e]
		if !ok {
			t.Errorf("missing zone event %v", e)
		}
	}
}

func Test_Map_Spawn_Npc(t *testing.T) {
	z := zone{}

	err := z.load()

	if err != nil {
		t.Fatal(err)
	}

	go z.run()

	zm, err := loadMap(1)

	if err != nil {
		t.Fatal(err)
	}

	zm.spawnNPCs()

	if len(zm.entities.npcs.active) == 0 {
		t.Fatal("there should be at least one npc active")
	}

	if len(zm.entities.npcs.active) != 30 {
		t.Fatalf("incorrect npc amount %v", len(zm.entities.npcs.active))
	}

	for _, npc := range zm.entities.npcs.active {
		if npc.nType == npcNoRole {
			t.Errorf("unexpected npcType %v %v", npc.nType, npc.data.mobInfo.InxName)
		}
	}

	zm1, err := loadMap(4)

	if err != nil {
		t.Fatal(err)
	}

	zm1.spawnNPCs()

	if len(zm1.entities.npcs.active) == 0 {
		t.Fatal("there should be at least one npc active")
	}

	if len(zm1.entities.npcs.active) != 8 {
		t.Fatalf("incorrect npc amount %v", len(zm1.entities.npcs.active))
	}

	for _, npc := range zm1.entities.npcs.active {
		if npc.nType == npcNoRole {
			t.Errorf("unexpected npcType %v %v", npc.nType, npc.data.mobInfo.InxName)
		}
	}
}

func Test_Map_Spawn_Monster_Npc(t *testing.T) {
	z := zone{}

	err := z.load()

	if err != nil {
		t.Fatal(err)
	}

	go z.run()

	zm, err := loadMap(1)

	if err != nil {
		t.Fatal(err)
	}

	zm.spawnMobs()

	if len(zm.entities.npcs.active) == 0 {
		t.Fatal("there should be at least one npc active")
	}

	if len(zm.entities.npcs.active) != 48 {
		t.Fatalf("incorrect npc amount %v", len(zm.entities.npcs.active))
	}
	// 48
}

func Test_Map_Path_A_B_astar(t *testing.T)  {
	
}

func Test_Map_Intermitent_Speed_Change_Path_A_B_AStar(t *testing.T)  {
	t.Fail()
	// start moving entity from point A to point B using speed X
	// midway, change speed to Y, point A will now be current point
	// assert entity arrives at destination in Z seconds ( time = distance / speed )
}

func Test_Entity_Chase(t *testing.T)  {

}

func Test_Path_A_B_Speed_Race(t *testing.T)  {

}

func Test_Entity_Within_Range(t *testing.T) {
	// Roumen
	zm, err := loadMap(1)

	if err != nil {
		t.Fatal(err)
	}

	// load players A and B in two distinct places
	pA := &player{
		baseEntity: &baseEntity{
			handle:    1,
			proximity: &entityProximity{
				entities: make(entitiesmap),
			},
		},
	}

	pB := &player{
		baseEntity: &baseEntity{
			handle:    2,
			proximity: &entityProximity{
				entities: make(entitiesmap),
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

func Test_Entity_Out_Of_Range(t *testing.T) {
	// Roumen
	zm, err := loadMap(1)

	if err != nil {
		t.Fatal(err)
	}

	// load players A and B in two distinct places
	pA := &player{
		baseEntity: &baseEntity{
			handle:    1,
			proximity: &entityProximity{
				entities: make(entitiesmap),
			},
		},
	}

	pB := &player{
		baseEntity: &baseEntity{
			handle:    2,
			proximity: &entityProximity{
				entities: make(entitiesmap),
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



func Test_Entity_Cast_(t *testing.T) {
	
}