package zone

import (
	"github.com/shine-o/shine.engine.emulator/internal/pkg/crypto"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/data"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/errors"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/persistence"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/structs"
	"math/rand"
	"sync"
	"time"
)

type playerInventories struct {
	equipped  itemBox
	inventory itemBox
	miniHouse itemBox
	reward    itemBox
	premium   itemBox
	*sync.RWMutex
}

func (pi *playerInventories) get(inventoryType persistence.InventoryType, slot int) *item {
	pi.RLock()
	defer pi.RUnlock()
	switch inventoryType {
	case persistence.EquippedInventory:
		item, ok := pi.equipped.items[slot]
		if ok {
			return item
		}
	case persistence.BagInventory:
		item, ok := pi.inventory.items[slot]
		if ok {
			return item
		}
	default:
		log.Error(errors.Err{
			Code:    errors.ZoneItemUnknownInventoryType,
			Details: errors.ErrDetails{
				"inventoryType": inventoryType,
				"slot": slot,
			},
		})
	}
	return nil
}

// move item from one inventory/slot to another inventory/slot
// the input from/to are values sent by the client
func (pi *playerInventories) moveItem(from, to uint16) (itemSlotChange, error) {
	var (
		change            = itemSlotChange{}
		fromInventoryType = persistence.InventoryType(from >> 10)
		toInventoryType   = persistence.InventoryType(to >> 10)
		fromInventorySlot = int(from & 1023)
		toInventorySlot   = int(to& 1023)
	)

	change.gameFrom = from
	change.gameTo = to

	switch fromInventoryType {
	case persistence.BagInventory:
		item := pi.get(fromInventoryType, fromInventorySlot)
		if item == nil {
			return change, errors.Err{
				Code: errors.ZoneItemSlotChangeNoItem,
				Details: errors.ErrDetails{
					"from": from,
					"to": to,
				},
			}
		}
		change.from.item = item
		change.from.slot = fromInventorySlot
		change.from.inventoryType = fromInventoryType
		break
	}

	switch toInventoryType {
	case persistence.BagInventory:
		item := pi.get(toInventoryType, toInventorySlot)
		change.to.item = item
		change.to.slot = toInventorySlot
		change.to.inventoryType = toInventoryType
		break
	}

	otherPItem, err := change.from.item.pItem.MoveTo(toInventoryType, toInventorySlot)

	if err != nil {
		return change, err
	}

	switch toInventoryType {
	case persistence.BagInventory:
		pi.Lock()
		//pi.inventory.items[change.from.slot] = nil
		delete(pi.inventory.items, change.from.slot)
		pi.inventory.items[change.to.slot] = change.from.item
		pi.Unlock()
		break
	}

	if change.to.item != nil {
		change.to.item.Lock()
		change.to.item.pItem = otherPItem
		change.to.item.Unlock()
		switch fromInventoryType {
		case persistence.BagInventory:
			pi.Lock()
			//delete(pi.inventory.items, change.to.slot)
			pi.inventory.items[change.from.slot] = change.to.item
			pi.Unlock()
			break
		}
	}

	return change, nil
}

type itemBox struct {
	box   int
	items map[int]*item
}

type item struct {
	pItem     *persistence.Item
	itemData  *itemData
	stats     itemStats
	amount    int
	stackable bool
	*sync.RWMutex
}

// static is a prefix for items that have stats defined in static files
// strength, dexterity, intelligence, endurance and spirit are an exception, as they can also be static and randomly generated
type itemStats struct {
	strength        itemStat
	dexterity       itemStat
	intelligence    itemStat
	endurance       itemStat
	spirit          itemStat
	physicalAttack  itemStat
	magicalAttack   itemStat
	physicalDefense itemStat
	magicalDefense  itemStat
	aim             itemStat
	evasion         itemStat
	maxHP           itemStat
	maxSP           itemStat
	// unsure what these stats are
	critical        itemStat
	criticalEvasion itemStat

	staticAttackSpeed       itemStat
	staticMinPAttack        itemStat
	staticMaxPAttack        itemStat
	staticMinMAttack        itemStat
	staticMaxMAttack        itemStat
	staticMinPACriticalRate itemStat
	staticMaxPACriticalRate itemStat
	staticMinMACriticalRate itemStat
	staticMaxMACriticalRate itemStat
	staticMAttackRate       itemStat
	staticPAttackRate       itemStat
	staticMDefenseRate      itemStat
	staticPDefenseRate      itemStat
	staticShieldDefenseRate itemStat

	staticEvasionRate  itemStat
	staticAimRate      itemStat
	staticCriticalRate itemStat
	staticPResistance  itemStat
	staticDResistance  itemStat
	staticCResistance  itemStat
	staticMResistance  itemStat
}

