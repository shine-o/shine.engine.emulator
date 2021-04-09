package zone

import (
	"github.com/shine-o/shine.engine.emulator/internal/pkg/crypto"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/data"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/errors"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/persistence"
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
	sync.RWMutex
}

type itemBox struct {
	box   uint8
	items map[int]*item
}

type item struct {
	pItem     *persistence.Item
	itemData  *itemData
	stats     itemStats
	amount    int
	stackable bool
	sync.RWMutex
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
	staticMDefense          itemStat
	staticPDefense          itemStat
	staticAim               itemStat
	staticEvasion           itemStat
	staticMaxHP             itemStat
	staticMaxSP             itemStat
	staticEvasionRate       itemStat
	staticAimRate           itemStat
	staticCriticalRate      itemStat
	staticPResistance       itemStat
	staticDResistance       itemStat
	staticCResistance       itemStat
	staticMResistance       itemStat
}

func (i *item) generateStats() {
	// first check if there are any random stats using (RandomOption / RandomOptionCount)
	// apply those first, after that check GradeItemOption for fixed stats
	// RNG for the number of stats that should be generated
	amount := amountStats(i.itemData)
	types := chosenStatTypes(amount, i.itemData)

	value := func(t data.RandomOptionType) int {
		ro := i.itemData.randomOption[t]
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
		case data.ROT_CON:
			if is.endurance.base > 0 {
				is.endurance.extra = value(t)
			} else {
				is.endurance.base = value(t)
			}
		case data.ROT_DEX:
			if is.dexterity.base > 0 {
				is.dexterity.extra = value(t)
			} else {
				is.dexterity.base = value(t)
			}
		case data.ROT_INT:
			if is.intelligence.base > 0 {
				is.intelligence.extra = value(t)
			} else {
				is.intelligence.base = value(t)
			}
		case data.ROT_MEN:
			if is.spirit.base > 0 {
				is.spirit.base = value(t)
			} else {
				is.spirit.base = value(t)
			}
		}
	}

	i.Lock()
	i.stats = is
	i.Unlock()
}

