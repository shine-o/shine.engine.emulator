package networking

import (
	"context"
)

// ContextKey identifier for values of common use within the Context
type ContextKey int

// ShineParameters is used by the Shine services to give extra parameters to their handlers
type Parameters struct {
	Command *Command
	NetVars *Network
	// anything that the Shine service wants to be sent to all handles, an alternative to using global variables
	ServiceParams interface{}
}

// HandleWarden utility struct for triggering functions implemented by the calling shine service
type ShineHandler map[uint16]func(context.Context, *Parameters)

// Read packet data from segments
func (ss *ShineService) handleInboundSegments(ctx context.Context, n *Network) {
	var (
		data      []byte
		offset    int
		xorOffset uint16
	)

	ctx = context.WithValue(ctx, XorOffset, &xorOffset)
	ctx, cancel := context.WithCancel(ctx)

	for i := 0; i < 10; i++ { // todo: number of routines to be put in config
		go ss.handlerWorker(ctx, n)
	}

	n.Commands.Send <- &Command{
		Base: CommandBase{
			OperationCode: 2055,
		},
	}

	defer cancel()

	offset = 0

	for {
		select {
		case <-ctx.Done():
			return
		case b := <- n.InboundSegments.Recv:
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

				if pLen == uint16(65535) {
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

				log.Infof("[inbound] metadata %v", c.Base.String())

				n.Commands.Send <- &c

				offset = 0
				data = nil
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
		case data := <- n.OutboundSegments.Recv:
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

func (ss ShineService) handlerWorker(ctx context.Context, n *Network) {
	for {
		select {
		case <-ctx.Done():
			return
		case c := <- n.Commands.Recv:
			if callback, ok := ss.ShineHandler[c.Base.OperationCode]; ok {
				go callback(ctx, &Parameters{
					Command:       c,
					NetVars:       n,
					ServiceParams: ss.ExtraParameters,
				})
			} else {
				log.Errorf("non existent handler for operation code  %v", c.Base.OperationCode)
			}
		}
	}
}
