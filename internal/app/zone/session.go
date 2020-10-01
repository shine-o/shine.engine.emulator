package zone

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/networking"
)

type sessionFactory struct{}

type session struct {
	id            string
	characterID   uint64
	characterName string
	// TODO: check if its viable to add directly the pointers to map, player objects to increase performance
	mapID   int
	handle  uint16
}

func (s sessionFactory) New() networking.Session {
	return &session{
		id: fmt.Sprintf("zone-%v", uuid.New().String()),
	}
}

func (s *session) Identifier() string {
	return s.id
}
