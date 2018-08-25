package consumer

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

/**
* take the data from the channel and return it to the client
* traverse the channels being monitored
* for each channel to traverse the information inside, each chanel
* the reason for opening a goroutine is because there is a blockage inside
* the reason for using the for range loop is because for MSG := range channle as long as the channel has data
* always execute, no need to add for {}
* if there is no for elsewhere, it is executed once and not executed, such as select {} if there is no for in the outermost layer, it can only be executed once
 */

func parseConfig() (map[string]map[string]string, error) {
	var result map[string]map[string]string
	data, err := ioutil.ReadFile("../channel.conf")
	if err != nil {
		return nil, err
	}

	// json to map
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (b *Broker) monitorChannel() {
	config, err := parseConfig()
	if err != nil {
		panic(err)
	}
	for channelKey, channle := range b.Channles {
		go func(channle chan string, channelKey string) {
			for msg := range channle {
				// key:delete.device.pc msg:message
				for messageChan, clientKey := range b.Clients {
					// clientKey is the params in the connection
					//messageChan <- msg
					brokerMessage := BrokerMessage{ChannelKey: channelKey, Msg: msg, ClientKey: clientKey, MessageChan: messageChan, Config: &config}
					method := config[channelKey]["eventID"]
					if method == "pushMessage" {
						log.Println("-----pushMessage")
						go brokerMessage.pushMessage()
					} else {
						log.Println("-----pushNotification")
						go brokerMessage.pushNotification()
					}
				}
			}
		}(channle, channelKey)
	}
}
