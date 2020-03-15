package networking

//type Session struct {
//	id string
//	ServiceData interface {}
//}
//
//func NewSession() * Session  {
//	return &Session{
//		id:	uuid.New().String(),
//	}
//}
//
//func (s * Session) Identifier() string {
//	return s.id
//}

type SessionFactory interface {
	New() Session
}

type Session interface {
	Identifier() string
}