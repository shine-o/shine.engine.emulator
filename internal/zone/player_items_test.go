package zone

import (
	"reflect"
	"testing"

	"github.com/shine-o/shine.engine.emulator/internal/pkg/data"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/errors"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/persistence"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/structs"
)

func TestNewItemSuccess(t *testing.T) {
	persistence.CleanDB()

	char := persistence.NewDummyCharacter("mage", false, "dummy")

	player := &player{
		baseEntity: &baseEntity{},
		persistence: &playerPersistence{
			char: char,
		},
	}

	err := player.load(char.Name)
	if err != nil {
		t.Fatal(err)
	}

	// item is not persisted here, only in memory
	item, _, err := makeItem("ShortStaff", makeItemOptions{})
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

	if item1.data.itemInfo.InxName != "ShortStaff" {
		t.Fail()
	}
}

func TestNewItemWithAttributes(t *testing.T) {
	persistence.CleanDB()

	itemInxName := "KarenStaff"
	char := persistence.NewDummyCharacter("mage", false, "dummy")

	player := &player{
		baseEntity: &baseEntity{},
		persistence: &playerPersistence{
			char: char,
		},
	}

	err := player.load(char.Name)
	if err != nil {
		t.Fatal(err)
	}

	// item is not persisted here, only in memory
	item, _, err := makeItem(itemInxName, makeItemOptions{})
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

	if item1.data.itemInfo.InxName != itemInxName {
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

func TestLoadItemWithAttributes(t *testing.T) {
	persistence.CleanDB()

	itemInxName := "KarenStaff"
	char := persistence.NewDummyCharacter("mage", false, "dummy")

	player := &player{
		baseEntity: &baseEntity{},
		persistence: &playerPersistence{
			char: char,
		},
	}

	err := player.load(char.Name)
	if err != nil {
		t.Fatal(err)
	}

	// item is not persisted here, only in memory
	item, _, err := makeItem(itemInxName, makeItemOptions{})
	if err != nil {
		t.Fatal(err)
	}

	// item is persisted here
	err = player.newItem(item)

	if err != nil {
		t.Fatal(err)
	}

	item1 := loadItem(item.pItem)

	if item1.data.itemInfo.InxName != itemInxName {
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

func TestNewItemCreateAllItems(t *testing.T) {
	persistence.CleanDB()

	for _, row := range itemsData.ItemInfo.ShineRow {
		_, _, err := makeItem(row.InxName, makeItemOptions{})
		if err != nil {
			t.Error(err)
		}
	}
}

func TestNewItemBadItemIndex(t *testing.T) {
	persistence.CleanDB()

	char := persistence.NewDummyCharacter("mage", false, "dummy")

	player := &player{
		baseEntity: &baseEntity{},
		persistence: &playerPersistence{
			char: char,
		},
	}

	err := player.load(char.Name)
	if err != nil {
		t.Fatal(err)
	}

	// item is not persisted here, only in memory
	_, _, err = makeItem("badindex", makeItemOptions{})

	if err == nil {
		t.Fatal("expected error, got null")
	}
}

func Test_AllItemsNC(t *testing.T) {
	persistence.CleanDB()

	for _, row := range itemsData.ItemInfo.ShineRow {
		item, _, err := makeItem(row.InxName, makeItemOptions{})
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

func Test_AllItemsWithAttributesNC(t *testing.T) {
	persistence.CleanDB()
	for _, row := range itemsData.ItemInfo.ShineRow {
		item, icd, err := makeItem(row.InxName, makeItemOptions{})
		if err != nil {
			continue
		}

		inc, err := protoItemPacketInformation(item)
		if err != nil {
			continue
		}

		if item.data.randomOption != nil && item.data.randomOptionCount != nil {
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
							"rotIndex":        item.data.itemInfoServer.RandomOptionDropGroup,
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
							"rotIndex":        item.data.itemInfoServer.RandomOptionDropGroup,
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
							"rotIndex":        item.data.itemInfoServer.RandomOptionDropGroup,
							"creationDetails": icd,
						},
					})
				}
				break
			}
		}
	}
}

func TestNewItemStackSuccess(t *testing.T) {
	persistence.CleanDB()

	char := persistence.NewDummyCharacter("mage", false, "dummy")

	player := &player{
		baseEntity: &baseEntity{},
		persistence: &playerPersistence{
			char: char,
		},
	}

	err := player.load(char.Name)
	if err != nil {
		t.Fatal(err)
	}

	// item is not persisted here, only in memory
	item, _, err := makeItem("ShortStaff", makeItemOptions{})
	if err != nil {
		t.Fatal(err)
	}

	// item is persisted here
	err = player.newItem(item)
}

func TestNewItemStackItemNotStackable(t *testing.T) {
	t.Fail()
}

func TestSplitItemStackSuccess(t *testing.T) {
	t.Fail()
}

func TestSplitItemStackNCSuccess(t *testing.T) {
	t.Fail()
}

func TestSplitItemStackBadDivision(t *testing.T) {
	t.Fail()
}

func TestSplitItemStackItemNotStackable(t *testing.T) {
	t.Fail()
}

func TestSplitItemStackChangeSlotWhileSplitting(t *testing.T) {
	// should fail with error
	t.Fail()
}

func TestSoftDeleteItemSuccess(t *testing.T) {
	t.Fail()
}

func TestLoadNewPlayerMageEquippedItems(t *testing.T) {
	// should have 1 staff in slot 12
	persistence.CleanDB()

	char := persistence.NewDummyCharacter("mage", true, "dummy")

	player := &player{
		baseEntity: &baseEntity{},
		persistence: &playerPersistence{
			char: char,
		},
	}

	err := player.load(char.Name)
	if err != nil {
		t.Fatal(err)
	}

	item, ok := player.inventories.equipped.items[12]

	if !ok {
		t.Fatal("expected to have an item, got nil")
	}

	if item.data.itemInfo.InxName != "ShortStaff" {
		t.Fatal("unexpected item index")
	}
}

func TestLoadNewPlayerWarriorEquippedItems(t *testing.T) {
	// should have 1 staff in slot 12
	persistence.CleanDB()

	char := persistence.NewDummyCharacter("fighter", true, "dummy")

	player := &player{
		baseEntity: &baseEntity{},
		persistence: &playerPersistence{
			char: char,
		},
	}

	err := player.load(char.Name)
	if err != nil {
		t.Fatal(err)
	}

	item, ok := player.inventories.equipped.items[12]

	if !ok {
		t.Fatal("expected to have an item, got nil")
	}

	if item.data.itemInfo.InxName != "ShortSword" {
		t.Fatal("unexpected item index")
	}
}

func TestLoadNewPlayerArcherEquippedItems(t *testing.T) {
	persistence.CleanDB()

	char := persistence.NewDummyCharacter("archer", true, "dummy")

	player := &player{
		baseEntity: &baseEntity{},
		persistence: &playerPersistence{
			char: char,
		},
	}

	err := player.load(char.Name)
	if err != nil {
		t.Fatal(err)
	}

	item, ok := player.inventories.equipped.items[10]

	if !ok {
		t.Fatal("expected to have an item, got nil")
	}

	if item.data.itemInfo.InxName != "ShortBow" {
		t.Fatal("unexpected item index")
	}
}

func TestLoadNewPlayerClericEquippedItems(t *testing.T) {
	persistence.CleanDB()

	char := persistence.NewDummyCharacter("cleric", true, "dummy")

	player := &player{
		baseEntity: &baseEntity{},
		persistence: &playerPersistence{
			char: char,
		},
	}

	err := player.load(char.Name)
	if err != nil {
		t.Fatal(err)
	}

	item, ok := player.inventories.equipped.items[12]

	if !ok {
		t.Fatal("expected to have an item, got nil")
	}

	if item.data.itemInfo.InxName != "ShortMace" {
		t.Fatal("unexpected item index")
	}
}

func TestPlayerPicksUpItem(t *testing.T) {
	t.Fail()
}

func TestPlayerDropsItem(t *testing.T) {
	t.Fail()
}

func TestPlayerDeletesItem(t *testing.T) {
	t.Fail()
}

func TestItemEquipSuccess(t *testing.T) {
	persistence.CleanDB()

	char := persistence.NewDummyCharacter("mage", false, "dummy")

	player := &player{
		baseEntity: &baseEntity{},
		persistence: &playerPersistence{
			char: char,
		},
	}

	err := player.load(char.Name)
	if err != nil {
		t.Fatal(err)
	}

	// item is not persisted here, only in memory
	item, _, err := makeItem("ShortStaff", makeItemOptions{})
	if err != nil {
		t.Fatal(err)
	}

	err = player.newItem(item)

	if err != nil {
		t.Fatal(err)
	}

	nc := &structs.NcItemEquipReq{
		Slot: 0,
	}

	itemSlotChange, err := player.equip(int(nc.Slot))
	if err != nil {
		t.Fatal(err)
	}

	if itemSlotChange.gameFrom != 9216 {
		t.Fail()
	}

	if itemSlotChange.gameTo != 8204 {
		t.Fail()
	}

	if itemSlotChange.from.item == nil {
		t.Fatal("from item should not be nil")
	}

	if itemSlotChange.to.item != nil {
		t.Fatal("to item should be nil, as no item is equippedID")
	}

	equippedItem, ok := player.inventories.equipped.items[12]

	if !ok {
		t.Fatal("item is expected to be in slot")
	}

	_, ok = player.inventories.inventory.items[0]

	if ok {
		t.Fatal("item is NOT expected to be in slot")
	}

	if equippedItem.data.itemInfo.InxName != item.data.itemInfo.InxName || equippedItem.pItem.ID != item.pItem.ID {
		t.Fatalf("mismatched items")
	}

	if equippedItem.pItem.InventoryType != int(persistence.EquippedInventory) {
		t.Fatalf("unexpected inventory type")
	}

	if equippedItem.pItem.Slot != 12 {
		t.Fatalf("unexpected inventory type")
	}

	_, ok = player.inventories.inventory.items[0]

	if ok {
		t.Fatal("item is NOT expected to be in slot")
	}
}

func TestItemEquipSuccessReplaceItem(t *testing.T) {
	persistence.CleanDB()

	char := persistence.NewDummyCharacter("mage", false, "dummy")

	player := &player{
		baseEntity: &baseEntity{},
		persistence: &playerPersistence{
			char: char,
		},
	}

	err := player.load(char.Name)
	if err != nil {
		t.Fatal(err)
	}

	item, _, err := makeItem("ShortStaff", makeItemOptions{})
	if err != nil {
		t.Fatal(err)
	}

	err = player.newItem(item)

	if err != nil {
		t.Fatal(err)
	}

	// create new item that will take the slot 0 that was recently occupied
	item2, _, err := makeItem("ShortStaff", makeItemOptions{})
	if err != nil {
		t.Fatal(err)
	}

	err = player.newItem(item2)

	if err != nil {
		t.Fatal(err)
	}

	nc := &structs.NcItemEquipReq{
		Slot: 0,
	}

	_, err = player.equip(int(nc.Slot))

	if err != nil {
		t.Fatal(err)
	}

	nc2 := &structs.NcItemEquipReq{
		Slot: 1,
	}

	itemSlotChange2, err := player.equip(int(nc2.Slot))
	if err != nil {
		t.Fatal(err)
	}

	if itemSlotChange2.to.item == nil {
		t.Fatal("equippedID item should not be nil")
	}

	inventoryItem, ok := player.inventories.inventory.items[1]

	if !ok {
		t.Fatal("inventoryItem should not be nil")
	}

	// check that it matches the first equippedID item
	if inventoryItem.pItem.ID != item.pItem.ID {
		t.Fatal("items do not match")
	}

	equippedItem, ok := player.inventories.equipped.items[12]

	if !ok {
		t.Fatal("item is expected to be in slot")
	}

	// assert it matches the second equippedID item
	if equippedItem.pItem.ID != item2.pItem.ID {
		t.Fatalf("mismatched items")
	}
}

func TestItemEquipSuccessNC(t *testing.T) {
	persistence.CleanDB()

	char := persistence.NewDummyCharacter("mage", false, "dummy")

	player := &player{
		baseEntity: &baseEntity{},
		persistence: &playerPersistence{
			char: char,
		},
	}

	err := player.load(char.Name)
	if err != nil {
		t.Fatal(err)
	}

	// item is not persisted here, only in memory
	item, _, err := makeItem("ShortStaff", makeItemOptions{})
	if err != nil {
		t.Fatal(err)
	}

	err = player.newItem(item)

	if err != nil {
		t.Fatal(err)
	}

	nc := &structs.NcItemEquipReq{
		Slot: 0,
	}

	itemSlotChange, err := player.equip(int(nc.Slot))
	if err != nil {
		t.Fatal(err)
	}

	nc1, _, err := ncItemEquipChangeCmd(itemSlotChange)
	if err != nil {
		t.Fatal(err)
	}

	if nc1.EquipSlot != 12 {
		t.Fatal("unexpected equip slot")
	}

	if nc1.From.Inventory != 9216 {
		t.Fatal("unexpected from slot")
	}

	if nc1.ItemData.ItemID != 1750 {
		t.Fatal("unexpected item ID")
	}

	if len(nc1.ItemData.ItemAttr) == 0 {
		t.Fatal("unexpected itemattr length")
	}

	_, nc2, err := ncItemCellChangeCmd(itemSlotChange)
	if err != nil {
		t.Fatal(err)
	}

	if nc2.Exchange.Inventory != 8204 {
		t.Fatal("unexpected equip slot")
	}

	if nc2.Location.Inventory != 9216 {
		t.Fatal("unexpected inventory slot")
	}

	if len(nc2.Item.ItemAttr) != 0 {
		t.Fatal("unexpected length for item attributes")
	}
}

func TestItemEquipSuccessReplaceItemNC(t *testing.T) {
	persistence.CleanDB()

	char := persistence.NewDummyCharacter("mage", false, "dummy")

	player := &player{
		baseEntity: &baseEntity{},
		persistence: &playerPersistence{
			char: char,
		},
	}

	err := player.load(char.Name)
	if err != nil {
		t.Fatal(err)
	}

	item, _, err := makeItem("ShortStaff", makeItemOptions{})
	if err != nil {
		t.Fatal(err)
	}

	err = player.newItem(item)

	if err != nil {
		t.Fatal(err)
	}

	// create new item that will take the slot 0 that was recently occupied
	item2, _, err := makeItem("ShortStaff", makeItemOptions{})
	if err != nil {
		t.Fatal(err)
	}

	err = player.newItem(item2)

	if err != nil {
		t.Fatal(err)
	}

	nc := &structs.NcItemEquipReq{
		Slot: 0,
	}

	_, err = player.equip(int(nc.Slot))

	if err != nil {
		t.Fatal(err)
	}

	nc2 := &structs.NcItemEquipReq{
		Slot: 1,
	}

	itemSlotChange, err := player.equip(int(nc2.Slot))
	if err != nil {
		t.Fatal(err)
	}

	nc3, _, err := ncItemEquipChangeCmd(itemSlotChange)
	if err != nil {
		t.Fatal(err)
	}

	if nc3.EquipSlot != 12 {
		t.Fatal("unexpected equip slot")
	}

	if nc3.From.Inventory != 9217 {
		t.Fatal("unexpected from slot")
	}

	if nc3.ItemData.ItemID != 1750 {
		t.Fatal("unexpected item ID")
	}

	if len(nc3.ItemData.ItemAttr) == 0 {
		t.Fatal("unexpected itemattr length")
	}

	_, nc4, err := ncItemCellChangeCmd(itemSlotChange)
	if err != nil {
		t.Fatal(err)
	}

	if nc4.Exchange.Inventory != 8204 {
		t.Fatal("unexpected equip slot")
	}

	if nc4.Location.Inventory != 9217 {
		t.Fatal("unexpected inventory slot")
	}

	if len(nc4.Item.ItemAttr) == 0 {
		t.Fatal("unexpected length for item attributes")
	}
}

func TestItemEquipBadItemEquipOrClass(t *testing.T) {
	persistence.CleanDB()

	char := persistence.NewDummyCharacter("mage", false, "dummy")

	player := &player{
		baseEntity: &baseEntity{},
		persistence: &playerPersistence{
			char: char,
		},
	}

	err := player.load(char.Name)
	if err != nil {
		t.Fatal(err)
	}

	// item is not persisted here, only in memory
	item, _, err := makeItem("El5", makeItemOptions{})
	if err != nil {
		t.Fatal(err)
	}

	err = player.newItem(item)

	if err != nil {
		t.Fatal(err)
	}

	nc := &structs.NcItemEquipReq{
		Slot: 0,
	}

	_, err = player.equip(int(nc.Slot))

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	cErr, ok := err.(errors.Err)

	if !ok {
		t.Fatal("unexpected error type")
	}

	if cErr.Code != errors.ZoneItemEquipBadType {
		t.Fatal("unexpected error code")
	}
}

func TestItemEquipLowLevel(t *testing.T) {
	t.Fail()
}

func TestItemEquipNoItemInSlot(t *testing.T) {
	persistence.CleanDB()

	char := persistence.NewDummyCharacter("mage", false, "dummy")

	player := &player{
		baseEntity: &baseEntity{},
		persistence: &playerPersistence{
			char: char,
		},
	}

	err := player.load(char.Name)
	if err != nil {
		t.Fatal(err)
	}

	// item is not persisted here, only in memory
	item, _, err := makeItem("ShortStaff", makeItemOptions{})
	if err != nil {
		t.Fatal(err)
	}

	err = player.newItem(item)

	if err != nil {
		t.Fatal(err)
	}

	nc := &structs.NcItemEquipReq{
		Slot: 3,
	}

	_, err = player.equip(int(nc.Slot))

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	cErr, ok := err.(errors.Err)

	if !ok {
		t.Fatal("unexpected error type")
	}

	if cErr.Code != errors.ZoneItemSlotEquipNoItem {
		t.Fatal("unexpected error code")
	}
}

func TestItemUnEquipSuccess(t *testing.T) {
	// t.Fail()
	persistence.CleanDB()

	char := persistence.NewDummyCharacter("mage", false, "dummy")

	player := &player{
		baseEntity: &baseEntity{},
		persistence: &playerPersistence{
			char: char,
		},
	}

	err := player.load(char.Name)
	if err != nil {
		t.Fatal(err)
	}

	// item is not persisted here, only in memory
	item, _, err := makeItem("ShortStaff", makeItemOptions{})
	if err != nil {
		t.Fatal(err)
	}

	err = player.newItem(item)

	if err != nil {
		t.Fatal(err)
	}

	nc := &structs.NcItemEquipReq{
		Slot: 0,
	}

	_, err = player.equip(int(nc.Slot))

	if err != nil {
		t.Fatal(err)
	}

	nc1 := &structs.NcItemUnequipReq{
		SlotEquip: 12,
		SlotInven: 0,
	}

	itemSlotChange, err := player.unEquip(int(nc1.SlotEquip), int(nc1.SlotInven))
	if err != nil {
		t.Fatal(err)
	}

	if itemSlotChange.gameFrom != 8204 {
		t.Fatal("unexpected value")
	}

	if itemSlotChange.gameTo != 9216 {
		t.Fatal("unexpected value")
	}

	if itemSlotChange.to.item != nil {
		t.Fatal("from item should be nil")
	}

	if itemSlotChange.from.item == nil {
		t.Fatal("from item should not be nil")
	}
}

func TestItemUnEquipInventoryFull(t *testing.T) {
	persistence.CleanDB()

	char := persistence.NewDummyCharacter("mage", false, "dummy")

	player := &player{
		baseEntity: &baseEntity{},
		persistence: &playerPersistence{
			char: char,
		},
	}

	err := player.load(char.Name)
	if err != nil {
		t.Fatal(err)
	}

	// item is not persisted here, only in memory
	item, _, err := makeItem("ShortStaff", makeItemOptions{})
	if err != nil {
		t.Fatal(err)
	}

	err = player.newItem(item)

	if err != nil {
		t.Fatal(err)
	}

	nc := &structs.NcItemEquipReq{
		Slot: 0,
	}

	_, err = player.equip(int(nc.Slot))

	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < persistence.BagInventoryMax; i++ {
		item, _, err := makeItem("ShortStaff", makeItemOptions{})
		if err != nil {
			t.Fatal(err)
		}

		err = player.newItem(item)

		if err != nil {
			t.Fatal(err)
		}
	}

	nc1 := &structs.NcItemUnequipReq{
		SlotEquip: 12,
		SlotInven: 0,
	}

	_, err = player.unEquip(int(nc1.SlotEquip), int(nc1.SlotInven))

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	cErr, ok := err.(errors.Err)

	if !ok {
		t.Fatal("unexpected error type")
	}

	if cErr.Code != errors.ZoneItemSlotInUse {
		t.Fatalf("unexpected error code %v", cErr.Code)
	}
}

func TestItemUnEquipNonExistentSlot(t *testing.T) {
	persistence.CleanDB()

	char := persistence.NewDummyCharacter("mage", false, "dummy")

	player := &player{
		baseEntity: &baseEntity{},
		persistence: &playerPersistence{
			char: char,
		},
	}

	err := player.load(char.Name)
	if err != nil {
		t.Fatal(err)
	}

	// item is not persisted here, only in memory
	item, _, err := makeItem("ShortStaff", makeItemOptions{})
	if err != nil {
		t.Fatal(err)
	}

	err = player.newItem(item)

	if err != nil {
		t.Fatal(err)
	}

	nc := &structs.NcItemEquipReq{
		Slot: 0,
	}

	_, err = player.equip(int(nc.Slot))

	if err != nil {
		t.Fatal(err)
	}

	nc1 := &structs.NcItemUnequipReq{
		SlotEquip: 12,
		SlotInven: 255,
	}

	_, err = player.unEquip(int(nc1.SlotEquip), int(nc1.SlotInven))

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	cErr, ok := err.(errors.Err)

	if !ok {
		t.Fatalf("unexpected error type")
	}

	if cErr.Code != errors.PersistenceOutOfRangeSlot {
		t.Fatalf("unexpected error code")
	}
}

//func TestItemUnEquip_First_Slot_Available(t *testing.T) {
//	t.Fail()
//	// assert the item is moved to the first slot available in the inventory
//}
func TestOpenDepositInventorySuccess(t *testing.T) {
	persistence.CleanDB()

	char := persistence.NewDummyCharacter("mage", false, "dummy")

	player := &player{
		baseEntity: &baseEntity{},
		persistence: &playerPersistence{
			char: char,
		},
	}

	err := player.load(char.Name)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < persistence.DepositInventoryMax; i++ {
		item, _, err := makeItem("ShortStaff", makeItemOptions{
			overrideInventory: true,
			inventoryType:     persistence.DepositInventory,
		})
		if err != nil {
			t.Fatal(err)
		}

		err = player.newItem(item)

		if err != nil {
			t.Fatal(err)
		}
	}

	deposit := playerDeposit(player.inventories)

	for i, page := range deposit {
		if page.maxPages != 4 {
			t.Errorf("unexpected number of pages %v", page.maxPages)
		}

		if page.currentPage != i {
			t.Errorf("unexpected currentlySelected page %v", page.currentPage)
		}

		if len(page.items) != 36 {
			t.Errorf("unexpected item count %v", len(page.items))
		}
	}
}

func TestOpenDepositInventoryNCSuccess(t *testing.T) {
	persistence.CleanDB()

	char := persistence.NewDummyCharacter("mage", false, "dummy")

	player := &player{
		baseEntity: &baseEntity{},
		persistence: &playerPersistence{
			char: char,
		},
	}

	err := player.load(char.Name)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < persistence.DepositInventoryMax; i++ {
		item, _, err := makeItem("ShortStaff", makeItemOptions{
			overrideInventory: true,
			inventoryType:     persistence.DepositInventory,
		})
		if err != nil {
			t.Fatal(err)
		}

		err = player.newItem(item)

		if err != nil {
			t.Fatal(err)
		}
	}

	deposit := playerDeposit(player.inventories)
	// INFO : 2021/04/25 01:38:42.724707 handlers.go:272: 2021-04-25 01:38:42.715758 +0200 CEST 9120->40575 inbound NC_ACT_NPCMENUOPEN_REQ {"packetType":"small","length":4,"department":8,"command":"1C","opCode":8220,"data":"6c00","rawData":"041c206c00","friendlyName":""}
	// INFO : 2021/04/25 01:38:45.942966 handlers.go:272: 2021-04-25 01:38:45.932944 +0200 CEST 40575->9120 outbound NC_ACT_NPCMENUOPEN_ACK {"packetType":"small","length":3,"department":8,"command":"1D","opCode":8221,"data":"01","rawData":"031d2001","friendlyName":""}
	//  send as many of these packets as needed:
	//  inbound NC_MENU_OPENSTORAGE_CMD {"packetType":"big","length":1460,"department":15,"command":"8","opCode":15368,
	for i, page := range deposit {
		nc := ncMenuOpenStorageCmd(page)

		if nc.Cen != 0 {
			t.Errorf("unexpected value %v", nc.Cen)
		}

		if nc.CountItems != byte(len(page.items)) {
			t.Errorf("unexpected value %v", nc.CountItems)
		}

		if len(nc.Items) != len(page.items) {
			t.Errorf("unexpected value %v", len(nc.Items))
		}

		if nc.CurrentPage != byte(i) {
			t.Errorf("unexpected value %v", nc.CurrentPage)
		}

		if nc.MaxPage != byte(len(deposit)) {
			t.Errorf("unexpected value %v", nc.MaxPage)
		}

		if nc.OpenType != 0 {
			t.Errorf("unexpected value %v", nc.OpenType)
		}

		_, err := structs.Pack(&nc)
		if err != nil {
			t.Error(err)
		}
	}
}

func TestOpenRewardInventorySuccess(t *testing.T) {
	// INFO : 2021/04/25 15:47:00.663951 handlers.go:272: 2021-04-25 15:47:00.652328 +0200 CEST 30776->9120 outbound NC_ITEM_REWARDINVENOPEN_REQ {"packetType":"small","length":4,"department":12,"command":"2C","opCode":12332,"data":"0000","rawData":"042c300000","friendlyName":""}
	// INFO : 2021/04/25 15:47:00.817892 handlers.go:272: 2021-04-25 15:47:00.814858 +0200 CEST 9120->30776 inbound NC_ITEM_REWARDINVENOPEN_ACK {"packetType":"small","length":194,"department":12,"command":"2D","opCode":12333,"data":"0e0500087f0e02050108750e02050208840e020503087a0e02050408890e02050508930e010c06082c8100000000000000000c07082e8100000000000000000c08082f810000000000000000450908d03a00000000000000ffff00000000ffff00000000ffff00000000ffff000000000000000000000000000000000000000000ffffffffffffffffff00000000000000000c0a08298100000000000000000c0b082a810000000000000000050c088e75030c0d082981000000000000000036","rawData":"c22d300e0500087f0e02050108750e02050208840e020503087a0e02050408890e02050508930e010c06082c8100000000000000000c07082e8100000000000000000c08082f810000000000000000450908d03a00000000000000ffff00000000ffff00000000ffff00000000ffff000000000000000000000000000000000000000000ffffffffffffffffff00000000000000000c0a08298100000000000000000c0b082a810000000000000000050c088e75030c0d082981000000000000000036","friendlyName":""}
	// NcItemRewardInvenOpenAck

	// page 1
	// INFO : 2021/04/26 09:49:27.859090 handlers.go:272: 2021-04-26 09:49:27.848211 +0200 CEST 6233->9120 outbound NC_ITEM_REWARDINVENOPEN_REQ {"packetType":"small","length":4,"department":12,"command":"2C","opCode":12332,"data":"0000","rawData":"042c300000","friendlyName":""}
	// INFO : 2021/04/26 09:49:28.078648 handlers.go:272: 2021-04-26 09:49:28.074126 +0200 CEST 9120->6233 inbound NC_ITEM_REWARDINVENOPEN_ACK {"packetType":"big","length":424,"department":12,"command":"2D","opCode":12333,"data":"18450008023500000000000000ffff00000000ffff00000000ffff00000000ffff000000000000000000000000000000000000000000ffffffffffffffffff0000000000000000050108830e05450208e93600000000000000ffff00000000ffff00000000ffff00000000ffff000000000000000000000000000000000000000000ffffffffffffffffff00000000000000000503087e0e05050408920e05050508aa0e052a0608d10800000000000000000000000100000000000000000000000000000000000000000000000000002a0708d20800000000000000000000000100000000000000000000000000000000000000000000000000002a08084f090000000000000000000000010000000000000000000000000000000000000000000000000000050908ac0d052a0a0850090000000000000000000000010000000000000000000000000000000000000000000000000000050b08b10d05050c08840e04050d08491f0a050e08491f0a050f08491f0a051008491f0a051108491f0a051208431f02051308491f06051408f10b02051508ff0b23051608020c01051708040c0400","rawData":"00a8012d3018450008023500000000000000ffff00000000ffff00000000ffff00000000ffff000000000000000000000000000000000000000000ffffffffffffffffff0000000000000000050108830e05450208e93600000000000000ffff00000000ffff00000000ffff00000000ffff000000000000000000000000000000000000000000ffffffffffffffffff00000000000000000503087e0e05050408920e050

	// page 2
	// INFO : 2021/04/26 09:49:56.300297 handlers.go:272: 2021-04-26 09:49:56.294153 +0200 CEST 6233->9120 outbound NC_ITEM_REWARDINVENOPEN_REQ {"packetType":"small","length":4,"department":12,"command":"2C","opCode":12332,"data":"0100","rawData":"042c300100","friendlyName":""}
	// INFO : 2021/04/26 09:49:56.300297 handlers.go:272: 2021-04-26 09:49:56.299559 +0200 CEST 9120->6233 inbound NC_ACT_SOMEONEMOVERUN_CMD {"packetType":"small","length":24,"department":8,"command":"1A","opCode":8218,"data":"6b58f110000069120000951000004e120000db007265","rawData":"181a206b58f110000069120000951000004e120000db007265","friendlyName":""}
	// INFO : 2021/04/26 09:49:56.455715 handlers.go:272: 2021-04-26 09:49:56.452787 +0200 CEST 9120->6233 inbound NC_ITEM_REWARDINVENOPEN_ACK {"packetType":"small","length":76,"department":12,"command":"2D","opCode":12333,"data":"0a050008c40b32050108e90b22050208ce0b0c050308c40b0f050408ca0b32050508ca0b03050608b80b24110708f9a700000000000000000000000001050808af0f03050908340a120f","rawData":"4c2d300a050008c40b32050108e90b22050208ce0b0c050308c40b0f050408ca0b32050508ca0b03050608b80b24110708f9a700000000000000000000000001050808af0f03050908340a120f","friendlyName":""}

	t.Fail()
}

func TestOpenRewardInventoryNCSuccess(t *testing.T) {
	// INFO : 2021/04/25 15:47:00.663951 handlers.go:272: 2021-04-25 15:47:00.652328 +0200 CEST 30776->9120 outbound NC_ITEM_REWARDINVENOPEN_REQ {"packetType":"small","length":4,"department":12,"command":"2C","opCode":12332,"data":"0000","rawData":"042c300000","friendlyName":""}
	// INFO : 2021/04/25 15:47:00.817892 handlers.go:272: 2021-04-25 15:47:00.814858 +0200 CEST 9120->30776 inbound NC_ITEM_REWARDINVENOPEN_ACK {"packetType":"small","length":194,"department":12,"command":"2D","opCode":12333,"data":"0e0500087f0e02050108750e02050208840e020503087a0e02050408890e02050508930e010c06082c8100000000000000000c07082e8100000000000000000c08082f810000000000000000450908d03a00000000000000ffff00000000ffff00000000ffff00000000ffff000000000000000000000000000000000000000000ffffffffffffffffff00000000000000000c0a08298100000000000000000c0b082a810000000000000000050c088e75030c0d082981000000000000000036","rawData":"c22d300e0500087f0e02050108750e02050208840e020503087a0e02050408890e02050508930e010c06082c8100000000000000000c07082e8100000000000000000c08082f810000000000000000450908d03a00000000000000ffff00000000ffff00000000ffff00000000ffff000000000000000000000000000000000000000000ffffffffffffffffff00000000000000000c0a08298100000000000000000c0b082a810000000000000000050c088e75030c0d082981000000000000000036","friendlyName":""}
	t.Fail()
}

func TestOpenPremiumInventorySuccess(t *testing.T) {
	// INFO : 2021/04/25 17:39:35.246924 handlers.go:272: 2021-04-25 17:39:35.237462 +0200 CEST 50634->9120 outbound NC_ITEM_CHARGEDINVENOPEN_REQ {"packetType":"small","length":4,"department":12,"command":"20","opCode":12320,"data":"0000","rawData":"0420300000","friendlyName":""}
	// INFO : 2021/04/25 17:39:35.480942 handlers.go:272: 2021-04-25 17:39:35.477496 +0200 CEST 9120->50634 inbound NC_ITEM_CHARGEDINVENOPEN_ACK {"packetType":"big","length":391,"department":12,"command":"21","opCode":12321,
	t.Fail()
}

func TestOpenPremiumInventoryNCSuccess(t *testing.T) {
	// INFO : 2021/04/25 17:39:35.246924 handlers.go:272: 2021-04-25 17:39:35.237462 +0200 CEST 50634->9120 outbound NC_ITEM_CHARGEDINVENOPEN_REQ {"packetType":"small","length":4,"department":12,"command":"20","opCode":12320,"data":"0000","rawData":"0420300000","friendlyName":""}
	// INFO : 2021/04/25 17:39:35.480942 handlers.go:272: 2021-04-25 17:39:35.477496 +0200 CEST 9120->50634 inbound NC_ITEM_CHARGEDINVENOPEN_ACK {"packetType":"big","length":391,"department":12,"command":"21","opCode":12321,

	t.Fail()
}

func TestChangeItemSlotInventoryEmptySlotSuccess(t *testing.T) {
	persistence.CleanDB()

	char := persistence.NewDummyCharacter("mage", false, "dummy")

	player := &player{
		baseEntity: &baseEntity{},
		persistence: &playerPersistence{
			char: char,
		},
	}

	err := player.load(char.Name)
	if err != nil {
		t.Fatal(err)
	}

	// item is not persisted here, only in memory
	item, _, err := makeItem("ShortStaff", makeItemOptions{})
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

	itemSlotChange, err := player.inventories.moveItem(nc.From.Inventory, nc.To.Inventory)
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

	if itemSlotChange.from.item.pItem.InventoryType != int(persistence.BagInventory) {
		t.Fatalf("unexpected inventoryType %v", itemSlotChange.from.item.pItem.InventoryType)
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

func TestChangeItemSlotInventoryOccupiedSlotSuccess(t *testing.T) {
	persistence.CleanDB()

	char := persistence.NewDummyCharacter("mage", false, "dummy")

	player := &player{
		baseEntity: &baseEntity{},
		persistence: &playerPersistence{
			char: char,
		},
	}

	err := player.load(char.Name)
	if err != nil {
		t.Fatal(err)
	}

	item, _, err := makeItem("ShortStaff", makeItemOptions{})
	if err != nil {
		t.Fatal(err)
	}

	err = player.newItem(item)

	if err != nil {
		t.Fatal(err)
	}

	item1, _, err := makeItem("ShortStaff", makeItemOptions{})
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

	itemSlotChange, err := player.inventories.moveItem(nc.From.Inventory, nc.To.Inventory)
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

	if itemSlotChange.from.item.pItem.InventoryType != int(persistence.BagInventory) {
		t.Fatalf("unexpected inventoryType %v", itemSlotChange.from.item.pItem.InventoryType)
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

	if itemSlotChange.to.item.pItem.InventoryType != int(persistence.BagInventory) {
		t.Fatalf("unexpected inventoryType %v", itemSlotChange.to.item.pItem.InventoryType)
	}

	i1, ok := player.inventories.inventory.items[0]
	if !ok {
		t.Fatalf("expected an item in inventory, found none")
	}

	if i1.pItem.ID != item1.pItem.ID {
		t.Fatalf("distinct items were found")
	}
}

func TestChangeItemSlotInventoryEmptySlotNC(t *testing.T) {
	persistence.CleanDB()

	char := persistence.NewDummyCharacter("mage", false, "dummy")

	player := &player{
		baseEntity: &baseEntity{},
		persistence: &playerPersistence{
			char: char,
		},
	}

	err := player.load(char.Name)
	if err != nil {
		t.Fatal(err)
	}

	// item is not persisted here, only in memory
	item, _, err := makeItem("ShortStaff", makeItemOptions{})
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

	change, err := player.inventories.moveItem(nc.From.Inventory, nc.To.Inventory)
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

func TestChangeItemSlotInventoryOccupiedSlotNC(t *testing.T) {
	persistence.CleanDB()

	char := persistence.NewDummyCharacter("mage", false, "dummy")

	player := &player{
		baseEntity: &baseEntity{},
		persistence: &playerPersistence{
			char: char,
		},
	}

	err := player.load(char.Name)
	if err != nil {
		t.Fatal(err)
	}

	item, _, err := makeItem("ShortStaff", makeItemOptions{})
	if err != nil {
		t.Fatal(err)
	}

	err = player.newItem(item)

	if err != nil {
		t.Fatal(err)
	}

	item1, _, err := makeItem("ShortStaff", makeItemOptions{})
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

	change, err := player.inventories.moveItem(nc.From.Inventory, nc.To.Inventory)
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

func TestChangeItemInventoryNonExistentSlot(t *testing.T) {
	persistence.CleanDB()

	char := persistence.NewDummyCharacter("mage", false, "dummy")

	player := &player{
		baseEntity: &baseEntity{},
		persistence: &playerPersistence{
			char: char,
		},
	}

	err := player.load(char.Name)
	if err != nil {
		t.Fatal(err)
	}

	item, _, err := makeItem("ShortStaff", makeItemOptions{})
	if err != nil {
		t.Fatal(err)
	}

	err = player.newItem(item)

	if err != nil {
		t.Fatal(err)
	}

	item1, _, err := makeItem("ShortStaff", makeItemOptions{})
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
			Inventory: 9600,
		},
	}

	_, err = player.inventories.moveItem(nc.From.Inventory, nc.To.Inventory)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	cErr, ok := err.(errors.Err)

	if !ok {
		t.Fatalf("unexpected error type")
	}

	if cErr.Code != errors.PersistenceOutOfRangeSlot {
		t.Fatalf("unexpected error code")
	}
}

func TestChangeItemSlotInventoryNoItemInSlot(t *testing.T) {
	persistence.CleanDB()

	char := persistence.NewDummyCharacter("mage", false, "dummy")

	player := &player{
		baseEntity: &baseEntity{},
		persistence: &playerPersistence{
			char: char,
		},
	}

	err := player.load(char.Name)
	if err != nil {
		t.Fatal(err)
	}

	nc := &structs.NcitemRelocateReq{
		From: structs.ItemInventory{
			Inventory: 9250,
		},
		To: structs.ItemInventory{
			Inventory: 9218,
		},
	}

	_, err = player.inventories.moveItem(nc.From.Inventory, nc.To.Inventory)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	cErr, ok := err.(errors.Err)

	if !ok {
		t.Fatalf("unexpected error type")
	}

	if cErr.Code != errors.ZoneItemNoItemInSlot {
		t.Fatalf("unexpected error code %v", cErr.Code)
	}
}

func TestChangeItemSlotInventoryInDropState(t *testing.T) {
	t.Fail()
}

func TestChangeItemSlotInventoryToDepositSuccess(t *testing.T) {
	t.Fail()
}

func TestChangeItemSlotInventoryToMHInventorySuccess(t *testing.T) { t.Fail() }

func TestChangeItemSlotInventoryToRewardInventoryShouldFail(t *testing.T) { t.Fail() }

func TestChangeItemSlotInventoryToPremiumInventoryShould_Fail(t *testing.T) { t.Fail() }

func TestChangeItemSlotDepositEmptySlotSuccess(t *testing.T) {
	persistence.CleanDB()

	char := persistence.NewDummyCharacter("mage", false, "dummy")

	player := &player{
		baseEntity: &baseEntity{},
		persistence: &playerPersistence{
			char: char,
		},
	}

	err := player.load(char.Name)
	if err != nil {
		t.Fatal(err)
	}

	// item is not persisted here, only in memory
	item, _, err := makeItem("ShortStaff", makeItemOptions{
		overrideInventory: true,
		inventoryType:     persistence.DepositInventory,
	})
	if err != nil {
		t.Fatal(err)
	}

	err = player.newItem(item)

	if err != nil {
		t.Fatal(err)
	}

	nc := &structs.NcitemRelocateReq{
		From: structs.ItemInventory{
			Inventory: 6144,
		},
		To: structs.ItemInventory{
			Inventory: 6145,
		},
	}

	itemSlotChange, err := player.inventories.moveItem(nc.From.Inventory, nc.To.Inventory)
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

	if itemSlotChange.from.item.pItem.InventoryType != int(persistence.DepositInventory) {
		t.Fatalf("unexpected inventoryType %v", itemSlotChange.from.item.pItem.InventoryType)
	}

	if itemSlotChange.from.item.pItem.Slot != 1 {
		t.Fatalf("unexpected slot %v", itemSlotChange.from.item.pItem.Slot)
	}

	i, ok := player.inventories.deposit.items[1]
	if !ok {
		t.Fatalf("expected an item in inventory, found none")
	}

	if i.pItem.ID != item.pItem.ID {
		t.Fatalf("distinct items were found")
	}
}

func TestChangeItemSlotDepositOccupiedSlotSuccess(t *testing.T) {
	persistence.CleanDB()

	char := persistence.NewDummyCharacter("mage", false, "dummy")

	player := &player{
		baseEntity: &baseEntity{},
		persistence: &playerPersistence{
			char: char,
		},
	}

	err := player.load(char.Name)
	if err != nil {
		t.Fatal(err)
	}

	item, _, err := makeItem("ShortStaff", makeItemOptions{
		overrideInventory: true,
		inventoryType:     persistence.DepositInventory,
	})
	if err != nil {
		t.Fatal(err)
	}

	err = player.newItem(item)

	if err != nil {
		t.Fatal(err)
	}

	item1, _, err := makeItem("ShortStaff", makeItemOptions{
		overrideInventory: true,
		inventoryType:     persistence.DepositInventory,
	})
	if err != nil {
		t.Fatal(err)
	}

	err = player.newItem(item1)

	if err != nil {
		t.Fatal(err)
	}

	nc := &structs.NcitemRelocateReq{
		From: structs.ItemInventory{
			Inventory: 6144,
		},
		To: structs.ItemInventory{
			Inventory: 6145,
		},
	}

	itemSlotChange, err := player.inventories.moveItem(nc.From.Inventory, nc.To.Inventory)
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

	if itemSlotChange.from.item.pItem.InventoryType != int(persistence.DepositInventory) {
		t.Fatalf("unexpected inventoryType %v", itemSlotChange.from.item.pItem.InventoryType)
	}

	i, ok := player.inventories.deposit.items[1]
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

	if itemSlotChange.to.item.pItem.InventoryType != int(persistence.DepositInventory) {
		t.Fatalf("unexpected inventoryType %v", itemSlotChange.to.item.pItem.InventoryType)
	}

	i1, ok := player.inventories.deposit.items[0]
	if !ok {
		t.Fatalf("expected an item in inventory, found none")
	}

	if i1.pItem.ID != item1.pItem.ID {
		t.Fatalf("distinct items were found")
	}
}

func TestChangeItemSlotDepositNonExistentSlot(t *testing.T) {
	persistence.CleanDB()

	char := persistence.NewDummyCharacter("mage", false, "dummy")

	player := &player{
		baseEntity: &baseEntity{},
		persistence: &playerPersistence{
			char: char,
		},
	}

	err := player.load(char.Name)
	if err != nil {
		t.Fatal(err)
	}

	item, _, err := makeItem("ShortStaff", makeItemOptions{
		overrideInventory: true,
		inventoryType:     persistence.DepositInventory,
	})
	if err != nil {
		t.Fatal(err)
	}

	err = player.newItem(item)

	if err != nil {
		t.Fatal(err)
	}

	item1, _, err := makeItem("ShortStaff", makeItemOptions{
		overrideInventory: true,
		inventoryType:     persistence.DepositInventory,
	})
	if err != nil {
		t.Fatal(err)
	}

	err = player.newItem(item1)

	if err != nil {
		t.Fatal(err)
	}

	nc := &structs.NcitemRelocateReq{
		From: structs.ItemInventory{
			Inventory: 6144,
		},
		To: structs.ItemInventory{
			Inventory: 6900,
		},
	}

	_, err = player.inventories.moveItem(nc.From.Inventory, nc.To.Inventory)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	cErr, ok := err.(errors.Err)

	if !ok {
		t.Fatalf("unexpected error type")
	}

	if cErr.Code != errors.PersistenceOutOfRangeSlot {
		t.Fatalf("unexpected error code")
	}
}

func TestChangeItemSlotDepositNoItemInSlot(t *testing.T) {
	persistence.CleanDB()

	char := persistence.NewDummyCharacter("mage", false, "dummy")

	player := &player{
		baseEntity: &baseEntity{},
		persistence: &playerPersistence{
			char: char,
		},
	}

	err := player.load(char.Name)
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

	_, err = player.inventories.moveItem(nc.From.Inventory, nc.To.Inventory)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	cErr, ok := err.(errors.Err)

	if !ok {
		t.Fatalf("unexpected error type")
	}

	if cErr.Code != errors.ZoneItemNoItemInSlot {
		t.Fatalf("unexpected error code %v", cErr.Code)
	}
}

func TestChangeItemSlotDepositToInventorySuccess(t *testing.T) {
	persistence.CleanDB()

	char := persistence.NewDummyCharacter("mage", false, "dummy")

	player := &player{
		baseEntity: &baseEntity{},
		persistence: &playerPersistence{
			char: char,
		},
	}

	err := player.load(char.Name)
	if err != nil {
		t.Fatal(err)
	}

	// item is not persisted here, only in memory
	item, _, err := makeItem("ShortStaff", makeItemOptions{})
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
			Inventory: 6145,
		},
	}

	itemSlotChange, err := player.inventories.moveItem(nc.From.Inventory, nc.To.Inventory)
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

	if itemSlotChange.from.item.pItem.InventoryType != int(persistence.DepositInventory) {
		t.Fatalf("unexpected inventoryType %v", itemSlotChange.from.item.pItem.InventoryType)
	}

	if itemSlotChange.from.item.pItem.Slot != 1 {
		t.Fatalf("unexpected slot %v", itemSlotChange.from.item.pItem.Slot)
	}

	i, ok := player.inventories.deposit.items[1]
	if !ok {
		t.Fatalf("expected an item in inventory, found none")
	}

	if i.pItem.ID != item.pItem.ID {
		t.Fatalf("distinct items were found")
	}
}

func TestChangeItemSlotDepositToPremiumInventoryShouldFail(t *testing.T) { t.Fail() }

func TestChangeItemSlotDepositToRewardInventoryShouldFail(t *testing.T) { t.Fail() }

func TestChangeItemSlotMHInventoryToInventoryShouldFail(t *testing.T) { t.Fail() }

func TestChangeItemSlotRewardInventoryShouldFail(t *testing.T) {
	persistence.CleanDB()

	char := persistence.NewDummyCharacter("mage", false, "dummy")

	player := &player{
		baseEntity: &baseEntity{},
		persistence: &playerPersistence{
			char: char,
		},
	}

	err := player.load(char.Name)
	if err != nil {
		t.Fatal(err)
	}

	// item is not persisted here, only in memory
	item, _, err := makeItem("ShortStaff", makeItemOptions{
		overrideInventory: true,
		inventoryType:     persistence.RewardInventory,
	})
	if err != nil {
		t.Fatal(err)
	}

	err = player.newItem(item)

	if err != nil {
		t.Fatal(err)
	}

	nc := &structs.NcitemRelocateReq{
		From: structs.ItemInventory{
			Inventory: 2048,
		},
		To: structs.ItemInventory{
			Inventory: 2049,
		},
	}

	_, err = player.inventories.moveItem(nc.From.Inventory, nc.To.Inventory)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	cErr, ok := err.(errors.Err)
	if !ok {
		t.Fatalf("unexpected error type %v", reflect.TypeOf(cErr).String())
	}

	if cErr.Code != errors.ZoneItemSlotChangeConstraint {
		t.Fatalf("unexpected error code %v", cErr.Code)
	}
}

func TestChangeItemSlotRewardInventoryToInventorySuccess(t *testing.T) {
	persistence.CleanDB()

	char := persistence.NewDummyCharacter("mage", false, "dummy")

	player := &player{
		baseEntity: &baseEntity{},
		persistence: &playerPersistence{
			char: char,
		},
	}

	err := player.load(char.Name)
	if err != nil {
		t.Fatal(err)
	}

	// item is not persisted here, only in memory
	item, _, err := makeItem("ShortStaff", makeItemOptions{
		overrideInventory: true,
		inventoryType:     persistence.RewardInventory,
	})
	if err != nil {
		t.Fatal(err)
	}

	err = player.newItem(item)

	if err != nil {
		t.Fatal(err)
	}

	nc := &structs.NcitemRelocateReq{
		From: structs.ItemInventory{
			Inventory: 2048,
		},
		To: structs.ItemInventory{
			Inventory: 9217,
		},
	}

	itemSlotChange, err := player.inventories.moveItem(nc.From.Inventory, nc.To.Inventory)
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

	if itemSlotChange.from.item.pItem.InventoryType != int(persistence.BagInventory) {
		t.Fatalf("unexpected inventoryType %v", itemSlotChange.from.item.pItem.InventoryType)
	}

	if itemSlotChange.from.item.pItem.Slot != 1 {
		t.Fatalf("unexpected slot %v", itemSlotChange.from.item.pItem.Slot)
	}

	i, ok := player.inventories.inventory.items[1]

	if !ok {
		t.Fatalf("expected an item in inventory, found none")
	}

	if i.pItem.ID != item.pItem.ID {
		t.Fatalf("distinct items were found")
	}
}

func TestChangeItemSlotRewardInventoryToInventoryInventoryFull(t *testing.T) {
	persistence.CleanDB()

	char := persistence.NewDummyCharacter("mage", false, "dummy")

	player := &player{
		baseEntity: &baseEntity{},
		persistence: &playerPersistence{
			char: char,
		},
	}

	err := player.load(char.Name)
	if err != nil {
		t.Fatal(err)
	}

	// item is not persisted here, only in memory
	item, _, err := makeItem("ShortStaff", makeItemOptions{
		overrideInventory: true,
		inventoryType:     persistence.RewardInventory,
	})
	if err != nil {
		t.Fatal(err)
	}

	err = player.newItem(item)

	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < persistence.BagInventoryMax; i++ {
		item, _, err := makeItem("ShortStaff", makeItemOptions{})
		if err != nil {
			t.Fatal(err)
		}

		err = player.newItem(item)

		if err != nil {
			t.Fatal(err)
		}
	}

	nc := &structs.NcitemRelocateReq{
		From: structs.ItemInventory{
			Inventory: 2048,
		},
		To: structs.ItemInventory{
			Inventory: 9217,
		},
	}

	_, err = player.inventories.moveItem(nc.From.Inventory, nc.To.Inventory)

	cErr, ok := err.(errors.Err)
	if !ok {
		t.Fatal("unexpected error type", reflect.TypeOf(cErr).String())
	}

	if cErr.Code != errors.PersistenceInventoryFull {
		t.Fatal("unexpected error code", cErr.Code)
	}
}

func TestChangeItemSlotRewardInventoryToPremiumInventoryShouldFail(t *testing.T) {
	persistence.CleanDB()

	char := persistence.NewDummyCharacter("mage", false, "dummy")

	player := &player{
		baseEntity: &baseEntity{},
		persistence: &playerPersistence{
			char: char,
		},
	}

	err := player.load(char.Name)
	if err != nil {
		t.Fatal(err)
	}

	// item is not persisted here, only in memory
	item, _, err := makeItem("ShortStaff", makeItemOptions{
		overrideInventory: true,
		inventoryType:     persistence.RewardInventory,
	})
	if err != nil {
		t.Fatal(err)
	}

	err = player.newItem(item)

	if err != nil {
		t.Fatal(err)
	}

	nc := &structs.NcitemRelocateReq{
		From: structs.ItemInventory{
			Inventory: 2048,
		},
		To: structs.ItemInventory{
			Inventory: 2049,
		},
	}

	_, err = player.inventories.moveItem(nc.From.Inventory, nc.To.Inventory)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	cErr, ok := err.(errors.Err)
	if !ok {
		t.Fatalf("unexpected error type %v", reflect.TypeOf(cErr).String())
	}

	if cErr.Code != errors.ZoneItemSlotChangeConstraint {
		t.Fatalf("unexpected error code %v", cErr.Code)
	}
}

func TestChangeItemSlotPremiumInventoryToInventorySuccess(t *testing.T) {
	t.Fail()
	// INFO : 2021/04/25 17:42:10.691919 handlers.go:272: 2021-04-25 17:42:10.682141 +0200 CEST 50634->9120 outbound NC_ITEM_CHARGED_WITHDRAW_REQ {"packetType":"small","length":6,"department":12,"command":"22","opCode":12322,"data":"c449b701","rawData":"062230c449b701","friendlyName":""}
	// INFO : 2021/04/25 17:42:10.877594 handlers.go:272: 2021-04-25 17:42:10.874376 +0200 CEST 9120->50634 inbound NC_ITEM_CELLCHANGE_CMD {"packetType":"small","length":9,"department":12,"command":"1","opCode":12289,"data":"152415249a7501","rawData":"090130152415249a7501","friendlyName":""}
	// INFO : 2021/04/25 17:42:10.877604 handlers.go:272: 2021-04-25 17:42:10.874376 +0200 CEST 9120->50634 inbound NC_ITEM_CHARGED_WITHDRAW_ACK {"packetType":"small","length":8,"department":12,"command":"23","opCode":12323,"data":"c449b7014110","rawData":"082330c449b7014110","friendlyName":""}
}

func TestChangeItemSlotPremiumInventoryToInventoryNCSuccess(t *testing.T) {
	t.Fail()
	// page 1
	// INFO : 2021/04/25 17:39:35.246924 handlers.go:272: 2021-04-25 17:39:35.237462 +0200 CEST 50634->9120 outbound NC_ITEM_CHARGEDINVENOPEN_REQ {"packetType":"small","length":4,"department":12,"command":"20","opCode":12320,"data":"0000","rawData":"0420300000","friendlyName":""}
	// page 2
	// INFO : 2021/04/25 18:20:31.407229 handlers.go:272: 2021-04-25 18:20:31.392034 +0200 CEST 50634->9120 outbound NC_ITEM_CHARGEDINVENOPEN_REQ {"packetType":"small","length":4,"department":12,"command":"20","opCode":12320,"data":"0100","rawData":"0420300100","friendlyName":""}

	// page 1
	// INFO : 2021/04/25 18:18:08.338040 handlers.go:272: 2021-04-25 18:18:08.335947 +0200 CEST 9120->50634 inbound NC_ITEM_CHARGEDINVENOPEN_ACK {"packetType":"big","length":391,"department":12,"command":"21","opCode":12321,"data":"4110031800d2861100c2280000000000000de23808336d5700d3280000010000000e8c315c98e89101b028000001000000108790049b0b2202c12700000100000011854176ab5f2202c22700000100000011854176cc072302c42700000100000011854900dc5b2302c527000001000000118549001fac24026d2800000100000011854902300025026e280000010000001185490455aa4202ab5500000100000013e4483456aa4202ab5500000100000013e4483457aa4202ab5500000100000013e448343bdc4202aa2800000100000013a45136a10443027b2700000100000013045230c2044302a12900000100000013045238c3044302a85500000100000013045238c4044302a85500000100000013045238c5044302842b00000100000013045238c7044302542800000100000013045238c8044302a1290000010000001304523ac9044302822b0000010000001304523aca044302842b0000010000001304523ccd044302d4290000010000001304523ed3044302a85500000100000013045244","rawData":"00870121304110031800d2861100c2280000000000000de23808336d5700d3280000010000000e8c315c98e89101b028000001000000108790049b0b2202c12700000100000011854176ab5f2202c22700000100000011854176cc072302c42700000100000011854900dc5b2302c527000001000000118549001fac24026d2800000100000011854902300025026e280000010000001185490455aa4202ab5500000100000013e4483456aa4202ab5500000100000013e4483457aa4202ab5500000100000013e448343bdc4202aa2800000100000013a45136a10443027b2700000100000013045230c2044302a12900000100000013045238c3044302a85500000100000013045238c4044302a85500000100000013045238c5044302842b00000100000013045238c7044302542800000100000013045238c8044302a1290000010000001304523ac9044302822b0000010000001304523aca044302842b0000010000001304523ccd044302d4290000010000001304523ed3044302a85500000100000013045244","friendlyName":""}
	// page 2 ? how the fuck does this work?
	// INFO : 2021/04/25 18:20:31.654372 handlers.go:272: 2021-04-25 18:20:31.652235 +0200 CEST 9120->50634 inbound NC_ITEM_CHARGEDINVENOPEN_ACK {"packetType":"small","length":231,"department":12,"command":"21","opCode":12321,"data":"4110030e00d40443028c5300000100000013045246d6044302a85500000100000013045246d8044302822b00000100000013045248da044302a85500000100000013045248db044302822b00000100000013045248e0044302842b00000100000013045248e20443028c530000010000001304524ae5044302a8550000010000001304524af4044302822b00000100000013045250f9044302842b00000100000013045252fe044302a8550000010000001304525407054302a8550000010000001304525a4e054302ce2b00000100000013045a0c2c7e5d02317500000100000014085276","rawData":"e721304110030e00d40443028c5300000100000013045246d6044302a85500000100000013045246d8044302822b00000100000013045248da044302a85500000100000013045248db044302822b00000100000013045248e0044302842b00000100000013045248e20443028c530000010000001304524ae5044302a8550000010000001304524af4044302822b00000100000013045250f9044302842b00000100000013045252fe044302a8550000010000001304525407054302a8550000010000001304525a4e054302ce2b00000100000013045a0c2c7e5d02317500000100000014085276","friendlyName":""}

	// INFO : 2021/04/25 17:42:10.691919 handlers.go:272: 2021-04-25 17:42:10.682141 +0200 CEST 50634->9120 outbound NC_ITEM_CHARGED_WITHDRAW_REQ {"packetType":"small","length":6,"department":12,"command":"22","opCode":12322,"data":"c449b701","rawData":"062230c449b701","friendlyName":""}
	// INFO : 2021/04/25 17:42:10.877594 handlers.go:272: 2021-04-25 17:42:10.874376 +0200 CEST 9120->50634 inbound NC_ITEM_CELLCHANGE_CMD {"packetType":"small","length":9,"department":12,"command":"1","opCode":12289,"data":"152415249a7501","rawData":"090130152415249a7501","friendlyName":""}
	// INFO : 2021/04/25 17:42:10.877604 handlers.go:272: 2021-04-25 17:42:10.874376 +0200 CEST 9120->50634 inbound NC_ITEM_CHARGED_WITHDRAW_ACK {"packetType":"small","length":8,"department":12,"command":"23","opCode":12323,"data":"c449b7014110","rawData":"082330c449b7014110","friendlyName":""}

	// INFO : 2021/04/25 17:59:04.123989 handlers.go:272: 2021-04-25 17:59:04.122346 +0200 CEST 50634->9120 outbound NC_ITEM_CHARGED_WITHDRAW_REQ {"packetType":"small","length":6,"department":12,"command":"22","opCode":12322,"data":"ed206302","rawData":"062230ed206302","friendlyName":""}
	// INFO : 2021/04/25 17:59:04.279049 handlers.go:272: 2021-04-25 17:59:04.270294 +0200 CEST 9120->50634 inbound NC_ITEM_CHARGED_WITHDRAW_ACK {"packetType":"small","length":8,"department":12,"command":"23","opCode":12323,"data":"ed2063024110","rawData":"082330ed2063024110","friendlyName":""}
	// INFO : 2021/04/25 17:59:04.279049 handlers.go:272: 2021-04-25 17:59:04.270294 +0200 CEST 9120->50634 inbound NC_ITEM_CELLCHANGE_CMD {"packetType":"small","length":16,"department":12,"command":"1","opCode":12289,"data":"16241624e0ff0000000000000000","rawData":"10013016241624e0ff0000000000000000","friendlyName":""}
}

func TestChangeItemSlotPremiumInventoryToInventoryInventoryFull(t *testing.T) {
	t.Fail()
}

func TestChangeItemSlotPremiumInventoryToRewardInventoryShouldFail(t *testing.T) { t.Fail() }

func TestChangeItemSlotPremiumInventoryShouldFail(t *testing.T) { t.Fail() }

func TestChangeItemSlotMHInventoryToRewardInventoryShouldFail(t *testing.T) { t.Fail() }

func TestChangeItemSlotMHInventoryToPremiumInventoryShouldFail(t *testing.T) { t.Fail() }

func TestChangeItemSlotMHInventoryShouldFail(t *testing.T) { t.Fail() }

func TestChangeItemSlotMHInventoryEquip(t *testing.T) {
	t.Fail()
	// INFO : 2021/04/25 01:07:35.262832 handlers.go:267: 2021-04-25 01:07:35.250157 +0200 CEST 39878->9120 outbound NC_MINIHOUSE_ACTIV_REQ {"packetType":"small","length":3,"department":35,"command":"1","opCode":35841,"data":"0b","rawData":"03018c0b","friendlyName":""}
	// INFO : 2021/04/25 01:07:35.371279 handlers.go:267: 2021-04-25 01:07:35.365524 +0200 CEST 9120->39878 inbound NC_MINIHOUSE_ACTIV_ACK {"packetType":"small","length":4,"department":35,"command":"2","opCode":35842,"data":"0110","rawData":"04028c0110","friendlyName":""}
	// INFO : 2021/04/25 01:07:35.526755 handlers.go:267: 2021-04-25 01:07:35.524136 +0200 CEST 9120->39878 inbound NC_ITEM_RELOC_ACK {"packetType":"small","length":4,"department":12,"command":"C","opCode":12300,"data":"4102","rawData":"040c304102","friendlyName":""}
	// INFO : 2021/04/25 01:07:35.527276 handlers.go:267: 2021-04-25 01:07:35.524136 +0200 CEST 9120->39878 inbound NC_ITEM_CELLCHANGE_CMD {"packetType":"small","length":12,"department":12,"command":"1","opCode":12289,"data":"0b3000303d79ffecbb76","rawData":"0c01300b3000303d79ffecbb76","friendlyName":""}
	// INFO : 2021/04/25 01:07:35.527285 handlers.go:267: 2021-04-25 01:07:35.524136 +0200 CEST 9120->39878 inbound NC_ITEM_CELLCHANGE_CMD {"packetType":"small","length":12,"department":12,"command":"1","opCode":12289,"data":"00300b301879ffecbb76","rawData":"0c013000300b301879ffecbb76","friendlyName":""}
}

func TestSellItemSuccess(t *testing.T) {
	t.Fail()
}

func TestSellItemNonExistingItem(t *testing.T) {
	t.Fail()
}

func TestBuyItemSuccess(t *testing.T) {
	t.Fail()
}

func TestOneUseItemSuccess(t *testing.T) {
	t.Fail()
}

// Like mounts, quest items
func TestMultipleUseItemSuccess(t *testing.T) {
	t.Fail()
}
