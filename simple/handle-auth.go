package simple

import (
	"github.com/gin-gonic/gin"
	"github.com/qf0129/ginz"
	"github.com/qf0129/ginz/crud"
	"github.com/qf0129/ginz/pkg/secures"
)

type AuthRequestBody struct {
	Username string
	Password string
}

// 用户登录接口
func UserLoginHandler(secretKey string, expiredSeconds int64) ginz.ApiHandler {
	return func(c *gin.Context) (data any, err *ginz.Err) {
		var req AuthRequestBody

		if er := c.ShouldBindJSON(&req); er != nil {
			err = ginz.ErrInvalidParams.Add(er.Error())
			return
		}

		if req.Username == "" {
			err = ginz.ErrInvalidParams.Add("Username")
			return
		}

		if req.Password == "" {
			err = ginz.ErrInvalidParams.Add("Password")
			return
		}

		existUser, er := crud.QueryOneByMap[User](map[string]any{"username": req.Username})
		if er != nil {
			err = ginz.ErrUserNotFound
			return
		}

		if !secures.VerifyPassword(req.Password, existUser.PasswordHash) {
			err = ginz.ErrIncorrectPassword
			return
		}

		token, er := secures.CreateToken(existUser.Id, secretKey)
		if er != nil {
			err = ginz.ErrCreateToken
			return
		}
		c.SetCookie("tk", token, int(expiredSeconds), "/", "*", true, true)
		c.SetCookie("uid", existUser.Id, int(expiredSeconds), "/", "*", true, false)
		data = map[string]any{"Token": token}
		return
	}
}

// 用户注册接口
func UserRegisterHandler() ginz.ApiHandler {
	return func(c *gin.Context) (data any, err *ginz.Err) {
		var req AuthRequestBody

		if er := c.ShouldBindJSON(&req); er != nil {
			err = ginz.ErrInvalidParams.Add(er.Error())
			return
		}

		if req.Username == "" {
			err = ginz.ErrInvalidParams.Add("username")
			return
		}

		if req.Password == "" {
			err = ginz.ErrInvalidParams.Add("password")
			return
		}

		existUser, _ := crud.QueryOneByMap[User](map[string]any{"username": req.Username})
		if existUser.Id != "" {
			err = ginz.ErrUserAlreadyExists
			return
		}

		psdHash, er := secures.HashPassword(req.Password)
		if er != nil {
			err = ginz.ErrHashPassword
			return
		}

		u := &User{
			Username:     req.Username,
			PasswordHash: psdHash,
		}

		if er = crud.CreateOne[User](u); er != nil {
			err = ginz.ErrCreateUser.Add(er.Error())
			return
		}
		data = map[string]any{"Id": u.Id}
		return
	}
}
