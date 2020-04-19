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
	"github.com/google/gopacket/reassembly"
	"github.com/google/logger"
	"github.com/google/uuid"
	"github.com/segmentio/ksuid"
	"github.com/shine-o/shine.engine.core/networking"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"
	"syscall"
	"time"
)

func init() {
	lf, err := os.OpenFile("streams.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		logger.Fatalf("Failed to open log file: %v", err)
	}
	log = logger.Init("SnifferLogger", true, false ,lf)
	log.Info("sniffer logger init()")
}

type shineStreamFactory struct {
	shineContext   context.Context
	localAddresses []pcap.InterfaceAddress
}

type decodedPacket struct {
	seen   time.Time
	packet networking.Command
	direction string
}

type shineStream struct {
	flowID         string
	net, transport gopacket.Flow
	client         chan<- shineSegment
	server         chan<- shineSegment
	packets        chan<- decodedPacket
	xorKey         chan<- uint16
	cancel         context.CancelFunc
	isServer       bool
	mu 	sync.Mutex
}

type shineSegment struct {
	data []byte
	seen time.Time
	direction string
}

var (
	iface   string
	snaplen int
	filter  string
	log     *logger.Logger
)

func config() {
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

	startPort := viper.GetString("network.portRange.start")
	endPort := viper.GetString("network.portRange.end")
	portRange := fmt.Sprintf("%v-%v", startPort, endPort)
	//filter = fmt.Sprintf("tcp and (dst portrange %v or src portrange %v)", portRange, portRange)
	filter = fmt.Sprintf("tcp and portrange %v", portRange)

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

func (ss *shineStream) handleDecodedPackets(ctx context.Context, decodedPackets <-chan decodedPacket) {
	for {
		select {
		case <-ctx.Done():
			return
		case dp := <-decodedPackets:
			go ss.logPacket(dp)
		}
	}
}

func (ss *shineStream) logPacket(dp decodedPacket) {
	packetID, err := ksuid.NewRandomWithTime(dp.seen)
	if err != nil {
		log.Error(err)
	}

	pv := PacketView{
		PacketID:      packetID.String(),
		TimeStamp:     dp.seen.String(),
		IPEndpoints:   ss.net.String(),
		PortEndpoints: ss.transport.String(),
		Direction:     dp.direction,
		PacketData:    dp.packet.Base.JSON(),
	}

	nr, err := ncStructRepresentation(dp.packet.Base.OperationCode, dp.packet.Base.Data)
	if err == nil {
		pv.NcRepresentation = nr
		//b, _ := json.Marshal(pv.ncRepresentation)
		//log.Info(string(b))
	} else {
		//log.Error(err)
	}

	var tPorts string
	if dp.direction == "inbound" {
		tPorts = ss.transport.Reverse().String()
	} else {
		tPorts = ss.transport.String()
	}
	if viper.GetBool("protocol.log.verbose") {
		log.Infof("\n%v\n%v\n%v\n%v\n%v\nunpacked data: %v \n%v",dp.packet.Base.ClientStructName, dp.seen, tPorts, dp.direction, dp.packet.Base.String(),  pv.NcRepresentation.UnpackedData, hex.Dump(dp.packet.Base.Data))
	} else {
		log.Infof("%v %v %v %v", dp.seen, tPorts, dp.direction,  dp.packet.Base.String())
	}
	pv.ConnectionKey = fmt.Sprintf("%v %v", ss.net.String(), ss.transport.String())
	ocs.mu.Lock()
	ocs.structs[dp.packet.Base.OperationCode] = dp.packet.Base.ClientStructName
	ocs.mu.Unlock()
	go sendPacketToUI(pv)
}

func (ss *shineStream) Accept(tcp *layers.TCP, ci gopacket.CaptureInfo, dir reassembly.TCPFlowDirection, nextSeq reassembly.Sequence, start *bool, ac reassembly.AssemblerContext) bool {
	// todo: save it to pcap file
	return true
}

func (ssf *shineStreamFactory) New(net, transport gopacket.Flow, tcp *layers.TCP, ac reassembly.AssemblerContext) reassembly.Stream {
	ctx, cancel := context.WithCancel(ssf.shineContext)

	xorKey := make(chan uint16)
	xorKeyFound := make(chan bool)

	s := &shineStream{
		flowID:    uuid.New().String(),
		net:       net,
		transport: transport,
		xorKey:    xorKey,
		cancel:    cancel,
		isServer:  false,
	}

	srcPort, _ := strconv.Atoi(transport.Src().String())
	if srcPort >= 9000 && srcPort <= 9600 {
		// server - client
		s.isServer = true
	}
	client := make(chan shineSegment, 512)
	server := make(chan shineSegment, 512)
	packets := make(chan decodedPacket, 512)
	s.client = client
	s.server = server
	s.packets = packets

	go s.decodeServerPackets(ctx, server, xorKeyFound, xorKey)
	go s.decodeClientPackets(ctx, client, xorKeyFound, xorKey)
	go s.handleDecodedPackets(ctx, packets)

	log.Infof("new stream from => [ %v ] [ %v ]", net, transport)
	return s
}

func (ss *shineStream) ReassembledSG(sg reassembly.ScatterGather, ac reassembly.AssemblerContext) {
	dir, _, _, _ := sg.Info()
	length, _ := sg.Lengths()
	seg := shineSegment{
		data: sg.Fetch(length),
		seen: ac.GetCaptureInfo().Timestamp,
	}
	//log.Info(dir, ss.net.String())
	ss.mu.Lock()
	if dir == reassembly.TCPDirClientToServer && !ss.isServer {
		seg.direction = "outbound"
		ss.client <- seg
	} else {
		seg.direction = "inbound"
		ss.server <- seg
	}
	ss.mu.Unlock()
}

func (ss *shineStream) ReassemblyComplete(ac reassembly.AssemblerContext) bool {
	log.Warningf("reassembly complete for stream [ %v - %v]", ss.net.String(), ss.transport.String()) // ip of the stream, port of the stream
	ss.cancel()
	//
	//go uiCompletedFlow(completedFlow{
	//	FlowCompleted: true,
	//	FlowID:        ss.flowID,
	//})
	return false
}

// handle stream data flowing from the client
func (ss *shineStream) decodeClientPackets(ctx context.Context, segments <-chan shineSegment, xorKeyFound <-chan bool, xorKey <-chan uint16) {
	var (
		d          []byte
		offset     int
		xorOffset  uint16
		hasXorKey  bool
		shouldQuit bool
	)
	offset = 0
	logActivated := viper.GetBool("protocol.log.client")

loop:
	for {
		select {
		case <-ctx.Done():
			log.Warningf("[%v %v] decodeClientPackets(): context was canceled", ss.net, ss.transport)
			shouldQuit = true
			return
		case <-xorKeyFound:
			log.Info("xor key found, waiting for it in a select")
			for {
				select {
				// retrial mechanism so it doesn't end up in infinite loop
				case xorOffset = <-xorKey:
					hasXorKey = true
					continue loop
				}
			}
		case segment := <-segments:
			d = append(d, segment.data...)

			if offset >= len(d) {
				log.Warningf("not enough data, next offset is %v ", offset)
				break
			}

			for offset < len(d) {
				if !hasXorKey {
					break
				}
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

				if logActivated {
					ss.packets <- decodedPacket{
						seen:   segment.seen,
						packet: pc,
						direction: segment.direction,
					}
				}
				offset += skipBytes + pLen
			}
			if shouldQuit {
				return
			}
		}
	}
}

// handle stream data flowing from the server
func (ss *shineStream) decodeServerPackets(ctx context.Context, segments <-chan shineSegment, xorKeyFound chan<- bool, xorKey chan<- uint16) {
	var (
		d              []byte
		offset         int
		xorOffsetFound bool
		shouldQuit     bool
	)
	xorOffsetFound = false
	offset = 0

	logActivated := viper.GetBool("protocol.log.server")
	for {
		select {
		case <-ctx.Done():
			log.Warningf("[%v %v] decodeServerPackets(): context was canceled", ss.net, ss.transport)
			shouldQuit = true
			return
		case segment := <-segments:
			d = append(d, segment.data...)
			if offset >= len(d) {
				log.Warningf("not enough data, next offset is %v ", offset)
				break
			}
			for offset < len(d) {
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
					log.Info("xor offset not found")
					if pc.Base.OperationCode == 2055 {
						var xorOffset uint16
						buf := bytes.NewBuffer(pc.Base.Data)
						if err := binary.Read(buf, binary.LittleEndian, &xorOffset); err != nil {
							log.Error(err)
							return
						}
						xorOffsetFound = true
						xorKeyFound <- true
						xorKey <- xorOffset
					}
				}

				if logActivated {
					ss.packets <- decodedPacket{
						seen:   segment.seen,
						packet: pc,
						direction: segment.direction,
					}
				}
				offset += skipBytes + pLen
			}
			if shouldQuit {
				return
			}
		}
	}
}

