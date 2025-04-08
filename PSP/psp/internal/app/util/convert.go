package util

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"strings"

	"github.com/yuansuan/ticp/PSP/psp/internal/app/dao/model"
	"github.com/yuansuan/ticp/PSP/psp/internal/app/dto"
	"github.com/yuansuan/ticp/PSP/psp/internal/common"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/app"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/license"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/monitor"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
	"github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
)

func ConvertAppDto(app *model.App, queues []*dto.QueueInfo, licenses []*dto.LicenseInfo) (*dto.App, error) {
	if app == nil {
		return nil, nil
	}

	appDto := &dto.App{
		ID:                app.ID.String(),
		OutAppID:          app.OutAppID,
		CloudOutAppID:     app.OutID,
		CloudOutAppName:   app.OutName,
		Name:              app.Name,
		Type:              app.Type,
		Version:           app.VersionNum,
		ComputeType:       app.ComputeType,
		Queues:            make([]*dto.QueueInfo, 0),
		Licenses:          make([]*dto.LicenseInfo, 0),
		State:             app.State,
		Image:             app.Image,
		EnableResidual:    app.EnableResidual,
		ResidualLogParser: app.ResidualLogParser,
		EnableSnapshot:    app.EnableSnapshot,
		Script:            app.Script,
		Icon:              app.Icon,
		Description:       app.Description,
		HelpDoc: &dto.HelpDoc{
			Type:  app.DocType,
			Value: app.DocContent,
		},
		BinPath:        GetKeyValue(app.BinPath),
		SchedulerParam: GetKeyValue(app.SchedulerParam),
	}

	if len(app.QueueNames) != 0 && len(queues) != 0 {
		queueNameMap := make(map[string]struct{}, len(app.QueueNames))
		for _, v := range app.QueueNames {
			queueNameMap[v] = struct{}{}
		}

		queueCopy := make([]*dto.QueueInfo, 0, len(queues))
		for _, v := range queues {
			queueCopy = append(queueCopy, &dto.QueueInfo{
				QueueName: v.QueueName,
				CPUNumber: v.CPUNumber,
				Select:    v.Select,
			})
		}

		queueSelected := make([]*dto.QueueInfo, 0, len(queueCopy))
		for _, v := range queueCopy {
			if _, ok := queueNameMap[v.QueueName]; ok {
				v.Select = true
			}
			queueSelected = append(queueSelected, v)
		}

		appDto.Queues = queueSelected
	}

	if app.LicenseManagerId != "" && len(licenses) != 0 {
		licenseCopy := make([]dto.LicenseInfo, 0, len(licenses))
		for _, lic := range licenses {
			licenseCopy = append(licenseCopy, *lic)
		}

		licenseSelected := make([]*dto.LicenseInfo, 0, len(licenseCopy))
		for i, lic := range licenseCopy {
			if app.LicenseManagerId == lic.Id {
				licenseCopy[i].Select = true
			}
			licenseSelected = append(licenseSelected, &licenseCopy[i])
		}

		appDto.Licenses = licenseSelected
	}

	subForm := &dto.SubForm{}
	err := yaml.Unmarshal([]byte(app.Content), subForm)
	if err != nil {
		return nil, fmt.Errorf("unmarshal the app content err: %v, app: [%+v]", err, app)
	}
	appDto.SubForm = subForm

	return appDto, nil
}

func ConvertAppModel(app *dto.App) (*model.App, error) {
	if app == nil {
		return nil, nil
	}

	appModel := &model.App{
		Name:        app.Name,
		Type:        app.Type,
		VersionNum:  app.Version,
		State:       app.State,
		Image:       app.Image,
		Script:      ReplaceWindowsLineBreak(app.Script),
		Description: app.Description,
		Icon:        app.Icon,
		BinPath:     GetKeyValueMap(app.BinPath),
	}

	if app.HelpDoc != nil {
		appModel.DocContent = app.HelpDoc.Value
		appModel.DocType = app.HelpDoc.Type
	}

	content, err := yaml.Marshal(app.SubForm)
	if err != nil {
		return nil, fmt.Errorf("marshal the app content err: %v, app: [%+v]", err, app)
	}
	appModel.Content = string(content)
	return appModel, nil
}

func ConvertCloudAppToAppModel(ID snowflake.ID, app *v20230530.Application) *model.App {
	if app == nil {
		return nil
	}
	appModel := &model.App{
		ID:                ID,
		OutAppID:          app.AppID,
		Name:              app.Name,
		Type:              app.Type,
		ComputeType:       common.Cloud,
		VersionNum:        app.Version,
		QueueNames:        make([]string, 0),
		State:             common.Unpublished,
		EnableResidual:    app.ResidualEnable,
		EnableSnapshot:    app.SnapshotEnable,
		ResidualLogParser: app.ResidualLogParser,
	}
	return appModel
}

