package impl

import (
	"encoding/json"
	"fmt"
	"os"

	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/cobra"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/admin/app/add"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/admin/app/update"
)

const appExampleFile = "app_example.json"

const appExampleFileContent = `{
	"Name": "MyApp 1.0",
	"Type": "Application",
	"Version": "1.0",
	"AppParamsVersion": 1,
	"Image": "myapp-image:v1",
	"Endpoint": "endpoint",
	"Command": "myapp start",
	"Description": "My Sample App",
	"IconUrl": "http://example.com/myapp.png",
	"CoresMaxLimit": 10,
	"CoresPlaceholder": "placeholder",
	"FileFilterRule": "{\"result\": \"\\\\\\\\.dat$\",\"model\": \"\\\\\\\\.(jou|cas)$\",\"log\": \"\\\\\\\\.(sta|dat|msg|out|log)$\",\"middle\": \"\\\\\\\\.(com|prt)$\"}",
	"ResidualEnable": false,
	"ResidualLogRegexp": "stdout.log",
	"ResidualLogParser": "fluent",
	"MonitorChartEnable": false,
	"MonitorChartRegexp": ".*\\\\.out",
	"MonitorChartParser": "fluent",
	"LicenseVars": "LicenseVars",
	"SnapshotEnable": false,
	"BinPath": "{\"az-jinan\":\"real jinan app bin path\",\"az-zhigu\":\"1212\",\"az-zhigu-pbs\":\"wqwq\"}",
	"ExtentionParams": "{\"YS_MAIN_FILE\":{\"Type\":\"String\",\"ReadableName\":\"主文件\",\"Must\":true}}",
	"LicManagerId": "",
	"PublishStatus": "published",
	"NeedLimitCore": false,
	"SpecifyQueue": {
		"az-jinan": "jinan-queue",
		"az-zhigu": "zhigu-queue"
	}
}
`

// AppOptions app options
type AppOptions struct {
	BaseOptions
}

// AddFlags add flags
func (o *AppOptions) AddFlags(cmd *cobra.Command) {
	o.BaseOptions.AddBaseOptions(cmd)
}

func init() {
	RegisterCmd(NewAppCommand())
}

// NewAppCommand new app command
func NewAppCommand() *cobra.Command {
	o := AppOptions{}
	cmd := &cobra.Command{
		Use:   "application",
		Short: "计算应用管理工具",
		Long:  "计算应用管理工具, 用于管理计算应用",
		RunE:  helpRun,
	}

	cmd.AddCommand(
		newAppAddCmd(o),
		newAppGetCmd(o),
		newAppListCmd(o),
		newAppDeleteCmd(o),
		newAppPutCmd(o),
		newAppExampleCmd(o),
		newQuotaCmd(o),
		newAppAllowCmd(o),
		newPublishCmd(o),
	)

	return cmd
}

