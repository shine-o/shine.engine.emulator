package zone

import (
	"github.com/shine-o/shine.engine.emulator/internal/pkg/data"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/errors"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/persistence"
	"sync"
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
	stats itemStats
	amount    int
	stackable bool
}

type itemStats struct {
	strength itemStat
}

type itemStat struct {
	base    int
	extra int
}

type itemData struct {
	itemInfo          *data.ItemInfo
	itemInfoServer    *data.ItemInfoServer
	gradeItemOption   *data.GradeItemOption
	randomOption      *data.RandomOption
	randomOptionCount *data.RandomOptionCount
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
	i.stats = itemStats{}

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
		if row.DropItemIndex == dropItemIndex {
			id.Lock()
			id.randomOption = &itemsData.RandomOption.ShineRow[i]
			id.Unlock()
			return
		}
	}
}

func addRandomOptionCountRow(dropItemIndex string, id *itemData, wg *sync.WaitGroup) {
	defer wg.Done()
	for i, row := range itemsData.RandomOptionCount.ShineRow {
		if row.DropItemIndex == dropItemIndex {
			id.Lock()
			id.randomOptionCount = &itemsData.RandomOptionCount.ShineRow[i]
			id.Unlock()
			return
		}
	}
}