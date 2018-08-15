package consumer

import (
	"fmt"
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

func (b *Broker) channleFunc() map[string]func(chanK string, msg string, clientChan chan string, comUDId string) {
	channelFuncMap := make(map[string]func(chanK string, msg string, clientChan chan string, comUDId string))
	channelFuncMap["delete.device.pc"] = deleteDevicePc
	channelFuncMap["delete.device.mobile"] = deleteDeviceMobile
	channelFuncMap["disable.user.pc"] = deleteDevicePc
	channelFuncMap["delete.user.pc"] = deleteDevicePc
	channelFuncMap["user.strategy.pc"] = deleteDevicePc
	channelFuncMap["client.update.pc"] = deleteDevicePc
	channelFuncMap["generalChannel.pc"] = generalChannelPc
	return channelFuncMap
}

func (b *Broker) monitorChannel() {
	channelFuncMap := b.channleFunc()
	for channelKey, channel := range b.Channles {
		go func(channel chan string, channelKey string) {
			for msg := range channel {
				// key:delete.device.pc msg:message
				fmt.Println("----key:", channelKey, "----msg:", msg)
				for messageChan, clientKey := range b.Clients {
					// key 建立连接的
					fmt.Println("---Clients key", clientKey)
					//messageChan <- msg
					channelFuncMap[channelKey](channelKey, msg, messageChan, clientKey)
				}
			}
		}(channel, channelKey)
	}
}
