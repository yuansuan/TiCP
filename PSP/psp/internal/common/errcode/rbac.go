package errcode

import (
	"google.golang.org/grpc/codes"
)

// RBAC service error codes from 14001 to 15000
const (
	NoPermission     = "no permission"
	NoSuchPermission = "no such permission"
	RBACUnknown      = "unknown"

	// db error
	ErrRBACInsertError        codes.Code = 14001
	ErrRBACQueryError         codes.Code = 14002
	ErrRBACDeleteError        codes.Code = 14003
	ErrRBACUpdateError        codes.Code = 14004
	ErrRBACInvalidOrderError  codes.Code = 14005
	ErrRBACPageShouldPositive codes.Code = 14006

	// business error
	ErrRBACRoleNotFound              codes.Code = 14010
	ErrRBACPermissionNotFound        codes.Code = 14011
	ErrRBACRoleAlreadyExist          codes.Code = 14012
	ErrRBACPermissionAlreadyExist    codes.Code = 14013
	ErrRBACNotAllowSuperAdminRole    codes.Code = 14014
	ErrRBACFailedGetUserID           codes.Code = 14015
	ErrRBACNotAllowOpSuperAdminRole  codes.Code = 14016
	ErrRBACCantRemoveInternalPerm    codes.Code = 14017
	ErrRBACInternalPermOnlyGiveAdmin codes.Code = 14018
	ErrRBACCantAddInternalPerm       codes.Code = 14019
	ErrRBACChildPerm                 codes.Code = 14020
	ErrRBACNoPermission              codes.Code = 14021
	ErrRBACUnknownResourceIdentify   codes.Code = 14022
	ErrRBACUnknown                   codes.Code = 14023
	ErrRBACAddRoleError              codes.Code = 14024
	ErrRBACAddRolePerm               codes.Code = 14025
	ErrRBACGetRole                   codes.Code = 14026
	ErrRBACTokenInvalid              codes.Code = 14027
	ErrRBACRoleUsed                  codes.Code = 14028
	ErrRBACRoleUsedByLdap            codes.Code = 14029
	ErrRBACSetSuperAdmin             codes.Code = 14030
	ErrRBACRoleNameExist             codes.Code = 14031
)

var RBACCodeMsg = map[codes.Code]string{
	ErrRBACAddRoleError:              "添加用户失败",
	ErrRBACNotAllowOpSuperAdminRole:  "不能修改超级管理员角色",
	ErrRBACInternalPermOnlyGiveAdmin: "不能操作内部权限",
	ErrRBACPermissionNotFound:        "所选权限不存在",
	ErrRBACNotAllowSuperAdminRole:    "不能删除超级管理员角色",
	ErrRBACAddRolePerm:               "添加角色权限失败",
	ErrRBACQueryError:                "角色查询失败",
	ErrRBACGetRole:                   "角色详情查询失败",
	ErrRBACUpdateError:               "修改角色失败",
	ErrRBACDeleteError:               "删除角色失败",
	ErrRBACRoleUsed:                  "该角色正在被使用，请解除与用户的绑定后再删除",
	ErrRBACRoleUsedByLdap:            "该角色被设置为ldap用户默认角色，请修改设置后再删除",
	ErrRBACRoleNotFound:              "角色不存在",
	ErrRBACSetSuperAdmin:             "不能将超级管理员设置为默认角色",
	ErrRBACRoleNameExist:             "角色名已存在",
}
