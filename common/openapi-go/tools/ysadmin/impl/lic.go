package impl

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	licenseinfo "github.com/yuansuan/ticp/common/project-root-api/license/v1/license_info"
	licmanager "github.com/yuansuan/ticp/common/project-root-api/license/v1/license_manager"
	moduleconfig "github.com/yuansuan/ticp/common/project-root-api/license/v1/module_config"
)

const (
	licManagerExampleFile        = "lic_manager_example.json"
	licManagerExampleFileContent = `{
	"Id": "52ZFMSvjWKb",
	"Os": 1,
	"Desc": "123222",
	"AppType": "StarCCM+",
	"ComputeRule": "echo '{\"ccmppower\": 1}'",
	"Status": 1
}`

	licInfoExampleFile        = "lic_info_example.json"
	licInfoExampleFileContent = `{
	"Id": "52ZFR2D2efG",
	"ManagerId": "52ZFMSvjWKb",
	"Provider": "XX超算",
	"MacAddr": "00:00:00:00:00:00",
	"ToolPath": "/home/yuansuan/lic_tool/lmutil",
	"LicenseUrl": "115.159.149.167",
	"Port": 29000,
	"licenseProxies": {
	  "http://10.0.10.3:8080": {
		"Url": "1.1.1.1",
		"Port": 29000
	  }
	},
	"LicenseNum": "4E30 AFEF CB5D D27D 353F 5B3D 7D",
	"Weight": 100,
	"BeginTime": "2021-11-06 00:05:35",
	"EndTime": "2025-11-07 00:05:35",
	"Auth": true,
	"LicenseEnvVar": "CDLMD_LICENSE_FILE",
	"LicenseType": 2,
	"HpcEndpoint": "http://10.0.10.3:8080",
	"AllowableHpcEndpoints": [
	  "http://10.0.10.3:8080"
	],
	"CollectorType": "flex"
}`

	moduleConfigExampleFile        = "module_config_example.json"
	moduleConfigExampleFileContent = `{
	"Id": "52ZFS6SDzCL",
	"LicenseId": "52ZFR2D2efG",
	"ModuleName": "ccmppower",
	"Total": 21
}`
)

type LicOptions struct {
	BaseOptions
}

func init() {
	RegisterCmd(NewLicCommand())
}

// NewLicCommand license server配置管理
func NewLicCommand() *cobra.Command {
	o := LicOptions{}
	cmd := &cobra.Command{
		Use:   "license",
		Short: "license server配置管理",
		Long:  "license server配置管理, 可以配置远算自有license和外部license",
		RunE:  helpRun,
	}

	cmd.AddCommand(
		newLicAddCommand(o),
		newLicGetCommand(o),
		newLicListCommand(o),
		newLicDeleteCommand(o),
		newLicPutCommand(o),
		newLicExampleCommand(),
	)
	return cmd
}

func newLicAddCommand(o LicOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "添加license server配置",
		Long:  "添加license server配置",
		RunE:  helpRun,
	}

	cmd.AddCommand(
		newLicAddLicManagerCommand(o),
		newLicAddLicInfoCommand(o),
		newLicAddModuleConfigCommand(o),
	)
	return cmd
}

func newLicAddLicManagerCommand(o LicOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "licmanager",
		Short: "添加license manager",
		Long:  "添加license manager",
		Args:  cobra.ExactArgs(0),
		Example: `- 添加license manager, 参数文件为lic_manager.json
  - ysadmin license add licmanager -F lic_manager.json`,
	}

	cmd.Flags().StringVarP(&o.JsonFile, "file", "F", "", "json文件")
	cmd.MarkFlagRequired("file")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return licAddLicManager(cmd, args, o)
	}

	return cmd
}

func newLicAddLicInfoCommand(o LicOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "licinfo",
		Short: "添加license info",
		Long:  "添加license info",
		Args:  cobra.ExactArgs(0),
		Example: `- 添加license info, 参数文件为lic_info.json
  - ysadmin license add licinfo -F lic_info.json`,
	}

	cmd.Flags().StringVarP(&o.JsonFile, "file", "F", "", "json文件")
	cmd.MarkFlagRequired("file")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return licAddLicInfo(cmd, args, o)
	}

	return cmd
}

