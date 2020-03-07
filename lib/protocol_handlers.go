package lib

import (
	"bytes"
	"context"
	"encoding/binary"
	"github.com/spf13/viper"
	"sync"
)

type handleWarden struct {
	handlers map[uint16]func(ctx context.Context, command *ProtocolCommand)
	mu       sync.Mutex
}

var (
	hw *handleWarden
)

func initHandlers() {
	hw = &handleWarden{
		handlers: make(map[uint16]func(ctx context.Context, command *ProtocolCommand)),
	}
	hw.handlers[2055] = miscSeedAck
	hw.handlers[3173] = userClientVersionCheckReq
	hw.handlers[3175] = userClientVersionCheckAck
	hw.handlers[3162] = userUsLoginReq
	hw.handlers[3082] = userLoginAck
	hw.handlers[3076] = userXtrapReq
	hw.handlers[3077] = userXtrapAck
	hw.handlers[3099] = userWorldStatusReq
	hw.handlers[3100] = userWorldStatusAck
}

func miscSeedAck(ctx context.Context, pc *ProtocolCommand) {
	select {
	case <-ctx.Done():
		return
	default:
		cwv := ctx.Value("connWriter")
		xov := ctx.Value("xorOffset")

		cw := cwv.(*clientWriter)
		xo := xov.(*uint16)

		xorOffset := randomXorKey(viper.GetInt("crypt.xorLimit"))

		*xo = xorOffset

		nc := ncMiscSeedAck{
			seed: xorOffset,
		}

		log.Infof("XorKey: %v", xorOffset)

		buf := new(bytes.Buffer)

		if err := binary.Write(buf, binary.LittleEndian, nc); err != nil {
			log.Fatal(err)
			return
		}

		pc.pcb.data = buf.Bytes()
		cw.mu.Lock()
		if _, err := cw.w.Write(pc.pcb.RawData()); err != nil {
			log.Error(err)
		} else {
			if err = cw.w.Flush(); err != nil {
				log.Error(err)
			}
		}
		cw.mu.Unlock()
		log.Infof("Handling packet %v", pc.pcb.String())
	}
}

func userClientVersionCheckReq(ctx context.Context, pc *ProtocolCommand) {
	select {
	case <-ctx.Done():
		return
	default:
		buf := bytes.NewBuffer(pc.pcb.data)

		nc := &ncUserClientVersionCheckReq{}

		if err := binary.Read(buf, binary.LittleEndian, nc); err != nil {
			log.Error(err)
		}
		// [...] future client version checking logic
		go userClientVersionCheckAck(ctx, &ProtocolCommand{})
		log.Infof("Handling packet %v", pc.pcb.String())
	}
}

func userClientVersionCheckAck(ctx context.Context, pc *ProtocolCommand) {
	select {
	case <-ctx.Done():
		return
	default:
		cwv := ctx.Value("connWriter")
		cw := cwv.(*clientWriter)

		pc.pcb = ProtocolCommandBase{
			operationCode: 3175,
		}

		cw.mu.Lock()
		if _, err := cw.w.Write(pc.pcb.RawData()); err != nil {
			log.Error(err)
		} else {
			if err = cw.w.Flush(); err != nil {
				log.Error(err)
			}
		}
		cw.mu.Unlock()
		log.Infof("Handle packet %v", pc.pcb.String())
	}
}

func userUsLoginReq(ctx context.Context, pc *ProtocolCommand) {
	select {
	case <-ctx.Done():
		return
	default:
		buf := bytes.NewBuffer(pc.pcb.data)

		nc := &ncUserUsLoginReq{}

		if err := binary.Read(buf, binary.LittleEndian, nc); err != nil {
			log.Error(err)
		}
		// [...] future user login checking logic
		go userLoginAck(ctx, &ProtocolCommand{})
		log.Infof("Handling packet %v - %v", pc.pcb.String(), nc)
	}
}

func userLoginAck(ctx context.Context, pc *ProtocolCommand) {
	select {
	case <-ctx.Done():
		return
	default:
		cwv := ctx.Value("connWriter")
		cw := cwv.(*clientWriter)

		pc.pcb = ProtocolCommandBase{
			operationCode: 3082,
		}

		w1 := WorldInfo{
			WorldNumber: 0,
			WorldName:   ComplexName{},
			WorldStatus: 0,
		}

		copy(w1.WorldName.Name[:], "EPITH")
		copy(w1.WorldName.NameCode[:], []uint16{262, 16720, 17735, 76})

		nc := &ncUserLoginAck{
			NumOfWorld: byte(2),
			Worlds: [1]WorldInfo{
				w1,
			},
		}
		buf := new(bytes.Buffer)

		if err := binary.Write(buf, binary.LittleEndian, nc); err != nil {
			log.Fatal(err)
			return
		}
		pc.pcb.data = buf.Bytes()
		cw.mu.Lock()
		if _, err := cw.w.Write(pc.pcb.RawData()); err != nil {
			log.Error(err)
		} else {
			if err = cw.w.Flush(); err != nil {
				log.Error(err)
			}
		}
		cw.mu.Unlock()
		log.Infof("Handling packet %v", pc.pcb.String())
	}
}

func userXtrapReq(ctx context.Context, pc *ProtocolCommand) {}

func userXtrapAck(ctx context.Context, pc *ProtocolCommand) {}

func userWorldStatusReq(ctx context.Context, pc *ProtocolCommand) {
	select {
	case <-ctx.Done():
		return
	default:
		go userWorldStatusAck(ctx, &ProtocolCommand{})
	}
}

func userWorldStatusAck(ctx context.Context, pc *ProtocolCommand) {
	select {
	case <-ctx.Done():
		return
	default:
		cwv := ctx.Value("connWriter")
		cw := cwv.(*clientWriter)

		pc.pcb = ProtocolCommandBase{
			operationCode: 3100,
		}

		cw.mu.Lock()
		if _, err := cw.w.Write(pc.pcb.RawData()); err != nil {
			log.Error(err)
		} else {
			if err = cw.w.Flush(); err != nil {
				log.Error(err)
			}
		}
		cw.mu.Unlock()
		log.Infof("Handling packet %v", pc.pcb.String())

	}
}