func newAppAddCmd(o AppOptions) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "add -F file.json",
		Short: "添加计算应用",
		Long:  `添加计算应用, 指定一个.json文件作为参数文件, 可使用'ysadmin application example'生成参考参数文件`,
		Example: `- 添加应用, 参数文件为app.json
  - ysadmin application add -F app.json
- 添加应用, 参数文件为app.json, command文件为script.sh
  - ysadmin application add -F app.json --sh script.sh`,
	}
	var shFile string
	cmd.Flags().StringVarP(&shFile, "sh", "", "", "脚本文件路径,若指定则会覆盖json文件中的command字段")
	cmd.Flags().StringVarP(&o.JsonFile, "file", "F", "", "JSON 文件路径")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		data, err := checkReqJsonFileAndRead(o.JsonFile)
		if err != nil {
			return err
		}
		req := new(add.Request)
		err = jsoniter.Unmarshal(data, req)
		if err != nil {
			return err
		}

		if shFile != "" {
			script, err := os.ReadFile(shFile)
			if err != nil {
				return fmt.Errorf("read script file failed: %w", err)
			}

			req.Command = string(script)
		}

		res, err := GetYsClient().Job.AdminAddAPP(
			GetYsClient().Job.AdminAddAPP.Name(req.Name),
			GetYsClient().Job.AdminAddAPP.Type(req.Type),
			GetYsClient().Job.AdminAddAPP.Version(req.Version),
			GetYsClient().Job.AdminAddAPP.AppParamsVersion(req.AppParamsVersion),
			GetYsClient().Job.AdminAddAPP.Image(req.Image),
			GetYsClient().Job.AdminAddAPP.LicManagerId(req.LicManagerId),
			GetYsClient().Job.AdminAddAPP.Endpoint(req.Endpoint),
			GetYsClient().Job.AdminAddAPP.Command(req.Command),
			GetYsClient().Job.AdminAddAPP.Description(req.Description),
			GetYsClient().Job.AdminAddAPP.IconUrl(req.IconUrl),
			GetYsClient().Job.AdminAddAPP.CoresMaxLimit(req.CoresMaxLimit),
			GetYsClient().Job.AdminAddAPP.CoresPlaceholder(req.CoresPlaceholder),
			GetYsClient().Job.AdminAddAPP.FileFilterRule(req.FileFilterRule),
			GetYsClient().Job.AdminAddAPP.ResidualEnable(req.ResidualEnable),
			GetYsClient().Job.AdminAddAPP.ResidualLogRegexp(req.ResidualLogRegexp),
			GetYsClient().Job.AdminAddAPP.ResidualLogParser(req.ResidualLogParser),
			GetYsClient().Job.AdminAddAPP.MonitorChartEnable(req.MonitorChartEnable),
			GetYsClient().Job.AdminAddAPP.MonitorChartRegexp(req.MonitorChartRegexp),
			GetYsClient().Job.AdminAddAPP.MonitorChartParser(req.MonitorChartParser),
			GetYsClient().Job.AdminAddAPP.LicenseVars(req.LicenseVars),
			GetYsClient().Job.AdminAddAPP.SnapshotEnable(req.SnapshotEnable),
			GetYsClient().Job.AdminAddAPP.BinPath(ToBinPathMap(req.BinPath)),
			GetYsClient().Job.AdminAddAPP.ExtentionParams(req.ExtentionParams),
			GetYsClient().Job.AdminAddAPP.NeedLimitCore(req.NeedLimitCore),
			GetYsClient().Job.AdminAddAPP.SpecifyQueue(req.SpecifyQueue),
		)
		PrintResp(res, err, "Admin Add APP")
		return nil
	}

	return cmd
}

func newAppGetCmd(o AppOptions) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "get -I id",
		Short: "获取计算应用",
		Long:  "获取计算应用",
		Example: `- 获取应用ID为 52XkfrGM9vE 的应用
  - ysadmin application get -I 52XkfrGM9vE
- 获取应用ID为 52XkfrGM9vE 的应用, 并输出到当前目录下的 52XkfrGM9vE.json 文件中
  - ysadmin application get -I 52XkfrGM9vE -O`,
	}
	cmd.Flags().StringVarP(&o.Id, "id", "I", "", "应用 ID")
	cmd.MarkFlagRequired("id")

	cmd.Flags().BoolP("output", "O", false, "输出到当前目录下的 [appID].json 文件中, 用于put命令使用")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		res, err := GetYsClient().Job.AdminGetAPP(
			GetYsClient().Job.AdminGetAPP.AppID(o.Id),
		)
		PrintResp(res, err, "Get App")

		if err == nil && cmd.Flags().Changed("output") {
			name := o.Id + ".json"
			data, err := json.MarshalIndent(res.Data, "", "  ")
			if err != nil {
				fmt.Println("MarshalError: \n", err.Error())
				return nil
			}

			f, err := os.Create(name)
			if err != nil {
				return err
			}
			defer f.Close()

			if _, err = f.Write(data); err != nil {
				return err
			}
			fmt.Println("输出文件:", name)
		}
		return nil
	}

	return cmd
}

func newAppListCmd(o AppOptions) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "list",
		Short: "列出计算应用",
		Long:  "列出计算应用, 可以指定用户 ID 查询该用户有配额的应用列表",
		Example: `- 列出所有应用
  - ysadmin application list
- 列出用户ID为 52XmLrdCkHb 的用户有配额的应用列表
  - ysadmin application list -U 52XmLrdCkHb`,
	}
	// userID
	cmd.Flags().StringP("user_id", "U", "", "用户 ID, 查询该用户有配额的应用列表")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		res, err := GetYsClient().Job.AdminListAPP(
			GetYsClient().Job.AdminListAPP.AllowUserID(cmd.Flag("user_id").Value.String()),
		)
		PrintResp(res, err, "List App")
		return nil
	}

	return cmd
}

