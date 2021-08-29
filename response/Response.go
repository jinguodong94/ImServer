package response

const (
	SUCCESS = iota
	ERROR
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type UserInfo struct {
	Account  string `json:"account"`
	NickName string `json:"nick_name"`
	Icon     string `json:"icon"`
	Token    string `json:"token"`
}

func NewResponse(code int, msg string, data interface{}) (resp *Response) {
	resp = new(Response)
	resp.Code = code
	resp.Msg = msg
	resp.Data = data
	return
}

func NewSuccessResponse(msg string, data interface{}) (resp *Response) {
	resp = NewResponse(SUCCESS, msg, data)
	return
}

func NewErrorResponse(msg string) (resp *Response) {
	resp = NewResponse(ERROR, msg, nil)
	return
}
