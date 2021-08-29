package constant

type ReqType int

const (
	Req_heartbeat   ReqType = 0 //心跳
	Req_login       ReqType = 1 //1 请求登录
	Req_sendMessage ReqType = 2 //2 发送消息
)
