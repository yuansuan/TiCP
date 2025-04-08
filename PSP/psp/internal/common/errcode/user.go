package errcode

import (
	"google.golang.org/grpc/codes"
)

// User service error codes from 16001 to 17000
const (
	ErrUserNotFound                   codes.Code = 16001
	ErrUserInvalidPassword            codes.Code = 16002
	ErrUserUpdatePasswordFailed       codes.Code = 16003
	ErrUserDeleted                    codes.Code = 16004
	ErrUserCreatedFailed              codes.Code = 16005
	ErrUserGetFailed                  codes.Code = 16006
	ErrUserUpdateFailed               codes.Code = 16007
	ErrUserNameExist                  codes.Code = 16008
	ErrUserNameEmpty                  codes.Code = 16009
	ErrUserCreateFailedAlreadyExist   codes.Code = 16010
	ErrUserQueryInvalidOrderCondition codes.Code = 16011
	ErrUserQueryFailed                codes.Code = 16012
	ErrUserCantRemoveInternalUser     codes.Code = 16013
	ErrUserActiveFailed               codes.Code = 16015
	ErrUserInActiveFailed             codes.Code = 16016
	ErrUserAddFailed                  codes.Code = 16017
	ErrAuthUserDisabled               codes.Code = 16018
	ErrAuthUserFailed                 codes.Code = 16019
	ErrUserNotExist                   codes.Code = 16020
	ErrOnlineListFailed               codes.Code = 16021
	ErrLDAPConnectFailed              codes.Code = 16022
	ErrUserLDAPUserInvalidPasswd      codes.Code = 16023
	ErrUserLDAPDefRoleFailed          codes.Code = 16024
	ErrUserInvalidOldPassword         codes.Code = 16025
	ErrUserOptionList                 codes.Code = 16026
	ErrUserNameCollideBuiltin         codes.Code = 16027
	ErrUserResetPass                  codes.Code = 16028

	ErrAuthGetMachineIDFailed          codes.Code = 16033
	ErrAuthGetLicenseInfoFailed        codes.Code = 16034
	ErrAuthUpdateLicenseInfoFailed     codes.Code = 16035
	ErrAuthLicenseHasExpiredOrNotExist codes.Code = 16036
	ErrAuthGetMachineIDEmpty           codes.Code = 16037
	ErrUserDefaultRoleNotExist         codes.Code = 16038
	ErrUserCantDeleteDefaultUser       codes.Code = 16039
	ErrOrgCreatedFailed                codes.Code = 16040
	ErrOrgUpdateFailed                 codes.Code = 16041
	ErrOrgDeleted                      codes.Code = 16042
	ErrOrgMemberAddFailed              codes.Code = 16043
	ErrOrgMemberDeletedFailed          codes.Code = 16044
	ErrOrgMemberUpdateFailed           codes.Code = 16045
	ErrOrgMemberListFailed             codes.Code = 16046

	ErrUserOpenapiCertAlreadyExist  codes.Code = 16050
	ErrUserOpenapiCertCreatedFailed codes.Code = 16051
	ErrUserOpenapiCertDeleteFailed  codes.Code = 16052
	ErrUserOpenapiCertGetFailed     codes.Code = 16053
	ErrUserOpenapiCertGenFailed     codes.Code = 16054
	ErrUserOpenapiCertDisable       codes.Code = 16055
	ErrUserOpenapiCertNotEnable     codes.Code = 16056
)

var UserCodeMsg = map[codes.Code]string{
	ErrAuthUserFailed:                  "账号有误，登录失败",
	ErrUserNameEmpty:                   "用户名不能为空",
	ErrAuthUserDisabled:                "用户未启用",
	ErrUserNameExist:                   "用户名已存在",
	ErrUserAddFailed:                   "用户添加失败",
	ErrUserQueryFailed:                 "用户查询失败",
	ErrUserGetFailed:                   "用户获取失败",
	ErrUserActiveFailed:                "用户启用失败",
	ErrUserInActiveFailed:              "用户禁用失败",
	ErrUserDeleted:                     "用户删除失败",
	ErrUserUpdateFailed:                "用户修改失败",
	ErrUserNotExist:                    "用户不存在",
	ErrAuthGetMachineIDFailed:          "获取机器标识失败",
	ErrAuthGetLicenseInfoFailed:        "获取系统许可证信息失败",
	ErrAuthUpdateLicenseInfoFailed:     "更新系统许可证信息失败",
	ErrAuthLicenseHasExpiredOrNotExist: "系统许可证已过期或不存在",
	ErrAuthGetMachineIDEmpty:           "获取的机器标识为空",
	ErrLDAPConnectFailed:               "连接LDAP/AD域失败",
	ErrUserUpdatePasswordFailed:        "修改用户密码失败",
	ErrUserInvalidPassword:             "账号有误，登录失败",
	ErrOnlineListFailed:                "获取登录用户列表失败",
	ErrUserLDAPDefRoleFailed:           "设置ldap用户默认角色失败",
	ErrUserInvalidOldPassword:          "原密码有误",
	ErrUserNameCollideBuiltin:          "用户名不能为public、project",
	ErrUserDefaultRoleNotExist:         "未设置ldap用户首次登录默认角色，请先使用管理员账号进入系统设置设定默认角色!",
	ErrUserCantDeleteDefaultUser:       "无法删除或禁用默认审批人，请先在系统全局设置中移除此审批人",
	ErrUserResetPass:                   "重置用户密码失败",
	ErrOrgCreatedFailed:                "创建新组织架构失败",
	ErrOrgUpdateFailed:                 "修改组织架构失败",
	ErrOrgDeleted:                      "删除组织架构失败",
	ErrOrgMemberAddFailed:              "添加组织架构成员失败",
	ErrOrgMemberDeletedFailed:          "删除组织架构成员失败",
	ErrOrgMemberUpdateFailed:           "修改组织架构成员失败",
	ErrOrgMemberListFailed:             "组织架构成员列表获取失败",

	ErrUserOpenapiCertAlreadyExist:  "openapi凭证已存在",
	ErrUserOpenapiCertCreatedFailed: "openapi凭证创建失败",
	ErrUserOpenapiCertDeleteFailed:  "openapi凭证删除失败",
	ErrUserOpenapiCertGetFailed:     "openapi凭证获取失败",
	ErrUserOpenapiCertGenFailed:     "openapi凭证生成失败",
	ErrUserOpenapiCertDisable:       "未开启openapi访问权限",
	ErrUserOpenapiCertNotEnable:     "禁用用户不允许调用openapi",
}
