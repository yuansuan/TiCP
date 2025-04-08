/*
 * Copyright (C) 2019 LambdaCal Inc.
 */

package util

import (
	"fmt"
	"strings"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/rbac"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/user"
	"github.com/yuansuan/ticp/PSP/psp/internal/rbac/consts"
	"github.com/yuansuan/ticp/PSP/psp/internal/rbac/dao/boring"
	"github.com/yuansuan/ticp/PSP/psp/internal/rbac/dao/model"
	"github.com/yuansuan/ticp/PSP/psp/internal/rbac/dto"
)

type toUserManagement struct{}

var ToUserManagement = toUserManagement{}

func (toUserManagement) UserID(id *rbac.ObjectID) *user.UserIdentity {
	return &user.UserIdentity{
		Id: id.Id,
	}
}

type fromUserManagement struct{}

var FromUserManagement = fromUserManagement{}

type toGRPC struct{}

var ToGRPC = toGRPC{}

func (toGRPC) Role(r *model.Role) *rbac.Role {
	if r == nil {
		return nil
	}

	return &rbac.Role{
		Id:      r.Id,
		Name:    r.Name,
		Comment: r.Comment,
		Type:    rbac.RoleType(r.Type),
	}
}

func (toGRPC) Resource(res *model.Resource) *rbac.Resource {
	if res == nil {
		return nil
	}

	return &rbac.Resource{
		Id:           res.Id,
		DisplayName:  res.DisplayName,
		ResourceType: res.Type,
		ResourceName: res.Name,
		ResourceId:   res.ExternalId,
		Action:       res.Action,
		ParentId:     res.ParentId,
	}
}

func (toGRPC) Roles(roles []*model.Role) []*rbac.Role {
	result := make([]*rbac.Role, 0, len(roles))
	for _, role := range roles {
		result = append(result, ToGRPC.Role(role))
	}
	return result
}

func (toGRPC) Resources(perms []*model.Resource) []*rbac.Resource {
	result := make([]*rbac.Resource, 0, len(perms))
	for _, perm := range perms {
		result = append(result, ToGRPC.Resource(perm))
	}
	return result
}

type fromGRPC struct{}

var FromGPRC = fromGRPC{}

func (fromGRPC) Role(r *rbac.Role) *model.Role {
	if r == nil {
		return nil
	}

	return &model.Role{
		Id:      r.Id,
		Name:    r.Name,
		Comment: r.Comment,
		Type:    int32(r.Type),
	}
}

func (fromGRPC) Resource(res *rbac.Resource) *model.Resource {
	if res == nil {
		return nil
	}

	custom := int8(1)
	if res.Custom == -1 {
		custom = -1
	}

	return &model.Resource{
		Id:          res.Id,
		Name:        res.ResourceName,
		DisplayName: res.DisplayName,
		Type:        res.ResourceType,
		ExternalId:  res.ResourceId,
		Action:      res.Action,
		Custom:      custom,
		ParentId:    res.ParentId,
	}
}

func (fromGRPC) ListQuery(q *rbac.ListQuery) *boring.ListRequest {
	if q == nil {
		return nil
	}

	return &boring.ListRequest{
		NameFilter: q.NameFilter,
		Page:       q.Page,
		PageSize:   q.PageSize,
		Desc:       q.Desc,
		OrderBy:    q.OrderBy,
	}
}

type toDTO struct{}

var ToDTO = toDTO{}

func (toDTO) RoleInfos(roles []*model.Role, defRoleID int64) []*dto.RoleInfo {
	result := make([]*dto.RoleInfo, 0, len(roles))
	for _, role := range roles {
		result = append(result, ToDTO.RoleInfo(role, defRoleID))
	}
	return result
}

func (toDTO) RoleInfo(r *model.Role, defRoleID int64) *dto.RoleInfo {
	if r == nil {
		return nil
	}

	info := &dto.RoleInfo{
		Id:      r.Id,
		Name:    r.Name,
		Comment: r.Comment,
		Type:    r.Type,
	}

	info.IsInternal = r.Type == consts.RoleTypeSuperAdmin
	info.IsDefault = r.Id == defRoleID
	return info
}

func (toDTO) Resources(resList []*model.Resource) []*dto.Resource {

	result := make([]*dto.Resource, 0, len(resList))
	for _, res := range resList {
		result = append(result, ToDTO.Resource(res))
	}
	return result
}

func (toDTO) Resource(res *model.Resource) *dto.Resource {
	if res == nil {
		return nil
	}

	return &dto.Resource{
		Id:           res.Id,
		DisplayName:  res.DisplayName,
		ResourceType: res.Type,
		ResourceName: res.Name,
		ExternalId:   res.ExternalId,
		Action:       res.Action,
		ParentId:     res.ParentId,
	}
}

type fromDTO struct{}

var FromDTO = fromDTO{}

func (fromDTO) Role(r *dto.RoleInfo) *model.Role {
	if r == nil {
		return nil
	}

	return &model.Role{
		Id:      r.Id,
		Name:    r.Name,
		Comment: r.Comment,
		Type:    r.Type,
	}
}

func PermLog(permList []string, typeName string) string {
	return fmt.Sprintf("%s:[%s]", typeName, strings.Join(permList, ","))
}
