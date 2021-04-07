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
	pItem    *persistence.Item
	itemData *itemData
	stats     itemStats
	amount    int
	stackable bool
	sync.RWMutex
}

type itemStats struct {
	strength itemStat
	dexterity itemStat
	intelligence itemStat
	endurance itemStat
	spirit itemStat
	aim itemStat
	critical itemStat
	physicalAttack itemStat
	magicalAttack itemStat
	physicalDefense itemStat
	magicalDefense itemStat
	evasion itemStat
	hp itemStat
}

func (i * item) generateStats() {
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
	for _, t := range types {
		switch t {
		case data.ROT_STR:
			if is.strength.base > 0 {
				is.strength.extra = value(t)
			} else {
				is.strength.base = value(t)
			}
		case data.ROT_CON:
			is.endurance.base = value(t)
		case data.ROT_DEX:
			is.dexterity.base = value(t)
		case data.ROT_INT:
			is.intelligence.base = value(t)
		case data.ROT_MEN:
			is.spirit.base = value(t)
		}
	}

	i.Lock()
	i.stats = is
	i.Unlock()
}

// todo: extend this beyond RNG using the player's session data for deciding amount of stats, e.g: how much damage he did to the mob that dropped it, how many kills before, etc
func amountStats(id * itemData) int {
	var keys []int
	for k, _ := range id.randomOptionCount {
		keys = append(keys, int(k))
	}
	return crypto.RandomIntBetween(0, len(keys))
}

func chosenStatTypes(amount int, id * itemData) []data.RandomOptionType  {
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
	base    int
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
			Code:    errors.ZoneItemMissingData,
			Details: errors.ErrDetails{
				"itemIndex": itemIndex,
				"type": "ItemInfo",
			},
		}
	}

	if itemData.itemInfoServer == nil {
		return i, errors.Err{
			Code:    errors.ZoneItemMissingData,
			Details: errors.ErrDetails{
				"itemIndex": itemIndex,
				"type": "ItemInfoServer",
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
				Code:    errors.ZoneInventorySlotOccupied,
				Details: errors.ErrDetails{
					"itemID": ei.pItem.ID,
					"slot": ei.pItem.Slot,
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
		pItem:     pItem,
		itemData:  getItemData(pItem.ShnInxName),
		stats: itemStats{
			strength: itemStat{
				base:  pItem.Attributes.StrengthBase,
				extra: pItem.Attributes.StrengthExtra,
			},
		},
		amount:    pItem.Amount,
		stackable: pItem.Stackable,
	}
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

	if id.itemInfoServer.RandomOptionDropGroup != "" {
		wg.Add(2)
		go addRandomOptionRow(id.itemInfoServer.RandomOptionDropGroup, id, wg)
		go addRandomOptionCountRow(id.itemInfoServer.RandomOptionDropGroup, id, wg)
	}

	wg.Wait()

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