package main

// this file is for pushing messages to the NSP, the producer
import (
	"go-nsq-sse/redis-client"
	"log"

	nsq "github.com/nsqio/go-nsq"
)

//producer
//nsq.NewProducer func NewProducer(addr string, config *Config) (*Producer, error)
//Producer is a struct，Publish is its method
//1.redisClient.SubscribeMessage()
func Producer(channels []string) {
	producer, err := nsq.NewProducer("127.0.0.1:4150", nsq.NewConfig())
	if err != nil {
		log.Println("NewProducer", err)
		panic(err)
	}
	channelMap := make(map[string]chan string)
	for _, v := range channels {
		channelMap[v] = make(chan string, 1000)
	}
	redisClient.SubscribeMessage(channelMap)
	for key, channle := range channelMap {
		go func(key string, channle chan string) {
			for {
				select {
				case message := <-channle:
					//fmt.Println("---server go message----:", message)
					//when publish to topics
					//user.strategy.14519b4e-4418-491d-8368-14278bf615e6
					if err := producer.Publish(key, []byte(message)); err != nil {
						log.Println("Publish nsq err:", err)
						panic(err)
					}
				}
			}
		}(key, channle)
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
