package packet_sniffer

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcapgo"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"sync"
)

// Flows utility struct for storing raw packets
type Flows struct {
	pfm map[string][]gopacket.Packet
	m   sync.Mutex
}

// write to disk pcapng files for each flow in the map
func (pf *Flows) persist() {
	pf.m.Lock()
	for k, v := range pf.pfm {
		pathName, err := filepath.Abs(fmt.Sprintf("%v%v", "output/", k))
		if err != nil {
			log.Fatal(err)
		}
		f, err := os.OpenFile(pathName, os.O_WRONLY|os.O_CREATE, 0666)

		r, err := pcapgo.NewNgWriter(f, layers.LinkTypeEthernet)

		for _, p := range v {
			err = r.WritePacket(p.Metadata().CaptureInfo, p.Data())
		}

		r.Flush()
		f.Close()
	}
	pf.m.Unlock()
}

func Decode(cmd *cobra.Command, args []string) {

}