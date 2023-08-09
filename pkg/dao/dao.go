package dao

import (
	"gorm.io/gorm"
)

type GormModel any

const FixedKeyPage = "Page"
const FixedKeyPageSize = "PageSize"
const FixedKeyPreload = "Preload"
const FixedKeyClosePaging = "ClosePaging"

var FIXED_KEYS = []string{FixedKeyPage, FixedKeyPageSize, FixedKeyPreload, FixedKeyClosePaging}

// 固定查询选项
type FixedOption struct {
	Page        int    // 页数，默认1
	PageSize    int    // 每页数量，默认10
	Preload     string // 预加载关联表名，若多个以英文逗号分隔
	ClosePaging bool   // 关闭分页，默认false
}

type PageBody[T any] struct {
	List     []T
	Page     int
	PageSize int
	Total    int64
}

var (
	DB              *gorm.DB
	DefaultPageSize = 10
	QueryPrimaryKey = "id"
)

type DaoOption struct {
	DB              *gorm.DB
	DefaultPageSize int
	QueryPrimaryKey string
}

func Init(option *DaoOption) {
	DB = option.DB
	DefaultPageSize = option.DefaultPageSize
	QueryPrimaryKey = option.QueryPrimaryKey
}