func newAppDeleteCmd(o AppOptions) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "delete -I id",
		Short: "删除计算应用",
		Long:  "删除计算应用",
		Example: `- 删除应用ID为 52XkfrGM9vE 的应用
  - ysadmin application delete -I 52XkfrGM9vE`,
	}
	cmd.Flags().StringVarP(&o.Id, "id", "I", "", "应用 ID")
	cmd.MarkFlagRequired("id")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		res, err := GetYsClient().Job.AdminDeleteAPP(
			GetYsClient().Job.AdminDeleteAPP.AppID(o.Id),
		)
		PrintResp(res, err, "Delete App "+o.Id)
		return nil
	}
	return cmd
}

func newAppPutCmd(o AppOptions) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "put -F file.json -I id",
		Short: "更新计算应用",
		Long:  `更新计算应用, 是对应用的全量更新, 未指定的字段会被置空, 指定一个.json文件作为参数文件, 可通过'ysadmin application get -O'获取当前应用的参数文件`,
		Example: `- 更新应用ID为 52XkfrGM9vE 的应用, 参数文件为app.json
  - ysadmin application put -F app.json -I 52XkfrGM9vE
- 更新应用ID为 52XkfrGM9vE 的应用, 参数文件为app.json, command文件为script.sh
  - ysadmin application put -F app.json -I 52XkfrGM9vE --sh script.sh`,
	}

	var shFile string
	cmd.Flags().StringVarP(&shFile, "sh", "", "", "脚本文件路径,若指定则会覆盖json文件中的command字段")
	cmd.Flags().StringVarP(&o.JsonFile, "file", "F", "", "JSON 文件路径")
	cmd.Flags().StringVarP(&o.Id, "id", "I", "", "应用 ID")
	cmd.MarkFlagRequired("id")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		data, err := checkReqJsonFileAndRead(o.JsonFile)
		if err != nil {
			return err
		}
		req := &update.Request{}
		err = jsoniter.Unmarshal(data, req)
		if err != nil {
			return err
		}

		if shFile != "" {
			script, err := os.ReadFile(shFile)
			if err != nil {
				return fmt.Errorf("read script file failed: %w", err)
			}

			req.Command = string(script)
		}

		res, err := GetYsClient().Job.AdminUpdateAPP(
			GetYsClient().Job.AdminUpdateAPP.AppID(o.Id),
			GetYsClient().Job.AdminUpdateAPP.Name(req.Name),
			GetYsClient().Job.AdminUpdateAPP.Type(req.Type),
			GetYsClient().Job.AdminUpdateAPP.Version(req.Version),
			GetYsClient().Job.AdminUpdateAPP.AppParamsVersion(req.AppParamsVersion),
			GetYsClient().Job.AdminUpdateAPP.Image(req.Image),
			GetYsClient().Job.AdminUpdateAPP.LicManagerId(req.LicManagerId),
			GetYsClient().Job.AdminUpdateAPP.Endpoint(req.Endpoint),
			GetYsClient().Job.AdminUpdateAPP.Command(req.Command),
			GetYsClient().Job.AdminUpdateAPP.Description(req.Description),
			GetYsClient().Job.AdminUpdateAPP.IconUrl(req.IconUrl),
			GetYsClient().Job.AdminUpdateAPP.CoresMaxLimit(req.CoresMaxLimit),
			GetYsClient().Job.AdminUpdateAPP.CoresPlaceholder(req.CoresPlaceholder),
			GetYsClient().Job.AdminUpdateAPP.FileFilterRule(req.FileFilterRule),
			GetYsClient().Job.AdminUpdateAPP.ResidualEnable(req.ResidualEnable),
			GetYsClient().Job.AdminUpdateAPP.ResidualLogParser(req.ResidualLogParser),
			GetYsClient().Job.AdminUpdateAPP.ResidualLogRegexp(req.ResidualLogRegexp),
			GetYsClient().Job.AdminUpdateAPP.MonitorChartEnable(req.MonitorChartEnable),
			GetYsClient().Job.AdminUpdateAPP.MonitorChartRegexp(req.MonitorChartRegexp),
			GetYsClient().Job.AdminUpdateAPP.MonitorChartParser(req.MonitorChartParser),
			GetYsClient().Job.AdminUpdateAPP.LicenseVars(req.LicenseVars),
			GetYsClient().Job.AdminUpdateAPP.SnapshotEnable(req.SnapshotEnable),
			GetYsClient().Job.AdminUpdateAPP.BinPath(ToBinPathMap(req.BinPath)),
			GetYsClient().Job.AdminUpdateAPP.ExtentionParams(req.ExtentionParams),
			GetYsClient().Job.AdminUpdateAPP.PublishStatus(req.PublishStatus),
			GetYsClient().Job.AdminUpdateAPP.NeedLimitCore(req.NeedLimitCore),
			GetYsClient().Job.AdminUpdateAPP.SpecifyQueue(req.SpecifyQueue),
		)
		PrintResp(res, err, "Put App")
		return nil
	}

	return cmd
}

