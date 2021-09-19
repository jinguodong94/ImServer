package dao

import "gorm.io/gorm"

const (
	TableName_GroupMessages         = "group_messages"
	TableName_GroupMessagesRelation = "group_messages_relations"
	TableName_UserGroupRelations    = "user_group_relations"
	TableName_FriendRelations       = "friend_relations"
)

//用户表
type Users struct {
	gorm.Model
	Account  string `gorm:"type:varchar(32);not null;unique" json:"account"`
	Pwd      string `gorm:"type:varchar(32);not null" json:"pwd"`
	NickName string `gorm:"type:varchar(32)" json:"nick_name"`
	Icon     string `json:"icon"`
}

//单聊临时消息表(未发送成功的)
type Messages struct {
	gorm.Model
	//聊天类型  0 点对点聊 , 1 群聊  , 2 聊天室
	ChatType    int `gorm:"type:int;not null;comment:'聊天类型  0 点对点聊 , 1 群聊  , 2 聊天室'" json:"chat_type"`
	MessageType int `gorm:"type:int;not null;comment:'消息类型  0 文本 , 1 图片'" json:"message_type"`
	//(不存表) 0 普通消息类型 , 1 指令类型 只针对在线用户
	MsgType int `gorm:"-" json:"msg_type"`

	MsgFlag int `gorm:"-" json:"msg_flag"`

	//离线发送 0 是  1不是
	IsOffLine int `gorm:"-" json:"is_off_line"`

	FromUid uint `gorm:"index;not null" json:"from_uid"` //发送人id
	ToUid   uint `gorm:"index;not null" json:"to_uid"`   //接收者id

	//消息状态 0 未发送 1 已发送待确定  2 发送成功
	Status         int    `gorm:"type:int;not null;comment:'消息状态 0 未发送 1 已发送待确定'" json:"status"`
	ExtendInfo     string `gorm:"comment:'扩展信息'" json:"extend_info"`     //扩展信息
	MessageContent string `gorm:"comment:'消息内容'" json:"message_content"` //消息内容
}

//单聊消息记录表(发送成功的)
type MessageHistory struct {
	gorm.Model
	//聊天类型  0 点对点聊 , 1 群聊  , 2 聊天室
	ChatType       int    `gorm:"type:int;not null;comment:'聊天类型  0 点对点聊 , 1 群聊  , 2 聊天室'" json:"chat_type"`
	MessageType    int    `gorm:"type:int;not null;comment:'消息类型  0 文本 , 1 图片'" json:"message_type"`
	FromUid        uint   `gorm:"index;not null" json:"from_uid"`        //发送人id
	ToUid          uint   `gorm:"not null" json:"to_uid"`                //接收者id
	ExtendInfo     string `gorm:"comment:'扩展信息'" json:"extend_info"`     //扩展信息
	MessageContent string `gorm:"comment:'消息内容'" json:"message_content"` //消息内容
}

//群聊消息表
type GroupMessages struct {
	gorm.Model
	//(不存表)消息类型  0 普通消息类型 , 1 指令类型 只针对在线用户
	MsgType int `gorm:"-" json:"msg_type"`
	//离线发送 0 是  1不是
	IsOffLine int `gorm:"-" json:"is_off_line"`

	MessageType    int    `gorm:"type:int;not null;comment:'消息类型  0 文本 , 1 图片'" json:"message_type"`
	FromUid        uint   `gorm:"index;not null" json:"from_uid"`        //发送人id
	ToGroupId      uint   `gorm:"index;not null" json:"to_group_id"`     //群id
	ExtendInfo     string `gorm:"comment:'扩展信息'" json:"extend_info"`     //扩展信息
	MessageContent string `gorm:"comment:'消息内容'" json:"message_content"` //消息内容
}

//群聊消息临时关联表
type GroupMessagesRelation struct {
	gorm.Model
	Uid   uint `gorm:"index;not null" json:"uid"`    //接收者id
	MsgId uint `gorm:"index;not null" json:"msg_id"` //消息id
	//消息状态 0 未发送 1 已发送待确定  2 发送成功
	Status int `gorm:"type:int;not null;comment:'消息状态 0 未发送 1 已发送待确定'" json:"status"`
}

//群聊消息记录关联表(发送成功的)
type GroupMessagesHistoryRelation struct {
	gorm.Model
	Uid   uint `gorm:"index;not null" json:"uid"`    //接收者id
	MsgId uint `gorm:"index;not null" json:"msg_id"` //消息id
}

//群组表
type Groups struct {
	gorm.Model
	GroupName string `gorm:"type:varchar(32);not null;comment:'群名'" json:"group_name"` //群名
	Notice    string `gorm:"comment:'公告'" json:"notice"`                               //公告
	Icon      string `gorm:"comment:'头像'" json:"icon"`                                 //群头像
	Extend    string `json:"extend"`                                                   //扩展字段
}

//用户和群组关系表
type UserGroupRelation struct {
	gorm.Model
	Uid     uint `gorm:"index;not null" json:"uid"`
	GroupId uint `gorm:"index;not null" json:"group_id"`
	Role    int  `gorm:"not null;default 0;comment:'成员角色 0 普通成员,1 群主 ,2 管理员'" json:"role"`
}

//好友关系表
type FriendRelation struct {
	gorm.Model
	Uid      uint `gorm:"index;not null" json:"uid"`
	FriendId uint `gorm:"index;not null" json:"friend_id"`
	Status   int  `gorm:"not null;comment:'0 正常关系  1 已拉黑'" json:"status"`
}

//好友申请表
type FriendApply struct {
	gorm.Model
	ToUid    uint   `gorm:"index;not null" json:"to_uid"`
	ApplyUid uint   `gorm:"index;not null;comment:'申请人id'" json:"apply_uid"`
	Status   int    `gorm:"not null;comment:'0 未处理  1 同意  2 拒绝'" json:"status"`
	Remarks  string `gorm:"not null;comment:'备注'" json:"remarks"`
}
