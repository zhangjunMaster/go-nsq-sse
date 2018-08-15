package consumer

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type Broker struct {
	Clients        map[chan string]string
	NewClients     chan chan string
	DefunctClients chan chan string
	Messages       chan string
	UserId         chan string
	Channles       map[string]chan string
}

type Data struct {
	Name      string            `json:"name"`
	Content   map[string]string `json:"content"`
	TimeStamp time.Time         `json:"timeStamp"`
}

type Message struct {
	EventID string `json:"eventID"`
	Data    Data   `json:"data"`
}

func (b *Broker) pushDeleteDeivice(channelKey string, msg string, clientChannle chan string, companyUserDeviceId string) {
	if msg != companyUserDeviceId {
		return
	}
	deviceId := strings.Split(msg, "_")[2]
	content := make(map[string]string)
	content["deviceId"] = deviceId
	data := Data{"removeDevice", content, time.Now()}
	message := Message{"pushMessage", data}
	structJson, _ := json.Marshal(message)
	string := string(structJson)
	fmt.Println("----message:", string)
	clientChannle <- string
}

// start是开启一个goroutine,遍历 b.clients这个map，一旦有message则将message给
// b.clients中的channel
// start是存信息到channel中
func (b *Broker) Start() {
	channels := []string{"delete.device.pc", "delete.device.mobile"}
	for _, v := range channels {
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
	// 从channel拿数据返回成client
	/**
	* 遍历监听的channels
	* 针对每个channel 遍历里面的信息，每个chanel
	* 开一个goroutine的原因是因为里面有阻塞
	* 用 for range循环的原因是因为，for msg := range channel 只要channel 有数据就会
	* 一直执行，不用 再加 for {}
	* 而其他地方没有for则是执行一遍就不执行，例如 select {} 如果没有for在最外层就是只能执行一遍
	 */
	go func() {
		for key, channel := range b.Channles {
			go func(channel chan string, key string) {
				for msg := range channel {
					fmt.Println("----key:", key, "----msg:", msg)
					for messageChan, key := range b.Clients {
						fmt.Println("---Clients key", key)
						//messageChan <- msg
						b.pushDeleteDeivice(key, msg, messageChan, key)
					}
				}
			}(channel, key)
		}
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