func newLicAddModuleConfigCommand(o LicOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "moduleconfig",
		Short: "添加module config",
		Long:  "添加module config",
		Args:  cobra.ExactArgs(0),
		Example: `- 添加module config, 参数文件为module_config.json
  - ysadmin license add moduleconfig -F module_config.json`,
	}

	cmd.Flags().StringVarP(&o.JsonFile, "file", "F", "", "json文件")
	cmd.MarkFlagRequired("file")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return licAddModuleConfig(cmd, args, o)
	}

	return cmd
}

func newLicGetCommand(o LicOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "获取license server配置",
		Long:  "获取license server配置",
		RunE:  helpRun,
	}

	cmd.AddCommand(
		newLicGetLicManagerCommand(o),
	)
	return cmd
}

func newLicGetLicManagerCommand(o LicOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "licmanager",
		Short: "获取license manager",
		Long:  "获取license manager",
		Args:  cobra.ExactArgs(0),
		Example: `- 获取ID为52Zzc3ycfEU的license manager
  - ysadmin license get licmanager -I 52Zzc3ycfEU`,
	}

	cmd.Flags().StringVarP(&o.Id, "id", "I", "", "license manager id")
	cmd.MarkFlagRequired("id")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return licGetLicManager(cmd, args, o)
	}

	return cmd
}

func newLicListCommand(o LicOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "列出license server配置",
		Long:  "列出license server配置",
		RunE:  helpRun,
	}

	cmd.AddCommand(
		newLicListLicManagerCommand(o),
	)
	return cmd
}

func newLicListLicManagerCommand(o LicOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "licmanager",
		Short: "列出license manager",
		Long:  "列出license manager",
		Args:  cobra.ExactArgs(0),
		Example: `- 列出所有license manager
  - ysadmin license list licmanager`,
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return licListLicManager(cmd, args)
	}

	return cmd
}

func newLicDeleteCommand(o LicOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "删除license server配置",
		Long:  "删除license server配置",
		RunE:  helpRun,
	}

	cmd.AddCommand(
		newLicDeleteLicManagerCommand(o),
		newLicDeleteLicInfoCommand(o),
		newLicDeleteModuleConfigCommand(o),
	)
	return cmd
}

func newLicDeleteLicManagerCommand(o LicOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "licmanager",
		Short: "删除license manager",
		Long:  "删除license manager",
		Args:  cobra.ExactArgs(0),
		Example: `- 删除ID为52Zzc3ycfEU的license manager
  - ysadmin license delete licmanager -I 52Zzc3ycfEU`,
	}

	cmd.Flags().StringVarP(&o.Id, "id", "I", "", "license manager id")
	cmd.MarkFlagRequired("id")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return licDeleteLicManager(cmd, args, o)
	}

	return cmd
}

func newLicDeleteLicInfoCommand(o LicOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "licinfo",
		Short: "删除license info",
		Long:  "删除license info",
		Args:  cobra.ExactArgs(0),
		Example: `- 删除ID为52Zzc3ycfEU的license info
  - ysadmin license delete licinfo -I 52Zzc3ycfEU`,
	}

	cmd.Flags().StringVarP(&o.Id, "id", "I", "", "license info id")
	cmd.MarkFlagRequired("id")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return licDeleteLicInfo(cmd, args, o)
	}

	return cmd
}

func newLicDeleteModuleConfigCommand(o LicOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "moduleconfig",
		Short: "删除module config",
		Long:  "删除module config",
		Args:  cobra.ExactArgs(0),
		Example: `- 删除ID为52Zzc3ycfEU的module config
  - ysadmin license delete moduleconfig -I 52Zzc3ycfEU`,
	}

	cmd.Flags().StringVarP(&o.Id, "id", "I", "", "module config id")
	cmd.MarkFlagRequired("id")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return licDeleteModuleConfig(cmd, args, o)
	}

	return cmd
}

func newLicPutCommand(o LicOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "put",
		Short: "修改license server配置",
		Long:  "修改license server配置",
		RunE:  helpRun,
	}

	cmd.AddCommand(
		newLicPutLicManagerCommand(o),
		newLicPutLicInfoCommand(o),
		newLicPutModuleConfigCommand(o),
	)
	return cmd
}

