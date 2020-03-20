package service

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcapgo"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

// Flows utility struct for storing raw packets
type Flows struct {
	pfm map[string][]gopacket.Packet
	m   sync.Mutex
}

// asses the package flow, either Client-Service or Service-Client
func (pf *Flows) add(p gopacket.Packet) {

	src, _ := strconv.Atoi(p.TransportLayer().TransportFlow().Src().String())
	dst, _ := strconv.Atoi(p.TransportLayer().TransportFlow().Dst().String())

	var fkey string
	shine.mu.Lock()
	if src >= 9000 && src <= 9600 {
		// server - client
		fkey = fmt.Sprintf("%v-Client.pcapng", shine.knownServices[src].name)
	} else {
		fkey = fmt.Sprintf("Client-%v.pcapng", shine.knownServices[dst].name)
	}
	shine.mu.Unlock()
	pf.m.Lock()
	pf.pfm[fkey] = append(pf.pfm[fkey], p)
	pf.m.Unlock()
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