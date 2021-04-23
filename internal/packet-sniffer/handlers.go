package packet_sniffer

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/segmentio/ksuid"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/crypto"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/networking"
	"github.com/spf13/viper"
	"time"
)

type shineSegment struct {
	data      []byte
	seen      time.Time
	direction string
}

type decodedPacket struct {
	seen      time.Time
	packet    *networking.Command
	direction string
}

// handle stream data flowing from the client
func (ss *shineStream) decodeClientPackets(ctx context.Context, segments <-chan shineSegment, xorKeyFound <-chan bool, xorKey <-chan uint16) {
	var (
		data       []byte
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
					crypto.XorCipher(packetData, xorBytes, &xorOffset, 499)
				}

				p, _ := networking.DecodePacket(packetData)

				if logActivated {
					ss.packets <- decodedPacket{
						seen:      segment.seen,
						packet:    &p,
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
		data           []byte
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
						seen:      segment.seen,
						packet:    &pc,
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

type CapturedPacket struct {
	Command   networking.Command
	Seen      time.Time
	Direction string
}

func (ss *shineStream) handleDecodedPackets(ctx context.Context, decodedPackets <-chan decodedPacket) {
	for {
		select {
		case <-ctx.Done():
			return
		case dp := <-decodedPackets:
			go func(p decodedPacket) {

				if params == nil {
					return
				}
				_, ok := params.WatchCommands[p.packet.Base.OperationCode]
				if ok {
					params.Send <- CapturedPacket{
						Command:   *p.packet,
						Seen:      p.seen,
						Direction: p.direction,
					}
				}
			}(dp)

			if params == nil {
				go ss.logPacket(dp)
			}
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

	//nr, err := ncStructRepresentation(dp.packet.Base.OperationCode, dp.packet.Base.Data)
	//if err == nil {
	//	pv.NcRepresentation = nr
	//	//b, _ := json.Marshal(pv.ncRepresentation)
	//	//log.Info(string(b))
	//} else {
	//	//log.Error(err)
	//}

	var tPorts string

	if dp.direction == "inbound" {
		tPorts = ss.transport.Reverse().String()
	} else {
		tPorts = ss.transport.String()
	}

	if viper.GetBool("protocol.log.verbose") {
		log.Infof("\n%v\n%v\n%v\n%v\n%v\nunpacked data: %v \n%v", dp.packet.Base.OperationCodeName, dp.seen, tPorts, dp.direction, dp.packet.Base.String(), pv.NcRepresentation.UnpackedData, hex.Dump(dp.packet.Base.Data))
	} else {
		log.Infof("%v %v %v %v %v", dp.seen, tPorts, dp.direction, dp.packet.Base.OperationCodeName, dp.packet.Base.String())
	}

	pv.ConnectionKey = fmt.Sprintf("%v %v", ss.net.String(), ss.transport.String())
	ocs.mu.Lock()
	ocs.structs[dp.packet.Base.OperationCode] = fmt.Sprint(dp.packet.Base.OperationCodeName)
	ocs.mu.Unlock()

	persistMovement(dp)
	persistPacketData(dp)
	if viper.GetBool("websocket.active") {
		sendPacketToUI(pv)
	}

	// log packet

}
