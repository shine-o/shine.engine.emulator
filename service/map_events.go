package service

type playerHandleEvent struct {
	player  *player
	session *session
	done chan bool
	err     chan error
}
