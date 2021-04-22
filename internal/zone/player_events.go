package zone

import "github.com/shine-o/shine.engine.emulator/internal/pkg/persistence"

// 	eduPosition
//	eduState
//	eduStats
//	eduEquipItem
//	eduUnEquipItem

type eduPositionEvent struct {
	x, y int
	zm * zoneMap
	err chan error
}

type eduStateEvent struct {}

type eduStatsEvent struct {}

type eduEquipItemEvent struct {
	slot int
	change itemSlotChange
	err chan error
}

type eduUnEquipItemEvent struct {
	slot int
	inventory persistence.InventoryType
	change itemSlotChange
	err chan error
}