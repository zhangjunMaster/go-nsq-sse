/**
1.Consumer: 消费者监听topic,channel.当有信息时，Consumer的consumer.AddHandler(&ConsumerHandler{b})，会调用
  (c *ConsumerHandler) HandleMessage，将信息传给 b.Messages


2.ServeHTTP中将messageChan给b,在Start中放message到b,在ServeHTTP的另外一个go goroutine中取信息
	1.建立长连接，messageChan := make(chan string) 空的messageChan给b.NewClients
	2.start 函数， 采用go func + for{} => s := <-b.NewClients 遍历的时候，如果发现有空的messageChan，则显示建立连接
	3.并同时将这个空的messageChan，挂在clients map上
	4.start函数中监听 b.Messages channel, 根据message的信息和clients中messageChan的key做条件对比，符合条件
	  将 messageChan <- Messages
	5.ServeHTTP中监听messageChan channel,如果有数据就发送

3.  Clients    放：map  UserId和messageChan,根据UserId判断是否将信息 messageChan <- message
	NewClients 放messageChan， 然后将messageChan给clients，组合形成map
	Messages   放给客户端的信息
	UserId     放客户端的id，用于匹配clients中map，因为message中会有UserId
*/
//这个文件是建立长连接，然后从nsq中拿信息
//本系统流程图，并发多个情况下适合用go func(){}()
//serveHttp返回信息
//   ^
//messageChan <-> NewClients <-> Clients 建立连接对应的client,还是用的map
//   								^
//								 messages<-channel<- nsq

// 取的时候，goroutine 以channel为标准
// 向client发的收 goroutine 以channel为标准

package main

import (
	"go-nsq/consumer"
	"go-nsq/handler"
	"net/http"
	_ "net/http/pprof"
)

// Main routine
//
func main() {
	b := &consumer.Broker{
		make(map[chan string]string),
		make(chan (chan string)),
		make(chan (chan string)),
		make(chan string),
		make(chan string),
		make(map[string]chan string),
		make([]string, 0),
	}
	// Start processing events
	// Start会产生一个goroutine，这个goroutine会监听各个channel的数据

	b.Start()
	go func() {
		consumer.Consumer(b)
	}()
	http.Handle("/client/v3/push/register/", b)
	http.Handle("/", http.HandlerFunc(handler.MainPageHandler))
	http.ListenAndServe(":8001", nil)
}
