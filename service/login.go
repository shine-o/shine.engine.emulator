package service

import (
	"context"
	//protocol "shine.engine.packet-protocol"
	protocol "github.com/shine-o/shine.engine.protocol"
	"sync"
)

type handleWarden struct {
	handlers map[uint16]func(ctx context.Context, command *protocol.Command)
	mu       sync.Mutex
}

var (
	hw *handleWarden
)

func commandHandlers() {
	hw = &handleWarden{
		handlers: make(map[uint16]func(ctx context.Context, command *protocol.Command)),
	}

	hw.handlers[2055] = miscSeedAck

	hw.handlers[3173] = userClientVersionCheckReq
	hw.handlers[3175] = userClientVersionCheckAck
	hw.handlers[3162] = userUsLoginReq
	hw.handlers[3081] = userLoginFailAck
	hw.handlers[3082] = userLoginAck
	hw.handlers[3076] = userXtrapReq
	hw.handlers[3077] = userXtrapAck
	hw.handlers[3099] = userWorldStatusReq
	hw.handlers[3100] = userWorldStatusAck
	hw.handlers[3083] = userWorldSelectReq
	hw.handlers[3084] = userWorldSelectAck
	hw.handlers[3096] = userNormalLogoutCmd
}

// Read packet data from segments
func handleLoginSegments(ctx context.Context, segment <-chan []byte) {
	var (
		data      []byte
		offset    int
		xorOffset uint16
	)

	ctx = context.WithValue(ctx, "xorOffset", &xorOffset)

	hw.mu.Lock()
	sendXorOffset := hw.handlers[2055]
	base := protocol.CommandBase{}
	base.SetOperationCode(2055)
	go sendXorOffset(ctx, &protocol.Command{
		Base: base,
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

				pLen, pType = protocol.PacketBoundary(offset, data)

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
				protocol.XorCipher(pd, &xorOffset)

				pc, _ := protocol.DecodePacket(pType, pLen, pd)

				log.Infof("Inbound packet %v", pc.Base.String())

				go handlePacket(ctx, &pc)

				offset += skipBytes + pLen
			}
		}
	}
}

// match operation code with handler if it exists
func handlePacket(ctx context.Context, command *protocol.Command) {
	select {
	case <-ctx.Done():
		return
	default:
		hw.mu.Lock()
		if callback, ok := hw.handlers[command.Base.OperationCode()]; ok {
			callback(ctx, command)
		} else {
			log.Errorf("non existent operation code from the client %v", command.Base.OperationCode())
		}
		hw.mu.Unlock()
	}
}
