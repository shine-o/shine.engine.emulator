package service

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/tcpassembly"
	"github.com/google/logger"
	networking "github.com/shine-o/shine.engine.networking"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"sync"
	"time"
)

func init() {
	log = logger.Init("SnifferLogger", true, false, ioutil.Discard)
	log.Info("sniffer logger init()")
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

type shineStreamFactory struct {
	shineContext context.Context
}

type shineStream struct {
	net, transport gopacket.Flow
	segments       chan<- shineSegment
	xorKey         chan<- uint16 // only used by decodeClientPackets()
	cancel         context.CancelFunc
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

	log *logger.Logger
)

var shine Shine

func captureConfig() {
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

	serverIP := viper.GetString("network.serverIP")
	startPort := viper.GetString("network.portRange.start")
	endPort := viper.GetString("network.portRange.end")
	portRange := fmt.Sprintf("%v-%v", startPort, endPort)

	filter = fmt.Sprintf("(dst net %v or src net %v) and (dst portrange %v or src portrange %v)", serverIP, serverIP, portRange, portRange)

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
	log.Error(filepath.Abs(viper.GetString("protocol.commands")))
	if path, err := filepath.Abs(viper.GetString("protocol.commands")); err != nil {
		log.Error(err)
	} else {
		s.CommandsFilePath = path
	}
	s.Set()
}

func (ssf *shineStreamFactory) New(net, transport gopacket.Flow) tcpassembly.Stream {
	ctx, cancel := context.WithCancel(ssf.shineContext)

	srcIP := net.Src().String()
	srcPort, _ := strconv.Atoi(transport.Src().String())

	// create new cancel func from context [will be called from reassembly complete]
	// assign context to shinestream
	segments := make(chan shineSegment, 512)
	xorKey := make(chan uint16)

	s := &shineStream{
		net:       net,
		transport: transport,
		segments:  segments,
		xorKey:    xorKey,
		cancel:    cancel,
	}
	key := fmt.Sprintf("%v:%v", srcIP, srcPort)
	if srcPort >= 9000 && srcPort <= 9600 {
		log.Infof("new server stream from => [ %v - %v]", srcIP, srcPort) // ip of the stream, port of the stream
		// server - client
		ss, ok := ssf.shineContext.Value(activeShineStreams).(*shineStreams)
		ss.mu.Lock()
		if !ok {
			log.Errorf("unexpected struct type: %v", reflect.TypeOf(ss).String())
			return s
		}
		key := fmt.Sprintf("%v:%v", srcIP, srcPort)
		ss.toClient[key] = s
		ss.mu.Unlock()
		go s.decodeServerPackets(ctx, segments)
	} else {
		log.Infof("new client stream from => [ %v - %v]", srcIP, srcPort) // ip of the stream, port of the stream
		ss, ok := ssf.shineContext.Value(activeShineStreams).(*shineStreams)
		ss.mu.Lock()
		if !ok {
			log.Errorf("unexpected struct type: %v", reflect.TypeOf(ss).String())
			return s
		}
		ss.fromClient[key] = s
		ss.mu.Unlock()
		go s.decodeClientPackets(ctx, segments, xorKey)
		return s
	}

	return s
}

func (fs *shineStream) Reassembled(reassemblies []tcpassembly.Reassembly) {
	for _, reassembly := range reassemblies {
		if len(reassembly.Bytes) == 0 {
			continue
		}
		seg := shineSegment{data: reassembly.Bytes, seen: reassembly.Seen}
		fs.segments <- seg
	}
}

func (fs *shineStream) ReassemblyComplete() {
	log.Warningf("reassembly complete for stream [ %v - %v]", fs.net.String(), fs.transport.String()) // ip of the stream, port of the stream

}

func (fs *shineStream) decodeClientPackets(ctx context.Context, segments <-chan shineSegment, xorKey <-chan uint16) {
	var d []byte
	var offset int
	var xorOffset uint16
	offset = 0
	xorOffset = 1500 // impossible value
	// block until xorKey is found
	for {
		select {
		case <-ctx.Done():
			log.Error("context was canceled")
			return
		case xorOffset = <-xorKey:
			break
		case segment := <-segments:
			d = append(d, segment.data...)
			if offset > len(d) {
				log.Info("not enough data, next offset is %v ", offset)
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
					log.Info("not enough data, next offset is %v ", nextOffset)
					break
				}

				rs = append(rs, d[offset+skipBytes:nextOffset]...)
				networking.XorCipher(rs, &xorOffset)
				pc, err := networking.DecodePacket(pType, pLen, rs)
				if err != nil {
					log.Error(err)
				}
				log.Infof("[%v] %v", segment.seen, pc.Base.String())
				offset += skipBytes + pLen
			}
		}
	}
}

