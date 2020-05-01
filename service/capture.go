package service

import (
	"context"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/reassembly"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

type Context struct {
	ci gopacket.CaptureInfo
}

func (c Context) GetCaptureInfo() gopacket.CaptureInfo {
	return c.ci
}

// Capture packets and decode them
func Capture(cmd *cobra.Command, args []string) {
	//p := profile.Start(profile.CPUProfile, profile.ProfilePath("."),profile.NoShutdownHook)
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
	go capturePackets(ctx, a)
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Kill, syscall.SIGTERM) // subscribe to system signals
	for {
		select {
		case <-c:
			cancel()
			//generateOpCodeSwitch()
		}
	}
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

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	for {
		select {
		case <-ctx.Done():
			log.Warningf("capture canceled")
			return
		case packet := <-packetSource.Packets():
			if tcp, ok := packet.TransportLayer().(*layers.TCP); ok {
				c := Context{
					ci: packet.Metadata().CaptureInfo,
				}
				a.AssembleWithContext(packet.NetworkLayer().NetworkFlow(), tcp, c)
			}
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
