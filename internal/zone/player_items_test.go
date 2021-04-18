package zone

import (
	"github.com/shine-o/shine.engine.emulator/internal/pkg/data"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/errors"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/persistence"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/structs"
	"reflect"
	"sync"
	"testing"
)

//
func TestNewItem_Success(t *testing.T) {
	persistence.CleanDB()

	char := persistence.NewDummyCharacter("mage")

	player := &player{
		baseEntity: baseEntity{},
		persistence: &playerPersistence{
			char: char,
		},
	}

	err := player.load(char.Name)

	if err != nil {
		t.Fatal(err)
	}

	// item is not persisted here, only in memory
	item, _, err := makeItem("ShortStaff")

	if err != nil {
		t.Fatal(err)
	}

	// item is persisted here
	err = player.newItem(item)

	if err != nil {
		t.Fatal(err)
	}

	// check if item is in player inventory
	item1, ok := player.inventories.inventory.items[0]
	if !ok {
		t.Fail()
	}

	if item1.itemData.itemInfo.InxName != "ShortStaff" {
		t.Fail()
	}

}

func TestNewItem_WithAttributes(t *testing.T) {
	persistence.CleanDB()

	itemInxName := "KarenStaff"
	char := persistence.NewDummyCharacter("mage")

	player := &player{
		persistence: &playerPersistence{
			char: char,
		},
	}

	err := player.load(char.Name)

	if err != nil {
		t.Fatal(err)
	}

	// item is not persisted here, only in memory
	item, _, err := makeItem(itemInxName)

	if err != nil {
		t.Fatal(err)
	}

	// item is persisted here
	err = player.newItem(item)

	if err != nil {
		t.Fatal(err)
	}

	// check if item is in player inventory
	item1, ok := player.inventories.inventory.items[0]
	if !ok {
		t.Fail()
	}

	if item1.itemData.itemInfo.InxName != itemInxName {
		t.Fail()
	}

	amount := 0

	if item1.stats.strength.base > 0 || item1.stats.dexterity.base > 0 || item1.stats.endurance.base > 0 || item1.stats.intelligence.base > 0 || item1.stats.spirit.base > 0 {
		amount++
	}

	if amount == 0 {
		t.Fail()
	}

	// should have 2 static stats (97 int, 500 HP through GradeItemOption)
	if item1.stats.intelligence.base != 97 {
		t.Fatal("unexpected stat value")
	}

	if item1.stats.maxHP.base != 500 {
		t.Fatal("unexpected stat value")
	}

	if item1.stats.staticAttackSpeed.base != 1300 {
		t.Fatal("unexpected stat value")
	}

	if item1.stats.staticMinPAttack.base != 438 {
		t.Fatal("unexpected stat value")
	}

	if item1.stats.staticMaxPAttack.base != 673 {
		t.Fatal("unexpected stat value")
	}

	if item1.stats.staticMinMAttack.base != 2773 {
		t.Fatal("unexpected stat value")
	}

	if item1.stats.staticMaxMAttack.base != 4265 {
		t.Fatal("unexpected stat value")
	}

	if item1.stats.staticMAttackRate.base != 1000 {
		t.Fatal("unexpected stat value")
	}

	if item1.stats.staticPAttackRate.base != 1000 {
		t.Fatal("unexpected stat value")
	}

	if item1.stats.staticMDefenseRate.base != 1000 {
		t.Fatal("unexpected stat value")
	}

	if item1.stats.staticPDefenseRate.base != 1000 {
		t.Fatal("unexpected stat value")
	}

	if item1.stats.staticShieldDefenseRate.base != 0 {
		t.Fatal("unexpected stat value")
	}

	if item1.stats.staticCriticalRate.base != 6 {
		t.Fatal("unexpected stat value")
	}

	if item1.stats.staticMinPACriticalRate.base != 0 {
		t.Fatal("unexpected stat value")
	}

	if item1.stats.staticMaxPACriticalRate.base != 0 {
		t.Fatal("unexpected stat value")
	}

	if item1.stats.staticMinMACriticalRate.base != 0 {
		t.Fatal("unexpected stat value")
	}

	if item1.stats.staticMaxMACriticalRate.base != 0 {
		t.Fatal("unexpected stat value")
	}

	if item1.stats.physicalDefense.base != 0 {
		t.Fatal("unexpected stat value")
	}

	if item1.stats.magicalDefense.base != 0 {
		t.Fatal("unexpected stat value")
	}

	if item1.stats.aim.base != 1326 {
		t.Fatal("unexpected stat value")
	}

	if item1.stats.evasion.base != 0 {
		t.Fatal("unexpected stat value")
	}

	if item1.stats.staticAimRate.base != 0 {
		t.Fatal("unexpected stat value")
	}

	if item1.stats.staticEvasionRate.base != 0 {
		t.Fatal("unexpected stat value")
	}

	if item1.stats.staticPResistance.base != 0 {
		t.Fatal("unexpected stat value")
	}

	if item1.stats.staticDResistance.base != 0 {
		t.Fatal("unexpected stat value")
	}

	if item1.stats.staticCResistance.base != 0 {
		t.Fatal("unexpected stat value")
	}

	if item1.stats.staticMResistance.base != 0 {
		t.Fatal("unexpected stat value")
	}
}

