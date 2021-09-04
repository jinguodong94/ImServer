package constant

const (
	//请求
	Req_heartbeat   int = 0 //心跳
	Req_login       int = 1 //1 请求登录
	Req_sendMessage int = 2 //2 发送消息
	Req_ACK_Message int = 3 //3 消息确认

	//回应
	Resp_heartbeat_success         int = 0 //0 心跳成功
	Resp_heartbeat_error           int = 1 //1 心跳失败
	Resp_login_success             int = 2 //2 登录成功
	Resp_login_error               int = 3 //3 登录失败
	Resp_get_or_update_friend_list int = 4 //4 获取或者更新好友列表
	Resp_user_info_update          int = 5 //5 用户信息变更
	Resp_send_message_success      int = 6 //6 发送消息成功
	Resp_send_message_error        int = 7 //7 发送消息失败
	Resp_recive_message            int = 8 //8 接受消息
	Resp_ack_message_success       int = 9 //9 消息确认成功

	Resp_parm_parse_error int = 999 // 参数解析失败

	//redis key
	Redis_key_user_im_login_info = "user_im_login_info"

	Redis_key_user_login_token = "user_login_token"

	Token = "token"
)
