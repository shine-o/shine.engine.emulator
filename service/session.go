package service

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/google/uuid"
	networking "github.com/shine-o/shine.engine.networking"
	"github.com/spf13/viper"
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
	host := viper.GetString("session.redis.host")
	port := viper.GetString("session.redis.port")
	db := viper.GetInt("session.redis.port")

	addr := fmt.Sprintf("%v:%v", host, port)
	log.Infof("initializing redis instance: %v", addr)
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // no password set
		DB:      db,  // use default DB
	})
	redisClient = client
	//err := client.Set("key", "value", 0).Err()
	//if err != nil {
	//	panic(err)
	//}
	//
	//val, err := client.Get("key").Result()
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println("key", val)
	//
	//val2, err := client.Get("key2").Result()
	//if err == redis.Nil {
	//	fmt.Println("key2 does not exist")
	//} else if err != nil {
	//	panic(err)
	//} else {
	//	fmt.Println("key2", val2)
	//}
}