type itemStat struct {
	isStatic   bool
	isUpgraded bool
	base       int
	extra      int
}

type itemData struct {
	itemInfo          *data.ItemInfo
	itemInfoServer    *data.ItemInfoServer
	gradeItemOption   *data.GradeItemOption
	randomOption      map[data.RandomOptionType]*data.RandomOption
	randomOptionCount map[uint16]*data.RandomOptionCount
	// itemUseEffect
	// ... etc
	sync.Mutex
}

type itemCreationDetails struct {
	randomTypes []data.RandomOptionType
	randomCount int
}

type itemSlotChange struct {
	gameFrom uint16
	gameTo   uint16
	from     itemSlot
	to       itemSlot
}

type itemSlot struct {
	slot int
	inventoryType persistence.InventoryType
	item *item
}

func (i *item) generateStats() (int, []data.RandomOptionType) {
	// first check if there are any random stats using (RandomOption / RandomOptionCount)
	// apply those first, after that check GradeItemOption for fixed stats
	// RNG for the number of stats that should be generated
	amount := amountStats(i.itemData)
	types := chosenStatTypes(amount, i.itemData)

	value := func(t data.RandomOptionType) int {
		ro := i.itemData.randomOption[t]
		if ro.Min == ro.Max {
			return int(ro.Min)
		}
		return int(crypto.RandomUint32Between(ro.Min, ro.Max))
	}

	is := itemStats{}

	is.staticStats(i.itemData)

	for _, t := range types {
		switch t {
		case data.ROT_STR:
			if is.strength.base > 0 {
				is.strength.extra = value(t)
			} else {
				is.strength.base = value(t)
			}
			break
		case data.ROT_CON:
			if is.endurance.base > 0 {
				is.endurance.extra = value(t)
			} else {
				is.endurance.base = value(t)
			}
			break
		case data.ROT_DEX:
			if is.dexterity.base > 0 {
				is.dexterity.extra = value(t)
			} else {
				is.dexterity.base = value(t)
			}
			break
		case data.ROT_INT:
			if is.intelligence.base > 0 {
				is.intelligence.extra = value(t)
			} else {
				is.intelligence.base = value(t)
			}
			break
		case data.ROT_MEN:
			if is.spirit.base > 0 {
				is.spirit.extra = value(t)
			} else {
				is.spirit.base = value(t)
			}
			break
		case data.ROT_TH:
			if is.aim.base > 0 {
				is.aim.extra = value(t)
			} else {
				is.aim.base = value(t)
			}
			break
		case data.ROT_CRI:
			if is.critical.base > 0 {
				is.critical.extra = value(t)
			} else {
				is.critical.base = value(t)
			}
			break
		case data.ROT_WC:
			if is.physicalAttack.base > 0 {
				is.physicalAttack.extra = value(t)
			} else {
				is.physicalAttack.base = value(t)
			}
			break
		case data.ROT_AC:
			if is.physicalDefense.base > 0 {
				is.physicalDefense.extra = value(t)
			} else {
				is.physicalDefense.base = value(t)
			}
			break
		case data.ROT_MA:
			if is.magicalAttack.base > 0 {
				is.magicalAttack.extra = value(t)
			} else {
				is.magicalAttack.base = value(t)
			}
			break
		case data.ROT_MR:
			if is.magicalDefense.base > 0 {
				is.magicalDefense.extra = value(t)
			} else {
				is.magicalDefense.base = value(t)
			}
			break
		case data.ROT_TB:
			if is.evasion.base > 0 {
				is.evasion.extra = value(t)
			} else {
				is.evasion.base = value(t)
			}
			break
		case data.ROT_CRITICAL_TB:
			if is.criticalEvasion.base > 0 {
				is.criticalEvasion.extra = value(t)
			} else {
				is.criticalEvasion.base = value(t)
			}
			break
		case data.ROT_DEMANDLVDOWN:
			if is.maxHP.base > 0 {
				is.maxHP.extra = value(t)
			} else {
				is.maxHP.base = value(t)
			}
			break
		case data.ROT_MAXHP:
			if is.maxHP.base > 0 {
				is.maxHP.extra = value(t)
			} else {
				is.maxHP.base = value(t)
			}
			break
		}
	}

	i.Lock()
	i.stats = is
	i.Unlock()

	return amount, types
}

