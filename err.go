package ginz

import "fmt"

type Errors struct {
	Code int
	Msg  string
}

func (err *Errors) String() string {
	return fmt.Sprintf("Code: %v, Message: %v", err.Code, err.Msg)
}

func (err *Errors) AddMsg(msg string) *Errors {
	return &Errors{Code: err.Code, Msg: err.Msg + "" + msg}
}

func (err *Errors) Args(args ...any) *Errors {
	return &Errors{Code: err.Code, Msg: fmt.Sprintf(err.Msg, args...)}
}

// common
var (
	ErrParams       = &Errors{Code: 1000001, Msg: "请求参数错误"}
	ErrAuthFailed   = &Errors{Code: 1000002, Msg: "认证失败"}
	ErrUserNotFound = &Errors{Code: 1000003, Msg: "用户不存在"}
)
