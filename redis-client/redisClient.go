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

func SubscribeMessage() chan string {
	Message := make(chan string, 1000)
	go func() {
		pubsub := Client.Subscribe("user.strategy")
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
			/*
				strategy := strings.Split(message, ":")
				key := strategy[0]
				fmt.Println("----分割后的-----", strings.Split(message, key+":"))
				strategyInfo := strings.Split(message, key+":")[1]
				fmt.Println("--------strategyInfo", strategyInfo)
				byt := []byte(strategyInfo)
				var dat map[string]interface{}
				var strategyItem = make(map[string]interface{})
				//解码
				if err := json.Unmarshal(byt, &dat); err != nil {
					panic(err)
				}
				strategyItem["key"] = key
				strategyItem["strategy"] = dat
				fmt.Printf("%v", dat)
			*/
			Message <- message
		}
	}()

	return Message
}
