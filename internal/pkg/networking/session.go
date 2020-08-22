package networking

// SessionFactory type for the shine service to implement if it needs session data
type SessionFactory interface {
	New() Session
}

// Session type for the shine service to implement if it needs session data
type Session interface {
	Identifier() string
}
