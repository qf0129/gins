package simple

import (
	"github.com/gin-gonic/gin"
	"github.com/qf0129/ginz"
	"github.com/qf0129/ginz/pkg/dao"
	"github.com/qf0129/ginz/pkg/errs"
	"github.com/qf0129/ginz/pkg/secures"
)

type AuthRequestBody struct {
	Username string
	Password string
}

// 用户登录接口
func UserLoginHandler() ginz.ApiHandler {
	return func(c *gin.Context) (data any, err *errs.Err) {
		var req AuthRequestBody

		if er := c.ShouldBindJSON(&req); er != nil {
			err = errs.ErrInvalidParams.Add(er.Error())
			return
		}

		if req.Username == "" {
			err = errs.ErrInvalidParams.Add("Username")
			return
		}

		if req.Password == "" {
			err = errs.ErrInvalidParams.Add("Password")
			return
		}

		existUser, er := dao.QueryOneByMap[User](map[string]any{"username": req.Username})
		if er != nil {
			err = errs.ErrUserNotFound
			return
		}

		if !secures.VerifyPassword(req.Password, existUser.PasswordHash) {
			err = errs.ErrIncorrectPassword
			return
		}

		token, er := secures.CreateToken(existUser.Id, ginz.Config.Secret)
		if er != nil {
			err = errs.ErrCreateToken
			return
		}
		c.SetCookie("tk", token, int(ginz.Config.TokenExpiredTime), "/", "*", true, true)
		c.SetCookie("uid", existUser.Id, int(ginz.Config.TokenExpiredTime), "/", "*", true, false)
		data = map[string]any{"Token": token}
		return
	}
}

// 用户注册接口
func UserRegisterHandler() ginz.ApiHandler {
	return func(c *gin.Context) (data any, err *errs.Err) {
		var req AuthRequestBody

		if er := c.ShouldBindJSON(&req); er != nil {
			err = errs.ErrInvalidParams.Add(er.Error())
			return
		}

		if req.Username == "" {
			err = errs.ErrInvalidParams.Add("username")
			return
		}

		if req.Password == "" {
			err = errs.ErrInvalidParams.Add("password")
			return
		}

		existUser, _ := dao.QueryOneByMap[User](map[string]any{"username": req.Username})
		if existUser.Id != "" {
			err = errs.ErrUserAlreadyExists
			return
		}

		psdHash, er := secures.HashPassword(req.Password)
		if er != nil {
			err = errs.ErrHashPassword
			return
		}

		u := &User{
			Username:     req.Username,
			PasswordHash: psdHash,
		}

		if er = dao.CreateOne[User](u); er != nil {
			err = errs.ErrCreateUser.Add(er.Error())
			return
		}
		data = map[string]any{"Id": u.Id}
		return
	}
}
