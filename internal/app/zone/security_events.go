package zone

import "github.com/shine-o/shine.engine.emulator/pkg/structs"

type securityEventError struct {
	code    int
	message string
}

func (e securityEventError) Error() string {
	return e.message
}

type playerSHNEvent struct {
	inboundNC structs.NcMapLoginReq
	ok        chan bool
	err       chan error
}


var errBadSHN = securityEventError{
	code:    0,
	message: "client and server SHN files do not match",
}
