package ginz

import (
	"github.com/gin-gonic/gin"
	"github.com/qf0129/ginz/crud"
)

func CreateOneHandler[T crud.GormModel](parentIdKeys ...string) ApiHandler {
	return func(c *gin.Context) (data any, err *Errors) {
		var er error
		var model T

		if er = c.ShouldBindJSON(&model); err != nil {
			err = ErrParams.AddMsg(er.Error())
			return
		}
		if len(parentIdKeys) > 0 {
			er = crud.CreateOneWithParentId[T](model, parentIdKeys[0], c.Param(parentIdKeys[0]))
		} else {
			// er = DB.Model(new(T)).Create(model).Error
			er = crud.CreateOne[T](&model)
		}

		if er != nil {
			err = ErrParams.AddMsg(er.Error())
			return
		}
		data = model
		return
	}
}

// func DeleteOneHandler[T crud.GormModel](parentIdKey string) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		err := DeleteOneByPk[T](c.Param(parentIdKey))
// 		if err != nil {
// 			if errMySQL, ok := err.(*mysql.MySQLError); ok {
// 				switch errMySQL.Number {
// 				case 1451:
// 					goo.RespFail(c, "无法删除有关联数据的项")
// 					return
// 				}
// 			} else {
// 				goo.RespFail(c, "DeleteOneFailed, "+err.Error())
// 				return
// 			}
// 		}
// 		goo.RespOk(c, true)
// 	}
// }

// func UpdateOneHandler[T crud.GormModel](parentIdKey string) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		if !ExistByPk[T](c.Param(parentIdKey)) {
// 			goo.RespFail(c, "FindNotItem")
// 			return
// 		}

// 		var objMap map[string]any
// 		if err := c.ShouldBindJSON(&objMap); err != nil {
// 			goo.RespFail(c, "InvalidData, "+err.Error())
// 			return
// 		}

// 		// gorm中updates结构体不支持更新空值，使用map不支持json类型
// 		// 因此遍历map，将子结构的map或slice转成json字符串
// 		for k, v := range objMap {
// 			valKind := reflect.ValueOf(v).Kind()
// 			if valKind == reflect.Map || valKind == reflect.Slice {
// 				bytes, err := json.Marshal(v)
// 				if err != nil {
// 					goo.RespFail(c, "InvalidJsonValue, "+err.Error())
// 					return
// 				}
// 				objMap[k] = string(bytes)
// 			}
// 		}

// 		err := UpdateOneByPk[T](c.Param(parentIdKey), &objMap)
// 		if err != nil {
// 			goo.RespFail(c, "UpdateOneFailed, "+err.Error())
// 			return
// 		}

// 		newModel := QueryOneByPk[T](c.Param(parentIdKey))
// 		goo.RespOk(c, &newModel)
// 	}
// }

// func QueryOneHandler[T crud.GormModel](parentIdKey string) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		ret := QueryOneByPkWithPreload[T](c.Param(parentIdKey), c.Query(FixedKeyPreload))
// 		goo.RespOk(c, ret)
// 	}
// }

// func QueryManyHandler[T crud.GormModel](parentIdKeys ...string) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		var fixedOptions FixedOption
// 		err := c.ShouldBind(&fixedOptions)
// 		if err != nil {
// 			goo.RespFail(c, "FixedOptionError, "+err.Error())
// 			return
// 		}

// 		filtersMap := GetUrlFilters(c.Request.URL.Query())
// 		if len(parentIdKeys) > 0 {
// 			filtersMap[parentIdKeys[0]] = c.Param(parentIdKeys[0])
// 		}

// 		var ret any
// 		if fixedOptions.ClosePaging {
// 			ret = QueryAll[T](&fixedOptions, filtersMap)
// 		} else {
// 			ret = QueryPage[T](&fixedOptions, filtersMap)
// 		}
// 		goo.RespOk(c, ret)
// 	}
// }
