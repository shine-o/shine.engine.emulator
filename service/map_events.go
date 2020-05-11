package service

type registerPlayerHandleEvent struct {
	player  *player
	session *session
	done chan bool
	err     chan error
}

type handleCleanUpEvent struct {
	player  *player
	session *session
	err     chan error
}

func (e *registerPlayerHandleEvent) erroneous() <-chan error {
	return e.err
}

func (e *handleCleanUpEvent) erroneous() <-chan error {
	return e.err
}