func (is *itemStats) staticStats(id *itemData) {

	is.staticAttackSpeed.base = int(id.itemInfo.AtkSpeed)
	is.staticMinPAttack.base = int(id.itemInfo.MinWC)
	is.staticMaxPAttack.base = int(id.itemInfo.MaxWC)
	is.staticMinMAttack.base = int(id.itemInfo.MinMA)
	is.staticMaxMAttack.base = int(id.itemInfo.MaxMA)
	is.staticPAttackRate.base = int(id.itemInfo.WCRate)
	is.staticMAttackRate.base = int(id.itemInfo.MARate)
	is.staticPDefenseRate.base = int(id.itemInfo.ACRate)
	is.staticMDefenseRate.base = int(id.itemInfo.MARate)
	is.staticCriticalRate.base = int(id.itemInfo.CriRate / 10)
	is.staticMinPACriticalRate.base = int(id.itemInfo.CriMinWc)
	is.staticMaxPACriticalRate.base = int(id.itemInfo.CriMaxWc)
	is.staticMinMACriticalRate.base = int(id.itemInfo.CriMinMa)
	is.staticMaxMACriticalRate.base = int(id.itemInfo.CriMaxMa)
	is.staticShieldDefenseRate.base = int(id.itemInfo.ShieldAC)

	if id.itemInfo.AC > 0 {
		is.physicalDefense.base = int(id.itemInfo.AC)
		is.physicalDefense.isStatic = true
	}

	if id.itemInfo.MR > 0 {
		is.magicalDefense.base = int(id.itemInfo.MR)
		is.magicalDefense.isStatic = true
	}

	if id.itemInfo.TH > 0 {
		is.aim.base = int(id.itemInfo.TH)
		is.aim.isStatic = true
	}

	if id.itemInfo.TB > 0 {
		is.evasion.base = int(id.itemInfo.TB)
		is.evasion.isStatic = true
	}

	if id.gradeItemOption != nil {
		if id.gradeItemOption.Strength > 0 {
			is.strength.base = int(id.gradeItemOption.Strength)
			is.strength.isStatic = true
		}

		if id.gradeItemOption.Endurance > 0 {
			is.endurance.base = int(id.gradeItemOption.Endurance)
			is.endurance.isStatic = true
		}

		if id.gradeItemOption.Dexterity > 0 {
			is.dexterity.base = int(id.gradeItemOption.Dexterity)
			is.dexterity.isStatic = true
		}

		if id.gradeItemOption.Intelligence > 0 {
			is.intelligence.base = int(id.gradeItemOption.Intelligence)
			is.intelligence.isStatic = true
		}

		if id.gradeItemOption.Spirit > 0 {
			is.spirit.base = int(id.gradeItemOption.Spirit)
			is.spirit.isStatic = true
		}

		if id.gradeItemOption.MaxHP > 0 {
			is.maxHP.base = int(id.gradeItemOption.MaxHP)
			is.spirit.isStatic = true
		}

		if id.gradeItemOption.MaxSP > 0 {
			is.maxSP.base = int(id.gradeItemOption.MaxSP)
			is.spirit.isStatic = true
		}

		if id.gradeItemOption.PoisonResistance > 0 {
			is.staticPResistance.base = int(id.gradeItemOption.PoisonResistance)
		}

		if id.gradeItemOption.DiseaseResistance > 0 {
			is.staticDResistance.base = int(id.gradeItemOption.DiseaseResistance)
		}

		if id.gradeItemOption.CurseResistance > 0 {
			is.staticCResistance.base = int(id.gradeItemOption.CurseResistance)
		}

		if id.gradeItemOption.MobilityResistance > 0 {
			is.staticMResistance.base = int(id.gradeItemOption.MobilityResistance)
		}

		if id.gradeItemOption.AimRate > 0 {
			is.staticAimRate.base = int(id.gradeItemOption.AimRate - 1000)
		}

		if id.gradeItemOption.EvasionRate > 0 {
			is.staticEvasionRate.base = int(id.gradeItemOption.EvasionRate - 1000)
		}
	}

}

