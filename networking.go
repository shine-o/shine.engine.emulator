package networking

import (
	"bufio"
	"context"
	"fmt"
	"github.com/google/logger"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
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

type ShineService struct {
	s  *Settings
	hw *HandleWarden
	sf SessionFactory
}

func NewShineService(s *Settings, hw *HandleWarden) *ShineService {
	return &ShineService{
		s:  s,
		hw: hw,
	}
}

// listen on  tcp socket
func (ss *ShineService) Listen(ctx context.Context, port string) {
	//log = logger.Init("NetworkingLogger", true, false, ioutil.Discard)
	//log.Info("Networking Logger init()")
	ss.s.Set()
	if l, err := net.Listen("tcp4", fmt.Sprintf(":%v", port)); err == nil {
		log.Infof("Listening for TCP connections on: %v", l.Addr())
		defer l.Close()
		rand.Seed(time.Now().Unix())
		for {
			if c, err := l.Accept(); err == nil {
				go ss.handleConnection(ctx, c)
			} else {
				log.Error(err)
			}
		}
	} else {
		log.Error(err)
	}
}

func (ss *ShineService) UseSessionFactory(sf SessionFactory)  {
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
		buffer  = make([]byte, 1024)
		segment = make(chan []byte)
		cr      = &clientReader{
			r: bufio.NewReader(c),
		}
		cw = &clientWriter{
			w: bufio.NewWriter(c),
		}
	)

	ctx = context.WithValue(ctx, "session", ss.sf.New())
	ctx = context.WithValue(ctx, "connWriter", cw)

	go ss.hw.handleSegments(ctx, segment)

	for {
		cr.mu.Lock()
		if n, err := cr.r.Read(buffer); err == nil {
			var data []byte
			data = append(data, buffer[:n]...)
			segment <- data
		} else {
			if err == io.EOF {
				cr.mu.Unlock()
				break
			} else {
				log.Error(err)
				cr.mu.Unlock()
				return
			}
		}
		cr.mu.Unlock()
	}
}

func WriteToClient(ctx context.Context, pc *Command) {
	select {
	case <-ctx.Done():
		return
	default:
		cwv := ctx.Value("connWriter")
		cw := cwv.(*clientWriter)
		log.Infof("Outbound packet: %v", pc.Base.String())
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
