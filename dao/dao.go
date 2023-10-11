package dao

import (
	"strings"
)

type GormModel any

const FixedKeyPage = "Page"
const FixedKeyPageSize = "PageSize"
const FixedKeyPreload = "Preload"
const FixedKeyClosePaging = "ClosePaging"

var FIXED_KEYS = []string{FixedKeyPage, FixedKeyPageSize, FixedKeyPreload, FixedKeyClosePaging}

// 查询结构体
type QueryBody struct {
	PageNum   int    // 页数，默认1
	PageSize  int    // 每页数量，默认10
	Preload   string // 预加载关联表名，若多个以英文逗号分隔
	NoPaging  bool   // 关闭分页，默认false
	Filter    string
	FilterMap map[string]any
}

func (query *QueryBody) ParseFilterToMap() {
	if query.FilterMap == nil {
		query.FilterMap = map[string]any{}
	}
	if query.Filter != "" {
		filterList := strings.Split(query.Filter, "|")
		for _, filter := range filterList {
			items := strings.Split(filter, ":")
			if len(items) == 2 {
				query.FilterMap[items[0]] = items[1]
			} else if len(items) >= 3 {
				query.FilterMap[items[0]+":"+items[1]] = items[2]
			}
		}
	}
}

type PageBody[T any] struct {
	List     []T
	PageNum  int
	PageSize int
	Total    int64
}
