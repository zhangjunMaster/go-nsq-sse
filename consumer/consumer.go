package consumer

// nsq建立连接，同时确定消费的是哪个topic,channel
import (
	"fmt"

	nsq "github.com/nsqio/go-nsq"
)

//nsq consumer处理消息
type ConsumerHandler struct {
	b       *Broker
	channel string
}

//func HandleMessage(message *Message) error
//func (h HandlerFunc) HandleMessage(m *Message) error
//    HandleMessage implements the Handler interface
func (c *ConsumerHandler) HandleMessage(msg *nsq.Message) error {
	fmt.Println("ConsumerHandler nsq.Message:", string(msg.Body))

	go func(c *ConsumerHandler, msg *nsq.Message) {
		redisChannel := c.channel
		// 如果这里没有被消费，就会阻塞 c.b.Channles[redisChannel] 这个channel
		c.b.Channles[redisChannel] <- string(msg.Body)
		fmt.Println("-----Consumer message:", redisChannel, string(msg.Body))
		//c.b.Messages <- string(msg.Body)
	}(c, msg)
	return nil
}

// nsq拿出数据，然后给客户端,消费者
// 注意调用的Handler是HandleMessage，不是httpServer
func Consumer(b *Broker) {
	channels := []string{"delete.device.pc", "delete.device.mobile"}
	for _, v := range channels {
		//b.Channles[v] = make(chan string, 10000)
		go func(v string) {
			consumer, err := nsq.NewConsumer(v, v, nsq.NewConfig())
			if err != nil {
				fmt.Println("NewConsumer:", err)
				panic(err)
			}
			consumer.AddHandler(&ConsumerHandler{b, v})
			if err := consumer.ConnectToNSQLookupd("127.0.0.1:4161"); err != nil {
				fmt.Println("ConnectToNSQLookupd", err)
				panic(err)
			}
		}(v)
	}
}
