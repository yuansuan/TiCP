package errcode

import (
	"google.golang.org/grpc/codes"
)

const (
	ErrProjectAdd                   codes.Code = 22001
	ErrProjectList                  codes.Code = 22002
	ErrProjectDetail                codes.Code = 22003
	ErrProjectDelete                codes.Code = 22004
	ErrProjectTerminated            codes.Code = 22005
	ErrProjectOwnerID               codes.Code = 22006
	ErrProjectEdit                  codes.Code = 22007
	ErrProjectModifyOwner           codes.Code = 22008
	ErrProjectMemberAdd             codes.Code = 22010
	ErrProjectMemberRemove          codes.Code = 22011
	ErrProjectAccessPermission      codes.Code = 22012
	ErrProjectSameName              codes.Code = 22013
	ErrProjectDelBeforeExistMembers codes.Code = 22014
	ErrProjectNotFound              codes.Code = 22015
	ErrProjectCurrentList           codes.Code = 22016
	ErrProjectEditState             codes.Code = 22017
	ErrProjectCurrentListForParam   codes.Code = 22018
	ErrProjectNameIsDefault         codes.Code = 22019
	ErrProjectDeleteState           codes.Code = 22020
	ErrProjectTerminateState        codes.Code = 22021
	ErrProjectMemberNotExist        codes.Code = 22022
)

var ProjectCodeMsg = map[codes.Code]string{
	ErrProjectAdd:                   "保存项目失败",
	ErrProjectList:                  "读取项目列表失败",
	ErrProjectDetail:                "读取项目详情失败",
	ErrProjectDelete:                "删除项目详情失败",
	ErrProjectTerminated:            "终止项目详情失败",
	ErrProjectOwnerID:               "项目管理员id不正确",
	ErrProjectEdit:                  "项目编辑失败",
	ErrProjectModifyOwner:           "修改项目管理员,失败",
	ErrProjectMemberAdd:             "添加项目成员失败",
	ErrProjectMemberRemove:          "移出项目成员失败",
	ErrProjectAccessPermission:      "无当前项目访问权限",
	ErrProjectSameName:              "项目记录名称已存在",
	ErrProjectDelBeforeExistMembers: "删除项目需先移除项目成员",
	ErrProjectNotFound:              "项目不存在",
	ErrProjectCurrentList:           "读取当前项目列表失败",
	ErrProjectEditState:             "终止或者结束项目不能编辑",
	ErrProjectCurrentListForParam:   "读取当前项目列表失败(条件参数)",
	ErrProjectNameIsDefault:         "项目名称不能为 `personal` 默认关键字",
	ErrProjectDeleteState:           "运行状态项目不能被删除",
	ErrProjectTerminateState:        "初始化或者运行状态项目不能被终止",
	ErrProjectMemberNotExist:        "项目成员不存在",
}
