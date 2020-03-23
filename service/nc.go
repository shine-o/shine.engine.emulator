package service

import "github.com/shine-o/shine.engine.networking/structs"

type NcWorldFactory struct{}

// New  creates new NC structs specific to this service, will be used by the networking service when unpacking data
func (nwf *NcWorldFactory) NewNc(opCode int) interface{} {
	switch opCode {
	case 3087:
		return structs.NcUserLoginWorldReq{}
	case 2061:
		return structs.NcMiscGameTimeAck{}
	default:
		log.Warningf("no struct was assigned for operation code %v", opCode)
		return nil
	}
}