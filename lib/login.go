package lib

import (
	"bytes"
	"context"
	"encoding/binary"
	//protocol "shine.engine.packet-protocol"
	protocol "github.com/shine-o/shine.engine.protocol"
	"sync"
)

// RE client struct:
// struct PROTO_NC_MISC_SEED_ACK
// {
//	unsigned __int16 seed;
// };
//
// xorKey offset used by client to encrypt data
// same offset is used on the server side to decrypt data sent by the client
type ncMiscSeedAck struct {
	seed uint16
}

// RE client struct:
// struct PROTO_NC_USER_CLIENT_VERSION_CHECK_REQ
// {
//  char sVersionKey[64];
// };
type ncUserClientVersionCheckReq struct {
	VersionKey [64]byte
}

// RE client struct:
// struct __cppobj PROTO_NC_USER_CLIENT_WRONGVERSION_CHECK_ACK
// {
// };
type ncUserClientRightversionCheckAck struct{}

// RE client struct:
// struct PROTO_NC_USER_US_LOGIN_REQ
// {
//  char sUserName[260];
//  char sPassword[36];
//  Name5 spawnapps;
// };
type ncUserUsLoginReq struct{}

// RE client struct:
// struct PROTO_NC_USER_XTRAP_REQ
// {
//  char XTrapClientKeyLength;
//  char XTrapClientKey[];
// };
type ncUserXtrapReq struct{}

// RE client struct:
// struct PROTO_NC_USER_XTRAP_ACK
// {
//  char bSuccess;
// };
type ncUserXtrapAck struct{}

// RE client struct:
// struct __unaligned __declspec(align(1)) PROTO_NC_USER_LOGIN_ACK
// {
//  char numofworld;
//  PROTO_NC_USER_LOGIN_ACK::WorldInfo worldinfo[];
// };
type ncUserLoginAck struct {
	NumOfWorld byte
	Worlds     [1]WorldInfo
}

// RE client struct:
// struct __cppobj PROTO_NC_USER_WORLD_STATUS_REQ
// {
// };
type ncUserWorldStatusReq struct{}

// OPERATION CODE ONLY 3100
type ncUserWorldStatusAck struct{}

// RE client struct:
//struct PROTO_NC_USER_WORLDSELECT_REQ
//{
//char worldno;
//};
type ncUserWorldSelectReq struct {
	WorldNo byte
}

// RE client struct:
// struct __unaligned __declspec(align(1)) PROTO_NC_USER_WORLDSELECT_ACK
// {
//	char worldstatus;
//	Name4 ip;
//	unsigned __int16 port;
//	unsigned __int16 validate_new[32];
//};
type ncUserWorldSelectAck struct {
	WorldStatus byte
	Ip          ComplexName
	Port        uint16
	ValidateNew [32]uint16
}

type handleWarden struct {
	handlers map[uint16]func(ctx context.Context, command *protocol.Command)
	mu       sync.Mutex
}

var (
	hw *handleWarden
)

func loginHandlers() {
	hw = &handleWarden{
		handlers: make(map[uint16]func(ctx context.Context, command *protocol.Command)),
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
	hw.handlers[3083] = userWorldSelectReq
	hw.handlers[3084] = userWorldSelectAck
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
				var skipBytes int
				var pLen int
				var pType string
				var pd []byte

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

func miscSeedAck(ctx context.Context, pc *protocol.Command) {
	select {
	case <-ctx.Done():
		return
	default:
		cwv := ctx.Value("connWriter")
		xov := ctx.Value("xorOffset")

		cw := cwv.(*clientWriter)
		xo := xov.(*uint16)

		xorOffset := protocol.RandomXorKey()

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

		pc.Base.SetData(buf.Bytes())

		cw.mu.Lock()
		if _, err := cw.w.Write(pc.Base.RawData()); err != nil {
			log.Error(err)
		} else {
			if err = cw.w.Flush(); err != nil {
				log.Error(err)
			}
		}
		cw.mu.Unlock()
		logOutboundPacket(pc)
	}
}

func userClientVersionCheckReq(ctx context.Context, pc *protocol.Command) {
	select {
	case <-ctx.Done():
		return
	default:
		buf := bytes.NewBuffer(pc.Base.Data())

		nc := &ncUserClientVersionCheckReq{}

		if err := binary.Read(buf, binary.LittleEndian, nc); err != nil {
			log.Error(err)
		}
		// [...] future client version checking logic
		go userClientVersionCheckAck(ctx, &protocol.Command{})
		logOutboundPacket(pc)
	}
}

func userClientVersionCheckAck(ctx context.Context, pc *protocol.Command) {
	select {
	case <-ctx.Done():
		return
	default:
		cwv := ctx.Value("connWriter")
		cw := cwv.(*clientWriter)
		base := protocol.CommandBase{}
		base.SetOperationCode(3175)
		pc.Base = base
		cw.mu.Lock()
		if _, err := cw.w.Write(pc.Base.RawData()); err != nil {
			log.Error(err)
		} else {
			if err = cw.w.Flush(); err != nil {
				log.Error(err)
			}
		}
		cw.mu.Unlock()
		logOutboundPacket(pc)
	}
}

func userUsLoginReq(ctx context.Context, pc *protocol.Command) {
	select {
	case <-ctx.Done():
		return
	default:
		buf := bytes.NewBuffer(pc.Base.Data())

		nc := &ncUserUsLoginReq{}

		if err := binary.Read(buf, binary.LittleEndian, nc); err != nil {
			log.Error(err)
		}
		// [...] future user login checking logic
		go userLoginAck(ctx, &protocol.Command{})
		logOutboundPacket(pc)
	}
}

func userLoginAck(ctx context.Context, pc *protocol.Command) {
	select {
	case <-ctx.Done():
		return
	default:
		cwv := ctx.Value("connWriter")
		cw := cwv.(*clientWriter)
		base := protocol.CommandBase{}
		base.SetOperationCode(3082)

		pc.Base = base

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
		pc.Base.SetData(buf.Bytes())
		cw.mu.Lock()
		if _, err := cw.w.Write(pc.Base.RawData()); err != nil {
			log.Error(err)
		} else {
			if err = cw.w.Flush(); err != nil {
				log.Error(err)
			}
		}
		cw.mu.Unlock()
		logOutboundPacket(pc)
	}
}

func userXtrapReq(ctx context.Context, pc *protocol.Command) {}

func userXtrapAck(ctx context.Context, pc *protocol.Command) {}

func userWorldStatusReq(ctx context.Context, pc *protocol.Command) {
	select {
	case <-ctx.Done():
		return
	default:
		go userWorldStatusAck(ctx, &protocol.Command{})
	}
}

func userWorldStatusAck(ctx context.Context, pc *protocol.Command) {
	select {
	case <-ctx.Done():
		return
	default:
		cwv := ctx.Value("connWriter")
		cw := cwv.(*clientWriter)

		base := protocol.CommandBase{}
		base.SetOperationCode(3100)
		pc.Base = base

		cw.mu.Lock()
		if _, err := cw.w.Write(pc.Base.RawData()); err != nil {
			log.Error(err)
		} else {
			if err = cw.w.Flush(); err != nil {
				log.Error(err)
			}
		}
		cw.mu.Unlock()
		logOutboundPacket(pc)
	}
}

func userWorldSelectAck(ctx context.Context, pc *protocol.Command) {}

func userWorldSelectReq(ctx context.Context, pc *protocol.Command) {}
