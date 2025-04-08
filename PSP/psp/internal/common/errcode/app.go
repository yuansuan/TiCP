package errcode

import (
	"google.golang.org/grpc/codes"
)

// App 错误码范围: 11001 ~ 12000
// 命名格式：Err + 服务名 + 具体错误
const (
	ErrAppTemplateHasExist                codes.Code = 11001
	ErrAppTemplateNotExist                codes.Code = 11002
	ErrAppTemplateHasPublished            codes.Code = 11003
	ErrAppTemplateHasUnpublished          codes.Code = 11004
	ErrAppBaseTemplateNotExist            codes.Code = 11011
	ErrAppGetAppInfoFailed                codes.Code = 11012
	ErrAppListAppFailed                   codes.Code = 11013
	ErrAppListTemplateFailed              codes.Code = 11014
	ErrAppGetAppScriptFailed              codes.Code = 11015
	ErrAppAddAppFailed                    codes.Code = 11016
	ErrAppUpdateAppFailed                 codes.Code = 11017
	ErrAppDeleteAppFailed                 codes.Code = 11018
	ErrAppPublishAppFailed                codes.Code = 11019
	ErrAppUnableOperateCloudAppData       codes.Code = 11020
	ErrAppSyncCloudAppFailed              codes.Code = 11021
	ErrAppUnableCreateNewAPPBaseCloudApp  codes.Code = 11022
	ErrAppUnableUpdateAppCompeteType      codes.Code = 11023
	ErrAppUnableCreateNewCloudApp         codes.Code = 11024
	ErrAppListZoneFailed                  codes.Code = 11025
	ErrAppParamLicenseManagerIdInvalid    codes.Code = 11026
	ErrAppParamBinPathOverLengthLimit     codes.Code = 11027
	ErrAppListQueueFailed                 codes.Code = 11028
	ErrAppLocalCompeteTypeNotSupported    codes.Code = 11029
	ErrAppCloudCompeteTypeNotSupported    codes.Code = 11030
	ErrAppSyncAppContentFailed            codes.Code = 11031
	ErrAppBindCloudTemplateNotExist       codes.Code = 11032
	ErrAppRelationTemplateNotExist        codes.Code = 11033
	ErrAppGetSchedulerResourceKeyFailed   codes.Code = 11034
	ErrAppGetSchedulerResourceValueFailed codes.Code = 11035
	ErrAppGetSchedulerResourceKeyNotFound codes.Code = 11036
	ErrAppListLicenseFailed               codes.Code = 11037
)

// AppCodeMsg ...
var AppCodeMsg = map[codes.Code]string{
	ErrAppTemplateHasExist:                "模版已存在",
	ErrAppTemplateNotExist:                "模版不存在",
	ErrAppTemplateHasPublished:            "模版已发布",
	ErrAppTemplateHasUnpublished:          "模版未发布",
	ErrAppBaseTemplateNotExist:            "基础模版不存在",
	ErrAppGetAppInfoFailed:                "获取应用信息失败",
	ErrAppListAppFailed:                   "获取应用列表失败",
	ErrAppListTemplateFailed:              "获取模版列表失败",
	ErrAppGetAppScriptFailed:              "获取应用脚本失败",
	ErrAppAddAppFailed:                    "添加应用失败",
	ErrAppUpdateAppFailed:                 "更新应用失败",
	ErrAppDeleteAppFailed:                 "删除应用失败",
	ErrAppPublishAppFailed:                "发布应用失败",
	ErrAppUnableOperateCloudAppData:       "无法操作云端计算应用数据",
	ErrAppSyncCloudAppFailed:              "同步云端计算应用模版失败",
	ErrAppUnableCreateNewAPPBaseCloudApp:  "不能基于云端计算应用模版创建新的应用模版",
	ErrAppUnableUpdateAppCompeteType:      "不能更新应用模版的计算类型",
	ErrAppUnableCreateNewCloudApp:         "不能创建新的云端计算应用",
	ErrAppListZoneFailed:                  "获取区域列表失败",
	ErrAppParamLicenseManagerIdInvalid:    "许可证管理 ID 参数错误, 示例: 4W8YD7Vmwbm",
	ErrAppParamBinPathOverLengthLimit:     "可执行文件路径参数少于或等于 255 字符",
	ErrAppListQueueFailed:                 "获取区域列表失败",
	ErrAppLocalCompeteTypeNotSupported:    "本地应用类型不支持当前操作",
	ErrAppCloudCompeteTypeNotSupported:    "云应用类型不支持当前操作",
	ErrAppSyncAppContentFailed:            "设置应用模版信息失败",
	ErrAppBindCloudTemplateNotExist:       "需要绑定的云应用不存在",
	ErrAppRelationTemplateNotExist:        "关联云模版不存在",
	ErrAppGetSchedulerResourceKeyFailed:   "获取调度器资源键信息失败",
	ErrAppGetSchedulerResourceValueFailed: "获取调度器资源值信息失败",
	ErrAppGetSchedulerResourceKeyNotFound: "调度器资源键信息不存在",
	ErrAppListLicenseFailed:               "获取许可证列表失败",
}
