package zone

import "github.com/prometheus/client_golang/prometheus"

// metrics specific to the zone service
type metrics struct {
	// every time a player enters or exit a map, update the gauge
	players prometheus.Gauge
	// every time a npc dies / respawns, update the gauge
	npcs prometheus.Gauge
}
