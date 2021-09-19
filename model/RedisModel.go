package model

type OnLineUserInfo struct {
	Token            string `json:"token"`
	UserId           uint   `json:"user_id"`
	HeartbeatExpires uint64 `json:"heartbeat_expires"` //心跳时间
	ServerId         string `json:"server_id"`         //所在的服务器ID
}

type RoomInfo struct {
	RoomId     uint   `json:"room_id"`
	RoomTitle  string `json:"room_title"`
	RoomIcon   string `json:"room_icon"`
	RoomNotice string `json:"room_notice"`
}
