package zone

import (
	"github.com/shine-o/shine.engine.emulator/internal/pkg/data"
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

func (p *player) itemData() {
	// for this character, load all items in each respective box
	// each item loaded should be validated so that, best way is to iterate all items and for each item launch a routine that validates it and returns the valid item through a channel
	// we also forward the error channel in case there is an error
	i := &playerInventories{
		equipped: itemBox{
			box: 8,
		},
		inventory: itemBox{
			box: 9,
		},
		miniHouse: itemBox{
			box: 12,
		},
		premium: itemBox{
			box: 15,
		},
	}
	p.Lock()
	p.inventories = i
	p.Unlock()
}
