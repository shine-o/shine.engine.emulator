package service

import (
	"context"
	"fmt"
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
		return &lw.WorldsInfo{}, status.Errorf(codes.Canceled, "context was canceled")
	default:
		var law activeWorlds
		aw.mu.Lock()
		law = *aw
		aw.mu.Unlock()

		worlds := make([]*lw.WorldInfo, 0)
		for _, w := range law.activeWorlds {
			id, err := strconv.Atoi(w.id)
			if err != nil {
				log.Error(err)
				continue
			}
			worlds = append(worlds, &lw.WorldInfo{
				WorldNumber:          int32(id),
				WorldName:            w.name,
				WorldStatus:          6,
			})
		}

		return &lw.WorldsInfo{
			Worlds: worlds,
		}, nil
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
		return &lw.WorldConnectionInfo{
			Name:                 w.name,
			IP:                   w.extIP,
			Port:                int32(port),
		}, nil
	}
}
