package consts

import (
	"github.com/yuansuan/ticp/PSP/psp/internal/common"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
)

const (
	LocalAppTemplateDir  = "config/template/local"
	CloudAppTemplateDir  = "config/template/cloud"
	IconDataBase64Prefix = "data:image/" + common.Png + ";base64,"
)

const (
	ResourceDataTypeDynamic = "dynamic"
	ResourceDataTypeKey     = "key"
	ResourceDataTypeValue   = "value"

	ResourceResolverTypePlatform = "platform"
	ResourceResolverTypeQueue    = "queue"
)

const (
	Cmd = "cmd"
)

const (
	InternalTemplateStarCCMId = snowflake.ID(1689929831401132032)
)

var PublishMap = map[string]string{
	common.Unpublished: "取消发布",
	common.Published:   "发布",
}