func TestLoadItem_WithAttributes(t *testing.T) {
	persistence.CleanDB()

	itemInxName := "KarenStaff"
	char := persistence.NewDummyCharacter("mage")

	player := &player{
		baseEntity: baseEntity{},
		persistence: &playerPersistence{
			char: char,
		},
	}

	err := player.load(char.Name)

	if err != nil {
		t.Fatal(err)
	}

	// item is not persisted here, only in memory
	item, _, err := makeItem(itemInxName)

	if err != nil {
		t.Fatal(err)
	}

	// item is persisted here
	err = player.newItem(item)

	if err != nil {
		t.Fatal(err)
	}

	item1 := loadItem(item.pItem)

	if item1.itemData.itemInfo.InxName != itemInxName {
		t.Fail()
	}

	amount := 0

	if item1.stats.strength.base > 0 || item1.stats.dexterity.base > 0 || item1.stats.endurance.base > 0 || item1.stats.intelligence.base > 0 || item1.stats.spirit.base > 0 {
		amount++
	}

	if amount == 0 {
		t.Fail()
	}

	// should have 2 static stats (97 int, 500 HP through GradeItemOption)
	if item1.stats.intelligence.base != 97 {
		t.Fatal("unexpected stat value")
	}

	if item1.stats.maxHP.base != 500 {
		t.Fatal("unexpected stat value")
	}

	if item1.stats.staticAttackSpeed.base != 1300 {
		t.Fatal("unexpected stat value")
	}

	if item1.stats.staticMinPAttack.base != 438 {
		t.Fatal("unexpected stat value")
	}

	if item1.stats.staticMaxPAttack.base != 673 {
		t.Fatal("unexpected stat value")
	}

	if item1.stats.staticMinMAttack.base != 2773 {
		t.Fatal("unexpected stat value")
	}

	if item1.stats.staticMaxMAttack.base != 4265 {
		t.Fatal("unexpected stat value")
	}

	if item1.stats.staticMAttackRate.base != 1000 {
		t.Fatal("unexpected stat value")
	}

	if item1.stats.staticPAttackRate.base != 1000 {
		t.Fatal("unexpected stat value")
	}

	if item1.stats.staticMDefenseRate.base != 1000 {
		t.Fatal("unexpected stat value")
	}

	if item1.stats.staticPDefenseRate.base != 1000 {
		t.Fatal("unexpected stat value")
	}

	if item1.stats.staticShieldDefenseRate.base != 0 {
		t.Fatal("unexpected stat value")
	}

	if item1.stats.staticCriticalRate.base != 6 {
		t.Fatal("unexpected stat value")
	}

	if item1.stats.staticMinPACriticalRate.base != 0 {
		t.Fatal("unexpected stat value")
	}

	if item1.stats.staticMaxPACriticalRate.base != 0 {
		t.Fatal("unexpected stat value")
	}

	if item1.stats.staticMinMACriticalRate.base != 0 {
		t.Fatal("unexpected stat value")
	}

	if item1.stats.staticMaxMACriticalRate.base != 0 {
		t.Fatal("unexpected stat value")
	}

	if item1.stats.physicalDefense.base != 0 {
		t.Fatal("unexpected stat value")
	}

	if item1.stats.magicalDefense.base != 0 {
		t.Fatal("unexpected stat value")
	}

	if item1.stats.aim.base != 1326 {
		t.Fatal("unexpected stat value")
	}

	if item1.stats.evasion.base != 0 {
		t.Fatal("unexpected stat value")
	}

	if item1.stats.staticAimRate.base != 0 {
		t.Fatal("unexpected stat value")
	}

	if item1.stats.staticEvasionRate.base != 0 {
		t.Fatal("unexpected stat value")
	}

	if item1.stats.staticPResistance.base != 0 {
		t.Fatal("unexpected stat value")
	}

	if item1.stats.staticDResistance.base != 0 {
		t.Fatal("unexpected stat value")
	}

	if item1.stats.staticCResistance.base != 0 {
		t.Fatal("unexpected stat value")
	}

	if item1.stats.staticMResistance.base != 0 {
		t.Fatal("unexpected stat value")
	}
}

func TestNewItem_CreateAllItems(t *testing.T) {
	persistence.CleanDB()

	for _, row := range itemsData.ItemInfo.ShineRow {
		_, _, err := makeItem(row.InxName)
		if err != nil {
			t.Error(errors.Err{
				Code:    errors.UnitTestError,
				Message: "error creating item",
				Details: errors.ErrDetails{
					"err":       err,
					"itemIndex": row.InxName,
				},
			})
		}
	}
}