func newLicPutLicManagerCommand(o LicOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "licmanager",
		Short: "修改license manager",
		Long:  "修改license manager",
		Args:  cobra.ExactArgs(0),
		Example: `- 修改ID为52Zzc3ycfEU的license manager, 参数文件为lic_manager.json
  - ysadmin license put licmanager -I 52Zzc3ycfEU -F lic_manager.json`,
	}

	cmd.Flags().StringVarP(&o.Id, "id", "I", "", "license manager id")
	cmd.MarkFlagRequired("id")

	cmd.Flags().StringVarP(&o.JsonFile, "file", "F", "", "json文件")
	cmd.MarkFlagRequired("file")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return licPutLicManager(cmd, args, o)
	}

	return cmd
}

func newLicPutLicInfoCommand(o LicOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "licinfo",
		Short: "修改license info",
		Long:  "修改license info",
		Args:  cobra.ExactArgs(0),
		Example: `- 修改ID为52Zzc3ycfEU的license info, 参数文件为lic_info.json
  - ysadmin license put licinfo -I 52Zzc3ycfEU -F lic_info.json`,
	}

	cmd.Flags().StringVarP(&o.Id, "id", "I", "", "license info id")
	cmd.MarkFlagRequired("id")

	cmd.Flags().StringVarP(&o.JsonFile, "file", "F", "", "json文件")
	cmd.MarkFlagRequired("file")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return licPutLicInfo(cmd, args, o)
	}

	return cmd
}

func newLicPutModuleConfigCommand(o LicOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "moduleconfig",
		Short: "修改module config",
		Long:  "修改module config",
		Args:  cobra.ExactArgs(0),
		Example: `- 修改ID为52Zzc3ycfEU的module config, 参数文件为module_config.json
  - ysadmin license put moduleconfig -I 52Zzc3ycfEU -F module_config.json`,
	}

	cmd.Flags().StringVarP(&o.Id, "id", "I", "", "module config id")
	cmd.MarkFlagRequired("id")

	cmd.Flags().StringVarP(&o.JsonFile, "file", "F", "", "json文件")
	cmd.MarkFlagRequired("file")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return licPutModuleConfig(cmd, args, o)
	}

	return cmd
}

func newLicExampleCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "example",
		Short: "输出配置文件示例",
		Long:  "输出配置文件示例",
		RunE:  helpRun,
	}

	cmd.AddCommand(
		newLicManagerExampleCommand(),
		newLicInfoExampleCommand(),
		newModuleConfigExampleCommand(),
	)
	return cmd
}

func newLicManagerExampleCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "licmanager",
		Short: "输出license manager配置文件示例",
		Long:  "输出license manager配置文件示例",
		Example: `- 输出license manager配置文件示例
  - ysadmin license example licmanager`,
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return licExample(cmd, args, licManagerExampleFile, licManagerExampleFileContent)
	}

	return cmd
}

func newLicInfoExampleCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "licinfo",
		Short: "输出license info配置文件示例",
		Long:  "输出license info配置文件示例",
		Example: `- 输出license info配置文件示例
  - ysadmin license example licinfo`,
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return licExample(cmd, args, licInfoExampleFile, licInfoExampleFileContent)
	}

	return cmd
}

func newModuleConfigExampleCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "moduleconfig",
		Short: "输出module config配置文件示例",
		Long:  "输出module config配置文件示例",
		Example: `- 输出module config配置文件示例
  - ysadmin license example moduleconfig`,
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return licExample(cmd, args, moduleConfigExampleFile, moduleConfigExampleFileContent)
	}

	return cmd
}

func licAddLicManager(cmd *cobra.Command, args []string, o LicOptions) error {
	req := new(licmanager.AddLicManagerRequest)
	err := ReadAndUnmarshal(o.JsonFile, req)
	if err != nil {
		fmt.Printf("read json file error: %v\n", err)
		return nil
	}

	res, err := GetYsClient().License.AddLicenseManager(
		GetYsClient().License.AddLicenseManager.Os(req.Os),
		GetYsClient().License.AddLicenseManager.AppType(req.AppType),
		GetYsClient().License.AddLicenseManager.Desc(req.Desc),
		GetYsClient().License.AddLicenseManager.ComputeRule(req.ComputeRule),
	)
	PrintResp(res, err, "Add License Manager")
	return nil
}

