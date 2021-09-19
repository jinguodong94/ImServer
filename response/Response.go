package response

const (
	SUCCESS = iota
	ERROR
	NOT_LOGIN_ERROR
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
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

func NewNotLoginErrorResponse() (resp *Response) {
	resp = NewResponse(NOT_LOGIN_ERROR, "未登录", nil)
	return
}
