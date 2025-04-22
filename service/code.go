package service

type Code int64

const (
	Success Code = 2000 + iota
	InvalidParam
	ServerBusy
)

var codeMsgMap = map[Code]string{
	Success:      "success",
	InvalidParam: "请求参数错误",
	ServerBusy:   "服务繁忙",
}

func (c Code) Msg() string {
	msg, ok := codeMsgMap[c]
	if !ok {
		msg = codeMsgMap[ServerBusy]
	}
	return msg
}
