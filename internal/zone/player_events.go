package zone

// 	eduPosition
//	eduState
//	eduStats
//	eduEquipItem
//	eduUnEquipItem

type eduPositionEvent struct {
	player * player
	x, y int
	zm * zoneMap
	err chan error
}

type eduStateEvent struct {}

type eduStatsEvent struct {}

type eduEquipItemEvent struct {}

type eduUnEquipItemEvent struct {}