func TestNewItem_BadItemIndex(t *testing.T) {
	persistence.CleanDB()

	char := persistence.NewDummyCharacter("mage")

	player := &player{
		baseEntity: baseEntity{},
		persistence: &playerPersistence{
			char: char,
		},
	}

	err := player.load(char.Name)

	if err != nil {
		t.Fatal(err)
	}

	// item is not persisted here, only in memory
	_, _, err = makeItem("badindex")

	if err == nil {
		t.Fatal("expected error, got null")
	}
}

func Test_AllItems_NC(t *testing.T) {
	persistence.CleanDB()

	for _, row := range itemsData.ItemInfo.ShineRow {
		item, _, err := makeItem(row.InxName)
		if err != nil {
			continue
		}

		inc, err := protoItemPacketInformation(item)
		if err != nil {
			t.Error(errors.Err{
				Code:    errors.UnitTestError,
				Message: "error creating item nc struct",
				Details: errors.ErrDetails{
					"err":       err,
					"itemIndex": row.InxName,
					"nc":        inc,
				},
			})
			continue
		}
		switch row.Class {
		case data.ItemClassByteLot:
		case data.ItemUpRed:
		case data.ItemUpBlue:
		case data.ItemUpGold:
		case data.ItemFeed:
		case data.ItemClassSkillScroll:
		case data.ItemClassRecallScroll:
		case data.ItemClassUpsource:
		case data.ItemClassWtLicence:
		case data.ItemKq:
		case data.ItemGbCoin:
		case data.ItemNoEffect:
		case data.ItemEnchant:
			var nc structs.ShineItemAttrByteLot
			err := structs.Unpack(inc.ItemAttr, &nc)
			if err != nil {
				t.Error(errors.Err{
					Code:    errors.UnitTestError,
					Message: "bad ItemAttr NC for item class",
					Details: errors.ErrDetails{
						"err":       err,
						"itemIndex": row.InxName,
						"data":      inc.ItemAttr,
						"ncType":    reflect.TypeOf(nc).String(),
					},
				})
			}
			break
		case data.ItemClassWordLot:
		case data.ItemKqStep:
		case data.ItemActiveSkill:
			var nc structs.ShineItemAttrWordLot
			err := structs.Unpack(inc.ItemAttr, &nc)
			if err != nil {
				t.Error(errors.Err{
					Code:    errors.UnitTestError,
					Message: "bad ItemAttr NC for item class",
					Details: errors.ErrDetails{
						"err":       err,
						"itemIndex": row.InxName,
						"data":      inc.ItemAttr,
						"ncType":    reflect.TypeOf(nc).String(),
					},
				})
			}
			break
		case data.ItemClassDwrdLot:
			var nc structs.ShineItemAttrDwrdLot
			err := structs.Unpack(inc.ItemAttr, &nc)
			if err != nil {
				t.Error(errors.Err{
					Code:    errors.UnitTestError,
					Message: "bad ItemAttr NC for item class",
					Details: errors.ErrDetails{
						"err":       err,
						"itemIndex": row.InxName,
						"data":      inc.ItemAttr,
						"ncType":    reflect.TypeOf(nc).String(),
					},
				})
			}
			break
		case data.ItemClassQuestItem:
			var nc structs.ShineItemAttrQuestItem
			err := structs.Unpack(inc.ItemAttr, &nc)
			if err != nil {
				t.Error(errors.Err{
					Code:    errors.UnitTestError,
					Message: "bad ItemAttr NC for item class",
					Details: errors.ErrDetails{
						"err":       err,
						"itemIndex": row.InxName,
						"data":      inc.ItemAttr,
						"ncType":    reflect.TypeOf(nc).String(),
					},
				})
			}
			break
		case data.ItemClassAmulet:
			var nc structs.ShineItemAttrAmulet
			err := structs.Unpack(inc.ItemAttr, &nc)
			if err != nil {
				t.Error(errors.Err{
					Code:    errors.UnitTestError,
					Message: "bad ItemAttr NC for item class",
					Details: errors.ErrDetails{
						"err":       err,
						"itemIndex": row.InxName,
						"data":      inc.ItemAttr,
						"ncType":    reflect.TypeOf(nc).String(),
					},
				})
			}
			break
		case data.ItemClassWeapon:
			var nc structs.ShineItemAttrWeapon
			err := structs.Unpack(inc.ItemAttr, &nc)
			if err != nil {
				t.Error(errors.Err{
					Code:    errors.UnitTestError,
					Message: "bad ItemAttr NC for item class",
					Details: errors.ErrDetails{
						"err":       err,
						"itemIndex": row.InxName,
						"data":      inc.ItemAttr,
						"ncType":    reflect.TypeOf(nc).String(),
					},
				})
			}
			break
		case data.ItemClassArmor:
			var nc structs.ShineItemAttrArmor
			err := structs.Unpack(inc.ItemAttr, &nc)
			if err != nil {
				t.Error(errors.Err{
					Code:    errors.UnitTestError,
					Message: "bad ItemAttr NC for item class",
					Details: errors.ErrDetails{
						"err":       err,
						"itemIndex": row.InxName,
						"data":      inc.ItemAttr,
						"ncType":    reflect.TypeOf(nc).String(),
					},
				})
			}
			break
		case data.ItemClassShield:
			var nc structs.ShineItemAttrShield
			err := structs.Unpack(inc.ItemAttr, &nc)
			if err != nil {
				t.Error(errors.Err{
					Code:    errors.UnitTestError,
					Message: "bad ItemAttr NC for item class",
					Details: errors.ErrDetails{
						"err":       err,
						"itemIndex": row.InxName,
						"data":      inc.ItemAttr,
						"ncType":    reflect.TypeOf(nc).String(),
					},
				})
			}
			break
		case data.ItemClassBoot:
			var nc structs.ShineItemAttrBoot
			err := structs.Unpack(inc.ItemAttr, &nc)
			if err != nil {
				t.Error(errors.Err{
					Code:    errors.UnitTestError,
					Message: "bad ItemAttr NC for item class",
					Details: errors.ErrDetails{
						"err":       err,
						"itemIndex": row.InxName,
						"data":      inc.ItemAttr,
						"ncType":    reflect.TypeOf(nc).String(),
					},
				})
			}
			break
		case data.ItemClassFurniture:
			var nc structs.ShineItemAttrFurniture
			err := structs.Unpack(inc.ItemAttr, &nc)
			if err != nil {
				t.Error(errors.Err{
					Code:    errors.UnitTestError,
					Message: "bad ItemAttr NC for item class",
					Details: errors.ErrDetails{
						"err":       err,
						"itemIndex": row.InxName,
						"data":      inc.ItemAttr,
						"ncType":    reflect.TypeOf(nc).String(),
					},
				})
			}
			break
		case data.ItemClassDecoration:
			var nc structs.ShineItemAttrDecoration
			err := structs.Unpack(inc.ItemAttr, &nc)
			if err != nil {
				t.Error(errors.Err{
					Code:    errors.UnitTestError,
					Message: "bad ItemAttr NC for item class",
					Details: errors.ErrDetails{
						"err":       err,
						"itemIndex": row.InxName,
						"data":      inc.ItemAttr,
						"ncType":    reflect.TypeOf(nc).String(),
					},
				})
			}
			break
		case data.ItemClassBindItem:
			var nc structs.ShineItemAttrBindItem
			err := structs.Unpack(inc.ItemAttr, &nc)
			if err != nil {
				t.Error(errors.Err{
					Code:    errors.UnitTestError,
					Message: "bad ItemAttr NC for item class",
					Details: errors.ErrDetails{
						"err":       err,
						"itemIndex": row.InxName,
						"data":      inc.ItemAttr,
						"ncType":    reflect.TypeOf(nc).String(),
					},
				})
			}
			break
		case data.ItemClassItemChest:
			var nc structs.ShineItemAttrItemChest
			err := structs.Unpack(inc.ItemAttr, &nc)
			if err != nil {
				t.Error(errors.Err{
					Code:    errors.UnitTestError,
					Message: "bad ItemAttr NC for item class",
					Details: errors.ErrDetails{
						"err":       err,
						"itemIndex": row.InxName,
						"data":      inc.ItemAttr,
						"ncType":    reflect.TypeOf(nc).String(),
					},
				})
			}
			break
		case data.ItemHouseSkin:
			var nc structs.ShineItemAttrMiniHouseSkin
			err := structs.Unpack(inc.ItemAttr, &nc)
			if err != nil {
				t.Error(errors.Err{
					Code:    errors.UnitTestError,
					Message: "bad ItemAttr NC for item class",
					Details: errors.ErrDetails{
						"err":       err,
						"itemIndex": row.InxName,
						"data":      inc.ItemAttr,
						"ncType":    reflect.TypeOf(nc).String(),
					},
				})
			}
			break
		case data.ItemRiding:
			var nc structs.ShineItemAttrRiding
			err := structs.Unpack(inc.ItemAttr, &nc)
			if err != nil {
				t.Error(errors.Err{
					Code:    errors.UnitTestError,
					Message: "bad ItemAttr NC for item class",
					Details: errors.ErrDetails{
						"err":       err,
						"itemIndex": row.InxName,
						"data":      inc.ItemAttr,
						"ncType":    reflect.TypeOf(nc).String(),
					},
				})
			}
			break
		case data.ItemAmount:
			var nc structs.ShineItemAttrAmount
			err := structs.Unpack(inc.ItemAttr, &nc)
			if err != nil {
				t.Error(errors.Err{
					Code:    errors.UnitTestError,
					Message: "bad ItemAttr NC for item class",
					Details: errors.ErrDetails{
						"err":       err,
						"itemIndex": row.InxName,
						"data":      inc.ItemAttr,
						"ncType":    reflect.TypeOf(nc).String(),
					},
				})
			}
			break
		case data.ItemCosWeapon:
			var nc structs.ShineItemAttrCostumeWeapon
			err := structs.Unpack(inc.ItemAttr, &nc)
			if err != nil {
				t.Error(errors.Err{
					Code:    errors.UnitTestError,
					Message: "bad ItemAttr NC for item class",
					Details: errors.ErrDetails{
						"err":       err,
						"itemIndex": row.InxName,
						"data":      inc.ItemAttr,
						"ncType":    reflect.TypeOf(nc).String(),
					},
				})
			}
			break
		case data.ItemActionItem:
			var nc structs.ShineItemAttrActionItem
			err := structs.Unpack(inc.ItemAttr, &nc)
			if err != nil {
				t.Error(errors.Err{
					Code:    errors.UnitTestError,
					Message: "bad ItemAttr NC for item class",
					Details: errors.ErrDetails{
						"err":       err,
						"itemIndex": row.InxName,
						"data":      inc.ItemAttr,
						"ncType":    reflect.TypeOf(nc).String(),
					},
				})
			}
			break
		case data.ItemCapsule:
			var nc structs.ShineItemAttrCapsule
			err := structs.Unpack(inc.ItemAttr, &nc)
			if err != nil {
				t.Error(errors.Err{
					Code:    errors.UnitTestError,
					Message: "bad ItemAttr NC for item class",
					Details: errors.ErrDetails{
						"err":       err,
						"itemIndex": row.InxName,
						"data":      inc.ItemAttr,
						"ncType":    reflect.TypeOf(nc).String(),
					},
				})
			}
			break
		case data.ItemClosedCard:
			var nc structs.ShineItemAttrMobCardCollectClosed
			err := structs.Unpack(inc.ItemAttr, &nc)
			if err != nil {
				t.Error(errors.Err{
					Code:    errors.UnitTestError,
					Message: "bad ItemAttr NC for item class",
					Details: errors.ErrDetails{
						"err":       err,
						"itemIndex": row.InxName,
						"data":      inc.ItemAttr,
						"ncType":    reflect.TypeOf(nc).String(),
					},
				})
			}
			break
		case data.ItemOpenCard:
			var nc structs.ShineItemAttrMobCardCollect
			err := structs.Unpack(inc.ItemAttr, &nc)
			if err != nil {
				t.Error(errors.Err{
					Code:    errors.UnitTestError,
					Message: "bad ItemAttr NC for item class",
					Details: errors.ErrDetails{
						"err":       err,
						"itemIndex": row.InxName,
						"data":      inc.ItemAttr,
						"ncType":    reflect.TypeOf(nc).String(),
					},
				})
			}
			break
		case data.ItemMoney:
			//
			break
		case data.ItemPup:
			var nc structs.ShineItemAttrPet
			err := structs.Unpack(inc.ItemAttr, &nc)
			if err != nil {
				t.Error(errors.Err{
					Code:    errors.UnitTestError,
					Message: "bad ItemAttr NC for item class",
					Details: errors.ErrDetails{
						"err":       err,
						"itemIndex": row.InxName,
						"data":      inc.ItemAttr,
						"ncType":    reflect.TypeOf(nc).String(),
					},
				})
			}
			break
		case data.ItemCosShield:
			var nc structs.ShineItemAttrCostumeShield
			err := structs.Unpack(inc.ItemAttr, &nc)
			if err != nil {
				t.Error(errors.Err{
					Code:    errors.UnitTestError,
					Message: "bad ItemAttr NC for item class",
					Details: errors.ErrDetails{
						"err":       err,
						"itemIndex": row.InxName,
						"data":      inc.ItemAttr,
						"ncType":    reflect.TypeOf(nc).String(),
					},
				})
			}
			break
		case data.ItemBracelet:
			var nc structs.ShineItemAttrBracelet
			err := structs.Unpack(inc.ItemAttr, &nc)
			if err != nil {
				t.Error(errors.Err{
					Code:    errors.UnitTestError,
					Message: "bad ItemAttr NC for item class",
					Details: errors.ErrDetails{
						"err":       err,
						"itemIndex": row.InxName,
						"data":      inc.ItemAttr,
						"ncType":    reflect.TypeOf(nc).String(),
					},
				})
			}
			break
		}
	}
}