func itemAttributesBytes(i *item) ([]byte, error) {
	i.RLock()
	defer i.RUnlock()
	var itemAttr []byte

	switch i.itemData.itemInfo.Class {
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
		attr := structs.ShineItemAttrByteLot(i.amount)
		bytes, err := structs.Pack(&attr)
		if err != nil {
			return itemAttr, err
		}
		itemAttr = bytes
		break
	case data.ItemClassWordLot:
	case data.ItemKqStep:
	case data.ItemActiveSkill:
		attr := structs.ShineItemAttrWordLot(i.amount)
		bytes, err := structs.Pack(&attr)
		if err != nil {
			return itemAttr, err
		}
		itemAttr = bytes
		break
	case data.ItemClassDwrdLot:
		attr := structs.ShineItemAttrDwrdLot(i.amount)
		bytes, err := structs.Pack(&attr)
		if err != nil {
			return itemAttr, err
		}
		itemAttr = bytes
		break
	case data.ItemClassQuestItem:
		attr := structs.ShineItemAttrQuestItem(i.amount)
		bytes, err := structs.Pack(&attr)
		if err != nil {
			return itemAttr, err
		}
		itemAttr = bytes
		break
	case data.ItemClassAmulet:
		attr := structs.ShineItemAttrAmulet{}
		attr.Option = itemOptionStorage(i.stats)
		bytes, err := structs.Pack(&attr)
		if err != nil {
			return itemAttr, err
		}
		itemAttr = bytes
		break
	case data.ItemClassWeapon:
		attr := structs.ShineItemAttrWeapon{
			// TODO: implement Licence feature
			Licences: [3]structs.ShineItemWeaponLicence{
				{
					MobID: 65535,
					BF2:   0,
				},
				{
					MobID: 65535,
					BF2:   0,
				},
				{
					MobID: 65535,
					BF2:   0,
				},
			},
			// TODO: implement Gems feature
			GemSockets: [3]structs.ShineItemWeaponGemSocket{
				{
					GemID:     65535,
					RestCount: 25,
				},
				{
					GemID:     65535,
					RestCount: 25,
				},
				{
					GemID:     65535,
					RestCount: 25,
				},
			},
			MaxSocketCount: 2,
		}
		attr.Option = itemOptionStorage(i.stats)
		bytes, err := structs.Pack(&attr)
		if err != nil {
			return itemAttr, err
		}
		itemAttr = bytes
		break
	case data.ItemClassArmor:
		attr := structs.ShineItemAttrArmor{}
		attr.Option = itemOptionStorage(i.stats)
		bytes, err := structs.Pack(&attr)
		if err != nil {
			return itemAttr, err
		}
		itemAttr = bytes
		break
	case data.ItemClassShield:
		attr := structs.ShineItemAttrShield{}
		bytes, err := structs.Pack(&attr)
		if err != nil {
			return itemAttr, err
		}
		itemAttr = bytes
		break
	case data.ItemClassBoot:
		attr := structs.ShineItemAttrBoot{}
		bytes, err := structs.Pack(&attr)
		if err != nil {
			return itemAttr, err
		}
		itemAttr = bytes
		break
	case data.ItemClassFurniture:
		attr := structs.ShineItemAttrFurniture{}
		bytes, err := structs.Pack(&attr)
		if err != nil {
			return itemAttr, err
		}
		itemAttr = bytes
		break
	case data.ItemClassDecoration:
		attr := structs.ShineItemAttrDecoration{}
		bytes, err := structs.Pack(&attr)
		if err != nil {
			return itemAttr, err
		}
		itemAttr = bytes
		break
	case data.ItemClassBindItem:
		attr := structs.ShineItemAttrBindItem{}
		bytes, err := structs.Pack(&attr)
		if err != nil {
			return itemAttr, err
		}
		itemAttr = bytes
		break
	case data.ItemClassItemChest:
		attr := structs.ShineItemAttrItemChest{}
		bytes, err := structs.Pack(&attr)
		if err != nil {
			return itemAttr, err
		}
		itemAttr = bytes
		break
	case data.ItemHouseSkin:
		attr := structs.ShineItemAttrMiniHouseSkin{}
		bytes, err := structs.Pack(&attr)
		if err != nil {
			return itemAttr, err
		}
		itemAttr = bytes
		break
	case data.ItemRiding:
		attr := structs.ShineItemAttrRiding{}
		bytes, err := structs.Pack(&attr)
		if err != nil {
			return itemAttr, err
		}
		itemAttr = bytes
		break
	case data.ItemAmount:
		attr := structs.ShineItemAttrAmount{}
		bytes, err := structs.Pack(&attr)
		if err != nil {
			return itemAttr, err
		}
		itemAttr = bytes
		break
	case data.ItemCosWeapon:
		attr := structs.ShineItemAttrCostumeWeapon{}
		bytes, err := structs.Pack(&attr)
		if err != nil {
			return itemAttr, err
		}
		itemAttr = bytes
		break
	case data.ItemActionItem:
		attr := structs.ShineItemAttrActionItem{}
		bytes, err := structs.Pack(&attr)
		if err != nil {
			return itemAttr, err
		}
		itemAttr = bytes
		break
	case data.ItemCapsule:
		attr := structs.ShineItemAttrCapsule{}
		bytes, err := structs.Pack(&attr)
		if err != nil {
			return itemAttr, err
		}
		itemAttr = bytes
		break
	case data.ItemClosedCard:
		attr := structs.ShineItemAttrMobCardCollectClosed{}
		bytes, err := structs.Pack(&attr)
		if err != nil {
			return itemAttr, err
		}
		itemAttr = bytes
		break
	case data.ItemOpenCard:
		attr := structs.ShineItemAttrMobCardCollect{}
		bytes, err := structs.Pack(&attr)
		if err != nil {
			return itemAttr, err
		}
		itemAttr = bytes
		break
	case data.ItemMoney:
		//
		break
	case data.ItemPup:
		attr := structs.ShineItemAttrPet{}
		bytes, err := structs.Pack(&attr)
		if err != nil {
			return itemAttr, err
		}
		itemAttr = bytes
		break
	case data.ItemCosShield:
		attr := structs.ShineItemAttrCostumeShield{}
		bytes, err := structs.Pack(&attr)
		if err != nil {
			return itemAttr, err
		}
		itemAttr = bytes
		break
	case data.ItemBracelet:
		attr := structs.ShineItemAttrBracelet{}
		bytes, err := structs.Pack(&attr)
		if err != nil {
			return itemAttr, err
		}
		itemAttr = bytes
		break
	default:
		return itemAttr, errors.Err{
			Code: errors.ZoneUnknownItemClass,
			Details: errors.ErrDetails{
				"itemID": i.pItem.ID,
			},
		}
	}

	return itemAttr, nil
}

