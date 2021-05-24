package packet_sniffer

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/networking"
	"github.com/spf13/viper"
)

// PacketView is used to represent data to the frontend UI
type PacketView struct {
	// time of capture
	PacketID         string                 `json:"packetID"`
	ConnectionKey    string                 `json:"connectionKey"`
	TimeStamp        string                 `json:"timestamp"`
	IPEndpoints      string                 `json:"ipEndpoints"`
	PortEndpoints    string                 `json:"portEndpoints"`
	Direction        string                 `json:"direction"`
	PacketData       networking.ExportedPcb `json:"packetData"`
	NcRepresentation ncRepresentation       `json:"ncRepresentation"`
}

type webSockets struct {
	cons map[*websocket.Conn]bool
	mu   sync.Mutex
}

var upgrader = websocket.Upgrader{} // use default options

var ws *webSockets // grrr, find other way to send packets to

func startUI(ctx context.Context) {
	select {
	case <-ctx.Done():
		return
	default:
		ws = &webSockets{
			cons: make(map[*websocket.Conn]bool),
		}

		addr := fmt.Sprintf("localhost:%v", viper.GetString("websocket.port"))
		log.Infof("starting websocket server on %v", addr)
		http.HandleFunc("/packets", packets)

		log.Error(http.ListenAndServe(addr, nil))
	}
}

func (pv *PacketView) String() string {
	sd, err := json.Marshal(&pv)
	if err != nil {
		log.Error(err)
	}
	return string(sd)
}

func sendPacketToUI(pv PacketView) {
	ws.mu.Lock()
	// check if it can be done with goroutine
	if len(ws.cons) == 0 {
		return
	}
	for c, active := range ws.cons {
		if active {
			err := c.WriteMessage(websocket.TextMessage, []byte(pv.String()))
			if err != nil {
				log.Error("write:", err)
				continue
			}
		}
		time.Sleep(time.Millisecond * 100)
	}
	ws.mu.Unlock()
}

type completedFlow struct {
	FlowCompleted bool   `json:"flow_completed"`
	FlowID        string `json:"flow_id"`
}

func (cf *completedFlow) String() string {
	sd, err := json.Marshal(&cf)
	if err != nil {
		log.Error(err)
	}
	return string(sd)
}

func uiCompletedFlow(cf completedFlow) {
	ws.mu.Lock()
	// check if it can be done with goroutine
	for c, active := range ws.cons {
		if active {
			err := c.WriteMessage(websocket.TextMessage, []byte(cf.String()))
			if err != nil {
				log.Error("write:", err)
				break
			}
		}
	}
	ws.mu.Unlock()
}

func packets(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Info("upgrade:", err)
		return
	}

	ws.mu.Lock()
	ws.cons[c] = true
	ws.mu.Unlock()

	defer closeWebSocket(c)
	log.Info("websocket connection made")
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Info("read:", err)
			break
		}
		log.Info("recv: %s", message)
	}
}

func closeWebSocket(c *websocket.Conn) {
	err := c.Close()
	if err != nil {
		log.Error()
	}
	ws.mu.Lock()
	ws.cons[c] = false
	ws.mu.Unlock()
}