func Test_AllItems_With_Attributes_NC(t *testing.T) {
	persistence.CleanDB()

	for _, row := range itemsData.ItemInfo.ShineRow {
		item, icd, err := makeItem(row.InxName)
		if err != nil {
			continue
		}

		inc, err := protoItemPacketInformation(item)
		if err != nil {
			continue
		}

		if item.itemData.randomOption != nil && item.itemData.randomOptionCount != nil {
			switch row.Class {
			case data.ItemClassArmor:
				attr := structs.ShineItemAttrArmor{}
				err := structs.Unpack(inc.ItemAttr, &attr)
				if err != nil {
					t.Error(errors.Err{
						Code:    errors.UnitTestError,
						Message: "error serializing item attr nc struct",
						Details: errors.ErrDetails{
							"err":             err,
							"itemIndex":       row.InxName,
							"creationDetails": icd,
						},
					})
				}
				if (attr.Option.AmountBit >> 1) == 0 {
					t.Error(errors.Err{
						Code:    errors.UnitTestError,
						Message: "missing stats, expected at least 1 or more stats",
						Details: errors.ErrDetails{
							"itemIndex":       row.InxName,
							"rotIndex":        item.itemData.itemInfoServer.RandomOptionDropGroup,
							"creationDetails": icd,
						},
					})
				}
				break
			case data.ItemClassAmulet:
				attr := structs.ShineItemAttrAmulet{}
				err := structs.Unpack(inc.ItemAttr, &attr)
				if err != nil {
					t.Error(errors.Err{
						Code:    errors.UnitTestError,
						Message: "error serializing item attr nc struct",
						Details: errors.ErrDetails{
							"err":             err,
							"itemIndex":       row.InxName,
							"creationDetails": icd,
						},
					})
				}
				if (attr.Option.AmountBit >> 1) == 0 {
					t.Error(errors.Err{
						Code:    errors.UnitTestError,
						Message: "missing stats, expected at least 1 or more stats",
						Details: errors.ErrDetails{
							"itemIndex":       row.InxName,
							"rotIndex":        item.itemData.itemInfoServer.RandomOptionDropGroup,
							"creationDetails": icd,
						},
					})
				}
				break
			case data.ItemClassWeapon:
				attr := structs.ShineItemAttrWeapon{}
				err := structs.Unpack(inc.ItemAttr, &attr)
				if err != nil {
					t.Error(errors.Err{
						Code:    errors.UnitTestError,
						Message: "error serializing item attr nc struct",
						Details: errors.ErrDetails{
							"err":             err,
							"itemIndex":       row.InxName,
							"creationDetails": icd,
						},
					})
				}
				if (attr.Option.AmountBit >> 1) == 0 {
					t.Error(errors.Err{
						Code:    errors.UnitTestError,
						Message: "missing stats, expected at least 1 or more stats",
						Details: errors.ErrDetails{
							"itemIndex":       row.InxName,
							"rotIndex":        item.itemData.itemInfoServer.RandomOptionDropGroup,
							"creationDetails": icd,
						},
					})
				}
				break
			}
		}
	}
}

