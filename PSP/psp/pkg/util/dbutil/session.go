package dbutil

import (
	"xorm.io/xorm"

	"github.com/yuansuan/ticp/PSP/psp/pkg/util/strutil"
	"github.com/yuansuan/ticp/PSP/psp/pkg/xtype"
)

// WrapSortSession 排序会话
func WrapSortSession(session *xorm.Session, orderSort *xtype.OrderSort) *xorm.Session {
	defaultOrderBy := "create_time"

	if orderSort == nil {
		orderSort = &xtype.OrderSort{
			OrderBy:   defaultOrderBy,
			SortByAsc: false,
		}
	}

	if strutil.IsEmpty(orderSort.OrderBy) {
		orderSort.OrderBy = defaultOrderBy
	}

	if orderSort.SortByAsc {
		session.Asc(orderSort.OrderBy)
	} else {
		session.Desc(orderSort.OrderBy)
	}

	return session
}

// WrapSortSessionWithTable 排序会话
func WrapSortSessionWithTable(session *xorm.Session, table string, orderSort *xtype.OrderSort) *xorm.Session {
	defaultOrderBy := table + ".create_time"

	if orderSort == nil {
		orderSort = &xtype.OrderSort{
			OrderBy:   defaultOrderBy,
			SortByAsc: false,
		}
	}

	if strutil.IsEmpty(orderSort.OrderBy) {
		orderSort.OrderBy = defaultOrderBy
	}

	if orderSort.SortByAsc {
		session.Asc(orderSort.OrderBy)
	} else {
		session.Desc(orderSort.OrderBy)
	}

	return session
}