// Capture packets and decode them
func Capture(cmd *cobra.Command, args []string) {
	//p := profile.Start(profile.CPUProfile, profile.ProfilePath("."),profile.NoShutdownHook)
	//p := profile.Start(profile.MemProfi	le, profile.ProfilePath("."),profile.NoShutdownHook)
	runtime.GOMAXPROCS(runtime.NumCPU())
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	config()

	ocs = &opCodeStructs{
		structs: make(map[uint16]string),
	}

	sf := &shineStreamFactory{
		shineContext: ctx,
	}
	sp := reassembly.NewStreamPool(sf)
	a := reassembly.NewAssembler(sp)

	go capturePackets(ctx, a)
	go startUI(ctx)

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM) // subscribe to system signals
	select {
	case <-c:
		//p.Stop()
		cancel()
		generateOpCodeSwitch()
	}
	<-c
}

type Context struct {
	ci  gopacket.CaptureInfo
}

func (c Context) GetCaptureInfo() gopacket.CaptureInfo {
	return c.ci
}

func capturePackets(ctx context.Context, a *reassembly.Assembler) {
	defer a.FlushAll()

	handle, err := pcap.OpenLive(iface, int32(snaplen), true, pcap.BlockForever)
	if err != nil {
		log.Fatal("error opening pcap handle: ", err)
	}
	if err := handle.SetBPFFilter(filter); err != nil {
		log.Fatal("error setting BPF filter: ", err)
	}
	var parser * gopacket.DecodingLayerParser
	var lb layers.Loopback
	var eth layers.Ethernet
	var ip4 layers.IPv4
	var tcp layers.TCP
	var payload gopacket.Payload

	if viper.GetBool("network.loopback") {
		parser = gopacket.NewDecodingLayerParser(layers.LayerTypeLoopback, &lb, &ip4, &tcp, &payload)

	} else {
		parser = gopacket.NewDecodingLayerParser(layers.LayerTypeEthernet, &eth, &ip4, &tcp, &payload)
	}

	decoded := make([]gopacket.LayerType, 0, 4)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			data, ci, err := handle.ZeroCopyReadPacketData()

			if err != nil {
				log.Errorf("error getting packet: %v", err)
				break
			}
			err = parser.DecodeLayers(data, &decoded)
			if err != nil {
				break
			}
			foundNetLayer := false
			var netFlow gopacket.Flow
			for _, typ := range decoded {
				switch typ {
				case layers.LayerTypeIPv4:
					netFlow = ip4.NetworkFlow()
					foundNetLayer = true
				case layers.LayerTypeTCP:
					if foundNetLayer {
						c := Context{
							ci: ci,
						}
						a.AssembleWithContext(netFlow, &tcp, c)
					} else {
						log.Error("could not find IPv4 or IPv6 layer, ignoring")
					}
					break
				}
			}
		}
	}
}
