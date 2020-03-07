package lib

import (
	"bufio"
	"bytes"
	"context"
	"encoding/binary"
	"github.com/spf13/viper"
	"sync"
)

type handleWarden struct {
	handlers map[int]func(ctx context.Context, command *ProtocolCommand)
	mu       sync.Mutex
}

var (
	hw *handleWarden
)

func initHandlers() {
	hw = &handleWarden{
		handlers: make(map[int]func(ctx context.Context, command *ProtocolCommand)),
	}
	hw.handlers[2055] = miscSeedAck
	hw.handlers[3173] = userClientVersionCheckReq
	hw.handlers[3175] = userClientVersionCheckAck
}

func miscSeedAck(ctx context.Context, pc *ProtocolCommand) {
	cwv := ctx.Value("connWriter")
	xov := ctx.Value("xorOffset")

	cw := cwv.(*bufio.Writer)
	xo := xov.(*uint16)

	xorOffset := randomXorKey(viper.GetInt("crypt.xorLimit"))

	*xo = xorOffset

	log.Infof("XorKey: %v", xorOffset)

	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.LittleEndian, xorOffset); err != nil {
		log.Fatal(err)
		return
	}

	pc.pcb.packetType = "small"
	pc.pcb.data = buf.Bytes()

	if _, err := cw.Write(pc.pcb.RawData()); err != nil {
		log.Error(err)
	} else {
		if err = cw.Flush(); err != nil {
			log.Error(err)
		}
	}
}

func userClientVersionCheckReq(ctx context.Context, pc *ProtocolCommand) {}

func userClientVersionCheckAck(ctx context.Context, pc *ProtocolCommand) {}
