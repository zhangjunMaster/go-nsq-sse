package consumer

import (
	"encoding/json"
	"log"
	"strings"
	"time"
)

type Data struct {
	Name      string            `json:"name"`
	Content   map[string]string `json:"content"`
	TimeStamp time.Time         `json:"timeStamp"`
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

func (bm *BrokerMessage) push() {
	eventID := (*bm.Config)[bm.ChannelKey]["eventID"]
	switch eventID {
	case "pushMessage":
		bm.pushMessage()
	case "pushNotification":
		bm.pushNotification()
	}
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

func (bm *BrokerMessage) createInfo() {
	name := (*bm.Config)[bm.ChannelKey]["name"]
	content := make(map[string]string)
	data := Data{name, content, time.Now()}
	message := Message{"pushNotification", data}
	bm.ClientMessage = message
}

func (bm *BrokerMessage) pushMessage() {
	log.Println("----pushMessage", bm.Msg, bm.ClientKey)
	if bm.Msg != bm.ClientKey {
		return
	}
	bm.createMessage()
	message := bm.ClientMessage
	log.Println("----message:", message)
	structJson, _ := json.Marshal(message)
	string := string(structJson)
	log.Println("----string:", string)
	bm.MessageChan <- string
}

// companyId.strategy:"companyId.strategy.userId1,companyId.strategy.userId2"
func (bm *BrokerMessage) pushNotification() {
	channelCID := strings.Split(bm.ClientKey, "_")[0]
	userId := strings.Split(bm.ClientKey, "_")[1]
	msgKey := strings.Split(bm.Msg, ":")[0]
	msgV := strings.Split(bm.Msg, ":")[1]
	if !strings.Contains(msgV, userId) {
		return
	}
	msgCID := strings.Split(msgKey, ".")[0]
	if channelCID != msgCID {
		return
	}

	bm.createInfo()
	message := bm.ClientMessage
	structJson, _ := json.Marshal(message)
	string := string(structJson)
	bm.MessageChan <- string
}
