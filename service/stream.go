package service

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/tcpassembly"
	"github.com/google/logger"
	"github.com/google/uuid"
	"github.com/shine-o/shine.engine.networking"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

func init() {
	log = logger.Init("SnifferLogger", true, false, ioutil.Discard)
	log.Info("sniffer logger init()")
}

type shineStreamFactory struct {
	shineContext context.Context
}

type shineStream struct {
	flowID         string
	net, transport gopacket.Flow
	flowName       string
	clientIP       string
	serverIP       string
	segments       chan<- shineSegment
	xorKey         chan<- uint16 // only used by decodeClientPackets()
	packetID       uint64
	cancel         context.CancelFunc
}

// Shine utility struct for storing service data
type Shine struct {
	knownServices map[int]*service // port => serviceName
	mu            sync.Mutex
}

type service struct {
	name string
}

type shineSegment struct {
	data []byte
	seen time.Time
}

type shineStreams struct {
	toClient   map[string]*shineStream
	fromClient map[string]*shineStream
	mu         sync.Mutex
}

type contextKey int

const (
	activeShineStreams contextKey = iota
)

var (
	iface   string
	snaplen int
	filter  string
	shine   Shine
	log     *logger.Logger
)

func config() {
	// remove output folder if exists, create it again
	dir, err := filepath.Abs("output/")
	if _, err := os.Stat(dir); os.IsNotExist(err) {

	} else {
		err = os.RemoveAll(dir)
		if err != nil {
			log.Error(err)
		}
	}

	err = os.Mkdir(dir, 0700)

	if err != nil {
		log.Error(err)
	}

	iface = viper.GetString("network.interface")
	snaplen = viper.GetInt("network.snaplen")

	//serverIP := viper.GetString("network.serverIP")
	startPort := viper.GetString("network.portRange.start")
	endPort := viper.GetString("network.portRange.end")
	portRange := fmt.Sprintf("%v-%v", startPort, endPort)

	//filter = fmt.Sprintf("(dst net %v or src net %v) and (dst portrange %v or src portrange %v)", serverIP, serverIP, portRange, portRange)
	filter = fmt.Sprintf("dst portrange %v or src portrange %v", portRange, portRange)

	shine.mu.Lock()
	shine.knownServices = make(map[int]*service)
	if viper.IsSet("protocol.services") {
		// snippet for loading yaml array
		services := make([]map[string]string, 0)
		var m map[string]string
		servicesI := viper.Get("protocol.services")
		servicesS := servicesI.([]interface{})
		for _, s := range servicesS {
			serviceMap := s.(map[interface{}]interface{})
			m = make(map[string]string)
			for k, v := range serviceMap {
				m[k.(string)] = v.(string)
			}
			services = append(services, m)
		}
		for _, v := range services {
			port, err := strconv.Atoi(v["port"])
			if err != nil {
				log.Error(err)
			}
			shine.knownServices[port] = &service{name: v["name"]}
		}
	} else {
		shine.knownServices[9000] = &service{name: "Account"}
		shine.knownServices[9311] = &service{name: "AccountLog"}
		shine.knownServices[9411] = &service{name: "Character"}
		shine.knownServices[9511] = &service{name: "GameLog"}
		shine.knownServices[9010] = &service{name: "Login"}
		shine.knownServices[9110] = &service{name: "WorldManager"}
		shine.knownServices[9210] = &service{name: "Zone00"}
		shine.knownServices[9212] = &service{name: "Zone01"}
		shine.knownServices[9214] = &service{name: "Zone02"}
		shine.knownServices[9216] = &service{name: "Zone03"}
		shine.knownServices[9218] = &service{name: "Zone04"}
	}
	shine.mu.Unlock()

	s := &networking.Settings{}

	xorKey, err := hex.DecodeString(viper.GetString("protocol.xorKey"))
	if err != nil {
		log.Fatal(err)
	}

	s.XorKey = xorKey

	xorLimit, err := strconv.Atoi(viper.GetString("protocol.xorLimit"))

	if err != nil {
		log.Fatal(err)
	}

	s.XorLimit = uint16(xorLimit)
	if path, err := filepath.Abs(viper.GetString("protocol.commands")); err != nil {
		log.Error(err)
	} else {
		s.CommandsFilePath = path
	}
	s.Set()
}

