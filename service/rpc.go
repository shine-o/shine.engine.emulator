package service

import (
	"context"
	"fmt"
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
		nc := structs.NcUserLoginAck{}
		nc.NumOfWorld = byte(1)

		w1 := structs.WorldInfo{
			WorldNumber: 0,
			WorldName:   structs.Name4{},
			// 1: behaviour -> cannot enter, message -> The server is under maintenance.
			// 2: behaviour -> cannot enter, message -> You cannot connect to an empty server.
			// 3: behaviour -> cannot enter, message -> The server has been reserved for a special use.
			// 4: behaviour -> cannot enter, message -> Login failed due to an unknown error.
			// 5: behaviour -> cannot enter, message -> The server is full.
			// 6: behaviour -> ok
			WorldStatus: 6,
		}
		copy(w1.WorldName.Name[:], "PAGEL")

		nc.Worlds = append(nc.Worlds, w1)

		data, err := structs.Pack(&nc)
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
			Ip: structs.Name4{},
			Port: uint16(port),
		}
		copy(nc.Ip.Name[:], w.extIP)

		data, err := structs.Pack(&nc)
		if err != nil {
			return nil, status.Errorf(codes.Unavailable, "something went horribly wrong")
		}
		return &lw.WorldConnectionInfo{Info: data}, nil
	}
}
