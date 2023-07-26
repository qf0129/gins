package ginz

import "fmt"

type Errors struct {
	err  error  `json:"-"`
	Code int    `json:"Code,omitempty"`
	Msg  string `json:"Msg,omitempty"`
}

func (err *Errors) String() string {
	return fmt.Sprintf("Code: %v, Message: %v,error:%v", err.Code, err.Msg, err.err)
}

func (err *Errors) Err() string {
	return err.Msg
}

func (err *Errors) AddMsg(msg string) *Errors {
	err.Msg += " " + msg
	return err
}

// common
var (
	ErrParams       = &Errors{Code: 1000001, Msg: "请求参数错误"}
	ErrAuthFailed   = &Errors{Code: 1000002, Msg: "认证失败"}
	ErrUserNotFound = &Errors{Code: 1000003, Msg: "用户不存在"}
)
