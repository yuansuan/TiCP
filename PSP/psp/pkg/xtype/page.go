package xtype

import (
	"errors"
	"fmt"
)

const (
	// MinPageIndex 最小分页Index是1
	MinPageIndex = int64(1)

	// MaxPageSize 每页最大容量是1000
	MaxPageSize = int64(1000)
)

// Page 分页请求数据
type Page struct {
	// 页数索引，范围大于等于1
	Index int64 `json:"index"`
	// 每页大小，范围[1, 1000]
	Size int64 `json:"size"`
}

// PageResp 分页响应数据
type PageResp struct {
	// 页数索引
	Index int64 `json:"index"`
	// 每页大小
	Size int64 `json:"size"`
	// 数据总量
	Total int64 `json:"total"`
}

// OrderSort 排序数据
type OrderSort struct {
	OrderBy   string `json:"order_by" form:"order_by"`       // 顺序字段
	SortByAsc bool   `json:"sort_by_asc" form:"sort_by_asc"` // 倒序字段
}

// GetPageOffset 获取分页偏移量
func GetPageOffset(index, size int64) (int64, error) {
	if index < MinPageIndex {
		return 0, errors.New(fmt.Sprintf("page index cannot less than '%v'", MinPageIndex))
	}

	return (index - 1) * size, nil
}

// CheckPageValid 检查分页信息是否合法
func CheckPageValid(page Page) bool {
	if page.Index < MinPageIndex || page.Size > MaxPageSize {
		return false
	}

	return true
}
