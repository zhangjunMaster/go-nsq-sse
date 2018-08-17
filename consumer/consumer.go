package consumer

// nsq建立连接，同时确定消费的是哪个topic,channle
import (
	"fmt"
	"log"

	nsq "github.com/nsqio/go-nsq"
)

//nsq consumer处理消息
type ConsumerHandler struct {
	b       *Broker
	channle string
}

//func HandleMessage(message *Message) error
//func (h HandlerFunc) HandleMessage(m *Message) error
//    HandleMessage implements the Handler interface
// 每一个channle都会new 一个ConsumerHandler，channle即这个channle的名称
func (c *ConsumerHandler) HandleMessage(msg *nsq.Message) error {
	log.Println("ConsumerHandler nsq.Message:", string(msg.Body))
	go func(c *ConsumerHandler, msg *nsq.Message) {
		channleKey := c.channle
		// 如果这里没有被消费，就会阻塞 c.b.Channles[redisChannel] 这个channel
		c.b.Channles[channleKey] <- string(msg.Body)
		//c.b.Messages <- string(msg.Body)
	}(c, msg)
	return nil
}

// nsq拿出数据，然后给客户端,消费者
// 注意调用的Handler是HandleMessage，不是httpServer
func Consumer(b *Broker) {
	//每一个channle与nsq建立一个连接
	for _, v := range b.ChannleTopics {
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
