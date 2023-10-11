package simple

import (
	"net/http"

	"github.com/qf0129/ginz"
	"github.com/qf0129/ginz/dao"
	"github.com/qf0129/ginz/pkg/errs"
	"github.com/qf0129/ginz/pkg/secures"
)

func RequireTokenFromCookie() ginz.Middleware {
	return func(c *ginz.Context) {
		tk, err := c.C.Cookie(ginz.Config.TokenKey)
		if err != nil {
			c.ReturnErr(errs.InvalidToken.Add(err.Error()))
			return
		}

		uid, err := secures.ParseToken(tk, ginz.Config.Secret, ginz.Config.TokenExpiredTime)
		if err != nil {
			c.ReturnErr(errs.InvalidToken.Add(err.Error()))
			return
		}

		existsUser, err := dao.QueryOneByPk[User](uid)
		if err != nil {
			c.ReturnErr(errs.UserNotFound.Add(err.Error()))
			return
		}
		c.C.Set("user", existsUser)
		c.C.Next()
	}
}

// 跨域请求
func CorsMiddleware() ginz.Middleware {
	return func(c *ginz.Context) {
		c.C.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.C.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.C.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.C.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.C.Request.Method == "OPTIONS" {
			c.C.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.C.Next()
	}
}
