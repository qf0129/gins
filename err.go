package ginz

import "fmt"

type Err struct {
	subErr error `json:"-"`
	Code   int
	Msg    string
}

func (err *Err) String() string {
	return fmt.Sprintf("Code: %v, Message: %v, Error:%v", err.Code, err.Msg, err.subErr)
}

func (err *Err) Add(msg string) *Err {
	return &Err{Code: err.Code, Msg: err.Msg + ":" + msg}
}

func (err *Err) Args(args ...any) *Err {
	return &Err{Code: err.Code, Msg: fmt.Sprintf(err.Msg, args...)}
}

// common
var (
	ErrInvalidParams = &Err{Code: 1000010, Msg: "无效的参数"}
	ErrInvalidToken  = &Err{Code: 1000011, Msg: "无效的令牌"}

	ErrAuthFailed        = &Err{Code: 1000020, Msg: "认证失败"}
	ErrUserAlreadyExists = &Err{Code: 1000021, Msg: "用户已存在"}
	ErrUserNotFound      = &Err{Code: 1000022, Msg: "用户不存在"}
	ErrIncorrectPassword = &Err{Code: 1000023, Msg: "密码不正确"}
	ErrHashPassword      = &Err{Code: 1000024, Msg: "哈希密码失败"}
	ErrCreateToken       = &Err{Code: 1000025, Msg: "创建令牌失败"}
	ErrCreateUser        = &Err{Code: 1000026, Msg: "创建用户失败"}

	ErrDBQueryFailed = &Err{Code: 1000030, Msg: "数据查询失败"}
)
