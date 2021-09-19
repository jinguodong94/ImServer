package response

type UserInfo struct {
	Account  string `json:"account"`
	NickName string `json:"nick_name"`
	Icon     string `json:"icon"`
	Token    string `json:"token"`
}
