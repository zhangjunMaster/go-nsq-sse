package main

//这个文件是往nsp上推送消息,生产者
import (
	"fmt"
	"go-nsq/redis-client"

	nsq "github.com/nsqio/go-nsq"
)

//生产者
//nsq.NewProducer func NewProducer(addr string, config *Config) (*Producer, error)
//Producer 是一个struct，Publish是其方法
//1.redisClient.SubscribeMessage()
func Producer(channels []string) {
	producer, err := nsq.NewProducer("127.0.0.1:4150", nsq.NewConfig())
	if err != nil {
		fmt.Println("NewProducer", err)
		panic(err)
	}
	channelMap := make(map[string]chan string)
	for _, v := range channels {
		fmt.Println("----v:", v)
		channelMap[v] = make(chan string, 1000)
	}
	redisClient.SubscribeMessage(channelMap)
	for key, channel := range channelMap {
		go func(key string, channel chan string) {
			for {
				select {
				case message := <-channel:
					//fmt.Println("---server go message----:", message)
					//publish的时候，分发n个topic
					//user.strategy.14519b4e-4418-491d-8368-14278bf615e6
					if err := producer.Publish(key, []byte(message)); err != nil {
						fmt.Println("Publish nsq err:", err)
						panic(err)
					}
				}
			}
		}(key, channel)
	}
}
func main() {
	channels := []string{
		"delete.device.pc", "delete.device.mobile",
		"disable.user.pc", "delete.user.pc",
		"user.strategy.pc", "client.update.pc",
		"generalChannel.pc",
	}
	//ConsumerA()
	done := make(chan bool)
	Producer(channels)
	<-done
}
