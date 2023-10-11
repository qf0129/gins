package ginz

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/qf0129/ginz/pkg/errs"
	"github.com/qf0129/ginz/pkg/maps"
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

func (c *Context) ShouldBindQuery(obj any) *errs.Err {
	return c.ShouldBindWith(obj, binding.Query)
}

var Validator = validator.New()

func (c *Context) Validate(data any) *errs.Err {
	err := Validator.Struct(data)
	if err != nil {
		if errList, ok := err.(validator.ValidationErrors); ok {
			for _, e := range errList {
				return errs.ValidateParamFailed.Add(e.Field())
			}
		} else {
			return errs.ValidateParamFailed.Add(err.Error())
		}
	}
	return nil
}

func (c *Context) ValidateMap(data, rules map[string]any) *errs.Err {
	errMap := Validator.ValidateMap(data, rules)
	k, err := maps.GetFirstOfMap(errMap)
	if err != nil {
		return errs.ValidateParamFailed.Add(k)
	}
	return nil
}

func (c *Context) ParseAndValidate(data any) *errs.Err {
	err := c.ShouldBindJSON(data)
	if err != nil {
		return err
	}
	return c.Validate(data)
}

func (c *Context) ParseAndValidateToMap(data, rules map[string]any) *errs.Err {
	err := c.ShouldBindJSON(&data)
	if err != nil {
		return err
	}
	return c.ValidateMap(data, rules)
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

func (c *Context) ReturnAuthErr(err *errs.Err) {
	c.C.JSON(http.StatusUnauthorized, RespBody{
		ReqId: c.ReqId,
		Code:  err.Code,
		Msg:   err.Msg,
		Data:  nil,
	})
	c.C.Abort()
}

func (c *Context) ReturnAnyErr(err any) {
	if er, ok := err.(*errs.Err); ok {
		c.ReturnErr(er)
	} else {
		c.ReturnErr(errs.RequestError.Add(fmt.Sprintf("%v", err)))
	}
}

func (c *Context) Panic(err *errs.Err) {
	panic(err)
}

func (c *Context) Next() {
	c.C.Next()
}

func (c *Context) SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool) {
	c.C.SetCookie(name, value, maxAge, path, domain, secure, httpOnly)
}

func (c *Context) Cookie(name string) (string, error) {
	return c.C.Cookie(name)
}

func (c *Context) Set(key string, value any) {
	c.C.Set(key, value)
}

func (c *Context) Get(key string) (value any, exists bool) {
	return c.C.Get(key)
}

func (c *Context) MustGet(key string) any {
	return c.C.MustGet(key)
}

func (c *Context) GetString(key string) (s string) {
	return c.C.GetString(key)
}

func (c *Context) GetBool(key string) (b bool) {
	return c.C.GetBool(key)
}

func (c *Context) GetInt(key string) (i int) {
	return c.C.GetInt(key)
}