func newPublishCmd(o AppOptions) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "publish -I id [-D]",
		Short: "计算应用发布/取消发布",
		Long:  "计算应用发布/取消发布",
		Example: `- 发布应用ID为 52XkfrGM9vE 的应用
  - ysadmin application publish -I 52XkfrGM9vE
- 取消发布应用ID为 52XkfrGM9vE 的应用
  - ysadmin application publish -I 52XkfrGM9vE -D`,
	}

	cmd.Flags().StringVarP(&o.Id, "id", "I", "", "应用 ID")
	cmd.MarkFlagRequired("id")

	cmd.Flags().BoolP("down", "D", false, "取消发布")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		res, err := GetYsClient().Job.AdminGetAPP(
			GetYsClient().Job.AdminGetAPP.AppID(o.Id),
		)
		if err != nil {
			fmt.Printf("Get App Error:\n%s\n", err.Error())
			return nil
		}

		app := res.Data

		publishStatus := func() update.Status {
			if cmd.Flags().Changed("down") {
				return update.Unpublished
			}

			return update.Published
		}()

		app.PublishStatus = string(publishStatus)
		binPath := app.BinPath

		_, err = GetYsClient().Job.AdminUpdateAPP(
			GetYsClient().Job.AdminUpdateAPP.AppID(app.AppID),
			GetYsClient().Job.AdminUpdateAPP.Name(app.Name),
			GetYsClient().Job.AdminUpdateAPP.Type(app.Type),
			GetYsClient().Job.AdminUpdateAPP.Version(app.Version),
			GetYsClient().Job.AdminUpdateAPP.AppParamsVersion(app.AppParamsVersion),
			GetYsClient().Job.AdminUpdateAPP.Image(app.Image),
			GetYsClient().Job.AdminUpdateAPP.LicManagerId(app.LicManagerId),
			GetYsClient().Job.AdminUpdateAPP.Endpoint(app.Endpoint),
			GetYsClient().Job.AdminUpdateAPP.Command(app.Command),
			GetYsClient().Job.AdminUpdateAPP.Description(app.Description),
			GetYsClient().Job.AdminUpdateAPP.IconUrl(app.IconUrl),
			GetYsClient().Job.AdminUpdateAPP.CoresMaxLimit(app.CoresMaxLimit),
			GetYsClient().Job.AdminUpdateAPP.CoresPlaceholder(app.CoresPlaceholder),
			GetYsClient().Job.AdminUpdateAPP.FileFilterRule(app.FileFilterRule),
			GetYsClient().Job.AdminUpdateAPP.ResidualEnable(app.ResidualEnable),
			GetYsClient().Job.AdminUpdateAPP.ResidualLogRegexp(app.ResidualLogRegexp),
			GetYsClient().Job.AdminUpdateAPP.ResidualLogParser(app.ResidualLogParser),
			GetYsClient().Job.AdminUpdateAPP.MonitorChartEnable(app.MonitorChartEnable),
			GetYsClient().Job.AdminUpdateAPP.MonitorChartRegexp(app.MonitorChartRegexp),
			GetYsClient().Job.AdminUpdateAPP.MonitorChartParser(app.MonitorChartParser),
			GetYsClient().Job.AdminUpdateAPP.LicenseVars(app.LicenseVars),
			GetYsClient().Job.AdminUpdateAPP.SnapshotEnable(app.SnapshotEnable),
			GetYsClient().Job.AdminUpdateAPP.BinPath(ToBinPathMap(binPath)),
			GetYsClient().Job.AdminUpdateAPP.ExtentionParams(app.ExtentionParams),
			GetYsClient().Job.AdminUpdateAPP.PublishStatus(update.Status(app.PublishStatus)),
			GetYsClient().Job.AdminUpdateAPP.NeedLimitCore(app.NeedLimitCore),
			GetYsClient().Job.AdminUpdateAPP.SpecifyQueue(app.SpecifyQueue),
		)
		if err != nil {
			fmt.Printf("%s App Error:\n%s\n", publishStatus, err.Error())
			return nil
		}

		fmt.Printf("%s App Success:\n", publishStatus)

		data, err := json.MarshalIndent(app, "", "  ")
		if err != nil {
			fmt.Println("MarshalError: \n", err.Error())
			return nil
		}
		fmt.Println(string(data))

		return nil
	}

	return cmd
}

