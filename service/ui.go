package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
	"html/template"
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
	PacketID uint32 `json:"packet_id"`
	// IP that is not the server
	// to be used in the frontend as an abstraction for the many flows between the client and the server
	// a clientIP with 0 active flows is considered inactive
	ClientIP string `json:"client_ip"`
	ServerIP string `json:"server_ip"`
	// time of capture
	TimeStamp  string `json:"timestamp"`
	PacketData string `json:"packet_data"`
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

		var addr = fmt.Sprintf("localhost:%v", viper.GetString("ui.port"))
		log.Infof("starting ui on: http://%v", addr)
		http.HandleFunc("/packets", packets)
		http.HandleFunc("/", home)

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

func notifySockets(pv PacketView) {
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
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Info("read:", err)
			break
		}
		log.Info("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Info("write:", err)
			break
		}
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

func home(w http.ResponseWriter, r *http.Request) {
	if err := homeTemplate.Execute(w, "ws://"+r.Host+"/packets"); err != nil {
		log.Error(err)
	}
}

var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<script>  
window.addEventListener("load", function(evt) {

    var output = document.getElementById("output");
    var input = document.getElementById("input");
    var ws;

    var print = function(message) {
        var d = document.createElement("div");
        d.textContent = message;
        output.appendChild(d);
    };

    document.getElementById("open").onclick = function(evt) {
        if (ws) {
            return false;
        }
        ws = new WebSocket("{{.}}");
        ws.onopen = function(evt) {
			console.log(evt);
        }
        ws.onclose = function(evt) {
            ws = null;
			console.log(evt);
        }
        ws.onmessage = function(evt) {
			console.log(evt.data);
        }
        ws.onerror = function(evt) {
			console.log(evt);
        }
        return false;
    };

    document.getElementById("send").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("SEND: " + input.value);
        ws.send(input.value);
        return false;
    };

    document.getElementById("close").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        ws.close();
        return false;
    };

});
</script>
</head>
<body>
<table>
<tr><td valign="top" width="50%">
<p>Click "Open" to create a connection to the server, 
"Send" to send a message to the server and "Close" to close the connection. 
You can change the message and send multiple times.
<p>
<form>
<button id="open">Open</button>
<button id="close">Close</button>
<p><input id="input" type="text" value="Hello world!">
<button id="send">Send</button>
</form>
</td><td valign="top" width="50%">
<div id="output"></div>
</td></tr></table>
</body>
</html>
`))
