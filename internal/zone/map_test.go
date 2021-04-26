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

	if zm.walkableX == nil {
		t.Fatal("value should not be nil")
	}

	if zm.walkableY == nil {
		t.Fatal("value should not be nil")
	}

	if zm.entities == nil {
		t.Fatal("value should not be nil")
	}

	if zm.metrics == nil {
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

	if len(zm.entities.npcs.active) == 0 {
		t.Fatal("there should be at least one npc active")
	}
}