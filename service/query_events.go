package service

type queryMapEvent struct {
	id  int
	zm  chan *zoneMap
	err chan error
}

type queryPlayerEvent struct {
	handle  uint16
	p  chan *player
	err chan error
}

func (e *queryMapEvent) erroneous() <-chan error {
	return e.err
}


func (e *queryPlayerEvent) erroneous() <-chan error {
	return e.err
}
