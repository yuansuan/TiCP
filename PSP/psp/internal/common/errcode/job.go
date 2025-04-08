package errcode

import (
	"google.golang.org/grpc/codes"
)

// Job service error codes from 12001 to 13000
const (
	ErrJobFailGet                         codes.Code = 12001
	ErrJobFailList                        codes.Code = 12002
	ErrJobFailTerminate                   codes.Code = 12003
	ErrJobFailAppNameList                 codes.Code = 12004
	ErrJobFailUserNameList                codes.Code = 12005
	ErrJobFailQueueNameList               codes.Code = 12006
	ErrJobFailSubmit                      codes.Code = 12007
	ErrJobExist                           codes.Code = 12008
	ErrJobFailCreateTempDir               codes.Code = 12009
	ErrJobFailComputeTypeList             codes.Code = 12010
	ErrJobFailCPUTimeTotal                codes.Code = 12011
	ErrJobFailStatisticsOverview          codes.Code = 12012
	ErrJobFailStatisticsDetail            codes.Code = 12013
	ErrJobFailStatisticsExport            codes.Code = 12014
	ErrJobFailAppJobNum                   codes.Code = 12015
	ErrJobFailNotEnoughBalance            codes.Code = 12016
	ErrJobResidualFailGet                 codes.Code = 12017
	ErrJobSnapshotFailList                codes.Code = 12018
	ErrJobSnapshotFailGet                 codes.Code = 12019
	ErrJobCurrentProjectNotRunning        codes.Code = 12020
	ErrJobCurrentProjectNotAccess         codes.Code = 12021
	ErrJobFailJobSetNameList              codes.Code = 12022
	ErrJobFailJobSetGet                   codes.Code = 12023
	ErrJobResubmitFailed                  codes.Code = 12024
	ErrJobNotExist                        codes.Code = 12025
	ErrJobStatusNotSupportResubmit        codes.Code = 12026
	ErrJobLastSubmitParamNotExist         codes.Code = 12027
	ErrJobResubmitAppNotExistOrNotPublish codes.Code = 12028
	ErrJobFailGetTop5ProjectInfo          codes.Code = 12029
)

// Job service error message
var JobCodeMsg = map[codes.Code]string{
	ErrJobFailGet:                         "获取作业信息失败",
	ErrJobFailList:                        "获取作业列表失败",
	ErrJobFailTerminate:                   "作业终止失败",
	ErrJobFailAppNameList:                 "获取应用名称失败",
	ErrJobFailUserNameList:                "获取用户名称失败",
	ErrJobFailQueueNameList:               "获取队列名称失败",
	ErrJobFailSubmit:                      "作业提交失败",
	ErrJobExist:                           "作业已存在",
	ErrJobFailCreateTempDir:               "创建作业临时目录失败",
	ErrJobFailComputeTypeList:             "获取计算类型失败",
	ErrJobFailCPUTimeTotal:                "获取作业总核时失败",
	ErrJobFailStatisticsOverview:          "获取作业总览数据失败",
	ErrJobFailStatisticsDetail:            "获取作业详情数据失败",
	ErrJobFailStatisticsExport:            "导出作业数据失败",
	ErrJobFailAppJobNum:                   "获取应用作业数失败",
	ErrJobFailNotEnoughBalance:            "账户余额不足，请及时充值",
	ErrJobResidualFailGet:                 "获取作业残差图失败",
	ErrJobSnapshotFailList:                "获取作业云图集失败",
	ErrJobSnapshotFailGet:                 "获取作业云图失败",
	ErrJobCurrentProjectNotRunning:        "该项目没有在[进行中]状态",
	ErrJobCurrentProjectNotAccess:         "该项目没有权限访问",
	ErrJobFailJobSetNameList:              "获取作业集名称失败",
	ErrJobFailJobSetGet:                   "获取作业集信息失败",
	ErrJobResubmitFailed:                  "重新提交作业失败",
	ErrJobNotExist:                        "作业不存在",
	ErrJobStatusNotSupportResubmit:        "作业状态不支持重新提交",
	ErrJobLastSubmitParamNotExist:         "作业的上一次提交参数不存在",
	ErrJobResubmitAppNotExistOrNotPublish: "重新提交作业的计算应用不存在或未发布",
	ErrJobFailGetTop5ProjectInfo:          "获取Top5项目信息失败",
}
