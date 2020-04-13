package net

import (
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

// Read packet data from segments
func (hw *HandleWarden) handleSegments(ctx context.Context, segment <-chan []byte) {
	var (
		data      []byte
		offset    int
		xorOffset uint16
	)

	ctx = context.WithValue(ctx, XorOffset, &xorOffset)

	hw.mu.Lock()
	sendXorOffset := hw.handlers[2055]
	go sendXorOffset(ctx, &Command{
		Base: CommandBase{
			OperationCode: 2055,
		},
	})
	hw.mu.Unlock()

	offset = 0
	for {
		select {
		case <-ctx.Done():
			return
		case b := <-segment:
			data = append(data, b...)

			if offset > len(data) {
				break
			}

			if offset != len(data) {
				var (
					skipBytes int
					pLen      int
					pType     string
					pd        []byte
				)

				pLen, pType = PacketBoundary(offset, data)

				if pType == "small" {
					skipBytes = 1
				} else {
					skipBytes = 3
				}

				nextOffset := offset + skipBytes + pLen
				if nextOffset > len(data) {
					break
				}

				pd = append(pd, data[offset+skipBytes:nextOffset]...)
				XorCipher(pd, &xorOffset)

				pc, _ := DecodePacket(pType, pLen, pd)

				log.Infof("[inbound] metadata %v", pc.Base.String())

				go hw.handleCommand(ctx, &pc)

				offset += skipBytes + pLen
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
