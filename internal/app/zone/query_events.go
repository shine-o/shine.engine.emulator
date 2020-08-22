package zone

type queryMapEvent struct {
	id  int
	zm  chan *zoneMap
	err chan error
}

type queryPlayerEvent struct {
	handle uint16
	p      chan *player
	err    chan error
}
