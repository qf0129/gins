package simple

import (
	"encoding/json"
	"reflect"

	"github.com/go-sql-driver/mysql"
	"github.com/qf0129/ginz"
	"github.com/qf0129/ginz/pkg/arrays"
	"github.com/qf0129/ginz/pkg/dao"
	"github.com/qf0129/ginz/pkg/errs"
	"github.com/sirupsen/logrus"
)

func CreateCrudRouter[T dao.GormModel](group *ginz.ApiGroup, methods ...string) {
	modelName := reflect.TypeOf(new(T)).Elem().Name()
	if modelName == "" {
		logrus.Fatalf("InvalidModelName: %s", modelName)
	}
	if len(methods) == 0 {
		methods = []string{"c", "r", "u", "d"}
	}
	if arrays.HasStrItem(methods, "c") {
		group.AddApi("Create"+modelName, CreateModelHandler[T]())
	}
	if arrays.HasStrItem(methods, "r") {
		group.AddApi("Query"+modelName, QueryModelHandler[T]())
	}
	if arrays.HasStrItem(methods, "u") {
		group.AddApi("Update"+modelName, UpdateModelHandler[T]())
	}
	if arrays.HasStrItem(methods, "d") {
		group.AddApi("Delete"+modelName, DeleteModelHandler[T]())
	}
}

func QueryModelHandler[T dao.GormModel]() ginz.ApiHandler {
	return func(c *ginz.Context) {
		var fixedOptions *dao.FixedOption
		err := c.ShouldBindJSON(&fixedOptions)
		if err != nil {
			c.ReturnErr(err)
			return
		}
		if fixedOptions.ClosePaging {
			data, er := dao.QueryAll[T](fixedOptions.Filters, fixedOptions.Preload)
			if er != nil {
				c.ReturnErr(errs.RetrieveDataFailed.Add(er.Error()))
				return
			}
			c.ReturnOk(data)
		} else {
			data, er := dao.QueryPage[T](fixedOptions.Filters, fixedOptions)
			if er != nil {
				c.ReturnErr(errs.RetrieveDataFailed.Add(er.Error()))
				return
			}
			c.ReturnOk(data)
		}
	}
}

func CreateModelHandler[T dao.GormModel]() ginz.ApiHandler {
	return func(c *ginz.Context) {
		var model T
		if er := c.ShouldBindJSON(&model); er != nil {
			c.ReturnErr(er)
			return
		}

		er := dao.CreateOne[T](&model)
		if er != nil {
			c.ReturnErr(errs.CreateDataFailed.Add(er.Error()))
			return
		}
		c.ReturnOk(model)
	}
}

func DeleteModelHandler[T dao.GormModel]() ginz.ApiHandler {
	return func(c *ginz.Context) {
		body := make(map[string]any)
		if err := c.ShouldBindJSON(&body); err != nil {
			c.ReturnErr(err)
			return
		}
		pk, ok := body[ginz.Config.QueryPrimaryKey]
		if !ok || pk == "" {
			c.ReturnErr(errs.InvalidParams.Add("主键为空"))
			return
		}
		if er := dao.DeleteOneByPk[T](pk); er != nil {
			if errMySQL, ok := er.(*mysql.MySQLError); ok {
				switch errMySQL.Number {
				case 1451:
					c.ReturnErr(errs.DeleteDataFailed.Add("无法删除有关联数据的项"))
					return
				}
			} else {
				c.ReturnErr(errs.DeleteDataFailed.Add(er.Error()))
				return
			}
		}
		c.ReturnOk(pk)
	}
}

func UpdateModelHandler[T dao.GormModel]() ginz.ApiHandler {
	return func(c *ginz.Context) {
		body := make(map[string]any)
		if err := c.ShouldBindJSON(&body); err != nil {
			c.ReturnErr(err)
			return
		}
		pk, ok := body[ginz.Config.QueryPrimaryKey]
		if !ok || pk == "" {
			c.ReturnErr(errs.InvalidParams.Add("主键为空"))
			return
		}
		if _, er := dao.QueryOneByPk[T](pk); er != nil {
			c.ReturnErr(errs.DataNotExists)
			return
		}

		// gorm中updates结构体不支持更新空值，使用map不支持json类型
		// 因此遍历map，将子结构的map或slice转成json字符串
		for k, v := range body {
			valKind := reflect.ValueOf(v).Kind()
			if valKind == reflect.Map || valKind == reflect.Slice {
				bytes, er := json.Marshal(v)
				if er != nil {
					c.ReturnErr(errs.InvalidParams.Add(er.Error()))
					return
				}
				body[k] = string(bytes)
			}
		}

		er := dao.UpdateOneByPk[T](pk, &body)
		if er != nil {
			c.ReturnErr(errs.UpdateDataFailed.Add(er.Error()))
			return
		}

		newModel, er := dao.QueryOneByPk[T](pk)
		if er != nil {
			c.ReturnErr(errs.RetrieveDataFailed.Add(er.Error()))
			return
		}
		c.ReturnOk(newModel)
	}
}