func (is *itemStats) staticStats(id *itemData) {

	is.staticAttackSpeed.base = int(id.itemInfo.AtkSpeed)
	is.staticMinPAttack.base = int(id.itemInfo.MinWC)
	is.staticMaxPAttack.base = int(id.itemInfo.MaxWC)
	is.staticMinMAttack.base = int(id.itemInfo.MinMA)
	is.staticMaxMAttack.base = int(id.itemInfo.MaxMA)
	is.staticPDefense.base = int(id.itemInfo.AC)
	is.staticMDefense.base = int(id.itemInfo.MR)
	is.staticAim.base = int(id.itemInfo.TH)
	is.staticEvasion.base = int(id.itemInfo.TB)
	is.staticPAttackRate.base = int(id.itemInfo.WCRate)
	is.staticMAttackRate.base = int(id.itemInfo.MARate)
	is.staticPDefenseRate.base = int(id.itemInfo.ACRate)
	is.staticMDefenseRate.base = int(id.itemInfo.MARate)
	is.staticPDefense.base = int(id.itemInfo.AC)
	is.staticMDefense.base = int(id.itemInfo.MR)
	is.staticCriticalRate.base = int(id.itemInfo.CriRate / 10)
	is.staticMinPACriticalRate.base = int(id.itemInfo.CriMinWc)
	is.staticMaxPACriticalRate.base = int(id.itemInfo.CriMaxWc)
	is.staticMinMACriticalRate.base = int(id.itemInfo.CriMinMa)
	is.staticMaxMACriticalRate.base = int(id.itemInfo.CriMaxMa)
	is.staticShieldDefenseRate.base = int(id.itemInfo.ShieldAC)

	if id.gradeItemOption != nil {
		if id.gradeItemOption.Strength > 0 {
			is.strength.base = int(id.gradeItemOption.Strength)
		}

		if id.gradeItemOption.Endurance > 0 {
			is.endurance.base = int(id.gradeItemOption.Endurance)
		}

		if id.gradeItemOption.Dexterity > 0 {
			is.dexterity.base = int(id.gradeItemOption.Dexterity)
		}

		if id.gradeItemOption.Intelligence > 0 {
			is.intelligence.base = int(id.gradeItemOption.Intelligence)
		}

		if id.gradeItemOption.Spirit > 0 {
			is.spirit.base = int(id.gradeItemOption.Spirit)
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

		if id.gradeItemOption.MaxHP > 0 {
			is.staticMaxHP.base = int(id.gradeItemOption.MaxHP)
		}

		if id.gradeItemOption.MaxSP > 0 {
			is.staticMaxSP.base = int(id.gradeItemOption.MaxSP)
		}
	}

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

type itemStat struct {
	base  int
	extra int
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

func makeItem(itemIndex string) (*item, error) {
	var i = &item{}

	itemData := getItemData(itemIndex)

	if itemData.itemInfo == nil {
		return i, errors.Err{
			Code: errors.ZoneItemMissingData,
			Details: errors.ErrDetails{
				"itemIndex": itemIndex,
				"type":      "ItemInfo",
			},
		}
	}

	if itemData.itemInfoServer == nil {
		return i, errors.Err{
			Code: errors.ZoneItemMissingData,
			Details: errors.ErrDetails{
				"itemIndex": itemIndex,
				"type":      "ItemInfoServer",
			},
		}
	}

	i.itemData = itemData
	i.pItem = &persistence.Item{}

	// first check if there are any random stats using (RandomOption / RandomOptionCount)
	// apply those first, after that check GradeItemOption for fixed stats
	if i.itemData.randomOption != nil && i.itemData.randomOptionCount != nil {
		i.generateStats()
	}

	if itemData.itemInfo.MaxLot > 1 {
		i.stackable = true
	}

	// will vary when created through the ItemDropTables
	// will vary when created through admin command with quantity parameter
	// will not vary if stackable is false
	i.amount = int(itemData.itemInfo.MaxLot)

	return i, nil
}

func (p *player) itemData() error {
	// for this character, load all items in each respective box
	// each item loaded should be validated so that, best way is to iterate all items and for each item launch a routine that validates it and returns the valid item through a channel
	// we also forward the error channel in case there is an error
	ivs := &playerInventories{
		equipped: itemBox{
			box: 8,
		},
		inventory: itemBox{
			box: persistence.BagInventory,
		},
		miniHouse: itemBox{
			box: 12,
		},
		premium: itemBox{
			box: 15,
		},
	}

	items, err := persistence.GetCharacterItems(int(p.char.ID), persistence.BagInventory)

	if err != nil {
		log.Error(err)
		return err
	}

	ivs.inventory.items = make(map[int]*item)

	for _, item := range items {

		ei, occupied := ivs.inventory.items[item.Slot]
		if occupied {
			log.Error(errors.Err{
				Code: errors.ZoneInventorySlotOccupied,
				Details: errors.ErrDetails{
					"itemID": ei.pItem.ID,
					"slot":   ei.pItem.Slot,
				},
			})
			continue
		}
		// load with goroutines and waitgroups
		ivs.inventory.items[item.Slot] = loadItem(item)
	}

	p.Lock()
	p.inventories = ivs
	p.Unlock()

	return nil
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
			dexterity:               itemStat{
				base:  pItem.Attributes.DexterityBase,
				extra: pItem.Attributes.DexterityExtra,
			},
			intelligence:            itemStat{
				base:  pItem.Attributes.IntelligenceBase,
				extra: pItem.Attributes.IntelligenceExtra,
			},
			endurance:               itemStat{
				base:  pItem.Attributes.EnduranceBase,
				extra: pItem.Attributes.EnduranceExtra,
			},
			spirit:                  itemStat{
				base:  pItem.Attributes.SpiritBase,
				extra: pItem.Attributes.SpiritExtra,
			},
			physicalAttack:          itemStat{
				base:  pItem.Attributes.PAttackBase,
				extra: pItem.Attributes.PAttackExtra,
			},
			magicalAttack:           itemStat{
				base:  pItem.Attributes.MAttackBase,
				extra: pItem.Attributes.MAttackExtra,
			},
			physicalDefense:         itemStat{
				base:  pItem.Attributes.PDefenseBase,
				extra: pItem.Attributes.PDefenseExtra,
			},
			magicalDefense:          itemStat{
				base:  pItem.Attributes.MDefenseBase,
				extra: pItem.Attributes.MDefenseExtra,
			},
			aim:                     itemStat{
				base:  pItem.Attributes.AimBase,
				extra: pItem.Attributes.AimExtra,
			},
			evasion:                 itemStat{
				base:  pItem.Attributes.EvasionBase,
				extra: pItem.Attributes.EvasionExtra,
			},
		},
		amount:    pItem.Amount,
		stackable: pItem.Stackable,
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
		if id.itemInfoServer.RandomOptionDropGroup != "" {
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
