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
func Producer() {
	producer, err := nsq.NewProducer("127.0.0.1:4150", nsq.NewConfig())
	if err != nil {
		fmt.Println("NewProducer", err)
		panic(err)
	}

	Message := redisClient.SubscribeMessage()
	go func() {
		for {
			select {
			case message := <-Message:
				//fmt.Println("---server go message----:", message)
				//publish的时候，分发n个topic
				//user.strategy.14519b4e-4418-491d-8368-14278bf615e6
				if err := producer.Publish("user.strategy", []byte(message)); err != nil {
					fmt.Println("Publish nsq err:", err)
					panic(err)
				}
			}
		}
	}()
}

func main() {
	//ConsumerA()
	done := make(chan bool)
	Producer()
	<-done
}
