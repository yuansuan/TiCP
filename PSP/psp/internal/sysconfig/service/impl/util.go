package impl

import (
	"context"
	"reflect"
	"strconv"
	"time"

	"github.com/spf13/viper"
	"github.com/yuansuan/ticp/common/go-kit/logging"

	sysconfig "github.com/yuansuan/ticp/PSP/psp/internal/sysconfig/config"
	"github.com/yuansuan/ticp/PSP/psp/internal/sysconfig/consts"
	"github.com/yuansuan/ticp/PSP/psp/internal/sysconfig/dao/model"
	"github.com/yuansuan/ticp/PSP/psp/internal/sysconfig/dto"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
)

func getNewViper() *viper.Viper {
	ctx := context.Background()
	newViper := viper.New()
	newViper.SetConfigType(consts.FileType)
	newViper.SetConfigName(consts.AlertManagerName)
	newViper.AddConfigPath(sysconfig.GetConfig().AlertManager.AlertManagerConfigPath)

	// Read the config file
	err := newViper.ReadInConfig()
	if err != nil {
		logging.GetLogger(ctx).Errorf("Failed to Read config in alertmanager.yml %v", err.Error())
		return nil
	}
	return newViper
}

func buildGlobalEmail(sid *snowflake.Node, in *dto.EmailConfig) []*model.AlertNotification {
	// 将结构体转换为map
	resultMap := make(map[string]interface{})
	resultMap[consts.KeyHost] = in.Host
	resultMap[consts.KeyPort] = in.Port
	resultMap[consts.KeyUsername] = in.UserName
	resultMap[consts.KeyPassword] = in.Password
	resultMap[consts.KeyFrom] = in.From
	resultMap[consts.KeyAdminAddr] = in.AdminAddr
	resultMap[consts.KeyUseTLS] = in.UseTLS

	// 将map转换为AlertNotification
	alertNotification := buildModel(sid, resultMap, consts.GlobalEmailType)

	return alertNotification
}

func buildAlertNotification(sid *snowflake.Node, in *dto.SetEmailConfigReq) []*model.AlertNotification {
	// 将结构体转换为map
	resultMap := make(map[string]interface{})
	resultMap[consts.KeyNodeBreakdown] = in.Notification.NodeBreakdown
	resultMap[consts.KeyDiskUsage] = in.Notification.DiskUsage
	resultMap[consts.KeyAgentBreakdown] = in.Notification.AgentBreakdown
	resultMap[consts.KeyJobFailNum] = in.Notification.JobFailNum

	// 将map转换为AlertNotification
	alertNotification := buildModel(sid, resultMap, consts.AlertManagerType)
	return alertNotification
}

func buildModel(sid *snowflake.Node, resultMap map[string]interface{}, alertType string) []*model.AlertNotification {
	var alertNotification []*model.AlertNotification

	now := time.Now()
	for k, v := range resultMap {
		var valueStr string
		//判断是否是bool类型
		if reflect.TypeOf(v).Kind() == reflect.Bool {
			valueStr = boolToString(v.(bool))
		} else if reflect.TypeOf(v).Kind() == reflect.Int {
			valueStr = strconv.Itoa(v.(int))
		} else {
			valueStr = v.(string)
		}
		alertNotification = append(alertNotification, &model.AlertNotification{
			Id:         sid.Generate(),
			Key:        k,
			Value:      valueStr,
			Type:       alertType,
			CreateTime: now,
			UpdateTime: now,
		})
	}
	return alertNotification
}

func boolToString(b bool) string {
	if b {
		return "1"
	}
	return "0"
}

func stringToBool(s string) bool {
	if s == "1" {
		return true
	}
	return false
}

func stringToInt(str string) int {
	num, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return num
}
