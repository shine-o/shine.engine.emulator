package zone

import (
	"testing"
)

func TestLoadNPCOk(t *testing.T) {
	//zm := zoneMap{}
	//
	//err := zm.load()
	t.Fail()
}

func TestLoadAllMapNpc(t *testing.T) {
	mapCount := 0
	npcCount := 0
	for _, npcs := range npcData.MapNPCs {
		mapCount++
		for _, npc := range npcs {
			n, err := loadBaseNpc(npc.MobIndex, isNPC)
			if err != nil {
				t.Fatal(err)
			}

			if n.data == nil {
				t.Errorf("value should not be nil %v", npc.MobIndex)
			}

			if n.ticks == nil {
				t.Errorf("value should not be nil %v", npc.MobIndex)
			}

			if n.state == nil {
				t.Errorf("value should not be nil %v", npc.MobIndex)
			}

			if n.stats == nil {
				t.Errorf("value should not be nil %v", npc.MobIndex)
			}

			if n.eType != isNPC {
				t.Errorf("unexpected entity type %v %v", n.eType, npc.MobIndex)
			}
			npcCount++
		}
	}
	log.Infof("npcs %v maps %v", npcCount, mapCount)
}

func TestLoadAllMapRegenMonsterNpc(t *testing.T) {
	// this is needed to avoid duplicate NPCs
	mapCount := 0
	mobCount := 0
	for _, mapRegen := range monsterData.MapRegens {
		mapCount++
		for _, group := range mapRegen.Groups {
			for _, mob := range group.Mobs {
				mn, err := loadBaseNpc(mob.Index, isMonster)
				if err != nil {
					t.Fatal(err)
				}

				if mn.data == nil {
					t.Errorf("value should not be nil %v", mob.Index)
				}

				if mn.ticks == nil {
					t.Errorf("value should not be nil %v", mob.Index)
				}

				if mn.state == nil {
					t.Errorf("value should not be nil %v", mob.Index)
				}

				if mn.stats == nil {
					t.Errorf("value should not be nil %v", mob.Index)
				}

				if mn.eType != isMonster {
					t.Errorf("unexpected entity type %v %v", mn.eType, mob.Index)
				}
				mobCount++
			}
		}
	}
	log.Infof("mobs %v maps %v", mobCount, mapCount)
}

func TestLoadAllMobInfoMonsterNpc(t *testing.T) {
	// this is needed to avoid duplicate NPCs
	//for , mapRegen := range monsterData.MapRegens {
	//	for , group := range mapRegen.Groups {
	//		for , mob := range group.Mobs {
	//			mob.
	//		}
	//	}
	//}
	t.Fail()
}

func TestLoadMonsterNPCOk(t *testing.T) {
	t.Fail()
}

func TestLoadNPCMissingData(t *testing.T) {
	t.Fail()
}

func TestLoadMonsterNPCMissingData(t *testing.T) {
	t.Fail()
}

func TestLoadVendorNPCMissingData(t *testing.T) {
	t.Fail()
}

func TestLoadVendorNPCOk(t *testing.T) {
	t.Fail()
}

func TestLoadVendorNPCLoadPages(t *testing.T) {
	t.Fail()
}

func TestLoadPortalNPCOk(t *testing.T) {
	t.Fail()
}

func TestLoadNPCOkNC(t *testing.T) {
	t.Fail()
}

func TestLoadMonsterNPCNC(t *testing.T) {
	t.Fail()
}

func TestLoadVendorNPCNC(t *testing.T) {
	t.Fail()
}

func TestTriggerNPCAction(t *testing.T) {
	t.Fail()
}
