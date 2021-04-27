package zone

import (
	"testing"
)

func Test_Load_NPC_Ok(t *testing.T) {
	//zm := zoneMap{}
	//
	//err := zm.load()
}

func Test_Load_All_Map_Npc(t *testing.T) {
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

			if n.nType == npcNoRole {
				t.Errorf("unexpected entity type %v %v", n.nType, npc.MobIndex)
			}

			if n.eType != isNPC {
				t.Errorf("unexpected entity type %v %v", n.eType, npc.MobIndex)
			}
			npcCount++
		}
	}
	log.Infof("npcs %v maps %v", npcCount, mapCount)
}

func Test_Load_All_MapRegen_Monster_Npc(t *testing.T) {
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

				if mn.nType != npcNoRole {
					t.Errorf("unexpected npcType %v", mob.Index)
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

func Test_Load_All_MobInfo_Monster_Npc(t *testing.T) {
	// this is needed to avoid duplicate NPCs
	//for _, mapRegen := range monsterData.MapRegens {
	//	for _, group := range mapRegen.Groups {
	//		for _, mob := range group.Mobs {
	//			mob.
	//		}
	//	}
	//}
}

func Test_Load_MonsterNPC_Ok(t *testing.T) {

}

func Test_Load_NPC_MissingData(t *testing.T) {

}

func Test_Load_MonsterNPC_MissingData(t *testing.T) {

}

func Test_Load_VendorNPC_MissingData(t *testing.T) {

}

func TestLoad_VendorNPC_Ok(t *testing.T) {

}

func TestLoad_VendorNPC_LoadPages(t *testing.T) {

}

func TestLoad_PortalNPC_Ok(t *testing.T) {

}

func TestLoadNPC_Ok_NC(t *testing.T) {

}

func TestLoad_MonsterNPC_NC(t *testing.T) {

}

func TestLoad_VendorNPC_NC(t *testing.T) {

}

func TestTrigger_NPC_Action(t *testing.T) {

}