func TestNewItemStack_Success(t *testing.T) {
	t.Fail()

}

func TestNewItemStack_ItemNotStackable(t *testing.T) {
	t.Fail()

}

func TestSplitItemStack_Success(t *testing.T) {
	t.Fail()

}

func TestSplitItemStack_NC_Success(t *testing.T) {
	t.Fail()

}

func TestSplitItemStack_BadDivision(t *testing.T) {
	t.Fail()

}

func TestSplitItemStack_ItemNotStackable(t *testing.T) {
	t.Fail()

}

func TestSoftDeleteItem_Success(t *testing.T) {
	t.Fail()

}

func TestLoadNewPlayer_Mage_EquippedItems(t *testing.T) {
	// should have 1 staff
	t.Fail()

}

func TestLoadNewPlayer_Warrior_EquippedItems(t *testing.T) {
	t.Fail()
}

func TestLoadNewPlayer_Archer_EquippedItems(t *testing.T) {
	t.Fail()

}

func TestLoadNewPlayer_Cleric_EquippedItems(t *testing.T) {
	t.Fail()

}

func TestPlayer_PicksUpItem(t *testing.T) {
	t.Fail()

}

func TestPlayer_DropsItem(t *testing.T) {
	t.Fail()

}

func TestPlayer_DeletesItem(t *testing.T) {
	t.Fail()

}

