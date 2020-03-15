package networking

type SessionFactory interface {
	New() Session
}

type Session interface {
	Identifier() string
}