package networking

import (
	"bufio"
	"context"
	"sync"
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

// HandleWarden utility struct for triggering functions implemented by the calling shine service
type HandleWarden struct {
	handlers map[uint16]func(ctx context.Context, command *Command)
	mu       sync.Mutex
}

// NewHandlerWarden handlers are callbacks to be called when an operationCode is detected in a packet.
func NewHandlerWarden(ch *CommandHandlers) *HandleWarden {
	hw := &HandleWarden{
		handlers: make(map[uint16]func(ctx context.Context, command *Command)),
	}
	hw.handlers[2055] = miscSeedAck
	for k, v := range *ch {
		hw.handlers[k] = v
	}
	return hw
}


func protocolCommandWorker(ctx context.Context, hw *HandleWarden, pc <- chan *Command) {
	for {
		select{
		case <- ctx.Done():
			return
		case c := <- pc:
			if callback, ok := hw.handlers[c.Base.OperationCode]; ok {
				go callback(ctx, c)
			} else {
				log.Errorf("non existent operation code from the client %v", c.Base.OperationCode)
			}
		}
	}
}

// Read packet data from segments
func handleInboundSegments(ctx context.Context, segment <-chan []byte, hw *HandleWarden) {
	var (
		data      []byte
		offset    int
		xorOffset uint16
	)

	ctx = context.WithValue(ctx, XorOffset, &xorOffset)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	pc := make(chan * Command, 4096)
	// as many as needed can be spawned
	go protocolCommandWorker(ctx, hw, pc)
	go protocolCommandWorker(ctx, hw, pc)
	go protocolCommandWorker(ctx, hw, pc)

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
				packetData := make([]byte, pLen)

			    copy(packetData, data[offset+skipBytes:nextOffset])

				XorCipher(packetData, &xorOffset)
				c, _ := DecodePacket(packetData)

				log.Infof("[inbound] metadata %v", c.Base.String())

				pc <- &c

				offset += skipBytes + int(pLen)
			}
		}
	}
}

func handleOutboundSegments(ctx context.Context, w *bufio.Writer, segment <-chan []byte) {
	for {
		select {
		case <-ctx.Done():
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

// match operation code with handler if it exists
func (hw *HandleWarden) handleCommand(ctx context.Context, command *Command) {
	select {
	case <-ctx.Done():
		return
	default:
		hw.mu.Lock()
		if callback, ok := hw.handlers[command.Base.OperationCode]; ok {
			callback(ctx, command)
		} else {
			log.Errorf("non existent operation code from the client %v", command.Base.OperationCode)
		}
		hw.mu.Unlock()
	}
}
