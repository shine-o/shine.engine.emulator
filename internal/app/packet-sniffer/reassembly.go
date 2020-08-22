package packet_sniffer

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/reassembly"
	"github.com/google/logger"
	"github.com/google/uuid"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/networking"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

func init() {
	dir, err := filepath.Abs("./output/")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.Mkdir(dir, 0700)

		if err != nil {
			log.Error(err)
		}
	} else {
		err = os.RemoveAll(dir)
		if err != nil {
			log.Error(err)
		}
	}

	lf, err := os.OpenFile("./output/streams.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0660)
	if err != nil {
		logger.Fatalf("Failed to open log file: %v", err)
	}
	log = logger.Init("SnifferLogger", true, false, lf)
	log.Info("sniffer logger init()")
}

type shineStreamFactory struct {
	shineContext   context.Context
	localAddresses []pcap.InterfaceAddress
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
	mu             sync.Mutex
}

var (
	iface             string
	snaplen           int
	filter            string
	log               *logger.Logger
	serverSideCapture bool
)

func config() {
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
