package casbin

import (
	gofmt "fmt"
	"strconv"
	"strings"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/rbac"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
)

const (
	rolePrefix = "r"
)

// format

type tFmt struct{}

var fmt tFmt

// RoleID RoleID
func (tFmt) RoleID(roleID int64) string {
	return gofmt.Sprintf("%s%d", rolePrefix, roleID)
}

// UserID UserID
func (tFmt) UserID(userID int64) string {
	return gofmt.Sprintf("%d", userID)
}

// PermID PermID
func (tFmt) PermID(permID int64) string {
	return gofmt.Sprintf("%d", permID)
}

// ObjectID ObjectID
func (tFmt) ObjectID(id *rbac.ObjectID) string {
	switch id.Type {
	case rbac.ObjectType_USER:
		return fmt.UserID(snowflake.MustParseString(id.Id).Int64())
	}
	return ""
}

// unformat

type tUnfmt struct{}

var unfmt tUnfmt

// RoleID RoleID
func (tUnfmt) RoleID(roleID string) int64 {
	result, err := strconv.ParseInt(strings.TrimPrefix(roleID, rolePrefix), 10, 64)
	if err != nil {
		panic(gofmt.Errorf("roleID in unknow format!!! roleID=%v", roleID))
	}
	return result
}

// UserID UserID
func (tUnfmt) UserID(userID string) int64 {
	result, err := snowflake.ParseString(userID)
	if err != nil {
		panic(gofmt.Errorf("userID in unknow format!!! userID=%v", userID))
	}
	return result.Int64()
}

// PermID PermID
func (tUnfmt) PermID(permID string) int64 {
	result, err := strconv.ParseInt(permID, 10, 64)
	if err != nil {
		panic(gofmt.Errorf("permID in unknow format!!! permID=%v", permID))
	}
	return result
}

// ObjectID ObjectID
func (tUnfmt) ObjectID(id string) *rbac.ObjectID {

	if len(id) >= 19 {
		i, _ := strconv.ParseInt(id, 10, 64)
		return &rbac.ObjectID{Id: snowflake.ID(i).String(), Type: rbac.ObjectType_USER}
	}
	return &rbac.ObjectID{Id: id, Type: rbac.ObjectType_USER}
}

// is

type tIs struct{}

var is tIs
