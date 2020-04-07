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
	"github.com/segmentio/ksuid"
	"github.com/shine-o/shine.engine.networking"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"os/signal"
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

type shineStreamFactory struct {
	shineContext   context.Context
	localAddresses []pcap.InterfaceAddress
}

type decodedPacket struct {
	seen   time.Time
	packet networking.Command
}

type shineStream struct {
	flowID         string
	net, transport gopacket.Flow
	direction      string
	segments       chan<- shineSegment
	packets        chan<- decodedPacket
	xorKey         chan<- uint16 // only used by decodeClientPackets()
	cancel         context.CancelFunc
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
			packetID, err := ksuid.NewRandomWithTime(dp.seen)
			if err != nil {
				log.Error(err)
				break
			}
			pv := PacketView{
				PacketID :	 	packetID.String(),
				TimeStamp:     dp.seen.String(),
				IpEndpoints:   ss.net.String(),
				PortEndpoints: ss.transport.String(),
				Direction:     ss.direction,
				PacketData:    dp.packet.Base.JSON(),
			}

			nr, err := ncStructRepresentation(dp.packet.Base.OperationCode, dp.packet.Base.Data)
			if err == nil {
				pv.NcRepresentation = nr
				//b, _ := json.Marshal(pv.NcRepresentation)
				//log.Info(string(b))
			}

			data, err := json.Marshal(&pv)

			if err != nil {
				log.Error(err)
				return
			}

			var connectionKey string
			if ss.direction == "inbound" {
				connectionKey = fmt.Sprintf("%v %v", ss.net.String(), ss.transport.String())
			} else {
				connectionKey = fmt.Sprintf("%v %v", ss.net.Reverse(), ss.transport.Reverse())
			}

			if viper.GetBool("protocol.log.verbose") {
				log.Infof("%v %v %v %v", ss.direction, ss.transport, dp.seen, string(data))
			} else {
				log.Infof("[ %v ] %v %v %v", connectionKey, ss.direction, dp.seen, dp.packet.Base.String())
			}
			pv.ConnectionKey = connectionKey
			ocs.mu.Lock()
			ocs.structs[dp.packet.Base.OperationCode] = dp.packet.Base.ClientStructName
			ocs.mu.Unlock()

			go sendPacketToUI(pv)
		}
	}
}

func isLocalAddress(ip string, addresses []pcap.InterfaceAddress) bool {
	for _, a := range addresses {
		if a.IP.String() == ip {
			return true
		}
	}
	return false
}

func (ssf *shineStreamFactory) New(net, transport gopacket.Flow) tcpassembly.Stream {

	var direction string

	ctx, cancel := context.WithCancel(ssf.shineContext)

	srcIP := net.Src().String()
	//dstIP := net.Dst().String()

	if isLocalAddress(srcIP, ssf.localAddresses) {
		direction = "outbound"
	} else {
		direction = "inbound"
	}

	srcPort, _ := strconv.Atoi(transport.Src().String())

	segments := make(chan shineSegment, 512)
	packets := make(chan decodedPacket, 100)

	xorKey := make(chan uint16)

	s := &shineStream{
		flowID:    uuid.New().String(),
		net:       net,
		transport: transport,
		direction: direction,
		segments:  segments,
		packets:   packets,
		xorKey:    xorKey,
		cancel:    cancel,
	}
	go s.handleDecodedPackets(ctx, packets)
	key := fmt.Sprintf("%v:%v", srcIP, srcPort)

	// 9000-9900 is known to be the server
	if srcPort >= 9000 && srcPort <= 9600 {
		// server - client

		if isLocalAddress(srcIP, ssf.localAddresses) {
			direction = ""
		}

		log.Infof("new stream from => [ %v ] [ %v ]", net, transport)
		log.Infof("reversed new stream from => [ %v ] [ %v ]", net.Reverse().String(), transport.Reverse())

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
		log.Infof("new stream from => [ %v ] [ %v ]", net, transport)
		log.Infof("reversed new stream from => [ %v ] [ %v ]", net.Reverse().String(), transport.Reverse())

		ss, ok := ssf.shineContext.Value(activeShineStreams).(*shineStreams)
		ss.mu.Lock()
		if !ok {
			log.Fatalf("unexpected struct type: %v", reflect.TypeOf(ss).String())
		}
		ss.fromClient[key] = s
		ss.mu.Unlock()
		go s.decodeClientPackets(ctx, segments, xorKey)
	}
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
			log.Warningf("[%v %v] decodeClientPackets(): context was canceled", ss.net, ss.transport)
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

				if logActivated {
					ss.packets <- decodedPacket{
						seen:   segment.seen,
						packet: pc,
					}
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
			log.Warningf("[%v %v] decodeServerPackets(): context was canceled", ss.net, ss.transport)
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

				if logActivated {
					ss.packets <- decodedPacket{
						seen:   segment.seen,
						packet: pc,
					}
					//go ss.handlePacket(ctx, segment.seen, pc)
				}
				offset += skipBytes + pLen
			}
		}
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

	var currentInterface pcap.Interface
	interfaces, _ := pcap.FindAllDevs()
	for _, i := range interfaces {
		log.Info(iface)
		if i.Name == iface {
			currentInterface = i
			break
		}
	}

	sf.localAddresses = currentInterface.Addresses

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
		case <-ctx.Done():
			return
		case packet := <-packetSource.Packets():
			if tcp, ok := packet.TransportLayer().(*layers.TCP); ok {

				a.AssembleWithTimestamp(packet.NetworkLayer().NetworkFlow(), tcp, packet.Metadata().Timestamp)
			}
		}
	}
}
