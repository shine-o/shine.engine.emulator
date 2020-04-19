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
	"sync"
	"time"
)

func init() {
	log = logger.Init("NetworkingLogger", true, false, ioutil.Discard)
	log.Info("networking logger init()")
}

type clientReader struct {
	r  *bufio.Reader
	mu sync.Mutex
}

type clientWriter struct {
	w  *bufio.Writer
	mu sync.Mutex
}

// CommandHandlers is a map of known operation codes linked to a caller
type CommandHandlers map[uint16]func(ctx context.Context, pc *Command)

// ShineService to be used by the calling shine service to:
// 	1. configure the settings for TCP socket
// 	2. assign the handlers for the operation codes handled by the shine service
//  3. use session factory specific to the shine service to create a session object in the context of each TCP connection
type ShineService struct {
	s  *Settings
	hw *HandleWarden
	sf SessionFactory
}

// NewShineService create new, the calling shine service must configure Settings and a HandlerWarden
func NewShineService(s *Settings, hw *HandleWarden) *ShineService {
	return &ShineService{
		s:  s,
		hw: hw,
	}
}

// Listen on TPC socket for connection on given port
func (ss *ShineService) Listen(ctx context.Context, port string) {
	ss.s.Set()
	if l, err := net.Listen("tcp4", fmt.Sprintf(":%v", port)); err == nil {
		log.Infof("listening for TCP connections on: %v", l.Addr())
		defer l.Close()
		rand.Seed(time.Now().Unix())
		for {
			select {
			case <-ctx.Done():
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

// UseSessionFactory given by the shine service
func (ss *ShineService) UseSessionFactory(sf SessionFactory) {
	ss.sf = sf
}

// for each connection launch go routine that handles tcp segment data
func (ss *ShineService) handleConnection(ctx context.Context, c net.Conn) {

	ctx, cancel := context.WithCancel(ctx)

	log.Infof("Serving %v", c.RemoteAddr().String())

	defer c.Close()
	defer cancel()
	defer log.Infof("Closing connection %v", c.RemoteAddr().String())

	var (
		buffer           = make([]byte, 1024)
		inboundSegments  = make(chan []byte, 4096)
		outboundSegments = make(chan []byte, 4096)
		r                *bufio.Reader
		w                *bufio.Writer
		//cw = &clientWriter{
		//	w: bufio.NewWriter(c),
		//}
	)
	r = bufio.NewReader(c)
	w = bufio.NewWriter(c)

	ctx = context.WithValue(ctx, ShineSession, ss.sf.New())
	ctx = context.WithValue(ctx, ConnectionWriter, outboundSegments)

	go ss.hw.handleInboundSegments(ctx, inboundSegments)

	go handleOutboundSegments(ctx, w, outboundSegments)

	for {
		if n, err := r.Read(buffer); err == nil {
			var data []byte
			data = append(data, buffer[:n]...)
			inboundSegments <- data
		} else {
			if err == io.EOF {
				break
			} else {
				log.Error(err)
				return
			}
		}
	}
}

func handleOutboundSegments(ctx context.Context, w *bufio.Writer, segment <-chan []byte) {
	for {
		select {
		case <-ctx.Done():
			log.Warning("handleOutboundSegments context canceled")
			return
		case data := <-segment:
			if _, err := w.Write(data); err != nil {
				log.Error(err)
			} else {
				if err = w.Flush(); err != nil {
					log.Error(err)
				}
			}
		}
	}
}

// Send bytes to the client
func (pc *Command) Send(ctx context.Context) {
	select {
	case <-ctx.Done():
		return
	default:
		cwv := ctx.Value(ConnectionWriter)
		cw := cwv.(chan []byte)

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
		//cw.mu.Lock()
		//if _, err := cw.w.Write(pc.Base.RawData()); err != nil {
		//	log.Error(err)
		//} else {
		//	if err = cw.w.Flush(); err != nil {
		//		log.Error(err)
		//	}
		//}
		//cw.mu.Unlock()
	}
}

func (pc *Command) SyncSend(ctx context.Context, err chan<- error) {
	select {
	case <-ctx.Done():
		return
	default:
		cwv := ctx.Value(ConnectionWriter)
		cw := cwv.(*clientWriter)

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

		cw.mu.Lock()
		if _, err := cw.w.Write(pc.Base.RawData()); err != nil {
			log.Error(err)
		} else {
			if err = cw.w.Flush(); err != nil {
				log.Error(err)
			}
		}
		cw.mu.Unlock()
	}
}