func ConvertRPCAppInfo(appInfo *dto.App) *app.App {
	if appInfo == nil {
		return nil
	}
	return &app.App{
		Id:                appInfo.ID,
		OutAppId:          appInfo.OutAppID,
		CloudOutAppId:     appInfo.CloudOutAppID,
		Name:              appInfo.Name,
		Type:              appInfo.Type,
		ComputeType:       appInfo.ComputeType,
		State:             appInfo.State,
		Version:           appInfo.Version,
		Script:            appInfo.Script,
		Icon:              appInfo.Icon,
		Description:       appInfo.Description,
		EnableResidual:    appInfo.EnableResidual,
		ResidualLogParser: appInfo.ResidualLogParser,
		EnableSnapshot:    appInfo.EnableSnapshot,
		HelpDoc:           ConvertRPCHelpDoc(appInfo.HelpDoc),
		SubForm:           ConvertRPCSections(appInfo.SubForm),
		SchedulerParam:    GetRPCKeyValue(appInfo.SchedulerParam),
	}
}

func ConvertRPCHelpDoc(doc *dto.HelpDoc) *app.HelpDoc {
	if doc == nil {
		return nil
	}
	return &app.HelpDoc{
		Type:  doc.Type,
		Value: doc.Value,
	}
}

func ConvertRPCSections(subForm *dto.SubForm) *app.SubForm {
	if subForm == nil {
		return nil
	}
	var sections []*app.Section
	for _, section := range subForm.Section {
		sections = append(sections, &app.Section{
			Name:  section.Name,
			Field: ConvertRPCFields(section.Field),
		})
	}
	return &app.SubForm{
		Section: sections,
	}
}

func ConvertRPCFields(fields []*dto.Field) []*app.Field {
	var result []*app.Field
	for _, field := range fields {
		result = append(result, &app.Field{
			Id:                      field.ID,
			Label:                   field.Label,
			Help:                    field.Help,
			Type:                    field.Type,
			Required:                field.Required,
			Hidden:                  field.Hidden,
			DefaultValue:            field.DefaultValue,
			DefaultValues:           field.DefaultValues,
			Value:                   field.Value,
			Values:                  field.Values,
			Action:                  field.Action,
			Options:                 field.Options,
			PostText:                field.PostText,
			FileFromType:            field.FileFromType,
			IsMasterSlave:           field.IsMasterSlave,
			MasterIncludeKeywords:   field.MasterIncludeKeywords,
			MasterIncludeExtensions: field.MasterIncludeExtensions,
			MasterSlave:             field.MasterSlave,
			OptionsFrom:             field.OptionsFrom,
			OptionsScript:           field.OptionsScript,
			CustomJsonValueString:   field.CustomJSONValueString,
			IsSupportMaster:         field.IsSupportMaster,
			MasterFile:              field.MasterFile,
			IsSupportWorkdir:        field.IsSupportWorkdir,
			Workdir:                 field.Workdir,
		})
	}
	return result
}

func ConvertQueueInfos(queueList *monitor.QueueListResponse) []*dto.QueueInfo {
	queues := make([]*dto.QueueInfo, 0, len(queueList.QueueNames))

	for _, v := range queueList.QueueNames {
		queues = append(queues, &dto.QueueInfo{
			QueueName: v,
		})
	}

	return queues
}

func ConvertLicenseInfos(licenseList *license.QueueLicenseTypeListResponse) []*dto.LicenseInfo {
	licenses := make([]*dto.LicenseInfo, 0, len(licenseList.LicenseTypes))

	for _, v := range licenseList.LicenseTypes {
		licenses = append(licenses, &dto.LicenseInfo{
			Id:           v.Id,
			Name:         v.TypeName,
			LicenceValid: v.LicenceValid,
		})
	}

	return licenses
}

func ConvertBinPathMap(binPath string) (map[string]string, error) {
	if binPath == "" {
		return nil, nil
	}

	binPathMap := make(map[string]string)
	err := json.Unmarshal([]byte(binPath), &binPathMap)
	if err != nil {
		return nil, fmt.Errorf("unmarshal the bin path err: %v, binPath: [%v]", err, binPath)
	}

	return binPathMap, nil
}

func GetKeyValueMap(kv []*dto.KeyValue) map[string]string {
	kvMap := make(map[string]string)

	for _, v := range kv {
		kvMap[v.Key] = v.Value
	}

	return kvMap
}

func GetKeyValue(kvMap map[string]string) []*dto.KeyValue {
	kv := make([]*dto.KeyValue, 0, len(kvMap))

	for k, v := range kvMap {
		kv = append(kv, &dto.KeyValue{Key: k, Value: v})
	}

	return kv
}

func GetRPCKeyValue(kvs []*dto.KeyValue) []*app.KeyValue {
	kv := make([]*app.KeyValue, 0, len(kvs))

	for _, v := range kvs {
		kv = append(kv, &app.KeyValue{Key: v.Key, Value: v.Value})
	}

	return kv
}

func ReplaceWindowsLineBreak(input string) string {
	return strings.ReplaceAll(input, "\r", "")
}

func GetSimpleAppValue(app *model.App) *model.App {
	if app != nil {
		app.Content = ""
		app.Icon = ""
	}
	return app
}