func protoItemPacketInformation(i *item) (*structs.ProtoItemPacketInformation, error) {
	var (
		nc       *structs.ProtoItemPacketInformation
		itemAttr []byte
	)

	nc = &structs.ProtoItemPacketInformation{}
	nc.Location.Inventory = uint16(i.pItem.InventoryType<<10 | i.pItem.Slot&1023)
	nc.ItemID = i.itemData.itemInfo.ID
	nc.DataSize = 4

	itemAttr, err := itemAttributesBytes(i)
	if err != nil {
		return nc, err
	}

	nc.ItemAttr = itemAttr
	nc.DataSize += byte(len(itemAttr))

	return nc, nil
}

func itemOptionStorage(stats itemStats) structs.ItemOptionStorage {
	var (
		storage  = structs.ItemOptionStorage{}
		elements []structs.ItemOptionStorageElement
	)

	if stats.strength.base > 0 || stats.strength.extra > 0 {
		iose := structs.ItemOptionStorageElement{
			ItemOptionType: byte(data.ROT_STR),
		}

		if stats.strength.isStatic {
			iose.ItemOptionValue = uint16(stats.strength.extra)
		} else {
			iose.ItemOptionValue = uint16(stats.strength.base)
		}

		elements = append(elements, iose)
	}

	if stats.endurance.base > 0 || stats.endurance.extra > 0 {
		iose := structs.ItemOptionStorageElement{
			ItemOptionType: byte(data.ROT_CON),
		}

		if stats.endurance.isStatic {
			iose.ItemOptionValue = uint16(stats.endurance.extra)
		} else {
			iose.ItemOptionValue = uint16(stats.endurance.base)
		}

		elements = append(elements, iose)
	}

	if stats.dexterity.base > 0 || stats.dexterity.extra > 0 {
		iose := structs.ItemOptionStorageElement{
			ItemOptionType: byte(data.ROT_DEX),
		}

		if stats.dexterity.isStatic {
			iose.ItemOptionValue = uint16(stats.dexterity.extra)
		} else {
			iose.ItemOptionValue = uint16(stats.dexterity.base)
		}

		elements = append(elements, iose)
	}

	if stats.intelligence.base > 0 || stats.intelligence.extra > 0 {
		iose := structs.ItemOptionStorageElement{
			ItemOptionType: byte(data.ROT_INT),
		}

		if stats.intelligence.isStatic {
			iose.ItemOptionValue = uint16(stats.intelligence.extra)
		} else {
			iose.ItemOptionValue = uint16(stats.intelligence.base)
		}

		elements = append(elements, iose)
	}

	if stats.spirit.base > 0 || stats.spirit.extra > 0 {
		iose := structs.ItemOptionStorageElement{
			ItemOptionType: byte(data.ROT_MEN),
		}

		if stats.spirit.isStatic {
			iose.ItemOptionValue = uint16(stats.spirit.extra)
		} else {
			iose.ItemOptionValue = uint16(stats.spirit.base)
		}

		elements = append(elements, iose)
	}

	if stats.aim.base > 0 || stats.aim.extra > 0 {
		iose := structs.ItemOptionStorageElement{
			ItemOptionType: byte(data.ROT_TH),
		}

		if stats.aim.isStatic {
			iose.ItemOptionValue = uint16(stats.aim.extra)
		} else {
			iose.ItemOptionValue = uint16(stats.aim.base)
		}

		elements = append(elements, iose)
	}

	if stats.critical.base > 0 || stats.critical.extra > 0 {
		iose := structs.ItemOptionStorageElement{
			ItemOptionType: byte(data.ROT_CRI),
		}

		if stats.critical.isStatic {
			iose.ItemOptionValue = uint16(stats.critical.extra)
		} else {
			iose.ItemOptionValue = uint16(stats.critical.base)
		}

		elements = append(elements, iose)
	}

	if stats.physicalAttack.base > 0 || stats.physicalAttack.extra > 0 {
		iose := structs.ItemOptionStorageElement{
			ItemOptionType: byte(data.ROT_WC),
		}

		if stats.physicalAttack.isStatic {
			iose.ItemOptionValue = uint16(stats.physicalAttack.extra)
		} else {
			iose.ItemOptionValue = uint16(stats.physicalAttack.base)
		}

		elements = append(elements, iose)
	}

	if stats.physicalDefense.base > 0 || stats.physicalDefense.extra > 0 {
		iose := structs.ItemOptionStorageElement{
			ItemOptionType: byte(data.ROT_AC),
		}

		if stats.physicalDefense.isStatic {
			iose.ItemOptionValue = uint16(stats.physicalDefense.extra)
		} else {
			iose.ItemOptionValue = uint16(stats.physicalDefense.base)
		}

		elements = append(elements, iose)
	}

	if stats.magicalAttack.base > 0 || stats.magicalAttack.extra > 0 {
		iose := structs.ItemOptionStorageElement{
			ItemOptionType: byte(data.ROT_MA),
		}

		if stats.magicalAttack.isStatic {
			iose.ItemOptionValue = uint16(stats.magicalAttack.extra)
		} else {
			iose.ItemOptionValue = uint16(stats.magicalAttack.base)
		}

		elements = append(elements, iose)
	}

	if stats.magicalDefense.base > 0 || stats.magicalDefense.extra > 0 {
		iose := structs.ItemOptionStorageElement{
			ItemOptionType: byte(data.ROT_MR),
		}

		if stats.magicalDefense.isStatic {
			iose.ItemOptionValue = uint16(stats.magicalDefense.extra)
		} else {
			iose.ItemOptionValue = uint16(stats.magicalDefense.base)
		}

		elements = append(elements, iose)
	}

	if stats.evasion.base > 0 || stats.evasion.extra > 0 {
		iose := structs.ItemOptionStorageElement{
			ItemOptionType: byte(data.ROT_TB),
		}

		if stats.evasion.isStatic {
			iose.ItemOptionValue = uint16(stats.evasion.extra)
		} else {
			iose.ItemOptionValue = uint16(stats.evasion.base)
		}

		elements = append(elements, iose)
	}

	if stats.criticalEvasion.base > 0 || stats.criticalEvasion.extra > 0 {
		iose := structs.ItemOptionStorageElement{
			ItemOptionType: byte(data.ROT_TB),
		}

		if stats.criticalEvasion.isStatic {
			iose.ItemOptionValue = uint16(stats.criticalEvasion.extra)
		} else {
			iose.ItemOptionValue = uint16(stats.criticalEvasion.base)
		}

		elements = append(elements, iose)
	}

	if stats.maxHP.base > 0 || stats.maxHP.extra > 0 {
		iose := structs.ItemOptionStorageElement{
			ItemOptionType: byte(data.ROT_MAXHP),
		}

		if stats.maxHP.isStatic {
			iose.ItemOptionValue = uint16(stats.maxHP.extra)
		} else {
			iose.ItemOptionValue = uint16(stats.maxHP.base)
		}

		elements = append(elements, iose)
	}

	storage.AmountBit = byte(len(elements)<<1 | 1)
	storage.Elements = elements
	return storage
}

