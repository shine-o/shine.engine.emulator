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
	"path/filepath"
	"runtime"
	"strconv"
	"sync"
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
	serverSideCapture bool
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
	serverSideCapture = viper.GetBool("network.serverSideCapture")
	snaplen = viper.GetInt("network.snaplen")

	if viper.GetBool("network.portRange.useThis") {
		startPort := viper.GetString("network.portRange.start")
		endPort := viper.GetString("network.portRange.end")
		portRange := fmt.Sprintf("%v-%v", startPort, endPort)
		filter = fmt.Sprintf("tcp and portrange %v", portRange)
		log.Infof("using bpf filter %v", filter)
	} else {
		specificPorts := viper.GetIntSlice("network.specificPorts.ports")
		for i, p := range specificPorts {
			if i == 0 {
				filter = fmt.Sprintf("tcp port %v", p)
			} else {
				filter += fmt.Sprintf(" or port %v", p)
			}
		}
		log.Infof("using bpf filter %v", filter)
	}

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
	length, _ := sg.Lengths()
	if length == 0 {
		return
	}

	dir, _, _, _ := sg.Info()

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
	return false
}

// handle stream data flowing from the client
func (ss *shineStream) decodeClientPackets(ctx context.Context, segments <-chan shineSegment, xorKeyFound <-chan bool, xorKey <-chan uint16) {
	var (
		data          []byte
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
			data = append(data, segment.data...)

			if offset >= len(data) {
				log.Warningf("not enough data, next offset is %v ", offset)
				break
			}

			for offset < len(data) {
				if !serverSideCapture {
					if !hasXorKey {
						break
					}
				}

				var skipBytes int
				var pLen uint16

				pLen, skipBytes = networking.PacketBoundary(offset, data)

				nextOffset := offset + skipBytes + int(pLen)

				if nextOffset > len(data) {
					log.Warningf("not enough data, next offset is %v ", nextOffset)
					break
				}

				if pLen == uint16(65535) {
					log.Errorf("bad length value %v", pLen)
					return
				}

				packetData := make([]byte, pLen)

				copy(packetData, data[offset+skipBytes:nextOffset])

				if !serverSideCapture {
					networking.XorCipher(packetData, &xorOffset)
				}

				p, _ := networking.DecodePacket(packetData)

				if logActivated {
					ss.packets <- decodedPacket{
						seen:   segment.seen,
						packet: p,
						direction: segment.direction,
					}
				}
				offset += skipBytes + int(pLen)
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
		data              []byte
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
			data = append(data, segment.data...)
			if offset >= len(data) {
				log.Warningf("not enough data, next offset is %v ", offset)
				break
			}

			for offset < len(data) {
				var skipBytes int
				var pLen uint16

				pLen, skipBytes = networking.PacketBoundary(offset, data)


				nextOffset := offset + skipBytes + int(pLen)

				if nextOffset > len(data) {
					log.Warningf("not enough data for stream %v, next offset is %v ", ss.transport, nextOffset)
					break
				}

				if pLen > uint16(32767) {
					log.Errorf("bad length value %v", pLen)
					return
				}

				packetData := make([]byte, pLen)

				copy(packetData, data[offset+skipBytes:nextOffset])

				pc, _ := networking.DecodePacket(packetData)

				if !serverSideCapture {
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
				}

				if logActivated {
					ss.packets <- decodedPacket{
						seen:   segment.seen,
						packet: pc,
						direction: segment.direction,
					}
				}
				offset += skipBytes + int(pLen)
			}

			if shouldQuit {
				return
			}
		}
	}
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

	go startUI(ctx)
	capturePackets(a)
	//go capturePackets(ctx, a)
	//c := make(chan os.Signal, 2)
	//signal.Notify(c, os.Interrupt, syscall.SIGTERM) // subscribe to system signals
	//select {
	//case <-c:
	//	//p.Stop()
	//	cancel()
	//	generateOpCodeSwitch()
	//}
	//<-c
}

type Context struct {
	ci  gopacket.CaptureInfo
}

func (c Context) GetCaptureInfo() gopacket.CaptureInfo {
	return c.ci
}

func capturePackets(a *reassembly.Assembler) {
	defer a.FlushAll()

	handle, err := pcap.OpenLive(iface, int32(snaplen), true, pcap.BlockForever)
	if err != nil {
		log.Fatal("error opening pcap handle: ", err)
	}
	if err := handle.SetBPFFilter(filter); err != nil {
		log.Fatal("error setting BPF filter: ", err)
	}

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	for packet := range packetSource.Packets() {
		if tcp, ok := packet.TransportLayer().(*layers.TCP); ok {
			c := Context{
				ci: packet.Metadata().CaptureInfo,
			}
			a.AssembleWithContext(packet.NetworkLayer().NetworkFlow(), tcp, c)
		}
	}

//
//	var parser * gopacket.DecodingLayerParser
//	var lb layers.Loopback
//	var eth layers.Ethernet
//	var ip4 layers.IPv4
//	var tcp layers.TCP
//	var payload gopacket.Payload
//
//	if viper.GetBool("network.loopback") {
//		parser = gopacket.NewDecodingLayerParser(layers.LayerTypeLoopback, &lb, &ip4, &tcp, &payload)
//	} else {
//		parser = gopacket.NewDecodingLayerParser(layers.LayerTypeEthernet, &eth, &ip4, &tcp, &payload)
//	}
//
//	decoded := make([]gopacket.LayerType, 4096, 4096)
//
//loop:
//	for {
//		data, ci, err := handle.ZeroCopyReadPacketData()
//
//		if err != nil {
//			log.Errorf("error getting packet: %v	", err)
//			continue
//		}
//		err = parser.DecodeLayers(data, &decoded)
//		if err != nil {
//			continue
//		}
//		foundNetLayer := false
//		var netFlow gopacket.Flow
//		for _, typ := range decoded {
//			switch typ {
//			case layers.LayerTypeIPv4:
//				netFlow = ip4.NetworkFlow()
//				foundNetLayer = true
//			case layers.LayerTypeTCP:
//				if foundNetLayer {
//					c := Context{
//						ci: ci,
//					}
//					a.AssembleWithContext(netFlow, &tcp, c)
//				} else {
//					log.Error("could not find IPv4 or IPv6 layer, ignoring")
//				}
//				continue loop
//			}
//		}
//	}
}
