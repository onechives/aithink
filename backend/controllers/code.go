package controllers

type ResCode int64

const (
	// 通用响应码约定（1000 为成功）
	CodeSuccess       ResCode = 1000
	CodeInvalidParams ResCode = 1001
	CodeNeedLogin     ResCode = 1002
	CodeInvalidToken  ResCode = 1003
	CodeServerBusy    ResCode = 1004
	CodeAuthFailed    ResCode = 1005
	CodeNotFound      ResCode = 1006
	CodeUserPending   ResCode = 1007
	CodeNicknameTaken ResCode = 1008
)

var codeMsgMap = map[ResCode]string{
	CodeSuccess:       "success",
	CodeInvalidParams: "invalid params",
	CodeNeedLogin:     "need login",
	CodeInvalidToken:  "invalid token",
	CodeServerBusy:    "server busy",
	CodeAuthFailed:    "auth failed",
	CodeNotFound:      "not found",
	CodeUserPending:   "user pending",
	CodeNicknameTaken: "nickname taken",
}

// Msg 将响应码映射为默认文案。
func (c ResCode) Msg() string {
	msg, ok := codeMsgMap[c]
	if ok {
		return msg
	}
	return codeMsgMap[CodeServerBusy]
}