func licAddLicInfo(cmd *cobra.Command, args []string, o LicOptions) error {
	req := new(licenseinfo.AddLicenseInfoRequest)
	err := ReadAndUnmarshal(o.JsonFile, req)
	if err != nil {
		fmt.Printf("read json file error: %v\n", err)
		return nil
	}
	res, err := GetYsClient().License.AddLicenseInfo(
		GetYsClient().License.AddLicenseInfo.ManagerId(req.ManagerId),
		GetYsClient().License.AddLicenseInfo.Provider(req.Provider),
		GetYsClient().License.AddLicenseInfo.MacAddr(req.MacAddr),
		GetYsClient().License.AddLicenseInfo.ToolPath(req.ToolPath),
		GetYsClient().License.AddLicenseInfo.LicenseUrl(req.LicenseUrl),
		GetYsClient().License.AddLicenseInfo.Port(req.Port),
		GetYsClient().License.AddLicenseInfo.LicenseAddresses(req.LicenseProxies),
		GetYsClient().License.AddLicenseInfo.LicenseNum(req.LicenseNum),
		GetYsClient().License.AddLicenseInfo.Weight(req.Weight),
		GetYsClient().License.AddLicenseInfo.BeginTime(req.BeginTime),
		GetYsClient().License.AddLicenseInfo.EndTime(req.EndTime),
		GetYsClient().License.AddLicenseInfo.Auth(req.Auth),
		GetYsClient().License.AddLicenseInfo.LicenseEnvVar(req.LicenseEnvVar),
		GetYsClient().License.AddLicenseInfo.LicenseType(req.LicenseType),
		GetYsClient().License.AddLicenseInfo.HpcEndpoint(req.HpcEndpoint),
		GetYsClient().License.AddLicenseInfo.AllowableHpcEndpoints(req.AllowableHpcEndpoints),
		GetYsClient().License.AddLicenseInfo.CollectorType(req.CollectorType),
		GetYsClient().License.AddLicenseInfo.EnableScheduling(req.EnableScheduling),
	)
	PrintResp(res, err, "Add License Info")
	return nil
}

func licAddModuleConfig(cmd *cobra.Command, args []string, o LicOptions) error {
	req := new(moduleconfig.AddModuleConfigRequest)
	err := ReadAndUnmarshal(o.JsonFile, req)
	if err != nil {
		fmt.Printf("read json file error: %v\n", err)
		return nil
	}
	res, err := GetYsClient().License.AddModuleConfig(
		GetYsClient().License.AddModuleConfig.LicenseId(req.LicenseId),
		GetYsClient().License.AddModuleConfig.ModuleName(req.ModuleName),
		GetYsClient().License.AddModuleConfig.Total(req.Total),
	)
	PrintResp(res, err, "Add Module Config")
	return nil
}

func licGetLicManager(cmd *cobra.Command, args []string, o LicOptions) error {
	res, err := GetYsClient().License.GetLicenseManager(
		GetYsClient().License.GetLicenseManager.Id(o.Id),
	)
	PrintResp(res, err, "Get License Manager")
	return nil
}

func licListLicManager(cmd *cobra.Command, args []string) error {
	res, err := GetYsClient().License.ListLicenseManager()
	PrintResp(res, err, "List License Manager")
	return nil
}

func licDeleteLicManager(cmd *cobra.Command, args []string, o LicOptions) error {
	res, err := GetYsClient().License.DeleteLicenseManager(
		GetYsClient().License.DeleteLicenseManager.Id(o.Id),
	)
	PrintResp(res, err, "Delete License Manager"+o.Id)
	return nil
}

func licDeleteLicInfo(cmd *cobra.Command, args []string, o LicOptions) error {
	res, err := GetYsClient().License.DeleteLicenseInfo(
		GetYsClient().License.DeleteLicenseInfo.Id(o.Id),
	)
	PrintResp(res, err, "Delete License Info"+o.Id)
	return nil
}

