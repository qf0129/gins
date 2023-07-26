package ginz

import (
	"github.com/gin-gonic/gin"
	"github.com/qf0129/ginz/pkg/strs"
)

type ApiHandler func(c *gin.Context) (data any, err *Errors)

type Api struct {
	Name    string
	Info    string
	Method  string
	Handler ApiHandler
}

type ApiGroup struct {
	BasePath    string
	RouterGroup *gin.RouterGroup
	Apis        []*Api
}

func (group *ApiGroup) Add(api *Api) {
	group.Apis = append(group.Apis, api)
	group.RouterGroup.Handle(api.Method, api.Name, func(c *gin.Context) {
		c.Set(REQUEST_KEY_ID, strs.UUID())
		data, err := api.Handler(c)
		if err == nil {
			RespOk(c, data)
		} else {
			RespErr(c, err)
		}
	})
}
