package ginz

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/qf0129/ginz/pkg/errs"
)

type CommonReq struct {
	Id       string
	Ids      []string
	Filter   map[string]any
	Page     int
	PageSize int
}

type RespBody struct {
	ReqId string
	Code  int
	Msg   string
	Data  any
}

type Context struct {
	C     *gin.Context
	ReqId string
}

func (c *Context) Param(key string) string {
	return c.C.Params.ByName(key)
}

func (c *Context) GetRequestData() (data *CommonReq) {
	c.ShouldBindJSON(&data)
	return
}

func (c *Context) ShouldBindWith(obj any, b binding.Binding) *errs.Err {
	if er := b.Bind(c.C.Request, obj); er != nil {
		return errs.ParseParamFailed.Add(er.Error())
	}
	return nil
}

func (c *Context) ShouldBindJSON(obj any) *errs.Err {
	return c.ShouldBindWith(obj, binding.JSON)
}

var Validator = validator.New()

func (c *Context) Validate(obj any) *errs.Err {
	err := Validator.Struct(obj)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return nil
		}
		for _, err := range err.(validator.ValidationErrors) {
			return errs.ValidateParamFailed.Add(err.Field())
		}
	}
	return nil
}

func (c *Context) ReturnOk(data any) {
	c.C.JSON(http.StatusOK, RespBody{
		ReqId: c.ReqId,
		Code:  0,
		Msg:   "ok",
		Data:  data,
	})
	c.C.Abort()
}

func (c *Context) ReturnErr(err *errs.Err) {
	c.C.JSON(http.StatusOK, RespBody{
		ReqId: c.ReqId,
		Code:  err.Code,
		Msg:   err.Msg,
		Data:  nil,
	})
	c.C.Abort()
}
