package simple

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qf0129/ginz"
	"github.com/qf0129/ginz/pkg/dao"
	"github.com/qf0129/ginz/pkg/errs"
	"github.com/qf0129/ginz/pkg/secures"
)

func RequireTokenFromCookie() ginz.Middleware {
	return func(c *gin.Context) {
		tk, err := c.Cookie(ginz.Config.TokenKey)
		if err != nil {
			ginz.RespErr(c, errs.ErrInvalidToken.Add(err.Error()))
			return
		}

		uid, err := secures.ParseToken(tk, ginz.Config.Secret, ginz.Config.TokenExpiredTime)
		if err != nil {
			ginz.RespErr(c, errs.ErrInvalidToken.Add(err.Error()))
			return
		}

		existsUser, err := dao.QueryOneByPk[User](uid)
		if err != nil {
			ginz.RespErr(c, errs.ErrUserNotFound.Add(err.Error()))
			return
		}
		c.Set("user", existsUser)
		c.Next()
	}
}

// 跨域请求
func CorsMiddleware() ginz.Middleware {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
