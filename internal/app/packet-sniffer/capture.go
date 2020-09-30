package packet_sniffer

import (
	"context"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/reassembly"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	ExtendedCapture(nil)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetConfigFile("configs/packet-sniffer-client.yml")


	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	required := []string{
		"network.interface",
	}

	for _, v := range required {
		if !viper.IsSet(v) {
			panic(fmt.Sprintf("required config parameter is missing: %v", v))
		}
	}

	viper.SetDefault("network.portRange.start", 9000)

	viper.SetDefault("network.portRange.end", 9600)

	viper.SetDefault("network.interface", 65536)

	viper.SetDefault("protocol.xorKey", "0759694a941194858c8805cba09ecd583a365b1a6a16febddf9402f82196c8e99ef7bfbdcfcdb27a009f4022fc11f90c2e12fba7740a7d78401e2ca02d06cba8b97eefde49ea4e13161680f43dc29ad486d7942417f4d665bd3fdbe4e10f50f6ec7a9a0c273d2466d322689c9a520be0f9a50b25da80490dfd3e77d156a8b7f40f9be80f5247f56f832022db0f0bb14385c1cba40b0219dff08becdb6c6d66ad45be89147e2f8910b89360d860def6fe6e9bca06c1759533cfc0b2e0cca5ce12f6e5b5b426c5b2184f2a5d261b654df545c98414dc7c124b189cc724e73c64ffd63a2cee8c8149396cb7dcbd94e232f7dd0afc020164ec4c940ab156f5c9a934de0f3827bc81300f7b3825fee83e29ba5543bf6b9f1f8a4952187f8af888245c4fe1a830878e501f2fd10cb4fd0abcdc1285e252ee4a5838abffc63db960640ab450d54089179ad585cfec0d7e817fe3c3040122ec27ccfa3e21a654c8de00b6df279ff625340785bfa7a5a5e0830c3d5d2040af60a36456f305c41c7d3798c3e85a6e5885a49a6b6af4a37b619b09401e604b32d951a4fef95d4e4afb4ad47c330233d59dce5baa5a7cd8f805fa1f2b8c725750ae6c1989ca01fcfc299b61126863654626c45b50aa2bbeef9a790223752c2013fdd95a7623f10bb5b859f99f7ae606e9a53ab450bf165898b39a6e36ee8deb")

	viper.SetDefault("protocol.xorLimit", 350)

	viper.SetDefault("protocol.log.client", true)

	viper.SetDefault("protocol.log.server", true)
}

type Params struct {
	// if command operation code is in this list, send command on the channel
	WatchCommands map[uint16]interface{}
	Send chan <- CapturedPacket
}

var params * Params

func ExtendedCapture(p * Params)  {
	initConfig()
	params = p

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

	if viper.GetBool("webstocket.active") {
		go startUI(ctx)
	}
	go capturePackets(ctx, a)

	em.Entities = make(map[uint16][]Movement)

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM) // subscribe to system signals
	for {
		select {
		case <-c:
			cancel()
			//generateOpCodeSwitch()
			exportEntitiesMovements()
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
