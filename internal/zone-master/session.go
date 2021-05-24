package zonemaster

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

func persist(ws *registeredMaps) error {
	sd, err := json.Marshal(ws)
	if err != nil {
		log.Error(err)
		return err
	}
	err = redisClient.Set("zone-master", sd, 0).Err()
	if err != nil {
		log.Error(err)
		return err
	}
	log.Infof("persisting maps %v ", string(sd))
	return nil
}
