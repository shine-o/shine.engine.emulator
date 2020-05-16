package service

type playerHandleEvent struct {
	player  *player
	session *session
	done chan bool
	err     chan error
}

func (e *playerHandleEvent) erroneous() <-chan error {
	return e.err
}