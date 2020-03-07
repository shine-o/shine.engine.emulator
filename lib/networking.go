package lib

import (
	"bufio"
	"context"
	"encoding/hex"
	"fmt"
	"github.com/google/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"sync"
	"time"
)

type clientWriter struct {
	w  *bufio.Writer
	mu sync.Mutex
}

var (
	log *logger.Logger
)

// open TCP socket on port 9010
// for each connection use Context
func Listen(cmd *cobra.Command, args []string) {
	serveConfig()
	initHandlers()
	log = logger.Init("LoginLogger", true, true, ioutil.Discard)

	if l, err := net.Listen("tcp4", fmt.Sprintf(":%v", viper.GetInt("serve.port"))); err == nil {
		log.Infof("Listening for TCP connections on: %v", l.Addr())
		defer l.Close()
		rand.Seed(time.Now().Unix())
		for {
			if c, err := l.Accept(); err == nil {
				go handleConnection(c)
			} else {
				logger.Fatal(err)
			}
		}
	} else {
		logger.Error(err)
	}
}

func serveConfig() {
	if b, err := hex.DecodeString(viper.GetString("crypt.xorKey")); err != nil {
		log.Error(err)
	} else {
		xorKey = b
	}
}

func handleConnection(c net.Conn) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	log.Infof("Serving %v", c.RemoteAddr().String())

	defer c.Close()
	defer cancel()
	defer log.Infof("Closing connection %v", c.RemoteAddr().String())

	var (
		buf       = make([]byte, 1024)
		r         = bufio.NewReader(c)
		xorOffset uint16
		segment   = make(chan []byte)
		cw        = &clientWriter{
			w: bufio.NewWriter(c),
		}
	)

	ctx = context.WithValue(ctx, "connWriter", cw)
	ctx = context.WithValue(ctx, "xorOffset", &xorOffset)

	hw.mu.Lock()
	sendXorOffset := hw.handlers[2055]
	go sendXorOffset(ctx, &ProtocolCommand{
		pcb: ProtocolCommandBase{operationCode: 2055},
	})
	hw.mu.Unlock()

	go handleSegments(ctx, segment, &xorOffset)

	for {
		if n, err := r.Read(buf); err == nil {
			log.Infof("Received %v bytes", n)
			var tmpBuf []byte
			tmpBuf = append(tmpBuf, buf[:n]...)
			segment <- tmpBuf
		} else {
			if err == io.EOF {
				break
			} else {
				log.Fatal(err)
				return
			}
		}
	}
}
