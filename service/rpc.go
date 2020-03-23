package service

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/shine-o/shine.engine.networking"
	"github.com/shine-o/shine.engine.networking/structs"
	lw "github.com/shine-o/shine.engine.protocol-buffers/login-world"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"
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
		nc := &structs.NcUserLoginAck{}
		nc.NumOfWorld = byte(1)

		var worlds [1]structs.WorldInfo
		w1 := structs.WorldInfo{
			WorldNumber: 0,
			WorldName:   structs.Name4{},
			WorldStatus: 1,
		}
		copy(w1.WorldName.Name[:], "INITIO")
		//copy(w1.WorldName.NameCode[:], []uint32{262, 16720, 17735, 76})
		//worlds[0] = w1

		nc.NumOfWorld = byte(1)
		nc.Worlds = worlds

		data, err := networking.WriteBinary(nc)
		if err != nil {
			return &lw.WorldsInfo{Info: []byte{0}}, err
		}
		return &lw.WorldsInfo{Info: data}, nil
	}
}

func (s *server) ConnectionInfo(ctx context.Context, req *lw.SelectedWorld) (*lw.WorldConnectionInfo, error) {
	select {
	case <-ctx.Done():
		return nil, status.Errorf(codes.Canceled, "context was canceled")
	default:
		aw.mu.Lock()
		inx := fmt.Sprintf("%v", req.Num)
		w, available := aw.activeWorlds[inx]
		aw.mu.Unlock()

		if !available {
			log.Error(aw)
			err := status.Errorf(codes.Unavailable, "requested world with id %v is not available", req.Num)
			log.Error(err)
			return nil, err
		}

		port, err := strconv.Atoi(w.port)
		if err != nil {
			log.Error(err)
			return nil, status.Errorf(codes.FailedPrecondition, "incorrect world port %v", port)
		}

		nc := structs.NcUserWorldSelectAck{
			WorldStatus: 6,
			Ip: structs.Name4{
				Name: [16]byte{},
			},
			Port: uint16(port),
		}

		copy(nc.Ip.Name[:], w.extIP)

		data := make([]byte, 0)
		data = append(data, nc.WorldStatus)

		if b, err := networking.WriteBinary(nc.Ip.Name); err == nil {
			data = append(data, b...)
		}

		if b, err := networking.WriteBinary(nc.Port); err == nil {
			data = append(data, b...)
		}
		log.Infof("sending server connection info to client %v", hex.EncodeToString(data))
		return &lw.WorldConnectionInfo{Info: data}, nil
	}
}
