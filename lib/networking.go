package lib

import (
	"bufio"
	"context"
	"encoding/hex"
	"fmt"
	"github.com/google/logger"
	//protocol "shine.engine.packet-protocol"
	protocol "github.com/shine-o/shine.engine.protocol"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type clientReader struct {
	r  *bufio.Reader
	mu sync.Mutex
}

type clientWriter struct {
	w  *bufio.Writer
	mu sync.Mutex
}

var (
	log *logger.Logger
)

func serveSetup() {
	settings := &protocol.Settings{}
	log = logger.Init("LoginLogger", true, true, ioutil.Discard)
	if xk, err := hex.DecodeString(viper.GetString("crypt.xorKey")); err != nil {
		log.Error(err)
		os.Exit(1)
	} else {
		settings.XorKey = xk
	}
	settings.XorLimit = uint16(viper.GetInt("crypt.xorLimit"))

	if path, err := filepath.Abs(viper.GetString("protocol.nc-data")); err != nil {
		log.Error(err)
	} else {
		settings.CommandsFilePath = path
	}
	settings.Set()
	loginHandlers()
}

// open tcp port
func Listen(cmd *cobra.Command, args []string) {
	serveSetup()

	if l, err := net.Listen("tcp4", fmt.Sprintf(":%v", viper.GetInt("serve.port"))); err == nil {
		log.Infof("Listening for TCP connections on: %v", l.Addr())
		defer l.Close()
		rand.Seed(time.Now().Unix())
		for {
			if c, err := l.Accept(); err == nil {
				go handleConnection(c)
			} else {
				log.Fatal(err)
			}
		}
	} else {
		log.Error(err)
	}
}

// for each connection launch go routine that handles tcp segment data
func handleConnection(c net.Conn) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	log.Infof("Serving %v", c.RemoteAddr().String())

	defer c.Close()
	defer cancel()
	defer log.Infof("Closing connection %v", c.RemoteAddr().String())

	var (
		buffer  = make([]byte, 1024)
		segment = make(chan []byte)
		cr      = &clientReader{
			r: bufio.NewReader(c),
		}
		cw = &clientWriter{
			w: bufio.NewWriter(c),
		}
	)

	ctx = context.WithValue(ctx, "connWriter", cw)

	go handleLoginSegments(ctx, segment)

	for {
		cr.mu.Lock()
		if n, err := cr.r.Read(buffer); err == nil {
			//log.Infof("Received %v bytes", n)
			var data []byte
			data = append(data, buffer[:n]...)
			segment <- data
		} else {
			if err == io.EOF {
				cr.mu.Unlock()
				break
			} else {
				log.Fatal(err)
				cr.mu.Unlock()
				return
			}
		}
		cr.mu.Unlock()
	}
}