func (ssf *shineStreamFactory) New(net, transport gopacket.Flow) tcpassembly.Stream {

	var (
		clientIP string
		serverIP string
	)

	// create new cancel func from context [will be called from reassembly complete]
	// assign context to shineStream
	ctx, cancel := context.WithCancel(ssf.shineContext)

	srcIP := net.Src().String()
	srcPort, _ := strconv.Atoi(transport.Src().String())

	segments := make(chan shineSegment, 512)
	xorKey := make(chan uint16)

	//decodedPackets := make(chan, 500)

	s := &shineStream{
		flowID:    uuid.New().String(),
		net:       net,
		transport: transport,
		segments:  segments,
		packetID:  0,
		xorKey:    xorKey,
		cancel:    cancel,
	}

	key := fmt.Sprintf("%v:%v", srcIP, srcPort)

	if srcPort >= 9000 && srcPort <= 9600 {
		// server - client

		serverIP = net.Src().String()
		clientIP = net.Dst().String()

		service, ok := shine.knownServices[srcPort]
		if !ok {
			//log.Fatal("something went horribly wrong")
			s.flowName = fmt.Sprintf("%v-client", "unknown")
		} else {
			s.flowName = fmt.Sprintf("%v-client", strings.ToLower(service.name))
		}
		log.Infof("new server stream from => [ %v - %v] [%v]", srcIP, srcPort, s.flowName)
		ss, ok := ssf.shineContext.Value(activeShineStreams).(*shineStreams)
		ss.mu.Lock()
		if !ok {
			log.Fatalf("unexpected struct type: %v", reflect.TypeOf(ss).String())
		}
		key := fmt.Sprintf("%v:%v", srcIP, srcPort)
		ss.toClient[key] = s
		ss.mu.Unlock()
		go s.decodeServerPackets(ctx, segments)
	} else {
		// client-server

		clientIP = net.Src().String()
		serverIP = net.Dst().String()

		dstPort, _ := strconv.Atoi(transport.Dst().String())

		service, ok := shine.knownServices[dstPort]
		if !ok {
			s.flowName = fmt.Sprintf("clien-%v", "unknown")
		} else {
			s.flowName = fmt.Sprintf("client-%v", strings.ToLower(service.name))
		}
		log.Infof("new server stream from => [ %v - %v] [%v]", srcIP, srcPort, s.flowName)
		ss, ok := ssf.shineContext.Value(activeShineStreams).(*shineStreams)
		ss.mu.Lock()
		if !ok {
			log.Fatalf("unexpected struct type: %v", reflect.TypeOf(ss).String())
		}
		ss.fromClient[key] = s
		ss.mu.Unlock()
		go s.decodeClientPackets(ctx, segments, xorKey)
	}
	s.clientIP = clientIP
	s.serverIP = serverIP
	return s
}

func (ss *shineStream) Reassembled(reassemblies []tcpassembly.Reassembly) {
	for _, reassembly := range reassemblies {
		if len(reassembly.Bytes) == 0 {
			continue
		}
		seg := shineSegment{data: reassembly.Bytes, seen: reassembly.Seen}
		ss.segments <- seg
	}
}

func (ss *shineStream) ReassemblyComplete() {
	log.Warningf("reassembly complete for stream [ %v - %v]", ss.net.String(), ss.transport.String()) // ip of the stream, port of the stream
	ss.cancel()
	cf := CompletedFlow{
		FlowCompleted: true,
		ClientIP:      ss.clientIP,
		FlowID:        ss.flowID,
	}
	go uiCompletedFlow(cf)
}

// handle stream data flowing from the client
func (ss *shineStream) decodeClientPackets(ctx context.Context, segments <-chan shineSegment, xorKey <-chan uint16) {
	var (
		d         []byte
		offset    int
		xorOffset uint16
	)
	offset = 0
	xorOffset = 1500 // impossible value
	logActivated := viper.GetBool("protocol.log.client")

	// block until xorKey is found
	for {
		select {
		case <-ctx.Done():
			log.Warningf("[%v] decodeClientPackets(): context was canceled", ss.flowName)
			return
		case xorOffset = <-xorKey:
			break
		case segment := <-segments:
			d = append(d, segment.data...)
			if offset > len(d) {
				log.Warningf("not enough data, next offset is %v ", offset)
				break
			}

			if xorOffset == 1500 { // impossible value
				break
			}

			if offset != len(d) {
				var skipBytes int
				var pLen int
				var pType string
				var rs []byte

				pLen, pType = networking.PacketBoundary(offset, d)

				if pType == "small" {
					skipBytes = 1
				} else {
					skipBytes = 3
				}

				nextOffset := offset + skipBytes + pLen
				if nextOffset > len(d) {
					log.Warningf("not enough data, next offset is %v ", nextOffset)
					break
				}

				rs = append(rs, d[offset+skipBytes:nextOffset]...)
				networking.XorCipher(rs, &xorOffset)
				pc, err := networking.DecodePacket(pType, pLen, rs)
				if err != nil {
					log.Error(err)
				}
				ss.packetID++

				if logActivated {
					go ss.handlePacket(ctx, segment.seen, pc)
				}
				offset += skipBytes + pLen
			}
		}
	}
}

