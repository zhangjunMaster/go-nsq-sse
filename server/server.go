/**
1.Consumer: 消费者监听topic,channel.当有信息时，Consumer的consumer.AddHandler(&ConsumerHandler{b})，会调用
  (c *ConsumerHandler) HandleMessage，将信息传给 b.messages


2.ServeHTTP中将messageChan给b,在Start中放message到b,在ServeHTTP的另外一个go goroutine中取信息
	1.建立长连接，messageChan := make(chan string) 空的messageChan给b.newClients
	2.start 函数， 采用go func + for{} => s := <-b.newClients 遍历的时候，如果发现有空的messageChan，则显示建立连接
	3.并同时将这个空的messageChan，挂在clients map上
	4.start函数中监听 b.messages channel, 根据message的信息和clients中messageChan的key做条件对比，符合条件
	  将 messageChan <- messages
	5.ServeHTTP中监听messageChan channel,如果有数据就发送

3.  clients    放：map  userId和messageChan,根据userId判断是否将信息 messageChan <- message
	newClients 放messageChan， 然后将messageChan给clients，组合形成map
	messages   放给客户端的信息
	userId     放客户端的id，用于匹配clients中map，因为message中会有userId
*/
//这个文件是建立长连接，然后从nsq中拿信息

package main

import (
	"fmt"
	"go-nsq/redis-client"
	"html/template"
	"log"
	"net/http"
	"strings"

	nsq "github.com/nsqio/go-nsq"
)

type Broker struct {
	clients        map[chan string]string
	newClients     chan chan string
	defunctClients chan chan string
	messages       chan string
	userId         chan string
}

// start是开启一个goroutine,遍历 b.clients这个map，一旦有message则将message给
// b.clients中的channel
// start是存信息到channel中
func (b *Broker) Start() {
	go func() {
		for {
			select {
			case s := <-b.newClients:
				go func() {
					key := <-b.userId
					b.clients[s] = key
				}()
				log.Println("Added new client")
			case s := <-b.defunctClients:
				delete(b.clients, s)
				close(s)
				log.Println("Removed client")
			case msg := <-b.messages:
				strategy := strings.Split(msg, ":")
				key := strategy[0]
				strategyInfo := strings.Split(msg, key+":")[1]
				go func() {
					len, err := redisClient.Client.LLen(key).Result()
					if err != nil {
						panic(err)
					}
					for i := 0; i < int(len); i++ {
						val, err := redisClient.Client.RPop(key).Result()
						if err != nil {
							panic(err)
						}
						for s, key := range b.clients {
							if key == val {
								s <- strategyInfo
							}
						}
					}

				}()
				log.Printf("Broadcast message to %d clients", len(b.clients))
			}
		}
	}()
}

func (b *Broker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println("userId:", r.URL.Query()["userId"][0])
	// Make sure that the writer supports flushing.
	// http.Flusher是接口，w.(http.Flusher) 判断是否符合这个接口
	f, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}
	messageChan := make(chan string)
	//将messageChan给b,在start中放message,在另外一个go goroutine中取信息
	b.newClients <- messageChan
	b.userId <- r.URL.Query()["userId"][0]
	notify := w.(http.CloseNotifier).CloseNotify()
	//监听是否断开连接
	go func() {
		<-notify
		// Remove this client from the map of attached clients when `EventHandler` exits.
		b.defunctClients <- messageChan
		log.Println("HTTP connection just closed.")
	}()
	//是建立sse的基础，设置头部信息
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	for {
		msg, open := <-messageChan
		if !open {
			// If our messageChan was closed, this means that the client has disconnected.
			break
		}
		//将信息返回到客户端
		fmt.Fprintf(w, "data: Message: %s\n\n", msg)
		// Flush the response.  This is only possible if the repsonse supports streaming.
		f.Flush()
	}
	log.Println("Finished HTTP request at ", r.URL.Path)
}

// route at "/" below in `main`.
// 返回template的Handler
func MainPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	t, err := template.ParseFiles("../templates/index.html")
	if err != nil {
		log.Fatal("WTF dude, error parsing your template.")
	}
	// Render the template, writing to `w`.
	t.Execute(w, "Duder")
	log.Println("Finished HTTP request at ", r.URL.Path)
}

//nsq consumer处理消息
type ConsumerHandler struct {
	b *Broker
}

func (c *ConsumerHandler) HandleMessage(msg *nsq.Message) error {
	fmt.Println("ConsumerHandler nsq.Message:", string(msg.Body))
	c.b.messages <- string(msg.Body)
	return nil
}

// nsq建立连接，同时确定消费的是哪个topic,channel
// nsq拿出数据，然后给客户端,消费者
// 注意调用的Handler是HandleMessage，不是httpServer
func Consumer(b *Broker) {
	consumer, err := nsq.NewConsumer("user.strategy", "user.strategy", nsq.NewConfig())
	if err != nil {
		fmt.Println("NewConsumer:", err)
		panic(err)
	}
	consumer.AddHandler(&ConsumerHandler{b})
	if err := consumer.ConnectToNSQLookupd("127.0.0.1:4161"); err != nil {
		fmt.Println("ConnectToNSQLookupd", err)
		panic(err)
	}
}

// Main routine
//
func main() {
	b := &Broker{
		make(map[chan string]string),
		make(chan (chan string)),
		make(chan (chan string)),
		make(chan string),
		make(chan string),
	}
	// Start processing events
	// Start会产生一个goroutine，这个goroutine会监听各个channel的数据
	b.Start()
	go func() {
		Consumer(b)
	}()
	http.Handle("/events", b)
	http.Handle("/", http.HandlerFunc(MainPageHandler))
	http.ListenAndServe(":8001", nil)
}
