package networking

import (
	"context"
	"sync"
)

type HandleWarden struct {
	handlers map[uint16]func(ctx context.Context, command *Command)
	mu       sync.Mutex
}

// handlers are callbacks to be called when an operationcode is detected in a packet.
func NewHandlerWarden(customHandlers map[uint16]func(ctx context.Context, command *Command)) *HandleWarden {
	hw := &HandleWarden{
		handlers: make(map[uint16]func(ctx context.Context, command *Command)),
	}
	hw.handlers[2055] = miscSeedAck

	for k, v := range customHandlers {
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

	ctx = context.WithValue(ctx, "xorOffset", &xorOffset)

	hw.mu.Lock()
	sendXorOffset := hw.handlers[2055]
	go sendXorOffset(ctx, &Command{
		CommandBase{
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

				log.Infof("Inbound packet %v", pc.Base.String())

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