func (fs *shineStream) decodeServerPackets(ctx context.Context, segments <-chan shineSegment) {
	// decode packets normally
	// if xorKey is found, use ctx.value.fiestaStreams to find opposite stream
	var (
		d              []byte
		offset         int
		xorOffsetFound bool
	)
	xorOffsetFound = false
	offset = 0

	for {
		select {
		case <-ctx.Done():
			log.Error("context was canceled")
			return
		case segment := <-segments:
			d = append(d, segment.data...)
			if offset > len(d) {
				log.Info("not enough data, next offset is %v ", offset)
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
					log.Info("not enough data, next offset is %v ", nextOffset)
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
						ss, ok := ctx.Value(activeShineStreams).(*shineStreams)
						if !ok {
							log.Errorf("unexpected struct type: %v", reflect.TypeOf(ss).String())
							return
						}
						ss.mu.Lock()

						//srcIP := fs.net.Src().String()
						dstIP := fs.net.Dst().String()
						//srcPort, _ := strconv.Atoi(fs.transport.Src().String())
						dstPort, _ := strconv.Atoi(fs.transport.Dst().String())
						key := fmt.Sprintf("%v:%v", dstIP, dstPort)

						if ss, ok := ss.fromClient[key]; ok {
							ss.xorKey <- xorOffset
						} else {
							log.Errorf("unexpected struct type: %v", reflect.TypeOf(ss).String())

						}
						ss.mu.Unlock()
						log.Errorf("xorOffset: %v found for client  %v", xorOffset, key)
					}
				}

				log.Infof("[%v] %v", segment.seen, pc.Base.String())
				offset += skipBytes + pLen
				// go fs.handlePacketDAta
			}
		}
	}
}

// Capture packets and decode them
func Capture(cmd *cobra.Command, args []string) {
	captureConfig()

	pf := &Flows{
		pfm: make(map[string][]gopacket.Packet),
	}

	ctx := context.Background()

	ss := &shineStreams{
		toClient:   make(map[string]*shineStream),
		fromClient: make(map[string]*shineStream),
	}

	ctx = context.WithValue(ctx, activeShineStreams, ss)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	sf := &shineStreamFactory{
		shineContext: ctx,
	}

	sp := tcpassembly.NewStreamPool(sf)
	a := tcpassembly.NewAssembler(sp)
	if handle, err := pcap.OpenLive(iface, int32(snaplen), true, pcap.BlockForever); err != nil {
		log.Fatal(err)
	} else if err := handle.SetBPFFilter(filter); err != nil { //
		log.Fatal(err)
	} else {
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
		for packet := range packetSource.Packets() {
			go pf.add(packet)
			tcp := packet.TransportLayer().(*layers.TCP)
			a.AssembleWithTimestamp(packet.NetworkLayer().NetworkFlow(), tcp, packet.Metadata().Timestamp)
		}
	}
}

func parseForFrontend() {

	// send to the client, the following
	//	{
	//		"clients": [
	//	{
	//		"uuid": "123e4567-e89b-12d3-a456-426655440000",
	//		"status": "alive",
	//		"ip" : "192.168.1.131",
	//		"flows": [
	//		{
	//			"key": "client-login",
	//			"status": "closed",
	//			"packets": [
	//			{
	//				"metadata": {
	//
	//				},
	//				"data": {
	//					"packetType": "small",
	//					"length": 24,
	//					"department": 8,
	//					"command": "1A",
	//					"opCode": 8218,
	//					"data": "a1086c0e0000560800006c0e0000560800006e006c20",
	//					"rawData": "181a20a1086c0e0000560800006c0e0000560800006e006c20",
	//					"friendlyName": "NC_ACT_SOMEONEMOVERUN_CMD"
	//				}
	//			}
	//		]
	//		},
	//		{
	//			"key": "client-zone01",
	//			"status": "active",
	//			"packets": [
	//			{
	//				"metadata": {
	//
	//				},
	//				"data": {
	//					"packetType": "small",
	//					"length": 24,
	//					"department": 8,
	//					"command": "1A",
	//					"opCode": 8218,
	//					"data": "a1086c0e0000560800006c0e0000560800006e006c20",
	//					"rawData": "181a20a1086c0e0000560800006c0e0000560800006e006c20",
	//					"friendlyName": "NC_ACT_SOMEONEMOVERUN_CMD"
	//				}
	//			}
	//		]
	//		}
	//	]
	//	}
	//]
	//}

}
