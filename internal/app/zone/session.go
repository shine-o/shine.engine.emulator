package zone

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/networking"
	"github.com/spf13/viper"
)

var redisClient *redis.Client

type sessionFactory struct{}

type session struct {
	id            string
	characterID   uint64
	characterName string
	// TODO: check if its viable to add directly the pointers to map, player objects to increase performance
	mapID int
	mapName int
	handle uint16
}

func (s sessionFactory) New() networking.Session {
	return &session{
		id: fmt.Sprintf("zone-%v", uuid.New().String()),
	}
}

func (s *session) Identifier() string {
	return s.id
}

func initRedis() {
	host := viper.GetString("session.redis.host")
	port := viper.GetString("session.redis.port")
	db := viper.GetInt("session.redis.db")

	addr := fmt.Sprintf("%v:%v", host, port)
	log.Infof("initializing redis instance: %v", addr)
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       db,
	})
	redisClient = client
}
