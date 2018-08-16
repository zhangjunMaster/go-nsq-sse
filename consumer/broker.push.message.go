package consumer

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type Data struct {
	Name      string            `json:"name"`
	Content   map[string]string `json:"content"`
	TimeStabm time.Time         `json:"timeStabm"`
}

type Message struct {
	EventID string `json:"eventID"`
	Data    Data   `json:"data"`
}

type BrokerMessage struct {
	ChannelKey    string
	Msg           string
	MessageChan   chan string
	Config        *map[string]map[string]string
	ClientKey     string
	ClientMessage Message
}

func (bm *BrokerMessage) createMessage() {
	deviceId := strings.Split(bm.Msg, "_")[2]
	userId := strings.Split(bm.Msg, "_")[1]
	content := make(map[string]string)
	content["deviceId"] = deviceId
	content["userId"] = userId
	name := (*bm.Config)[bm.ChannelKey]["name"]
	data := Data{name, content, time.Now()}
	message := Message{"pushMessage", data}
	bm.ClientMessage = message
}

func (bm *BrokerMessage) pushMessage() {
	fmt.Printf("%+v", bm)
	if bm.Msg != bm.ClientKey {
		return
	}
	bm.createMessage()
	message := bm.ClientMessage
	structJson, _ := json.Marshal(message)
	string := string(structJson)
	fmt.Println("----message:", string)
	bm.MessageChan <- string
}

func (bm *BrokerMessage) pushInfo() {
	if bm.Msg != bm.ChannelKey {
		return
	}
	message := bm.ClientMessage
	structJson, _ := json.Marshal(message)
	string := string(structJson)
	fmt.Println("----message:", string)
	bm.MessageChan <- string
}