func TestItemEquip_Success(t *testing.T) {
	t.Fail()

	//    SUCCESS = 641, // 0x0281
	//    FAILED = 645, // 0x0285
	//player := &player{
	//	baseEntity: baseEntity{
	//		handle: 1,
	//	},
	//}
	//
	//item := &item{}
	//
	//itemSlotChange, err := player.equip(item, data.ItemEquipHat)
	//
	//if err != nil {
	//	t.Fatal(err)
	//}
	//
	//if itemSlotChange.from != 0 {
	//	t.Fail()
	//}
	//
	//if itemSlotChange.to != 1 {
	//	t.Fail()
	//}
	//
	//equippedItem, ok := player.inventories.equipped.items[int(data.ItemEquipHat)]
	//
	//if !ok {
	//	t.Fail()
	//}
	//
	//if equippedItem.pItem.ID != item.pItem.ID {
	//	t.Fail()
	//}
	//
	//clauses := make(map[string]interface{})
	//
	//clauses["item_id"] = item.pItem.ID
	//clauses["character_id"] = player.char.ID
	//clauses["inventory_type"] = persistence.EquippedInventory
	//
	//_, err = persistence.GetItemWhere(clauses, false)
	//
	//if err != nil {
	//	t.Fail()
	//}
}

func TestItemEquip_NC_Success(t *testing.T) {
	t.Fail()
}

