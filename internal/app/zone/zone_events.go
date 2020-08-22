package zone

type playerLogoutStartEvent struct {
	sessionID string
	mapID int
	handle uint16
	err chan error
}

type playerLogoutCancelEvent struct {
	sessionID string
	err chan error
}

type playerLogoutConcludeEvent struct {
	sessionID string
	err chan error
}
