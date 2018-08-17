package redisClient

import (
	"log"
	"time"

	"github.com/go-redis/redis"
)

var Client = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "", // no password set
	DB:       0,  // use default DB
})

var dat map[string]interface{}

func SubscribeMessage(channelMap map[string]chan string) {
	for key, channle := range channelMap {
		go func(key string, channle chan string) {
			pubsub := Client.Subscribe(key)
			defer pubsub.Close()
			subscr, err := pubsub.ReceiveTimeout(time.Second)
			if err != nil {
				panic(err)
			}
			log.Println("subscr", subscr)
			for {
				msg, err := pubsub.ReceiveMessage()
				if err != nil {
					log.Println("err:", err)
					panic(err)
				}
				message := msg.Payload
				log.Println(message)
				channle <- message
			}
		}(key, channle)
	}

}
