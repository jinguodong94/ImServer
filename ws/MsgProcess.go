package ws

import (
	"gindemo/model"
	"log"
)

//处理心跳
func processHeartbeat(client *Client, heartbeat *model.Heartbeat) {
	//TODO 校验用户token 后返回信息给Client
	log.Println("处理心跳", heartbeat.Token)

}

//处理登录
func processLogin(client *Client, login *model.Login) {
	//TODO 校验用户账号密码，生成Token存入redis

}

//处理消息
func processMessage(client *Client, message *model.Message) {
	//TODO 发送消息

}
