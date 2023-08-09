package ginz

import (
	"encoding/json"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/qf0129/ginz/pkg/dao"
	"github.com/qf0129/ginz/pkg/errs"
)

func CreateOneHandler[T dao.GormModel](parentIdKeys ...string) ApiHandler {
	return func(c *gin.Context) (data any, err *errs.Err) {
		var er error
		var model T

		if er = c.ShouldBindJSON(&model); err != nil {
			err = errs.ErrInvalidParams.Add(er.Error())
			return
		}
		if len(parentIdKeys) > 0 {
			er = dao.CreateOneWithParentId[T](model, parentIdKeys[0], c.Param(parentIdKeys[0]))
		} else {
			er = dao.CreateOne[T](&model)
		}

		if er != nil {
			err = errs.ErrCreateDataFailed.Add(er.Error())
			return
		}
		data = model
		return
	}
}

func DeleteOneHandler[T dao.GormModel](parentIdKey string) ApiHandler {
	return func(c *gin.Context) (data any, err *errs.Err) {
		if er := dao.DeleteOneByPk[T](c.Param(parentIdKey)); er != nil {
			if errMySQL, ok := er.(*mysql.MySQLError); ok {
				switch errMySQL.Number {
				case 1451:
					err = errs.ErrDeleteDataFailed.Add("无法删除有关联数据的项")
					return
				}
			} else {
				err = errs.ErrDeleteDataFailed.Add(er.Error())
				return
			}
		}
		data = parentIdKey
		return
	}
}

func UpdateOneHandler[T dao.GormModel](parentIdKey string) ApiHandler {
	return func(c *gin.Context) (data any, err *errs.Err) {
		if _, er := dao.QueryOneByPk[T](c.Param(parentIdKey)); er != nil {
			err = errs.ErrNotExistsData
			return
		}

		var objMap map[string]any
		if er := c.ShouldBindJSON(&objMap); er != nil {
			err = errs.ErrInvalidParams.Add(er.Error())
			return
		}

		// gorm中updates结构体不支持更新空值，使用map不支持json类型
		// 因此遍历map，将子结构的map或slice转成json字符串
		for k, v := range objMap {
			valKind := reflect.ValueOf(v).Kind()
			if valKind == reflect.Map || valKind == reflect.Slice {
				bytes, er := json.Marshal(v)
				if er != nil {
					err = errs.ErrInvalidParams.Add(er.Error())
					return
				}
				objMap[k] = string(bytes)
			}
		}

		er := dao.UpdateOneByPk[T](c.Param(parentIdKey), &objMap)
		if er != nil {
			err = errs.ErrUpdateDataFailed.Add(er.Error())
			return
		}

		newModel, er := dao.QueryOneByPk[T](c.Param(parentIdKey))
		if er != nil {
			err = errs.ErrRetrieveDataFailed.Add(er.Error())
			return
		}
		data = newModel
		return
	}
}

// func QueryOneHandler[T dao.GormModel](parentIdKey string) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		ret := QueryOneByPkWithPreload[T](c.Param(parentIdKey), c.Query(FixedKeyPreload))
// 		goo.RespOk(c, ret)
// 	}
// }

func QueryManyHandler[T dao.GormModel](parentIdKeys ...string) ApiHandler {
	return func(c *gin.Context) (data any, err *errs.Err) {
		var fixedOptions *dao.FixedOption
		er := c.ShouldBind(&fixedOptions)
		if er != nil {
			err = errs.ErrInvalidParams.Add(er.Error())
			return
		}

		filtersMap := dao.ParseUrlFilters(c.Request.URL.Query())
		if len(parentIdKeys) > 0 {
			filtersMap[parentIdKeys[0]] = c.Param(parentIdKeys[0])
		}

		if fixedOptions.ClosePaging {
			data, er = dao.QueryAll[T](filtersMap, fixedOptions.Preload)
		} else {
			data, er = dao.QueryPage[T](filtersMap, fixedOptions)
		}

		if er != nil {
			err = errs.ErrInvalidParams.Add(er.Error())
			return
		}
		return
	}
}