// todo: extend this beyond RNG using the player's session data for deciding amount of stats, e.g: how much damage he did to the mob that dropped it, how many kills before, etc
func amountStats(id *itemData) int {
	var keys []int
	for k, _ := range id.randomOptionCount {
		keys = append(keys, int(k))
	}
	if len(keys) > 0 {
		return keys[crypto.RandomIntBetween(0, len(keys))]
	}
	return 0
}

func chosenStatTypes(amount int, id *itemData) []data.RandomOptionType {
	var types []data.RandomOptionType

	for rot, _ := range id.randomOption {
		types = append(types, rot)
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(types), func(i, j int) {
		types[i], types[j] = types[j], types[i]
	})

	return types[:amount]
}

func makeItem(itemIndex string) (*item, itemCreationDetails, error) {
	var (
		i = &item{
			RWMutex: &sync.RWMutex{},
		}
		icd = itemCreationDetails{}
	)

	itemData := getItemData(itemIndex)

	if itemData.itemInfo == nil {
		return i, icd, errors.Err{
			Code: errors.ZoneItemMissingData,
			Details: errors.ErrDetails{
				"itemIndex": itemIndex,
				"type":      "ItemInfo",
			},
		}
	}

	if itemData.itemInfoServer == nil {
		return i, icd, errors.Err{
			Code: errors.ZoneItemMissingData,
			Details: errors.ErrDetails{
				"itemIndex": itemIndex,
				"type":      "ItemInfoServer",
			},
		}
	}

	i.itemData = itemData
	i.pItem = &persistence.Item{
		InventoryType: int(persistence.BagInventory),
	}

	// first check if there are any random stats using (RandomOption / RandomOptionCount)
	// apply those first, after that check GradeItemOption for fixed stats
	if i.itemData.randomOption != nil && i.itemData.randomOptionCount != nil {
		count, types := i.generateStats()
		icd.randomCount = count
		icd.randomTypes = types
	}

	if itemData.itemInfo.MaxLot > 1 {
		i.stackable = true
	}

	// will vary when created through the ItemDropTables
	// will vary when created through admin command with quantity parameter
	// will not vary if stackable is false
	i.amount = int(itemData.itemInfo.MaxLot)

	return i, icd, nil
}

