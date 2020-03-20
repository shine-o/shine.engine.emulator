package service

import (
	"context"
	"encoding/json"
	"github.com/gorilla/websocket"
	networking "github.com/shine-o/shine.engine.networking"
	"strconv"
	"sync"
	"time"
)

type PacketView struct {
	// unique uuid
	// when a stream finished its reassembly, notify the frontend that the flow with uuid is closed
	FlowID string `json:"flow_id"`
	FlowName string `json:"flow_name"`
	// numerary for the packets processed in this flow, informational only
	PacketID uint32 `json:"packet_id"`
	// IP that is not the server
	// to be used in the frontend as an abstraction for the many flows between the client and the server
	// a clientIP with 0 active flows is considered inactive
	ClientIP string `json:"client_ip"`
	ServerIP string `json:"server_ip"`
	// time of capture
	TimeStamp string `json:"timestamp"`
	PacketData string `json:"packet_data"`
}


func (ss * shineStream) handlePacket(ctx context.Context, wg * sync.WaitGroup, seen time.Time, pc networking.Command) {
	defer wg.Done()

	select {
	case <- ctx.Done():
		return
	default:
		var (
			clientIP string
			serverIP string
		)
		srcPort, _ := strconv.Atoi(ss.transport.Src().String())

		if srcPort >= 9000 && srcPort <= 9600 {
			// server ip
			serverIP = ss.net.Src().String()
			clientIP = ss.net.Dst().String()
		} else {
			clientIP = ss.net.Src().String()
			serverIP = ss.net.Dst().String()
		}
		//log.Infof("[%v] [%v] %v", ss.flowName, seen, pc.Base.String())
		pv := PacketView{
			FlowID:           ss.flowID,
			FlowName: ss.flowName,
			PacketID:         0,
			ClientIP:         clientIP,
			ServerIP:         serverIP,
			TimeStamp:        seen.String(),
			PacketData:       pc.Base.String(),
		}

		//sd, err := json.Marshal(&pv)
		//
		//if err != nil {
		//	log.Error(err)
		//	return
		//}
		go notifySockets(pv)
		//log.Info(string(sd))
	}
}

func (pv * PacketView) String() string  {
	sd, err := json.Marshal(&pv)
	if err != nil {
		log.Error(err)
	}
	return string(sd)
}

func notifySockets(pv PacketView)  {
	ws.mu.Lock()
	// check if it can be done with goroutine
	for c, active := range ws.cons {
		if active {
			err := c.WriteMessage(websocket.TextMessage, []byte(pv.String()))
			if err != nil {
				log.Error("write:", err)
				break
			}
		}
	}
	ws.mu.Unlock()
}