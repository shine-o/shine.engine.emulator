package service

type queryMapEvent struct {
	id  int
	zm  chan *zoneMap
	err chan error
}

type queryPlayerEvent struct {
	id  int
	zm  chan *player
	err chan error
}

func (e *queryMapEvent) erroneous() <-chan error {
	return e.err
}
