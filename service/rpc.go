package service

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	gs "shine.engine.game_structs"
	networking "shine.engine.networking"
	lw "shine.engine.protocol-buffers/login-world"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	lw.UnimplementedWorldServer
}

// There can be many Worlds, each World with its own Zones
// for simplicity we use hardcoded one for now :)
func (s *server) AvailableWorlds(ctx context.Context, in *lw.ClientMetadata) (*lw.WorldsInfo, error) {
	select {
	case <-ctx.Done():
		return &lw.WorldsInfo{Info: []byte{0}}, status.Errorf(codes.Canceled, "context was canceled")
	default:
		nc := &gs.NcUserLoginAck{}
		nc.NumOfWorld = byte(1)

		var worlds [1]gs.WorldInfo
		w1 := gs.WorldInfo{
			WorldNumber: 0,
			WorldName:   gs.Name4{},
			WorldStatus: 1,
		}
		copy(w1.WorldName.Name[:], "INITIO")
		copy(w1.WorldName.NameCode[:], []uint32{262, 16720, 17735, 76})
		worlds[0] = w1

		nc.NumOfWorld = byte(1)
		nc.Worlds = worlds
		//

		// unaligned structures =/
		//return &lw.WorldsInfo{Info: []byte{1, 0, 1}}, nil

		if data, err := networking.WriteBinary(nc); err == nil {
			return &lw.WorldsInfo{Info: data}, nil
		} else {
			return &lw.WorldsInfo{Info: []byte{0}}, err
		}
	}
}

func (s *server) ConnectionInfo(ctx context.Context, req *lw.SelectedWorld) (*lw.WorldConnectionInfo, error) {
	select {
	case <-ctx.Done():
		return nil, status.Errorf(codes.Canceled, "context was canceled")
	default:
		//nc := gs.NcUserWorldSelectAck{
		//	WorldStatus: 6,
		//	Ip:          gs.Name4{},
		//	Port:        9110,
		//	ValidateNew: [32]uint16{},
		//}

		data := make([]byte, 0)
		data = append(data, byte(6))
		data = append(data, []byte("192.168.1.131")...)
		data = append(data, []byte{0, 0, 0}...)

		var port int16 = 9110
		b := make([]byte, 2)
		binary.LittleEndian.PutUint16(b, uint16(port))
		data = append(data, b...)
		log.Info(hex.EncodeToString(data))
		return &lw.WorldConnectionInfo{Info: data}, nil

		//if data, err := networking.WriteBinary(nc); err == nil {
		//	return &lw.WorldConnectionInfo{Info: data}, nil
		//} else {
		//	return &lw.WorldConnectionInfo{Info: []byte{0}}, err
		//}
	}
}