func loadItem(pItem *persistence.Item) *item {
	i := &item{
		pItem:    pItem,
		itemData: getItemData(pItem.ShnInxName),
		stats: itemStats{
			strength: itemStat{
				base:  pItem.Attributes.StrengthBase,
				extra: pItem.Attributes.StrengthExtra,
			},
			dexterity: itemStat{
				base:  pItem.Attributes.DexterityBase,
				extra: pItem.Attributes.DexterityExtra,
			},
			intelligence: itemStat{
				base:  pItem.Attributes.IntelligenceBase,
				extra: pItem.Attributes.IntelligenceExtra,
			},
			endurance: itemStat{
				base:  pItem.Attributes.EnduranceBase,
				extra: pItem.Attributes.EnduranceExtra,
			},
			spirit: itemStat{
				base:  pItem.Attributes.SpiritBase,
				extra: pItem.Attributes.SpiritExtra,
			},
			physicalAttack: itemStat{
				base:  pItem.Attributes.PAttackBase,
				extra: pItem.Attributes.PAttackExtra,
			},
			magicalAttack: itemStat{
				base:  pItem.Attributes.MAttackBase,
				extra: pItem.Attributes.MAttackExtra,
			},
			physicalDefense: itemStat{
				base:  pItem.Attributes.PDefenseBase,
				extra: pItem.Attributes.PDefenseExtra,
			},
			magicalDefense: itemStat{
				base:  pItem.Attributes.MDefenseBase,
				extra: pItem.Attributes.MDefenseExtra,
			},
			aim: itemStat{
				base:  pItem.Attributes.AimBase,
				extra: pItem.Attributes.AimExtra,
			},
			evasion: itemStat{
				base:  pItem.Attributes.EvasionBase,
				extra: pItem.Attributes.EvasionExtra,
			},
		},
		amount:    pItem.Amount,
		stackable: pItem.Stackable,
		RWMutex:   &sync.RWMutex{},
	}

	i.stats.staticStats(i.itemData)

	return i
}

func getItemData(itemIndex string) *itemData {
	var (
		id = &itemData{}
		wg = &sync.WaitGroup{}
	)

	wg.Add(3)

	go addItemInfoRow(itemIndex, id, wg)

	go addItemInfoServerRow(itemIndex, id, wg)

	go addGradeItemOptionRow(itemIndex, id, wg)

	wg.Wait()

	if id.itemInfo != nil && id.itemInfoServer != nil {
		if id.itemInfoServer.RandomOptionDropGroup != "" && id.itemInfoServer.RandomOptionDropGroup != "-" {
			wg.Add(2)
			go addRandomOptionRow(id.itemInfoServer.RandomOptionDropGroup, id, wg)
			go addRandomOptionCountRow(id.itemInfoServer.RandomOptionDropGroup, id, wg)
		}

		wg.Wait()
	}

	return id
}

