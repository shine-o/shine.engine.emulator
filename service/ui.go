package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
	"net/http"
	"sync"
)

// PacketView is used to represent data to the frontend UI
type PacketView struct {
	// unique uuid
	// when a stream finished its reassembly, notify the frontend that the flow with uuid is closed
	FlowID   string `json:"flow_id"`
	FlowName string `json:"flow_name"`
	// numerary for the packets processed in this flow, informational only
	PacketID uint64 `json:"packet_id"`
	// IP that is not the server
	// to be used in the frontend as an abstraction for the many flows between the client and the server
	// a clientIP with 0 active flows is considered inactive
	ClientIP string `json:"client_ip"`
	ServerIP string `json:"server_ip"`
	// time of capture
	TimeStamp        string           `json:"timestamp"`
	PacketData       string           `json:"packet_data"`
	NcRepresentation NcRepresentation `json:"nc_representation"`
}

type webSockets struct {
	//cons [] * websocket.Conn
	cons map[*websocket.Conn]bool
	mu   sync.Mutex
}

type ss map[*websocket.Conn]bool

var upgrader = websocket.Upgrader{} // use default options

var ws *webSockets

func startUI(ctx context.Context) {
	select {
	case <-ctx.Done():
		return
	default:

		ws = &webSockets{
			cons: make(map[*websocket.Conn]bool),
		}

		var addr = fmt.Sprintf("localhost:%v", viper.GetString("websocket.port"))
		log.Infof("starting ui on: http://%v", addr)
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
		//time.Sleep(time.Millisecond * 150)
	}
	ws.mu.Unlock()
}

type CompletedFlow struct {
	FlowCompleted bool   `json:"flow_completed"`
	ClientIP      string `json:"client_ip"`
	FlowID        string `json:"flow_id"`
}

func (cf *CompletedFlow) String() string {
	sd, err := json.Marshal(&cf)
	if err != nil {
		log.Error(err)
	}
	return string(sd)
}

func uiCompletedFlow(cf CompletedFlow) {
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
