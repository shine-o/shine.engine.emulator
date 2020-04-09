package service

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/google/uuid"
	"github.com/shine-o/shine.engine.networking"
	"github.com/spf13/viper"
)

var redisClient *redis.Client

type sessionFactory struct{}

type session struct {
	id       string
	userName string
}

func (s sessionFactory) New() networking.Session {
	return &session{
		id: fmt.Sprintf("login-%v", uuid.New().String()),
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
