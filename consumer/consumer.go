package consumer

// nsq建立连接，同时确定消费的是哪个topic,channle
import (
	"fmt"

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
// each channle will have a new ConsumerHandler, channle is the name of this channle

func (c *ConsumerHandler) HandleMessage(msg *nsq.Message) error {
	go func(c *ConsumerHandler, msg *nsq.Message) {
		channleKey := c.channle
		// if it's not consumed, it blocks the channel
		c.b.Channles[channleKey] <- (string(msg.Body))
		//c.b.Messages <- string(msg.Body)
	}(c, msg)
	return nil
}

// NSQ takes the data and gives it to the client, the consumer
// note that the Handler invoked is HandleMessage, not httpServer
func Consumer(b *Broker) {
	// each channle establishes a connection with NSQ
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