func addItemInfoRow(itemIndex string, id *itemData, wg *sync.WaitGroup) {
	defer wg.Done()
	for i, row := range itemsData.ItemInfo.ShineRow {
		if row.InxName == itemIndex {
			id.Lock()
			id.itemInfo = &itemsData.ItemInfo.ShineRow[i]
			id.Unlock()
			return
		}
	}
}

func addItemInfoServerRow(itemIndex string, id *itemData, wg *sync.WaitGroup) {
	defer wg.Done()
	for i, row := range itemsData.ItemInfoServer.ShineRow {
		if row.InxName == itemIndex {
			id.Lock()
			id.itemInfoServer = &itemsData.ItemInfoServer.ShineRow[i]
			id.Unlock()
			return
		}
	}
}

func addGradeItemOptionRow(itemIndex string, id *itemData, wg *sync.WaitGroup) {
	defer wg.Done()
	for i, row := range itemsData.GradeItemOptions.ShineRow {
		if row.ItemIndex == itemIndex {
			id.Lock()
			id.gradeItemOption = &itemsData.GradeItemOptions.ShineRow[i]
			id.Unlock()
			return
		}
	}
}

func addRandomOptionRow(dropItemIndex string, id *itemData, wg *sync.WaitGroup) {
	defer wg.Done()
	for i, row := range itemsData.RandomOption.ShineRow {
		if id.randomOption == nil {
			id.randomOption = make(map[data.RandomOptionType]*data.RandomOption)
		}
		if row.DropItemIndex == dropItemIndex {
			id.Lock()
			id.randomOption[row.RandomOptionType] = &itemsData.RandomOption.ShineRow[i]
			id.Unlock()
			//return
		}
	}
}

func addRandomOptionCountRow(dropItemIndex string, id *itemData, wg *sync.WaitGroup) {
	defer wg.Done()
	for i, row := range itemsData.RandomOptionCount.ShineRow {
		if id.randomOptionCount == nil {
			id.randomOptionCount = make(map[uint16]*data.RandomOptionCount)
		}
		if row.DropItemIndex == dropItemIndex {
			id.Lock()
			id.randomOptionCount[row.LimitCount] = &itemsData.RandomOptionCount.ShineRow[i]
			id.Unlock()
			//return
		}
	}
}

// NC_ITEM_CELLCHANGE_CMD
// info about a slot change between items
func ncItemCellChangeCmd(change itemSlotChange) (*structs.NcItemCellChangeCmd, *structs.NcItemCellChangeCmd, error) {
	var (
		nc1 = &structs.NcItemCellChangeCmd{}
		nc2 = &structs.NcItemCellChangeCmd{}
	)

	nc1.Exchange.Inventory = change.gameFrom
	nc1.Location.Inventory = change.gameTo
	nc1.Item.ItemID = change.from.item.itemData.itemInfo.ID
	itemAttr, err := itemAttributesBytes(change.from.item)

	if err != nil {
		return nc1, nc2, err
	}

	nc1.Item.ItemAttr = itemAttr

	nc2.Exchange.Inventory = change.gameTo
	nc2.Location.Inventory = change.gameFrom
	if change.to.item != nil {
		nc2.Item.ItemID = change.to.item.itemData.itemInfo.ID
		itemAttr1, err := itemAttributesBytes(change.to.item)

		if err != nil {
			return nc1, nc2, err
		}

		nc2.Item.ItemAttr = itemAttr1
	} else {
		nc2.Item.ItemID = 65535
	}

	return nc1, nc2, nil
}

// NC_ITEM_EQUIPCHANGE_CMD
// data about the recently equipped item
func ncItemEquipChangeCmd(change itemSlotChange) (structs.NcItemEquipChangeCmd, error) {
	nc := structs.NcItemEquipChangeCmd{
		From:      structs.ItemInventory{
			Inventory: change.gameFrom,
		},
		EquipSlot: byte(change.to.slot),
		ItemData:  structs.ShineItemVar{
			ItemID:   change.from.item.itemData.itemInfo.ID,
		},
	}

	itemAttr, err := itemAttributesBytes(change.from.item)

	if err != nil {
		return nc, err
	}

	nc.ItemData.ItemAttr = itemAttr

	return nc, nil
}