func licDeleteModuleConfig(cmd *cobra.Command, args []string, o LicOptions) error {
	res, err := GetYsClient().License.DeleteModuleConfig(
		GetYsClient().License.DeleteModuleConfig.Id(o.Id),
	)
	PrintResp(res, err, "Delete Module Config"+o.Id)
	return nil
}

func licPutLicManager(cmd *cobra.Command, args []string, o LicOptions) error {
	req := new(licmanager.PutLicManagerRequest)
	err := ReadAndUnmarshal(o.JsonFile, req)
	if err != nil {
		fmt.Printf("read json file error: %v\n", err)
		return nil
	}

	res, err := GetYsClient().License.PutLicenseManager(
		GetYsClient().License.PutLicenseManager.Id(o.Id),
		GetYsClient().License.PutLicenseManager.Os(req.Os),
		GetYsClient().License.PutLicenseManager.AppType(req.AppType),
		GetYsClient().License.PutLicenseManager.Desc(req.Desc),
		GetYsClient().License.PutLicenseManager.Status(req.Status),
		GetYsClient().License.PutLicenseManager.ComputeRule(req.ComputeRule),
	)
	PrintResp(res, err, "Put License Manager")
	return nil
}

func licPutLicInfo(cmd *cobra.Command, args []string, o LicOptions) error {
	req := new(licenseinfo.PutLicenseInfoRequest)
	err := ReadAndUnmarshal(o.JsonFile, req)
	if err != nil {
		fmt.Printf("read json file error: %v\n", err)
		return nil
	}

	res, err := GetYsClient().License.PutLicenseInfo(
		GetYsClient().License.PutLicenseInfo.Id(o.Id),
		GetYsClient().License.PutLicenseInfo.ManagerId(req.ManagerId),
		GetYsClient().License.PutLicenseInfo.Provider(req.Provider),
		GetYsClient().License.PutLicenseInfo.MacAddr(req.MacAddr),
		GetYsClient().License.PutLicenseInfo.ToolPath(req.ToolPath),
		GetYsClient().License.PutLicenseInfo.LicenseUrl(req.LicenseUrl),
		GetYsClient().License.PutLicenseInfo.Port(req.Port),
		GetYsClient().License.PutLicenseInfo.LicenseAddresses(req.LicenseProxies),
		GetYsClient().License.PutLicenseInfo.LicenseNum(req.LicenseNum),
		GetYsClient().License.PutLicenseInfo.Weight(req.Weight),
		GetYsClient().License.PutLicenseInfo.BeginTime(req.BeginTime),
		GetYsClient().License.PutLicenseInfo.EndTime(req.EndTime),
		GetYsClient().License.PutLicenseInfo.Auth(req.Auth),
		GetYsClient().License.PutLicenseInfo.LicenseEnvVar(req.LicenseEnvVar),
		GetYsClient().License.PutLicenseInfo.LicenseType(req.LicenseType),
		GetYsClient().License.PutLicenseInfo.HpcEndpoint(req.HpcEndpoint),
		GetYsClient().License.PutLicenseInfo.AllowableHpcEndpoints(req.AllowableHpcEndpoints),
		GetYsClient().License.PutLicenseInfo.CollectorType(req.CollectorType),
		GetYsClient().License.PutLicenseInfo.EnableScheduling(req.EnableScheduling),
	)
	PrintResp(res, err, "Put License Info")
	return nil
}

func licPutModuleConfig(cmd *cobra.Command, args []string, o LicOptions) error {
	req := new(moduleconfig.PutModuleConfigRequest)
	err := ReadAndUnmarshal(o.JsonFile, req)
	if err != nil {
		fmt.Printf("read json file error: %v\n", err)
		return nil
	}

	res, err := GetYsClient().License.PutModuleConfig(
		GetYsClient().License.PutModuleConfig.Id(o.Id),
		GetYsClient().License.PutModuleConfig.ModuleName(req.ModuleName),
		GetYsClient().License.PutModuleConfig.Total(req.Total),
	)
	PrintResp(res, err, "Put Module Config")
	return nil
}

func licExample(cmd *cobra.Command, args []string, file, content string) error {
	if checkFileExist(file) {
		fmt.Println("示例文件:", file)
		return nil
	}
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err = f.WriteString(content); err != nil {
		return err
	}
	fmt.Println("示例文件:", file)
	return nil
}
