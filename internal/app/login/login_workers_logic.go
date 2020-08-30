package login

import (
	"github.com/shine-o/shine.engine.emulator/internal/pkg/networking"
	"github.com/shine-o/shine.engine.emulator/pkg/structs"
	"google.golang.org/grpc/connectivity"
)

func clientVersionLogic(ev *clientVersionEvent) {
	err := checkClientVersion(ev.nc)
	if err != nil {
		log.Error(err)
		ncUserClientWrongVersionCheckAck(ev.np)
		return
	}
	ncUserClientRightVersionCheckAck(ev.np)
}

func worldManagerStatusLogic(ev *worldManagerStatusEvent) {
	conn, err := newRPCClient("world_master")

	if err != nil {
		log.Error(err)
		return
	}

	defer conn.Close()

	if conn.GetState() != connectivity.Ready {
		log.Error("connection with the world master was not Ready")
		return
	}

	ncUserWorldStatusAck(ev.np)
}

func credentialsLoginLogic(ev *credentialsLoginEvent, e event, l *login) {
	err := checkCredentials(ev.nc)
	if err != nil {
		log.Error(e)
		ncUserLoginFailAck(ev.np, 69)
		return
	}

	loginSuccessful(l, ev.np)
}

func tokenLoginLogic(l *login, ev *tokenLoginEvent) {
	_, err := redisClient.Get(ev.nc.Otp.Name).Result()
	if err != nil {
		log.Error(err)
		ncUserLoginFailAck(ev.np, 69)
		return
	}

	loginSuccessful(l, ev.np)
}

func loginSuccessful(l *login, np *networking.Parameters) {
	nc := structs.NcUserLoginAck{}
	for _, w := range l.worlds {
		nc.Worlds = append(nc.Worlds, structs.WorldInfo{
			WorldNumber: byte(w.id),
			WorldName: structs.Name4{
				Name: w.name,
			},
			WorldStatus: 6,
		})
	}
	nc.NumOfWorld = byte(len(l.worlds))
	ncUserLoginAck(np, nc)
}

func serverSelectLogic(l *login, ev *serverSelectEvent) {
	for _, w := range l.worlds {
		if byte(w.id) == ev.nc.WorldNo {
			nc := structs.NcUserWorldSelectAck{
				WorldStatus: 6,
				Ip: structs.Name4{
					Name: w.ip,
				},
				Port: uint16(w.port),
			}
			ncUserWorldSelectAck(ev.np, nc)
			return
		}
	}
	log.Errorf("failed to find world with ID %v in %v", ev.nc.WorldNo, l.worlds)
}
