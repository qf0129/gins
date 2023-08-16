package main

import (
	"github.com/qf0129/ginz"
	"github.com/qf0129/ginz/pkg/dao"
	"github.com/qf0129/ginz/pkg/errs"
	"github.com/qf0129/ginz/simple"
	"gorm.io/gorm/logger"
)

var App *ginz.Ginz

func main() {
	App = ginz.Init(&ginz.Option{
		LoadConfigFile:    true,
		ConnectDB:         true,
		AddHealthCheckApi: true,
		DBLogLevel:        logger.Error,
		Models:            []any{&simple.User{}},
		// Middlewares:       []ginz.Middleware{simple.RequireTokenFromCookie()},
	})

	group1 := App.Group("/api")
	group1.AddApi("login", simple.UserLoginHandler())
	group1.AddApi("register", simple.UserRegisterHandler())
	group1.AddApi("CreateUser", simple.CreateModelHandler[simple.User]())
	group1.AddApi("QueryUser", simple.QueryModelHandler[simple.User]())
	group1.AddApi("UpdateUser", simple.UpdateModelHandler[simple.User]())
	group1.AddApi("DeleteUser", simple.DeleteModelHandler[simple.User]())

	group2 := App.Group("/api")
	group2.Use(simple.RequireTokenFromCookie())
	group2.AddApi("test1", func(c *ginz.Context) {
		data, er := dao.QueryAll[simple.User](nil)
		if er != nil {
			c.ReturnErr(errs.RetrieveDataFailed.Add(er.Error()))
			return
		}
		c.ReturnOk(data)
	})
	App.Run()
}
