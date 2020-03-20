package service

import (
	"github.com/shine-o/shine.engine.networking"
	"time"
)

type PacketView struct {
	// unique uuid
	// when a stream finished its reassembly, notify the frontend that the flow with uuid is closed
	flowID string
	flowFriendlyName string
	// numerary for the packets processed in this flow, informational only
	packetID uint32
	// IP that is not the server
	// to be used in the frontend as an abstraction for the many flows between the client and the server
	// a clientIP with 0 active flows is considered inactive
	clientIP string
	// time of capture
	timestamp time.Time
	pc networking.Command
}