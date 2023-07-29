package redis

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"gome/utils"
	"gopkg.in/yaml.v3"
	"os"
)

var Conf *utils.MeConfig

func init() {
	confFile, _ := os.ReadFile("config.yaml")
	err := yaml.Unmarshal(confFile, &Conf)
	if err != nil {
		fmt.Println(err)
	}
}

func NewRedisClient() *redis.Client {
	host := Conf.CacheConf.Host
	port := Conf.CacheConf.Port
	//password := conf.CacheConf.Password
	cache := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return cache
}
