package service

import "github.com/shine-o/shine.engine.core/structs"

type securityEventError struct {
	code    int
	message string
}

func (e securityEventError) Error() string {
	return e.message
}

type clientSHNEvent struct {
	inboundNC structs.NcMapLoginReq
	ok chan bool
	err chan error
}

func (e * clientSHNEvent) erroneous() <- chan error {
	return e.err
}

var errBadSHN = securityEventError{
	code:    0,
	message: "client and server SHN files do not match",
}
