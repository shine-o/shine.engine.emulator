package networking

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/logger"
	"github.com/shine-o/shine.engine.core/structs"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"reflect"
	"time"
)

func init() {
	log = logger.Init("networking logger", true, false, ioutil.Discard)
	log.Info("networking logger init()")
}

type ShineService struct {
	Settings
	ShineHandler
	SessionFactory
	ExtraParameters interface {}
}

// Listen on TPC socket for connection on given port
func (ss * ShineService) Listen(ctx context.Context, port string) {
	ss.Settings.Set()
	if l, err := net.Listen("tcp4", fmt.Sprintf(":%v", port)); err == nil {
		log.Infof("listening for TCP connections on: %v", l.Addr())
		defer l.Close()
		var src cryptoSource
		rnd := rand.New(src)
		rand.Seed(rnd.Int63n(time.Now().Unix()))
		for {
			select {
			case <- ctx.Done():
				return
			default:
				if c, err := l.Accept(); err == nil {
					go ss.handleConnection(ctx, c)
				} else {
					log.Error(err)
				}
			}
		}
	} else {
		log.Error(err)
	}
}

// for each connection launch go routine that handles tcp segment data
func (ss * ShineService) handleConnection(ctx context.Context, c net.Conn) {

	ctx, cancel := context.WithCancel(ctx)

	log.Infof("serving %v", c.RemoteAddr().String())

	defer c.Close()
	defer cancel()
	defer log.Infof("closing connection %v", c.RemoteAddr().String())

	var (
		buffer           = make([]byte, 4096)
		inboundSegments  = make(chan []byte, 4096)
		outboundSegments = make(chan []byte, 4096)
		closeConnection  = make(chan bool)
		r                *bufio.Reader
		w                *bufio.Writer
	)
	r = bufio.NewReader(c)
	w = bufio.NewWriter(c)

	ctx = context.WithValue(ctx, ShineSession, ss.SessionFactory.New())
	ctx = context.WithValue(ctx, ConnectionWriter, outboundSegments)

	go ss.handleInboundSegments(ctx, inboundSegments, closeConnection)
	go ss.handleOutboundSegments(ctx, w, outboundSegments)
	go waitForClose(closeConnection, c)
	for {
		n, err := r.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				log.Error(err)
				cancel()
				return
			}
		}
		data := make([]byte, n)
		copy(data, buffer[:n])
		inboundSegments <- data
	}
}

func waitForClose(close <-chan bool, c net.Conn) {
	for {
		select {
		case <-close:
			c.Close()
			return
		}
	}
}

// Send bytes to the client
func (pc *Command) Send(ctx context.Context) {
	cwv := ctx.Value(ConnectionWriter)
	cw := cwv.(chan []byte) //maybe the handlers themselves should receive the outboundSegments channel as parameter

	if pc.NcStruct != nil {
		data, err := structs.Pack(pc.NcStruct)
		if err != nil {
			log.Error(err)
			return
		}
		pc.Base.Data = data
		sd, err := json.Marshal(pc.NcStruct)
		if err != nil {
			log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(pc.NcStruct).String(), err)
		}
		log.Infof("[outbound] structured packet data: %v %v", reflect.TypeOf(pc.NcStruct).String(), string(sd))
	}
	log.Infof("[outbound] metadata: %v", pc.Base.String())
	cw <- pc.Base.RawData()
}
