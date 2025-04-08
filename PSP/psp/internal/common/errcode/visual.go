package errcode

import (
	"google.golang.org/grpc/codes"
)

// Visual 错误码范围: 17001 ~ 18000
// 命名格式：Err + 服务名 + 具体错误
const (
	ErrVisualSessionNotFound                     codes.Code = 17001
	ErrVisualHardwareNotFound                    codes.Code = 17002
	ErrVisualSoftwareNotFound                    codes.Code = 17003
	ErrVisualRemoteAppNotFound                   codes.Code = 17004
	ErrVisualSessionNotStarted                   codes.Code = 17005
	ErrVisualHardwareHasUsed                     codes.Code = 17006
	ErrVisualSoftwareHasUsed                     codes.Code = 17007
	ErrVisualRemoteAppHasUsed                    codes.Code = 17008
	ErrVisualListSessionFailed                   codes.Code = 17009
	ErrVisualStartSessionFailed                  codes.Code = 17010
	ErrVisualCloseSessionFailed                  codes.Code = 17011
	ErrVisualReadySessionFailed                  codes.Code = 17012
	ErrVisualGetRemoteAppURLFailed               codes.Code = 17013
	ErrVisualListHardwareFailed                  codes.Code = 17014
	ErrVisualAddHardwareFailed                   codes.Code = 17015
	ErrVisualUpdateHardwareFailed                codes.Code = 17016
	ErrVisualDeleteHardwareFailed                codes.Code = 17017
	ErrVisualListSoftwareFailed                  codes.Code = 17018
	ErrVisualAddSoftwareFailed                   codes.Code = 17019
	ErrVisualUpdateSoftwareFailed                codes.Code = 17020
	ErrVisualDeleteSoftwareFailed                codes.Code = 17021
	ErrVisualGetSoftwarePresetsFailed            codes.Code = 17022
	ErrVisualSetSoftwarePresetsFailed            codes.Code = 17023
	ErrVisualAddRemoteAppFailed                  codes.Code = 17024
	ErrVisualUpdateRemoteAppFailed               codes.Code = 17025
	ErrVisualDeleteRemoteAppFailed               codes.Code = 17026
	ErrVisualDurationStatisticFailed             codes.Code = 17027
	ErrVisualListHistoryDurationFailed           codes.Code = 17028
	ErrVisualSessionRepeatStart                  codes.Code = 17029
	ErrVisualHardwareHasExist                    codes.Code = 17030
	ErrVisualSoftwareHasExist                    codes.Code = 17031
	ErrVisualPublishSoftwareFailed               codes.Code = 17032
	ErrVisualSoftwareHasPublishedFailed          codes.Code = 17033
	ErrVisualSoftwareHasUsedForPublish           codes.Code = 17034
	ErrVisualListSoftwareUsingStatusesFailed     codes.Code = 17035
	ErrVisualListUsedProjectNamesFailed          codes.Code = 17036
	ErrVisualGetMountInfoFailed                  codes.Code = 17037
	ErrVisualCurrentProjectNotRunning            codes.Code = 17038
	ErrVisualCurrentProjectNotAccess             codes.Code = 17039
	ErrVisualRebootSessionFailed                 codes.Code = 17040
	ErrVisualPowerOffSessionFailed               codes.Code = 17041
	ErrVisualPowerOnSessionFailed                codes.Code = 17042
	ErrVisualSessionUsageDurationStatisticFailed codes.Code = 17043
	ErrVisualSessionCreateNumberStatisticFailed  codes.Code = 17044
	ErrVisualSessionNumberStatusStatisticFailed  codes.Code = 17045
	ErrVisualExportSessionInfoFailed             codes.Code = 17046
	ErrVisualExportUsageDurationStatisticFailed  codes.Code = 17047
)

// VisualCodeMsg ...
var VisualCodeMsg = map[codes.Code]string{
	ErrVisualSessionNotFound:                     "会话不存在",
	ErrVisualHardwareNotFound:                    "实例不存在",
	ErrVisualSoftwareNotFound:                    "镜像不存在",
	ErrVisualRemoteAppNotFound:                   "远程应用不存在",
	ErrVisualSessionNotStarted:                   "会话未启动",
	ErrVisualHardwareHasUsed:                     "实例已被使用, 无法删除",
	ErrVisualSoftwareHasUsed:                     "镜像已被使用, 无法删除",
	ErrVisualRemoteAppHasUsed:                    "远程应用已被使用, 无法删除",
	ErrVisualListSessionFailed:                   "获取会话列表失败",
	ErrVisualStartSessionFailed:                  "启动会话失败",
	ErrVisualCloseSessionFailed:                  "关闭会话失败",
	ErrVisualReadySessionFailed:                  "准备会话失败",
	ErrVisualGetRemoteAppURLFailed:               "获取远程应用URL失败",
	ErrVisualListHardwareFailed:                  "获取实例列表失败",
	ErrVisualAddHardwareFailed:                   "添加实例失败",
	ErrVisualUpdateHardwareFailed:                "更新实例失败",
	ErrVisualDeleteHardwareFailed:                "删除实例失败",
	ErrVisualListSoftwareFailed:                  "获取镜像列表失败",
	ErrVisualAddSoftwareFailed:                   "添加镜像失败",
	ErrVisualUpdateSoftwareFailed:                "更新镜像失败",
	ErrVisualDeleteSoftwareFailed:                "删除镜像失败",
	ErrVisualGetSoftwarePresetsFailed:            "获取镜像预设失败",
	ErrVisualSetSoftwarePresetsFailed:            "设置镜像预设失败",
	ErrVisualAddRemoteAppFailed:                  "添加远程应用失败",
	ErrVisualUpdateRemoteAppFailed:               "更新远程应用失败",
	ErrVisualDeleteRemoteAppFailed:               "删除远程应用失败",
	ErrVisualDurationStatisticFailed:             "时长统计失败",
	ErrVisualListHistoryDurationFailed:           "获取历史时长失败",
	ErrVisualSessionRepeatStart:                  "项目 [%v] 会话「%v」已存在, 请先删除",
	ErrVisualHardwareHasExist:                    "实例已存在",
	ErrVisualSoftwareHasExist:                    "镜像已存在",
	ErrVisualPublishSoftwareFailed:               "发布镜像失败",
	ErrVisualSoftwareHasPublishedFailed:          "镜像已发布",
	ErrVisualSoftwareHasUsedForPublish:           "镜像已被使用, 请先关闭相关会话",
	ErrVisualListSoftwareUsingStatusesFailed:     "获取软件使用状态失败",
	ErrVisualListUsedProjectNamesFailed:          "获取已使用项目名称列表失败",
	ErrVisualGetMountInfoFailed:                  "获取挂载信息失败",
	ErrVisualCurrentProjectNotRunning:            "该项目没有在[进行中]状态",
	ErrVisualCurrentProjectNotAccess:             "该项目没有权限访问",
	ErrVisualRebootSessionFailed:                 "重启会话失败",
	ErrVisualPowerOffSessionFailed:               "会话关机失败",
	ErrVisualPowerOnSessionFailed:                "会话开机失败",
	ErrVisualSessionUsageDurationStatisticFailed: "会话使用时长统计失败",
	ErrVisualSessionCreateNumberStatisticFailed:  "会话创建数量统计失败",
	ErrVisualSessionNumberStatusStatisticFailed:  "会话数量状态统计失败",
	ErrVisualExportSessionInfoFailed:             "导出会话信息失败",
	ErrVisualExportUsageDurationStatisticFailed:  "导出会话使用时长统计失败",
}
