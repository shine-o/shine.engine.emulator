package service

import "github.com/shine-o/shine.engine.core/game/maps"

type zoneMap struct {
	data maps.MapData
	sectors sectorGrid
}

// load maps
