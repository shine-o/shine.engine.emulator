package networking

import (
	"bufio"
	"context"
)

// ContextKey identifier for values of common use within the Context
type ContextKey int

const (
	// XorOffset indicates what offset in the xor hex table to use to start decrypting client data
	XorOffset ContextKey = iota
	// ShineSession if used, shine service can access session data within their handler's context
	ShineSession
	// ConnectionWriter is a utility struct which contains the tcp connection object and a mutex
	// it is used to write data to the client from any shine service handler
	ConnectionWriter
)

// ShineParameters is used by the Shine services to give extra parameters to their handlers
type Parameters struct {
	Command *Command
	Extra interface{}
}

// HandleWarden utility struct for triggering functions implemented by the calling shine service
type ShineHandler map[uint16]func(context.Context, *Parameters)

// Read packet data from segments
func (ss * ShineService) handleInboundSegments(ctx context.Context, segment <-chan []byte, closeConnection chan <- bool) {

	var (
		data      []byte
		offset    int
		xorOffset uint16
	)

	ctx = context.WithValue(ctx, XorOffset, &xorOffset)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	pc := make(chan * Command, 4096)

	for i := 0; i < 10; i++ {// todo: number of routines to be put in config
		go ss.handlerWorker(ctx, pc)
	}

	pc <- &Command{
		Base: CommandBase{
			OperationCode: 2055,
		},
	}

	offset = 0

	for {
		select {
		case <-ctx.Done():
			return
		case b := <-segment:
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
					closeConnection <- true
					return
				}

				if nextOffset > len(data) {
					log.Errorf("next offset [%v] is bigger than current available data [%v], cannot proceed", nextOffset, len(data))
					closeConnection <- true
					return
				}

				packetData := make([]byte, pLen)

			    copy(packetData, data[offset+skipBytes:nextOffset])

				XorCipher(packetData, &xorOffset)
				c, _ := DecodePacket(packetData)

				log.Infof("[inbound] metadata %v", c.Base.String())

				pc <- &c

				offset = 0
				data = nil
				//offset += skipBytes + int(pLen)
			}
		}
	}
}

func (ss * ShineService) handleOutboundSegments(ctx context.Context, w *bufio.Writer, segment <-chan []byte) {
	for {
		select {
		case <- ctx.Done():
			log.Warning("handleOutboundSegments context canceled")
			return
		case data := <-segment:
			if _, err := w.Write(data); err != nil {
				log.Error(err)
			} else {
				if err = w.Flush(); err != nil {
					log.Error(err)
				}
			}
		}
	}
}


func (ss ShineService) handlerWorker(ctx context.Context, pc <- chan *Command) {
	for {
		select{
		case <- ctx.Done():
			return
		case c := <- pc:
			if callback, ok := ss.ShineHandler[c.Base.OperationCode]; ok {
				go callback(ctx, &Parameters{
					Command: c,
					Extra: ss.ExtraParameters,
				})
			} else {
				log.Errorf("non existent handler for operation code %v", c.Base.OperationCode)
			}
		}
	}
}