func TestItemEquip_Failed(t *testing.T) {
	// p.equipItem(item{}) (itemSlotChange{}, error)
	// err := error.(ErrorCodeZone)
	// err.Code = ItemEquipFailed
	// err.Details["pHandle"]
	//
	t.Fail()

}

func TestItemEquip_NC_Failed(t *testing.T) {
	//    FAILED = 645, // 0x0285
	// nc := itemEquipFailNc(err) structs.NcItemEquipFailNc ?
	//nc.Code == 645
	t.Fail()

}

func TestItemEquip_BadSlot(t *testing.T) {
	t.Fail()
}

func TestItemUnEquip_NC_Success(t *testing.T) {
	t.Fail()

}

func TestItemUnEquip_Success(t *testing.T) {
	t.Fail()

}

func TestChangeItemSlot_EmptySlot_Success(t *testing.T) {
	persistence.CleanDB()

	char := persistence.NewDummyCharacter("mage")

	player := &player{
		baseEntity: baseEntity{},
		persistence: &playerPersistence{
			char: char,
		},
		dz: &sync.RWMutex{},
	}

	err := player.load(char.Name)

	if err != nil {
		t.Fatal(err)
	}

	// item is not persisted here, only in memory
	item, _, err := makeItem("ShortStaff")

	if err != nil {
		t.Fatal(err)
	}

	err = player.newItem(item)

	if err != nil {
		t.Fatal(err)
	}

	nc := &structs.NcitemRelocateReq{
		From: structs.ItemInventory{
			Inventory: 9216,
		},
		To: structs.ItemInventory{
			Inventory: 9217,
		},
	}

	itemSlotChange, err := player.inventories.changeItemSlot(nc)

	if err != nil {
		t.Error(err)
	}

	// should be nil as I'm moving it to an empty slot
	if itemSlotChange.to.item != nil {
		t.Fatal(errors.Err{
			Code:    errors.UnitTestError,
			Message: "item should be nil",
			Details: errors.ErrDetails{
				"itemSlotChange": itemSlotChange,
			},
		})
	}

	if itemSlotChange.from.item.pItem.Slot != 1 {
		t.Fatalf("expected slot %v", 1)
	}

	i, ok := player.inventories.inventory.items[1]
	if !ok {
		t.Fatalf("expected an item in inventory, found none")
	}

	if i.pItem.ID != item.pItem.ID {
		t.Fatalf("distinct items were found")
	}
}

func TestChangeItemSlot_OccupiedSlot_Success(t *testing.T) {
	persistence.CleanDB()

	char := persistence.NewDummyCharacter("mage")

	player := &player{
		baseEntity: baseEntity{},
		persistence: &playerPersistence{
			char: char,
		},
		dz: &sync.RWMutex{},
	}

	err := player.load(char.Name)

	if err != nil {
		t.Fatal(err)
	}

	item, _, err := makeItem("ShortStaff")

	if err != nil {
		t.Fatal(err)
	}

	err = player.newItem(item)

	if err != nil {
		t.Fatal(err)
	}

	item1, _, err := makeItem("ShortStaff")

	if err != nil {
		t.Fatal(err)
	}

	err = player.newItem(item1)

	if err != nil {
		t.Fatal(err)
	}

	nc := &structs.NcitemRelocateReq{
		From: structs.ItemInventory{
			Inventory: 9216,
		},
		To: structs.ItemInventory{
			Inventory: 9217,
		},
	}

	itemSlotChange, err := player.inventories.changeItemSlot(nc)

	if err != nil {
		t.Error(err)
	}

	if itemSlotChange.to.item == nil {
		t.Fatal(errors.Err{
			Code:    errors.UnitTestError,
			Message: "item should not be nil",
			Details: errors.ErrDetails{
				"itemSlotChange": itemSlotChange,
			},
		})
	}

	if itemSlotChange.from.item.pItem.Slot != 1 {
		t.Fatalf("expected slot %v", 1)
	}

	i, ok := player.inventories.inventory.items[1]
	if !ok {
		t.Fatalf("expected an item in inventory, found none")
	}

	if i.pItem.ID != item.pItem.ID {
		t.Fatalf("distinct items were found")
	}

	//
	if itemSlotChange.to.item.pItem.Slot != 0 {
		t.Fatalf("expected slot %v", 0)
	}

	i1, ok := player.inventories.inventory.items[0]
	if !ok {
		t.Fatalf("expected an item in inventory, found none")
	}

	if i1.pItem.ID != item1.pItem.ID {
		t.Fatalf("distinct items were found")
	}

}

