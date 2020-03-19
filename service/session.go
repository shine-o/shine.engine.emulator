package service

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/google/uuid"
	networking "github.com/shine-o/shine.engine.networking"
	"github.com/spf13/viper"
)

type sessionFactory struct {
	worldId string
}

type session struct {
	Id string `json:"id"`
	WorldId string
	UserName string `json:"user_name"`
}

func (s sessionFactory) New() networking.Session  {
	return &session {
		Id:	uuid.New().String(),
		WorldId: s.worldId,
	}
}

func (s * session) Identifier() string  {
	return s.Id
}

var redisClient * redis.Client

func initRedis()  {
	host := viper.GetString("session.redis.host")
	port := viper.GetString("session.redis.port")
	db := viper.GetInt("session.redis.db")
	addr := fmt.Sprintf("%v:%v", host, port)
	log.Infof("initializing redis instance: %v", addr)
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // no password set
		DB:       db,  // use default DB
	})
	redisClient = client
}

func persistSession(ws *session) error {
	if sd, err := json.Marshal(ws); err != nil {
		log.Error(err)
		return err
	} else {
		key := fmt.Sprintf("%v-world", ws.UserName)
		if err := redisClient.Set(key, sd, 0).Err(); err != nil {
			log.Error(err)
			return err
		} else {
			log.Infof("persisting session %v -> %v", key, string(sd))
			return nil
		}
	}
}