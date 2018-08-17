/**
1.Consumer: consumers listen to topic,channle. When there is information, Consumer's ConsumerHandler (&ConsumerHandler{b}) will be called
(c *ConsumerHandler) HandleMessage, passing the information to b.Messages
2. Give messageChan to b in ServeHTTP, put message to b in Start, and get information in another go goroutine of ServeHTTP
	1. Establish long connection, messageChan := make(chan string) empty messageChan to b.NewClients
	2. Start function, use go func + for{} => s := < -b.n ewClients when traversing, if free messageChan is found, the connection is established
	3. Hang this empty messageChan on clients map at the same time
	4. Start function listens for b.Messages channle, which is qualified according to the information of message and the key of messageChan in clients
Will messageChan < - Messages
	5. Listen to messageChan channle in ServeHTTP and send if there is data
3. Clients put: map UserId and messageChan, and judge whether the information is messageChan < -message based on the UserId
NewClients puts messageChan and then gives messageChan to clients, which is combined to form a map
Messages is sent to the client
The UserId puts the client's id to match the map in clients because there will be a UserId in the message
// this file is to make a long connection, and then take the information from NSQ
// the system flow chart, which is suitable for the use of go func(){}()
//serveHttp returns information
//messageChan <-> NewClients <-> Clients sets up a concatenated client, again using a map
// messages < - channle < - NSQ
// when taking, goroutine takes channel as the standard
// the collection of goroutine sent to client is based on channel
*/
package main

import (
	"go-nsq-sse/consumer"
	"go-nsq-sse/handler"
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
	// Start will produce a goroutine, which will listen to the data on each channel

	b.Start()
	go func() {
		consumer.Consumer(b)
	}()
	http.Handle("/client/v3/push/register/", b)
	http.Handle("/", http.HandlerFunc(handler.MainPageHandler))
	http.ListenAndServe(":8001", nil)
}
