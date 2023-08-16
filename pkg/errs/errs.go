package errs

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
	CreateDataFailed   = &Err{Code: 1000101, Msg: "新建数据失败"}
	QueryDataFailed    = &Err{Code: 1000102, Msg: "查询数据失败"}
	RetrieveDataFailed = &Err{Code: 1000102, Msg: "查询数据失败"}
	UpdateDataFailed   = &Err{Code: 1000103, Msg: "更新数据失败"}
	DeleteDataFailed   = &Err{Code: 1000104, Msg: "删除数据失败"}
	DataNotExists      = &Err{Code: 1000105, Msg: "数据不存在"}

	InvalidParams       = &Err{Code: 1000201, Msg: "无效的参数"}
	ValidateParamFailed = &Err{Code: 1000202, Msg: "校验参数失败"}
	ParseParamFailed    = &Err{Code: 1000203, Msg: "解析参数失败"}

	AuthFailed        = &Err{Code: 1000301, Msg: "认证失败"}
	InvalidToken      = &Err{Code: 1000302, Msg: "无效的令牌"}
	UserAlreadyExists = &Err{Code: 1000303, Msg: "用户已存在"}
	UserNotFound      = &Err{Code: 1000304, Msg: "用户不存在"}
	IncorrectPassword = &Err{Code: 1000305, Msg: "密码不正确"}
	HashPassword      = &Err{Code: 1000306, Msg: "哈希密码失败"}
	CreateToken       = &Err{Code: 1000307, Msg: "创建令牌失败"}
	CreateUser        = &Err{Code: 1000308, Msg: "创建用户失败"}
)