func newAppExampleCmd(o AppOptions) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "example",
		Short: "创建计算应用示例文件",
		Long:  "创建计算应用示例文件",
		Example: `- 创建计算应用示例文件
  - ysadmin application example`,
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if checkFileExist(appExampleFile) {
			fmt.Println("示例文件:", appExampleFile)
			return nil
		}

		f, err := os.Create(appExampleFile)
		if err != nil {
			return err
		}
		defer f.Close()

		if _, err = f.WriteString(appExampleFileContent); err != nil {
			return err
		}
		fmt.Println("示例文件:", appExampleFile)
		return nil
	}

	return cmd
}

func newQuotaCmd(o AppOptions) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "quota",
		Short: "计算应用配额管理工具",
		Long:  "计算应用配额管理工具, 用于管理计算应用配额, 配额是指用户对应用的使用权限",
		RunE:  helpRun,
	}

	QuotaOptions := QuotaOptions{
		AppOptions: o,
	}
	cmd.AddCommand(
		newQuotaAddCmd(QuotaOptions),
		newQuotaGetCmd(QuotaOptions),
		newQuotaDeleteCmd(QuotaOptions),
	)

	return cmd
}

// QuotaOptions quota options
type QuotaOptions struct {
	AppOptions
	UserID string
}

func newQuotaAddCmd(o QuotaOptions) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "add -I app_id -U user_id",
		Short: "添加计算应用配额",
		Long:  "添加计算应用配额",
		Example: `- 为用户ID为 52XmLrdCkHb 的用户添加应用ID为 52XkfrGM9vE 的应用配额
  - ysadmin application quota add -I 52XkfrGM9vE -U 52XmLrdCkHb`,
	}
	cmd.Flags().StringVarP(&o.Id, "id", "I", "", "应用 ID")
	cmd.Flags().StringVarP(&o.UserID, "user_id", "U", "", "用户 ID")
	cmd.MarkFlagRequired("id")
	cmd.MarkFlagRequired("user_id")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		res, err := GetYsClient().Job.AdminAddAPPQuota(
			GetYsClient().Job.AdminAddAPPQuota.AppID(o.Id),
			GetYsClient().Job.AdminAddAPPQuota.UserID(o.UserID),
		)
		PrintResp(res, err, "Add App Quota")
		return nil
	}

	return cmd
}

func newQuotaGetCmd(o QuotaOptions) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "get -I app_id -U user_id",
		Short: "获取计算应用配额",
		Long:  "获取计算应用配额",
		Example: `- 查看用户ID为 52XmLrdCkHb 的用户是否有应用ID为 52XkfrGM9vE 的应用配额
  - ysadmin application quota get -I 52XkfrGM9vE -U 52XmLrdCkHb`,
	}
	cmd.Flags().StringVarP(&o.Id, "id", "I", "", "应用 ID")
	cmd.Flags().StringVarP(&o.UserID, "user_id", "U", "", "用户 ID")
	cmd.MarkFlagRequired("id")
	cmd.MarkFlagRequired("user_id")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		res, err := GetYsClient().Job.AdminGetAPPQuota(
			GetYsClient().Job.AdminGetAPPQuota.AppID(o.Id),
			GetYsClient().Job.AdminGetAPPQuota.UserID(o.UserID),
		)
		PrintResp(res, err, "Get App Quota")
		return nil
	}

	return cmd
}

