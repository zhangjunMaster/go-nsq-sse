package consumer

import (
	"encoding/json"
	"io/ioutil"
)

/**
* 从channel拿数据返回成client
* 遍历监听的channels
* 针对每个channel 遍历里面的信息，每个chanel
* 开一个goroutine的原因是因为里面有阻塞
* 用 for range循环的原因是因为，for msg := range channel 只要channel 有数据就会
* 一直执行，不用 再加 for {}
* 而其他地方没有for则是执行一遍就不执行，例如 select {} 如果没有for在最外层就是只能执行一遍
 */

func parseConfig() (map[string]map[string]string, error) {
	var result map[string]map[string]string
	data, err := ioutil.ReadFile("../config.json")
	if err != nil {
		return nil, err
	}

	//读取的数据为json格式，需要进行解码
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
	for channelKey, channel := range b.Channles {
		go func(channel chan string, channelKey string) {
			for msg := range channel {
				// key:delete.device.pc msg:message
				for messageChan, clientKey := range b.Clients {
					// clientKey 建立连接中带的东西
					//messageChan <- msg
					//考虑到并发，不能new一个object,只采用了传参的方法，所以使用了map
					brokerMessage := BrokerMessage{ChannelKey: channelKey, Msg: msg, ClientKey: clientKey, MessageChan: messageChan, Config: &config}
					go brokerMessage.pushMessage()
				}
			}
		}(channel, channelKey)
	}
}
