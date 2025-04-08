package errcode

import (
	"google.golang.org/grpc/codes"
)

// Billing 错误码范围: 22001 ~ 23000
// 命名格式：Err + 服务名 + 具体错误

const (
	ErrFailedLicenseManagerList       codes.Code = 22001
	ErrFailedLicenseManagerAdd        codes.Code = 22002
	ErrFailedLicenseManagerEdit       codes.Code = 22003
	ErrFailedLicenseManagerDelete     codes.Code = 22004
	ErrFailedLicenseManagerDeleteBind codes.Code = 22005

	ErrFailedLicenseInfoAdd  codes.Code = 22006
	ErrFailedLicenseInfoEdit codes.Code = 22007
	ErrFailedLicenseInfoDel  codes.Code = 22008

	ErrFailedConfigModuleList   codes.Code = 22009
	ErrFailedConfigModuleAdd    codes.Code = 22010
	ErrFailedConfigModuleEdit   codes.Code = 22011
	ErrFailedConfigModuleDelete codes.Code = 22012

	ErrFailedLicenseNameRepeat         codes.Code = 22013
	ErrFailedAppTypeRepeat             codes.Code = 22014
	ErrFailedLicenseManagerDeleteExist codes.Code = 22015
	ErrFailedModuleNameRepeat          codes.Code = 22016
)

const (
	MsgFailedLicenseManagerList       = "获取license列表异常"
	MsgFailedLicenseManagerAdd        = "保存license管理器异常"
	MsgFailedLicenseManagerEdit       = "编辑license管理器异常"
	MsgFailedLicenseManagerDelete     = "删除license管理器异常"
	MsgFailedLicenseManagerDeleteBind = "删除license管理器异常，当前license已被应用绑定"

	MsgFailedLicenseInfoAdd  = "保存license信息异常"
	MsgFailedLicenseInfoEdit = "编辑license信息异常"
	MsgFailedLicenseInfoDel  = "删除license信息异常"

	MsgFailedConfigModuleList   = "获取模块列表异常"
	MsgFailedConfigModuleAdd    = "保存模块异常"
	MsgFailedConfigModuleEdit   = "编辑模块异常"
	MsgFailedConfigModuleDelete = "删除模块异常"

	MsgFailedLicenseNameRepeat         = "许可证名字重复"
	MsgFailedAppTypeRepeat             = "许可证类型重复"
	MsgFailedModuleNameRepeat          = "许可证模块名字重复"
	MsgFailedLicenseManagerDeleteExist = "当前许可证下面存在许可证服务器信息，请先删除"
)
