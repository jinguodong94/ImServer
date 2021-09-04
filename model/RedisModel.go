package model

type OnLineUserInfo struct {
	Token            string `json:"token"`
	UserId           uint   `json:"user_id"`
	HeartbeatExpires uint64 `json:"heartbeat_expires"` //心跳时间
	ServerId         string `json:"server_id"`         //所在的服务器ID
}