func newQuotaDeleteCmd(o QuotaOptions) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "delete -I app_id -U user_id",
		Short: "删除计算应用配额",
		Long:  "删除计算应用配额",
		Example: `- 删除用户ID为 52XmLrdCkHb 的用户的应用ID为 52XkfrGM9vE 的应用配额
  - ysadmin application quota delete -I 52XkfrGM9vE -U 52XmLrdCkHb`,
	}
	cmd.Flags().StringVarP(&o.Id, "id", "I", "", "应用 ID")
	cmd.Flags().StringVarP(&o.UserID, "user_id", "U", "", "用户 ID")
	cmd.MarkFlagRequired("id")
	cmd.MarkFlagRequired("user_id")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		res, err := GetYsClient().Job.AdminDeleteAPPQuota(
			GetYsClient().Job.AdminDeleteAPPQuota.AppID(o.Id),
			GetYsClient().Job.AdminDeleteAPPQuota.UserID(o.UserID),
		)
		PrintResp(res, err, "Delete App Quota")
		return nil
	}

	return cmd
}

func newAppAllowCmd(o AppOptions) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "allow",
		Short: "计算应用白名单工具",
		Long:  "计算应用白名单工具, 用于管理计算应用白名单, 在白名单中的应用，所有用户都可以使用",
		RunE:  helpRun,
	}
	cmd.AddCommand(
		newAppAllowAddCmd(o),
		newAppAllowGetCmd(o),
		newAppAllowDeleteCmd(o),
	)
	return cmd
}

func newAppAllowAddCmd(o AppOptions) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "add -I app_id",
		Short: "添加计算应用白名单",
		Long:  "添加计算应用白名单",
		Example: `- 添加应用ID为 52XkfrGM9vE 的白名单
  - ysadmin application allow add -I 52XkfrGM9vEb`,
	}
	cmd.Flags().StringVarP(&o.Id, "id", "I", "", "应用 ID")
	cmd.MarkFlagRequired("id")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		res, err := GetYsClient().Job.AdminAddAppAllow(
			GetYsClient().Job.AdminAddAppAllow.AppID(o.Id),
		)
		PrintResp(res, err, "Add App Allow")
		return nil
	}

	return cmd
}

func newAppAllowGetCmd(o AppOptions) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "get -I app_id",
		Short: "获取计算应用白名单",
		Long:  "获取计算应用白名单",
		Example: `- 查看是否有应用ID为 52XkfrGM9vE 的应用白名单
  - ysadmin application allow get -I 52XkfrGM9vE`,
	}
	cmd.Flags().StringVarP(&o.Id, "id", "I", "", "应用 ID")
	cmd.MarkFlagRequired("id")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		res, err := GetYsClient().Job.AdminGetAppAllow(
			GetYsClient().Job.AdminGetAppAllow.AppID(o.Id),
		)
		PrintResp(res, err, "Get App Allow")
		return nil
	}

	return cmd
}

func newAppAllowDeleteCmd(o AppOptions) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "delete -I app_id",
		Short: "删除计算应用白名单",
		Long:  "删除计算应用白名单",
		Example: `- 删除应用ID为 52XkfrGM9vE 的应用白名单
  - ysadmin application allow delete -I 52XkfrGM9vE`,
	}
	cmd.Flags().StringVarP(&o.Id, "id", "I", "", "应用 ID")
	cmd.MarkFlagRequired("id")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		res, err := GetYsClient().Job.AdminDeleteAppAllow(
			GetYsClient().Job.AdminDeleteAppAllow.AppID(o.Id),
		)
		PrintResp(res, err, "Delete App Allow")
		return nil
	}

	return cmd
}

// ToBinPathMap binPath to map,if binPath not exist return nil
func ToBinPathMap(binPath string) map[string]string {
	if len(binPath) == 0 {
		return nil
	}
	binPathMap := make(map[string]string)
	err := json.Unmarshal([]byte(binPath), &binPathMap)
	if err != nil {
		panic(err)
	}
	return binPathMap
}