// handle stream data flowing from the server
func (ss *shineStream) decodeServerPackets(ctx context.Context, segments <-chan shineSegment) {
	var (
		d              []byte
		offset         int
		xorOffsetFound bool
	)
	xorOffsetFound = false
	offset = 0

	logActivated := viper.GetBool("protocol.log.server")

	for {
		select {
		case <-ctx.Done():
			log.Warningf("[%v] decodeServerPackets(): context was canceled", ss.flowName)
			return
		case segment := <-segments:
			d = append(d, segment.data...)
			if offset > len(d) {
				log.Warningf("not enough data, next offset is %v ", offset)
				break
			}

			if offset != len(d) {
				var skipBytes int
				var pLen int
				var pType string
				var rs []byte

				pLen, pType = networking.PacketBoundary(offset, d)

				if pType == "small" {
					skipBytes = 1
				} else {
					skipBytes = 3
				}

				nextOffset := offset + skipBytes + pLen
				if nextOffset > len(d) {
					log.Warningf("not enough data, next offset is %v ", nextOffset)
					break
				}

				rs = append(rs, d[offset+skipBytes:nextOffset]...)

				pc, err := networking.DecodePacket(pType, pLen, rs)
				if err != nil {
					log.Error(err)
				}

				if !xorOffsetFound {
					if pc.Base.OperationCode == 2055 {
						var xorOffset uint16
						buf := bytes.NewBuffer(pc.Base.Data)
						if err := binary.Read(buf, binary.LittleEndian, &xorOffset); err != nil {
							log.Error(err)
							return
						}
						xorOffsetFound = true
						// LOL
						ass, ok := ctx.Value(activeShineStreams).(*shineStreams)
						if !ok {
							log.Errorf("unexpected struct type: %v", reflect.TypeOf(ss).String())
							return
						}
						ass.mu.Lock()

						dstIP := ss.net.Dst().String()
						dstPort, _ := strconv.Atoi(ss.transport.Dst().String())

						key := fmt.Sprintf("%v:%v", dstIP, dstPort)

						if ss, ok := ass.fromClient[key]; ok {
							ss.xorKey <- xorOffset
						} else {
							log.Errorf("unexpected struct type: %v", reflect.TypeOf(ss).String())
						}
						ass.mu.Unlock()
						log.Warningf("xorOffset: %v found for client  %v", xorOffset, key)
					}
				}
				ss.packetID++

				if logActivated {
					go ss.handlePacket(ctx, segment.seen, pc)
				}
				offset += skipBytes + pLen
			}
		}
	}
}

// prepare for frontend and send it to the active socket connections
func (ss *shineStream) handlePacket(ctx context.Context, seen time.Time, pc networking.Command) {
	select {
	case <-ctx.Done():
		return
	default:
		pv := PacketView{
			FlowID:     ss.flowID,
			FlowName:   ss.flowName,
			PacketID:   ss.packetID,
			ClientIP:   ss.clientIP,
			ServerIP:   ss.serverIP,
			TimeStamp:  seen.String(),
			PacketData: pc.Base.String(),
		}

		nr, err := ncStructRepresentation(pc.Base.OperationCode, pc.Base.Data)
		if err == nil {
			pv.NcRepresentation = nr
			b, _ := json.Marshal(pv.NcRepresentation)
			log.Info(string(b))
		}

		data, err := json.Marshal(&pv)

		if err != nil {
			log.Error(err)
			return
		}

		if viper.GetBool("protocol.log.verbose") {
			log.Infof("[%v] [%v] [%v] %v", ss.packetID, ss.flowName, seen, string(data))
		} else {
			log.Infof("[%v] [%v] [%v] %v", ss.packetID, ss.flowName, seen, pc.Base.String())
		}
		ocs.mu.Lock()
		ocs.structs[pc.Base.OperationCode] = pc.Base.ClientStructName
		//fmt.Println(ocs.structs)
		ocs.mu.Unlock()
		go sendPacketToUI(pv)
		//generateStructCode(pc)

	}
}


// Capture packets and decode them
func Capture(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	config()

	ocs = &OpCodeStructs{
		structs: make(map[uint16]string),
	}
	ss := &shineStreams{
		toClient:   make(map[string]*shineStream),
		fromClient: make(map[string]*shineStream),
	}

	ctx = context.WithValue(ctx, activeShineStreams, ss)

	sf := &shineStreamFactory{
		shineContext: ctx,
	}
	sp := tcpassembly.NewStreamPool(sf)
	a := tcpassembly.NewAssembler(sp)

	go capturePackets(ctx, a)
	go startUI(ctx)


	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	select {
	case <-c:
		cancel()
		generateOpCodeSwitch()
	}
	<-c

}

func capturePackets(ctx context.Context, a *tcpassembly.Assembler) {
	handle, err := pcap.OpenLive(iface, int32(snaplen), true, pcap.BlockForever)
	if err != nil {
		log.Fatal(err)
	}
	err = handle.SetBPFFilter(filter)
	if err != nil {
		log.Fatal(err)
	}
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for {
		select {
		case <- ctx.Done():
			return
		case packet := <- packetSource.Packets():
			tcp := packet.TransportLayer().(*layers.TCP)
			a.AssembleWithTimestamp(packet.NetworkLayer().NetworkFlow(), tcp, packet.Metadata().Timestamp)
		}
	}
}