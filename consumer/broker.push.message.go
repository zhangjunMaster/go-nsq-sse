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
	TimeStamp time.Time         `json:"timeStamp"`
}

type Message struct {
	EventID string `json:"eventID"`
	Data    Data   `json:"data"`
}

func deleteDevicePc(chanKey string, msg string, clientChan chan string, comUDId string) {
	if msg != comUDId {
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
	clientChan <- string
}

func deleteDeviceMobile(chanKey string, msg string, clientChan chan string, comUDId string) {
	deleteDevicePc(chanKey, msg, clientChan, comUDId)
}

func generalChannelPc(chanKey string, msg string, clientChan chan string, comUDId string) {
	if msg != comUDId {
		return
	}
	data := Data{"pullGateway", make(map[string]string), time.Now()}
	message := Message{"pushNotification", data}
	structJson, _ := json.Marshal(message)
	string := string(structJson)
	fmt.Println("----message:", string)
	clientChan <- string
}
