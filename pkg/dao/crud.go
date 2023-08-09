package dao

import (
	"reflect"
	"strings"
)

// 查询分页数据
func QueryPage[T GormModel](filters map[string]any, fixedOptions ...*FixedOption) (result PageBody[T], err error) {
	query := DB.Model(new(T))
	for _, fc := range ParseFilters(filters) {
		query = fc(query)
	}
	if err = query.Count(&result.Total).Error; err != nil {
		return
	}

	var fixedOption *FixedOption
	if len(fixedOptions) > 0 {
		fixedOption = fixedOptions[0]
	} else {
		fixedOption = &FixedOption{
			Page:     1,
			PageSize: 10,
		}
	}

	for _, preload := range strings.Split(fixedOption.Preload, ",") {
		query = FiltePreloadFunc(preload)(query)
	}
	if fixedOption.Page < 1 {
		fixedOption.Page = 1
	}
	if fixedOption.PageSize < 1 {
		fixedOption.PageSize = DefaultPageSize
	}

	result.Page = fixedOption.Page
	result.PageSize = fixedOption.PageSize
	query = FiltePageFunc(result.Page, result.PageSize)(query)
	err = query.Find(&result.List).Error
	return
}

func QueryAll[T GormModel](filters map[string]any, preloads ...string) (result []T, err error) {
	query := DB.Model(new(T))
	for _, fc := range ParseFilters(filters) {
		query = fc(query)
	}
	for _, preloadStr := range preloads {
		for _, preload := range strings.Split(preloadStr, ",") {
			query = FiltePreloadFunc(preload)(query)
		}
	}
	err = query.Find(&result).Error
	return
}

func ExistByPk[T GormModel](pk any) (err error) {
	item := new(T)
	return DB.Model(new(T)).Where(map[string]any{QueryPrimaryKey: pk}).First(&item).Error
}

func QueryOneByPk[T GormModel](pk any) (result T, err error) {
	err = DB.Model(new(T)).Where(map[string]any{QueryPrimaryKey: pk}).First(&result).Error
	return
}

func QueryOneByPkWithPreload[T GormModel](pk any, preload string) (result T, err error) {
	query := DB.Model(new(T)).Where(map[string]any{QueryPrimaryKey: pk})
	for _, field := range strings.Split(preload, ",") {
		query = FiltePreloadFunc(field)(query)
	}
	err = query.Take(&result).Error
	return
}

func QueryOneByMap[T GormModel](filters map[string]any) (result T, err error) {
	query := DB.Model(new(T))
	for _, fc := range ParseFilters(filters) {
		query = fc(query)
	}
	err = query.First(&result).Error
	return
}

func QueryOneByMapWithPreload[T GormModel](filters map[string]any, preload string) (result T, err error) {
	query := DB.Model(new(T))
	for _, fc := range ParseFilters(filters) {
		query = fc(query)
	}

	for _, field := range strings.Split(preload, ",") {
		query = FiltePreloadFunc(field)(query)
	}

	err = query.First(&result).Error
	return
}

func CreateOne[T GormModel](obj any) error {
	return DB.Model(new(T)).Create(obj).Error
}

func CreateOneWithParentId[T GormModel](obj any, parentIdKey string, parentIdVal string) error {
	types := reflect.TypeOf(obj)
	vals := reflect.ValueOf(obj).Elem()
	for i := 0; i < types.NumField(); i++ {
		if types.Field(i).Name == parentIdKey {
			vals.Field(i).Set(reflect.ValueOf(parentIdVal))
		}
	}
	return DB.Model(new(T)).Create(&obj).Error
}

func UpdateOneByPk[T GormModel](pk any, data any) error {
	return DB.Model(new(T)).Where(map[string]any{QueryPrimaryKey: pk}).Updates(data).Error
}

func DeleteOneByPk[T GormModel](pk any) error {
	return DB.Where(map[string]any{QueryPrimaryKey: pk}).Delete(new(T)).Error
}

func HasField[T GormModel](field string) bool {
	return DB.Model(new(T)).Select(field).Take(new(T)).Error == nil
}
