package ginz

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// type PageBody struct {
// 	List     any   `json:"list"`
// 	Page     int   `json:"page"`
// 	PageSize int   `json:"page_size"`
// 	Total    int64 `json:"total"`
// }

type RespBody struct {
	ReqId string
	Code  int
	Msg   string
	Data  any
}

func RespOk(c *gin.Context, data any) {
	c.JSON(http.StatusOK, RespBody{
		ReqId: c.GetString(REQUEST_KEY_ID),
		Code:  0,
		Msg:   "ok",
		Data:  data,
	})
	c.Abort()
}

func RespErr(c *gin.Context, err *Errors) {
	c.JSON(http.StatusOK, RespBody{
		ReqId: c.GetString(REQUEST_KEY_ID),
		Code:  err.Code,
		Msg:   err.Msg,
		Data:  nil,
	})
	c.Abort()
}
