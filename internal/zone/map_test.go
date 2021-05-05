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

	if zm.walkableNodes == nil {
		t.Fatal("value should not be nil")
	}

	if zm.pathNodes == nil {
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

func Test_Map_Path_A_B_BIAStar(t *testing.T)  {

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