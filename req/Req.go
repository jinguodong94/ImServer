package req

type LoginReq struct {
	Account string `form:"account"`
	Pwd     string `form:"pwd"`
}

type RegisterReq struct {
	Account  string `form:"account"`
	Pwd      string `form:"pwd"`
	NickName string `form:"nickName"`
	Icon     string `form:"icon"`
}

type CreateGroupReq struct {
	GroupName string `form:"groupName"`
	Icon      string `form:"icon"`
}

type AddFriendReq struct {
	FriendId uint `form:"friendId"`
	//备注
	Remarks string `form:"remarks"`
}

type JoinOrExitGroupReq struct {
	GroupId string `form:"groupId"`
}

type OptionFriendApplyReq struct {
	OptionId   string `form:"optionId"`
	OptionType int    `form:"optionType"` //0 同意  1拒绝
}

type DeleteFriendReq struct {
	OptionUid uint `form:"optionUid"`
}

type CreateRoomReq struct {
	RoomTitle  string `form:"roomTitle"`
	RoomIcon   string `form:"roomIcon"`
	RoomNotice string `form:"roomNotice"`
}

type JoinOrExitRoomReq struct {
	RoomId uint `form:"roomId"`
}

type UserInfoReq struct {
}
