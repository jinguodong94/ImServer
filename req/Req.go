package req

type LoginReq struct {
	Account string `form:"account"`
	Pwd string `form:"pwd"`
}

type RegisterReq struct {
	Account string `form:"account"`
	Pwd string `form:"pwd"`
	NickName string	`form:"nickName"`
	Icon string	`form:"icon"`
}

type UserInfoReq struct {


}