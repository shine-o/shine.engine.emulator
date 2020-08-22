package login

import (
	"github.com/shine-o/shine.engine.emulator/internal/pkg/networking"
	"github.com/shine-o/shine.engine.emulator/pkg/structs"
)

type clientVersionEvent struct {
	nc *structs.NcUserClientVersionCheckReq
	np *networking.Parameters
}

//	credentialsLogin
//	serverList
//	serverSelect
//	tokenLogin
type credentialsLoginEvent struct {
	nc  *structs.NcUserUsLoginReq
	np  *networking.Parameters
}

type worldManagerStatusEvent struct {
	np  *networking.Parameters
}

type serverSelectEvent struct {
	nc * structs.NcUserWorldSelectReq
	np  *networking.Parameters
}

type tokenLoginEvent struct {
	nc * structs.NcUserLoginWithOtpReq
	np  *networking.Parameters
}