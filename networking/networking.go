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
	ExtraParameters interface{}
}

type InboundSegments struct {
	Recv <-chan []byte
	Send chan<- []byte
}
type OutboundSegments struct {
	Recv <-chan []byte
	Send chan<- []byte
}

type Commands struct {
	Send chan<- *Command
	Recv <-chan *Command
}

type Network struct {
	InboundSegments
	OutboundSegments
	Commands
	CloseConnection chan bool
	Conn            net.Conn
	Reader          *bufio.Reader
	Writer          *bufio.Writer
}

const (
	// XorOffset indicates what offset in the xor hex table to use to start decrypting client data
	XorOffset ContextKey = iota
	// ShineSession if used, shine service can access session data within their handler's context
	ShineSession
	// OutboundStream is a utility struct which contains the tcp connection object and a mutex
	// it is used to write data to the client from any shine service handler
	NetworkVariables
)

var logInboundPackets  chan <- * Command
var logOutboundPackets chan <- * Command


// Listen on TPC socket for connection on given port
func (ss *ShineService) Listen(ctx context.Context, port string) {
	ss.Settings.Set()


	in :=  make(chan * Command, 4096)
	out :=  make(chan * Command, 4096)

	logInboundPackets  = in
	logOutboundPackets = out

	go logPackets(ctx, in, out)
	go logPackets(ctx, in, out)
	go logPackets(ctx, in, out)
	go logPackets(ctx, in, out)

	if l, err := net.Listen("tcp4", fmt.Sprintf(":%v", port)); err == nil {
		log.Infof("listening for TCP connections on: %v", l.Addr())
		defer l.Close()
		var src cryptoSource
		rnd := rand.New(src)
		rand.Seed(rnd.Int63n(time.Now().Unix()))
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if c, err := l.Accept(); err == nil {
					go ss.handleConnection(c)
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
func (ss *ShineService) handleConnection(conn net.Conn) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	log.Infof("serving %v", conn.RemoteAddr().String())

	defer conn.Close()
	defer cancel()
	defer log.Infof("closing connection %v", conn.RemoteAddr().String())

	var (
		buffer   = make([]byte, 4096)
		inbound  = make(chan []byte, 4096)
		outbound = make(chan []byte, 4096)
		commands = make(chan *Command, 4096)

		n = &Network{
			InboundSegments: InboundSegments{
				Recv: inbound,
				Send: inbound,
			},
			OutboundSegments: OutboundSegments{
				Recv: outbound,
				Send: outbound,
			},
			Commands: Commands{
				Send: commands,
				Recv: commands,
			},
			CloseConnection: make(chan bool),
			Conn:            conn,
			Reader:          bufio.NewReader(conn),
			Writer:          bufio.NewWriter(conn),
		}
	)

	ctx = context.WithValue(ctx, ShineSession, ss.SessionFactory.New())
	ctx = context.WithValue(ctx, NetworkVariables, n)

	go ss.handleInboundSegments(ctx, n)
	go ss.handleOutboundSegments(ctx, n)
	go waitForClose(n)

	for {
		size, err := n.Reader.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				log.Error(err)
				return
			}
		}
		data := make([]byte, size)
		copy(data, buffer[:size])
		n.InboundSegments.Send <- data
	}
}

func waitForClose(n *Network) {
	for {
		select {
		case <-n.CloseConnection:
			err := n.Conn.Close()
			if err != nil {
				log.Error(err)
			}
			return
		}
	}
}

// Send bytes to the client
func (pc *Command) Send(ctx context.Context) {
	nv := ctx.Value(NetworkVariables)
	n := nv.(*Network)

	if pc.NcStruct != nil {
		data, err := structs.Pack(pc.NcStruct)
		if err != nil {
			log.Error(err)
			return
		}
		pc.Base.Data = data
	}
	n.OutboundSegments.Send <- pc.Base.RawData()
	logOutboundPackets <- pc
}

func (pc * Command) SendDirectly(outboundStream chan<- []byte) {
	if pc.NcStruct != nil {
		data, err := structs.Pack(pc.NcStruct)
		if err != nil {
			log.Error(err)
			return
		}
		pc.Base.Data = data
	}
	outboundStream <- pc.Base.RawData()
	logOutboundPackets <- pc
}


func logPackets(ctx context.Context, in <- chan* Command, out <-chan *Command) {
	for {
		select {
		case <- ctx.Done():
			return
		case ipc := <- in:
			logDirection(ipc, "inbound")
		case opc := <- out:
			logDirection(opc, "outbound")
		}
	}
}

func logDirection(pc * Command, direction string) {
	pc.RLock()
	defer pc.RUnlock()
	cn :=  commandName(pc)
	log.Infof("%v %v packet metadata: %v", direction, cn, pc.Base.String())
	if pc.NcStruct != nil {
		sd, err := json.Marshal(pc.NcStruct)
		if err != nil {
			log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(pc.NcStruct).String(), err)
		} else {
			log.Infof("%v %v packet structure data: %v %v", direction, cn, reflect.TypeOf(pc.NcStruct).String(), string(sd))
		}
	}
}

func commandName(pc *Command) string {
	commandList.mu.Lock()
	defer 		commandList.mu.Unlock()
	if (&PCList{}) != commandList {// should be commented out on production to increase performance
		opCode := pc.Base.OperationCode
		department := opCode >> 10
		command := opCode &1023
		if dpt, ok := commandList.Departments[uint8(department)]; ok {
			return  dpt.ProcessedCommands[fmt.Sprintf("%X", command)]
		} else {
			log.Warningf("Missing friendly name for command with: operationCode %v,  department %v, command %v, ", opCode, department, fmt.Sprintf("%X", command))
		}
	}
	return ""
}