package model

//数据返回
type MsgResponse struct {
	/*
		0 心跳回应	1 登录成功  2登录失败  3 获取或者更新好友列表   4 用户信息变更   5 消息
	*/
	RespType int
	Data     interface{} //返回的数据
}

//客户端请求数据
type MsgRequest struct {
	Type int         `json:"type"` //0 心跳  1 请求登录   2 发送消息
	Data interface{} `json:"data"` //数据
}

//消息
type Message struct {
	ChatType       int8   //聊天类型  0 点对点聊 , 1 群聊  , 2 聊天室
	MsgType        int8   //消息类型  0 普通消息类型 , 1 指令类型 只针对在线用户
	FromUid        string //发送人id
	ToUid          string //接收者id
	ExtendInfo     string //扩展信息
	MessageContent string //消息内容
}

//心跳
type Heartbeat struct {
	Token string `json:"token"`
}

//登录
type Login struct {
	UserId string `json:"user_id"`
	Pwd    string `json:"pwd"`
}
