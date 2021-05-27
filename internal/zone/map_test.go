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

func TestMapSpawnNpc(t *testing.T) {
	zm, err := loadMap(1)
	if err != nil {
		t.Fatal(err)
	}

	zm.spawnNPCs()

	if len(zm.entities.npc) == 0 {
		t.Fatal("there should be at least one npc active")
	}

	if len(zm.entities.npc) != 29 {
		t.Fatalf("incorrect npc amount %v", len(zm.entities.npc))
	}

	for npc := range zm.entities.allNpc() {
		if npc.nType == npcNoRole {
			t.Errorf("unexpected npcType %v %v", npc.nType, npc.data.mobInfo.InxName)
		}
	}

	zm1, err := loadMap(4)
	if err != nil {
		t.Fatal(err)
	}

	zm1.spawnNPCs()

	if len(zm1.entities.npc) == 0 {
		t.Fatal("there should be at least one npc active")
	}

	if len(zm1.entities.npc) != 8 {
		t.Fatalf("incorrect npc amount %v", len(zm1.entities.npc))
	}

	for npc := range zm1.entities.allNpc() {
		if npc.nType == npcNoRole {
			t.Errorf("unexpected npcType %v %v", npc.nType, npc.data.mobInfo.InxName)
		}
	}
}

func TestMapSpawnMonsterNpc(t *testing.T) {
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

	if len(zm.entities.npc) == 0 {
		t.Fatal("there should be at least one npc active")
	}

	if len(zm.entities.npc) != 48 {
		t.Fatalf("incorrect npc amount %v", len(zm.entities.npc))
	}
	// 48
}

func TestMapPathABAstar(t *testing.T) {
	t.Fail()
}

func TestMapIntermittentSpeedChangePathABAStar(t *testing.T) {
	t.Fail()
	// start moving entity from point A to point B using speed X
	// midway, change speed to Y, point A will now be currentlySelected point
	// assert entity arrives at destination in Z seconds ( time = distance / speed )
}

func TestEntityChase(t *testing.T) {
	t.Fail()
}

func Test_Path_A_B_Speed_Race(t *testing.T) {
	t.Fail()
}
