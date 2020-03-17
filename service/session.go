package service

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/google/logger"
	"github.com/google/uuid"
	"github.com/shine-o/shine.engine.networking"
	"github.com/spf13/viper"
	"io/ioutil"
)

type sessionFactory struct {}

type session struct {
	id string
	userName string
}

func (s sessionFactory) New() networking.Session  {
	return &session{
		id:	fmt.Sprintf("login-%v", uuid.New().String()),
	}
}

func (s * session) Identifier() string  {
	return s.id
}

var redisClient * redis.Client

func initRedis()  {
	log = logger.Init("LoginLogger", true, true, ioutil.Discard)
	log.Info("LoginLogger init()")
	host := viper.GetString("session.redis.host")
	port := viper.GetString("session.redis.port")
	db := viper.GetInt("session.redis.db")

	addr := fmt.Sprintf("%v:%v", host, port)
	log.Infof("initializing redis instance: %v", addr)
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // no password set
		DB:      db,  // use default DB
	})
	redisClient = client
}
