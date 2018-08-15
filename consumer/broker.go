package consumer

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

type Broker struct {
	Clients        map[chan string]string
	NewClients     chan chan string
	DefunctClients chan chan string
	Messages       chan string
	UserId         chan string
	Channles       map[string]chan string
	ChannleTopics  []string
}

// start是开启一个goroutine,遍历 b.clients这个map，一旦有message则将message给
// b.clients中的channel
// start是存信息到channel中
func (b *Broker) Start() {
	b.ChannleTopics = []string{
		"delete.device.pc", "delete.device.mobile",
		"disable.user.pc", "delete.user.pc",
		"user.strategy.pc", "client.update.pc",
		"generalChannel.pc",
	}
	for _, v := range b.ChannleTopics {
		b.Channles[v] = make(chan string, 10000)
	}
	go func() {
		for {
			select {
			case s := <-b.NewClients:
				go func() {
					key := <-b.UserId
					b.Clients[s] = key
				}()
				log.Println("Added new client")
			case s := <-b.DefunctClients:
				delete(b.Clients, s)
				close(s)
				log.Println("Removed client")
			}
		}
	}()

	go func() {
		b.monitorChannel()
	}()
}

//建立连接，没建立一个连接，创建一个messageChan
func (b *Broker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	userId := strings.TrimPrefix(r.URL.Path, "/client/v3/push/register/")
	fmt.Println("UserId:", userId)
	// Make sure that the writer supports flushing.
	// http.Flusher是接口，w.(http.Flusher) 判断是否符合这个接口
	f, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}
	messageChan := make(chan string)
	//将messageChan给b,在start中放message,在另外一个go goroutine中取信息
	b.NewClients <- messageChan
	b.UserId <- userId

	notify := w.(http.CloseNotifier).CloseNotify()
	//监听是否断开连接
	go func() {
		<-notify
		// Remove this client from the map of attached Clients when `EventHandler` exits.
		b.DefunctClients <- messageChan
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
		go func(msg string) {
			//将信息返回到客户端
			fmt.Fprintf(w, "data: Message: %s\n\n", msg)
			// Flush the response.  This is only possible if the repsonse supports streaming.
			f.Flush()
		}(msg)

	}
	log.Println("Finished HTTP request at ", r.URL.Path)
}
