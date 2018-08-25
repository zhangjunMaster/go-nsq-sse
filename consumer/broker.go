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

// Start is to open a goroutine, traverse the map of b.clients,
//and give the message to the channel in b.clients once there is a message

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
				log.Println("Added new client, total:", len(b.Clients)+1)
			case s := <-b.DefunctClients:
				delete(b.Clients, s)
				close(s)
				log.Println("Removed client, total:", len(b.Clients)+1)
			}
		}
	}()

	go func() {
		b.monitorChannel()
	}()
}

//Establish a connection, and create a messageChan for each connection
func (b *Broker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	userId := strings.TrimPrefix(r.URL.Path, "/client/v3/push/register/")
	log.Println("client key:", userId)
	// Make sure that the writer supports flushing.
	// http.Flusher is interface，w.(http.Flusher) assertions

	f, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}
	messageChan := make(chan string)
	//Send messageChan to b, put message in start, and get the message in another go goroutine
	b.NewClients <- messageChan
	b.UserId <- userId

	notify := w.(http.CloseNotifier).CloseNotify()
	//Listen to if the connection is broken
	go func() {
		<-notify
		// Remove this client from the map of attached Clients when `EventHandler` exits.
		b.DefunctClients <- messageChan
		log.Println("HTTP connection just closed.")
	}()

	//Establish the foundation of sse and set the header information
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	for {
		msg, open := <-messageChan
		if !open {
			// If our messageChan was closed, this means that the client has disconnected.
			break
		}
		fmt.Println("------msg:", msg)
		go func(msg string) {
			n, err := fmt.Fprintf(w, "%s", msg)
			fmt.Println("-----n, err:", n, err)
			//json.NewEncoder(w).Encode(msg)
			// Flush the response.  This is only possible if the repsonse supports streaming.
			f.Flush()
		}(msg)

	}
	log.Println("Finished HTTP request at ", r.URL.Path)
}
