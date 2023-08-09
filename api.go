package ginz

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qf0129/ginz/pkg/errs"
)

type ApiHandler func(c *gin.Context) (any, *errs.Err)

type Api struct {
	Name      string
	Info      string
	Method    string
	Handler   ApiHandler
	WithCache bool
}

type ApiGroup struct {
	BasePath    string
	RouterGroup *gin.RouterGroup
	Apis        []*Api
}

// 使用中间件
func (group *ApiGroup) Use(middleware Middleware) {
	group.RouterGroup.Use(func(ctx *gin.Context) {
		middleware(ctx)
	})
}

// 添加api对象
func (group *ApiGroup) Add(api *Api) {
	group.Apis = append(group.Apis, api)
	group.RouterGroup.Handle(api.Method, api.Name, func(c *gin.Context) {
		data, err := api.Handler(c)
		if err != nil {
			RespErr(c, err)
		} else {
			RespOk(c, data)
		}
	})
}

// 添加api对象
func (group *ApiGroup) AddApi(name string, handler ApiHandler) {
	group.Add(&Api{
		Name:    name,
		Method:  http.MethodPost,
		Handler: handler,
	})
}
