/*
 * Copyright (C) 2019 LambdaCal Inc.
 */

package boring

import (
	"context"
	"fmt"
	"reflect"

	"github.com/go-sql-driver/mysql"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/middleware"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/collection"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Codes Codes
type Codes struct {
	EntityNotFound     codes.Code
	EntityAlreadyExist codes.Code
	InsertError        codes.Code
	QueryError         codes.Code
	DeleteError        codes.Code
	UpdateError        codes.Code
	InvalidOrderError  codes.Code
	PageShouldPositive codes.Code
}

// ListRequest ListRequest
type ListRequest struct {
	NameFilter string
	Page       int64
	PageSize   int64
	Desc       bool
	OrderBy    string
}

// Dao is use for :
//
//	id column is `id`
//	id is int64
type Dao struct {
	Bean            reflect.Type
	Codes           Codes
	NameFilterField string
	UpdateOmit      []string
	AllowOrderBy    map[string]bool
}

// NewBoringDao NewBoringDao
func NewBoringDao(bean interface{}, codes Codes, nameFilterField string, updateOmit []string, allowOrderBy []string) *Dao {
	mapAllowOrderBy := map[string]bool{}
	for _, orderBy := range allowOrderBy {
		mapAllowOrderBy[orderBy] = true
	}

	return &Dao{
		Bean:            reflect.TypeOf(bean).Elem(),
		Codes:           codes,
		NameFilterField: nameFilterField,
		UpdateOmit:      updateOmit,
		AllowOrderBy:    mapAllowOrderBy,
	}
}

func isDuplicate(err error) bool {
	me, ok := err.(*mysql.MySQLError)
	if !ok {
		return false
	}
	return me.Number == 1062
}

// Add Add
// result is *models.Bean
func (dao *Dao) Add(ctx context.Context, result interface{}) (err error) {
	session := middleware.DefaultSession(ctx)
	defer session.Close()
	_, err = session.Omit("id").Insert(result)
	if err != nil {
		if isDuplicate(err) {
			return status.Error(dao.Codes.EntityAlreadyExist, err.Error())
		}
		return status.Error(dao.Codes.InsertError, err.Error())
	}
	return nil
}

// Get Get
// result is *models.Bean
func (dao *Dao) Get(ctx context.Context, id int64, result interface{}) (err error) {
	session := middleware.DefaultSession(ctx)
	defer session.Close()

	ok, err := session.ID(id).Get(result)
	if err != nil {
		return status.Error(dao.Codes.QueryError, err.Error())
	}

	if !ok {
		return status.Error(dao.Codes.EntityNotFound, "")
	}

	return nil
}

// GetByName GetByName
// result is *models.Bean
func (dao *Dao) GetByName(ctx context.Context, name string, result interface{}) (err error) {
	session := middleware.DefaultSession(ctx)
	defer session.Close()

	ok, err := session.Where(fmt.Sprintf("%s = ?", dao.NameFilterField), name).Get(result)
	if err != nil {
		return status.Error(dao.Codes.QueryError, err.Error())
	}

	if !ok {
		return status.Error(dao.Codes.EntityNotFound, "")
	}

	return nil
}

// Gets Gets
// `result` is *[]models.Bean or *[]*models.Bean
func (dao *Dao) Gets(ctx context.Context, ids []int64, result interface{}) (err error) {
	session := middleware.DefaultSession(ctx)
	defer session.Close()

	count, err := session.In("id", ids).FindAndCount(result)
	if err != nil {
		return status.Error(dao.Codes.QueryError, err.Error())
	}
	if count != int64(len(ids)) {
		return status.Error(dao.Codes.EntityNotFound, "")
	}

	return
}

// List List
// `result` is *[]models.Bean or *[]*models.Bean
func (dao *Dao) List(ctx context.Context, listReq *ListRequest, result interface{}) (total int64, err error) {
	sql := middleware.DefaultSession(ctx)
	defer sql.Close()

	if listReq.OrderBy != "" && len(dao.AllowOrderBy) != 0 && !dao.AllowOrderBy[listReq.OrderBy] {
		return 0, status.Errorf(dao.Codes.InvalidOrderError, "unsupported order [%v], only support: %v",
			listReq.OrderBy, dao.AllowOrderBy)
	}

	if listReq.NameFilter != "" {
		sql = sql.Where(fmt.Sprintf("%s LIKE ?", dao.NameFilterField), fmt.Sprintf("%%%s%%", listReq.NameFilter))
	}

	// if pageSize <= 0, no limit
	if listReq.PageSize > 0 {
		if listReq.Page < 1 {
			return 0, status.Error(dao.Codes.PageShouldPositive, "")
		}
		sql = sql.Limit(int(listReq.PageSize), int(listReq.Page-1)*int(listReq.PageSize))
	}

	if listReq.OrderBy != "" {
		if listReq.Desc {
			sql.Desc(listReq.OrderBy)
		} else {
			sql.Asc(listReq.OrderBy)
		}
	}

	total, err = sql.FindAndCount(result)
	if err != nil {
		return total, status.Error(dao.Codes.QueryError, err.Error())
	}
	return total, nil
}

// Update Update will update by fields in `model`
// `result.id` will be ignored, won't update
// result is *models.Bean
func (dao *Dao) Update(ctx context.Context, id int64, result interface{}) error {
	session := middleware.DefaultSession(ctx)
	defer session.Close()

	bean := reflect.New(dao.Bean).Interface()
	count, err := session.ID(id).Count(bean)
	if err != nil {
		return status.Error(dao.Codes.QueryError, err.Error())
	}

	if count == 0 {
		return status.Error(dao.Codes.EntityNotFound, "")
	}

	_, err = session.ID(id).Omit(dao.UpdateOmit...).AllCols().Update(result)

	// `no content found to be updated` means all fields will not be updated
	if err != nil && err.Error() != "No content found to be updated" {
		if isDuplicate(err) {
			return status.Error(dao.Codes.EntityAlreadyExist, err.Error())
		}
		return status.Error(dao.Codes.UpdateError, err.Error())
	}

	return nil
}

// Delete Delete
func (dao *Dao) Delete(ctx context.Context, id int64) error {
	session := middleware.DefaultSession(ctx)
	defer session.Close()

	bean := reflect.New(dao.Bean).Interface()
	affected, err := session.ID(id).Delete(bean)
	if err != nil {
		return status.Error(dao.Codes.DeleteError, err.Error())
	}
	if affected == 0 {
		return status.Error(dao.Codes.EntityNotFound, "")
	}
	return nil
}

// AllExists AllExists
func (dao *Dao) AllExists(ctx context.Context, ids ...int64) error {
	if len(ids) == 0 {
		return nil
	}

	session := middleware.DefaultSession(ctx)
	defer session.Close()

	ids = collection.UniqueInt64Array(ids)

	bean := reflect.New(dao.Bean).Interface()
	count, err := session.In("id", ids).Count(bean)
	if err != nil {
		return status.Error(dao.Codes.QueryError, err.Error())
	}
	if count != int64(len(ids)) {
		return status.Error(dao.Codes.EntityNotFound, fmt.Sprintf("[%d/%d] not found", int64(len(ids))-count, len(ids)))
	}
	return nil
}
