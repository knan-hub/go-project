package common

type Code int64

const (
	SUCCESS Code = 2000 + iota
	ERROR
	INVALID_PARAMS
	SERVER_BUSY
	REDIS_SET_ERROR
	REDIS_GET_ERROR
)

var codeMsgMap = map[Code]string{
	SUCCESS:         "success",
	ERROR:           "error",
	INVALID_PARAMS:  "请求参数错误",
	SERVER_BUSY:     "服务繁忙",
	REDIS_SET_ERROR: "redis set error",
	REDIS_GET_ERROR: "redis get error",
}

func (c Code) Msg() string {
	msg, ok := codeMsgMap[c]
	if !ok {
		msg = codeMsgMap[SERVER_BUSY]
	}
	return msg
}