func TestChangeItemSlot_EmptySlot_NC(t *testing.T) {
	persistence.CleanDB()

	char := persistence.NewDummyCharacter("mage")

	player := &player{
		baseEntity: baseEntity{},
		persistence: &playerPersistence{
			char: char,
		},
		dz: &sync.RWMutex{},
	}

	err := player.load(char.Name)

	if err != nil {
		t.Fatal(err)
	}

	// item is not persisted here, only in memory
	item, _, err := makeItem("ShortStaff")

	if err != nil {
		t.Fatal(err)
	}

	err = player.newItem(item)

	if err != nil {
		t.Fatal(err)
	}

	nc := &structs.NcitemRelocateReq{
		From: structs.ItemInventory{
			Inventory: 9216,
		},
		To: structs.ItemInventory{
			Inventory: 9217,
		},
	}

	change, err := player.inventories.changeItemSlot(nc)

	if err != nil {
		t.Fatal(err)
	}

	enc1, enc2, err := ncItemCellChangeCmd(change)
	if err != nil {
		t.Fatal(err)
	}

	if enc1.Exchange.Inventory != nc.From.Inventory {
		t.Fatalf("mismatched from inventory")
	}

	if enc1.Location.Inventory != nc.To.Inventory {
		t.Fatalf("mismatched to inventory")
	}

	if len(enc1.Item.ItemAttr) == 0 {
		t.Fatalf("item attributes length should not be 0")
	}

	if enc2.Exchange.Inventory != nc.To.Inventory {
		t.Fatalf("mismatched from inventory")
	}

	if enc2.Location.Inventory != nc.From.Inventory {
		t.Fatalf("mismatched to inventory")
	}

	if len(enc2.Item.ItemAttr) != 0 {
		t.Fatalf("item attributes length should be 0")
	}

}

func TestChangeItemSlot_OccupiedSlot_NC(t *testing.T) {
	persistence.CleanDB()

	char := persistence.NewDummyCharacter("mage")

	player := &player{
		baseEntity: baseEntity{},
		persistence: &playerPersistence{
			char: char,
		},
		dz: &sync.RWMutex{},
	}

	err := player.load(char.Name)

	if err != nil {
		t.Fatal(err)
	}

	item, _, err := makeItem("ShortStaff")

	if err != nil {
		t.Fatal(err)
	}

	err = player.newItem(item)

	if err != nil {
		t.Fatal(err)
	}

	item1, _, err := makeItem("ShortStaff")

	if err != nil {
		t.Fatal(err)
	}

	err = player.newItem(item1)

	if err != nil {
		t.Fatal(err)
	}

	nc := &structs.NcitemRelocateReq{
		From: structs.ItemInventory{
			Inventory: 9216,
		},
		To: structs.ItemInventory{
			Inventory: 9217,
		},
	}

	change, err := player.inventories.changeItemSlot(nc)

	if err != nil {
		t.Fatal(err)
	}

	enc1, enc2, err := ncItemCellChangeCmd(change)
	if err != nil {
		t.Fatal(err)
	}

	if enc1.Exchange.Inventory != nc.From.Inventory {
		t.Fatalf("mismatched from inventory")
	}

	if enc1.Location.Inventory != nc.To.Inventory {
		t.Fatalf("mismatched to inventory")
	}

	if len(enc1.Item.ItemAttr) == 0 {
		t.Fatalf("item attributes length should not be 0")
	}

	if enc2.Exchange.Inventory != nc.To.Inventory {
		t.Fatalf("mismatched from inventory")
	}

	if enc2.Location.Inventory != nc.From.Inventory {
		t.Fatalf("mismatched to inventory")
	}

	if len(enc2.Item.ItemAttr) == 0 {
		t.Fatalf("item attributes length should not be 0")
	}

}

func TestChangeItem_NonExistentSlot(t *testing.T) {
	//    ERR_BOUND = 586, // 0x024A

	t.Fail()

}

func TestChangeItemSlot_BadItemType(t *testing.T) {
	t.Fail()

}

func TestChangeItemSlot_NoItemInSlot(t *testing.T) {
	t.Fail()
	//    ERR_BOUND = 586, // 0x024A

}

func TestSellItem_Success(t *testing.T) {
	t.Fail()

}

func TestSellItem_NonExistingItem(t *testing.T) {
	t.Fail()

}

func TestBuyItem_Success(t *testing.T) {
	t.Fail()

}

func TestOneUseItem_Success(t *testing.T) {
	t.Fail()

}

// Like mounts, quest items
func TestMultipleUseItem_Success(t *testing.T) {
	t.Fail()

}
