package application

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/admin/app/update"
	schema "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/consts"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/util"
)

func validateName(name string) error {
	if len(name) == 0 {
		return errors.New("empty name")
	}
	if len(name) > 255 {
		return errors.New("name length limit")
	}
	return nil
}

func validateLicManagerID(licManagerID string) error {
	if len(licManagerID) != 0 {
		_, err := snowflake.ParseString(licManagerID)
		if err != nil {
			return fmt.Errorf("parse licManagerId to snowflake failed, err:%v", err)
		}
	}
	return nil
}

func validateType(tp string) error {
	if len(tp) == 0 {
		return errors.New("empty type")
	}
	return nil
}

func validateVersion(version string) error {
	if len(version) == 0 {
		return errors.New("empty version")
	}
	return nil
}

func validateDescription(description string) error {
	if len(description) > 255 {
		return errors.New("description length limit")
	}
	return nil
}

func validateBinPath(binPath string, zones schema.Zones) error {
	binPathLen := len(binPath)
	if binPathLen == 0 {
		return nil
	}

	if binPathLen > 65535 {
		return errors.New("binPath length limit")
	}

	binPathZones := util.ToStringMap(binPath)
	if binPathZones == nil {
		return errors.New("invalid bin path")
	}
	for key := range binPathZones {
		if !zones.IsZone(key) {
			return errors.New("invalid bin path")
		}
	}

	return nil
}

func validateExtentionParams(extentionParams string) error {
	ep := make(map[string]schema.ExtentionParam)
	err := json.Unmarshal([]byte(extentionParams), &ep)
	if err != nil {
		return errors.New("invalid extention params")
	}

	return nil
}

func validateSpecifyQueue(q map[string]string, zones schema.Zones) error {
	if len(q) == 0 {
		return nil
	}
	for key := range q {
		if !zones.IsZone(key) {
			msg := fmt.Sprintf("invalid specify queue, zone %s not exist", key)
			return errors.New(msg)
		}
	}

	return nil
}

func validateBinPathAndImage(binPath string, image string, zones schema.Zones) error {
	if len(binPath) == 0 && len(image) == 0 {
		return errors.New("image and binPath cannot be empty at the same time")
	}

	return validateBinPath(binPath, zones)
}

func validatePublishStatus(publishStatus update.Status) error {
	if len(publishStatus) != 0 && publishStatus != update.Published && publishStatus != update.Unpublished {
		return errors.New("invalid publish status")
	}

	return nil
}

func validateResidualLogParser(residualLogParser string) error {
	if residualLogParser != "" && residualLogParser != schema.ResidualLogParserTypeStarccm && residualLogParser != schema.ResidualLogParserTypeFluent {
		return errors.New("invalid residualLogParser, only support [starccm, fluent]")
	}
	return nil
}

func validateMonitorChartParser(monitorChartParser string) error {
	if monitorChartParser != "" && monitorChartParser != schema.MonitorChartParserTypeFluent && monitorChartParser != schema.MonitorChartParserTypeCfx {
		return errors.New("invalid monitorChartParser, only support [fluent, cfx]")
	}
	return nil
}

func validateCommand(command string) error {
	// Check if the command contains reserved command flag, return error if not
	if !strings.Contains(command, consts.AppPreparedFlag) {
		return fmt.Errorf("invalid command, must contain reserved command flag [%s]", consts.AppPreparedFlag)
	}

	return nil
}
