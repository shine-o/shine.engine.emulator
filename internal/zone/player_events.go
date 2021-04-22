package zone

// 	eduPosition
//	eduState
//	eduStats
//	eduEquipItem
//	eduUnEquipItem

type eduPositionEvent struct {
	x, y int
	zm   *zoneMap
	err  chan error
}

type eduStateEvent struct{}

type eduStatsEvent struct{}

type eduEquipItemEvent struct {
	slot   int
	change itemSlotChange
	err    chan error
}

type eduUnEquipItemEvent struct {
	from, to     int
	change    itemSlotChange
	err       chan error
}
