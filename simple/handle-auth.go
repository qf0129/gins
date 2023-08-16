package simple

import (
	"github.com/qf0129/ginz"
	"github.com/qf0129/ginz/pkg/dao"
	"github.com/qf0129/ginz/pkg/errs"
	"github.com/qf0129/ginz/pkg/secures"
)

type AuthRequestBody struct {
	Username string `validate:"gt=2"`
	Password string `validate:"gt=2"`
}

// 用户登录接口
func UserLoginHandler() ginz.ApiHandler {
	return func(c *ginz.Context) {
		var req AuthRequestBody
		if err := c.ShouldBindJSON(&req); err != nil {
			c.ReturnErr(err)
			return
		}

		if err := c.Validate(&req); err != nil {
			c.ReturnErr(err)
			return
		}

		existUser, er := dao.QueryOneByMap[User](map[string]any{"username": req.Username})
		if er != nil {
			c.ReturnErr(errs.UserNotFound)
			return
		}

		if !secures.VerifyPassword(req.Password, existUser.PasswordHash) {
			c.ReturnErr(errs.IncorrectPassword)
			return
		}

		token, er := secures.CreateToken(existUser.Id, ginz.Config.Secret)
		if er != nil {
			c.ReturnErr(errs.CreateToken)
			return
		}
		c.C.SetCookie("tk", token, int(ginz.Config.TokenExpiredTime), "/", "*", true, true)
		c.C.SetCookie("uid", existUser.Id, int(ginz.Config.TokenExpiredTime), "/", "*", true, false)
		c.ReturnOk(map[string]any{"Token": token})
	}
}

// 用户注册接口
func UserRegisterHandler() ginz.ApiHandler {
	return func(c *ginz.Context) {
		var req AuthRequestBody
		if err := c.ShouldBindJSON(&req); err != nil {
			c.ReturnErr(err)
			return
		}

		if err := c.Validate(&req); err != nil {
			c.ReturnErr(err)
			return
		}

		existUser, _ := dao.QueryOneByMap[User](map[string]any{"username": req.Username})
		if existUser.Id != "" {
			c.ReturnErr(errs.UserAlreadyExists)
			return
		}

		psdHash, er := secures.HashPassword(req.Password)
		if er != nil {
			c.ReturnErr(errs.HashPassword)
			return
		}

		u := &User{
			Username:     req.Username,
			PasswordHash: psdHash,
		}

		if er = dao.CreateOne[User](u); er != nil {
			c.ReturnErr(errs.CreateUser.Add(er.Error()))
			return
		}
		c.ReturnOk(map[string]any{"Id": u.Id})
	}
}
