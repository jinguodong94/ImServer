package model

import "imserver/dao"

//数据返回
type MsgResponse struct {
	/*
		0 心跳成功	1 心跳失败	2 登录成功  3 登录失败  4 获取或者更新好友列表   5 用户信息变更   6 消息
	*/
	RespType int         `json:"resp_type"`
	Data     interface{} `json:"data"` //返回的数据
}

//客户端请求数据
type MsgRequest struct {
	Type int         `json:"type"` //0 心跳  1 请求登录   2 发送消息
	Data interface{} `json:"data"` //数据
}

//消息
type Message struct {
	ChatType       int    `json:"chat_type"`       //聊天类型  0 点对点聊 , 1 群聊  , 2 聊天室
	MsgType        int    `json:"msg_type"`        //消息类型  0 普通消息类型 , 1 指令类型 只针对在线用户
	FromUid        string `json:"from_uid"`        //发送人id
	ToUid          string `json:"to_uid"`          //接收者id
	ExtendInfo     string `json:"extend_info"`     //扩展信息
	MessageContent string `json:"message_content"` //消息内容
}

//心跳
type Heartbeat struct {
	Token string `json:"token"`
}

//登录
type Login struct {
	Account string `json:"account"`
	Pwd     string `json:"pwd"`
}

type MqMessage struct {
	Type int         `json:"type"` //0 单聊消息   1 群组信息
	Data interface{} `json:"data"`
}

type AckMessage struct {
	ChatType  int `json:"chat_type"`  //0 单聊  1群聊  2聊天室
	MessageId int `json:"message_id"` //消息id
}

type GroupMessage struct {
	dao.GroupMessages
	ToUid uint `json:"to_uid"`
}
