package zone

import (
	"testing"
)

func Test_Load_NPC_Ok(t *testing.T) {
	//zm := zoneMap{}
	//
	//err := zm.load()



}

func Test_Load_All_Npc(t *testing.T) {
	for _, npcs := range npcData.MapNPCs {
		for _, npc := range npcs {
			n, err := loadNpc(npc.MobIndex, 1, 1)
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

			if n.nType == npcUnknownRole {
				t.Errorf("value should not be nil %v", npc.MobIndex)
			}
		}
	}
}

func Test_Load_All_Monster_Npc(t *testing.T) {

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
