package networking

import (
	"context"
)

// ContextKey identifier for values of common use within the Context
type ContextKey int

// ShineParameters is used by the Shine services to give extra parameters to their handlers
type Parameters struct {
	Command *Command
	*Network
	// anything that the Shine service wants to be sent to all handles, an alternative to using global variables
	ServiceParams interface{}
}

// HandleWarden utility struct for triggering functions implemented by the calling shine service
//type ShinePacketRegistry map[uint16]func(context.Context, *Parameters)
type ShinePacketRegistry map[OperationCode]ShinePacket

type ShinePacket struct {
	Handler  func(context.Context, *Parameters)
	NcStruct interface{}
}

// Read packet data from segments
func (ss *ShineService) handleInboundSegments(ctx context.Context, n *Network) {
	var (
		data            []byte
		offset          int
		xorOffset       uint16
		randomXorOffset = make(chan uint16)
	)

	ctx = context.WithValue(ctx, XorOffset, randomXorOffset)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		for i := 0; i < 10; i++ { // todo: number of routines to be put in config
			go ss.commandWorker(ctx, n)
		}

		n.Commands.Send <- &Command{
			Base: CommandBase{
				OperationCode: 2055,
			},
		}
	}()

	offset = 0

	select {
	case xorOffset = <-randomXorOffset:
		break
	}

	for {
		select {
		// probably not needed as nothing will be received if the connection drops and this routine will get collected
		case <-ctx.Done():
			return
		case b := <-n.InboundSegments.Recv:
			data = append(data, b...)

			if offset >= len(data) {
				break
			}

			for offset < len(data) {
				pLen, skipBytes := PacketBoundary(offset, data)

				nextOffset := offset + skipBytes + int(pLen)

				if nextOffset > len(data) {
					break
				}

				if pLen >= uint16(65535) {
					log.Error("max value reached for packet length")
					n.CloseConnection <- true
					return
				}

				if nextOffset > len(data) {
					log.Errorf("next offset [%v] is bigger than current available data [%v], cannot proceed", nextOffset, len(data))
					n.CloseConnection <- true
					return
				}

				packetData := make([]byte, pLen)

				copy(packetData, data[offset+skipBytes:nextOffset])
				XorCipher(packetData, &xorOffset)
				c, _ := DecodePacket(packetData)

				n.Commands.Send <- &c
				logInboundPackets <- &c

				offset += skipBytes + int(pLen)
			}
		}
	}
}

func (ss *ShineService) handleOutboundSegments(ctx context.Context, n *Network) {
	for {
		select {
		case <-ctx.Done():
			log.Warning("handleOutboundSegments context canceled")
			return
		case data := <-n.OutboundSegments.Recv:
			if _, err := n.Writer.Write(data); err != nil {
				log.Error(err)
			} else {
				if err = n.Writer.Flush(); err != nil {
					log.Error(err)
				}
			}
		}
	}
}

func (ss ShineService) commandWorker(ctx context.Context, n *Network) {
	for {
		select {
		case <-ctx.Done():
			log.Warning("commandWorker context canceled")
			return
		case c := <-n.Commands.Recv:
			if packet, ok := ss.ShinePacketRegistry[OperationCode(c.Base.OperationCode)]; ok {
				go packet.Handler(ctx, &Parameters{
					Command:       c,
					Network:       n,
					ServiceParams: ss.ExtraParameters,
				})
			} else {
				log.Errorf("non existent handler for operation code  %v %v", c.Base.OperationCode, c.Base.OperationCodeName)
			}
		}
	}
}
