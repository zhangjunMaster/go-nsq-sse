package redisClient

import (
	"fmt"
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
	for key, channel := range channelMap {
		go func(key string, channel chan string) {
			pubsub := Client.Subscribe(key)
			defer pubsub.Close()
			subscr, err := pubsub.ReceiveTimeout(time.Second)
			if err != nil {
				panic(err)
			}
			fmt.Println("subscr", subscr)
			for {
				msg, err := pubsub.ReceiveMessage()
				if err != nil {
					fmt.Println("err:", err)
					panic(err)
				}
				message := msg.Payload
				fmt.Println(message)
				channel <- message
			}
		}(key, channel)
	}

}
