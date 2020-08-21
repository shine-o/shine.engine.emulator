package world_master

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/spf13/viper"
)

var redisClient *redis.Client

func initRedis() {
	host := viper.GetString("session.redis.host")
	port := viper.GetString("session.redis.port")
	db := viper.GetInt("session.redis.db")
	addr := fmt.Sprintf("%v:%v", host, port)
	log.Infof("initializing redis instance: %v", addr)
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // no password set
		DB:       db, // use default DB
	})
	redisClient = client
}

func persist(rw *registeredWorlds) error {
	sd, err := json.Marshal(rw)
	if err != nil {
		log.Error(err)
		return err
	}

	err = redisClient.Set("world-master", sd, 0).Err()
	if err != nil {
		log.Error(err)
		return err
	}

	log.Infof("persisting worlds %v ", string(sd))
	return nil
}
