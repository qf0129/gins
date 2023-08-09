package simple

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qf0129/ginz"
	"github.com/qf0129/ginz/pkg/dao"
	"github.com/qf0129/ginz/pkg/secures"
)

func RequireTokenFromCookie(app *ginz.Ginz) ginz.Middleware {
	return func(c *gin.Context) {
		tk, err := c.Cookie(app.Config.TokenKey)
		if err != nil {
			ginz.RespErr(c, ginz.ErrInvalidToken.Add(err.Error()))
			return
		}

		uid, err := secures.ParseToken(tk, app.Config.Secret, app.Config.TokenExpiredTime)
		if err != nil {
			ginz.RespErr(c, ginz.ErrInvalidToken.Add(err.Error()))
			return
		}

		existsUser, err := dao.QueryOneByPk[User](uid)
		if err != nil {
			ginz.RespErr(c, ginz.ErrUserNotFound.Add(err.Error()))
			return
		}
		c.Set("user", existsUser)
		c.Next()
	}
}

// 跨域请求
func CorsMiddleware(app *ginz.Ginz) ginz.Middleware {
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
