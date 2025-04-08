package impl

import (
	"fmt"
	"os"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/cobra"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/admin/jobcreate"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/jobpreschedule"
)

const submitJobExampleFile = "submit_example.json"
const scheduleJobExampleFile = "preschedule_example.json"

const submitJobExampleFileContent = `
{
    "Name":"SampleJob",
    "Params":{
        "Application":{
            "Command":"some command; some other command; $SomeEnvVar; $YS_MAIN_FILE",
            "AppID":"solver123"
        },
        "Resource":{
            "Cores":8,
            "Memory":512
        },
        "EnvVars":{
            "SomeEnvVar":"some value, EnvVars are used for application commands. like $SomeEnvVar, $YS_MAIN_FILE etc.",
            "YS_MAIN_FILE":"some_file.sh"
        },
        "Input":{
            "Type":"cloud_storage",
            "Source":"https://storage.example.com/input_data",
            "Destination":"jobs/sample/input_data"
        },
        "Output":{
            "Type":"cloud_storage",
            "Address":"https://storage.example.com/output_data",
            "NoNeededPaths":"some_path"
        },
        "TmpWorkdir":true,
        "SubmitWithSuspend":false,
        "CustomStateRule":{
            "KeyStatement":"error",
            "ResultState":"failed"
        }
    },
    "Timeout":3600,
    "Zone":"az-zhigu",
    "Comment":"Testing job creation",
    "ChargeParam":{
        "ChargeType":"PostPaid",
        "PeriodType":"hour",
        "PeriodNum":1
    },
    "NoRound":false
}
`

const scheduleJobExampleFileContent = `
{
   "Params":{
       "Application":{
           "Command":"some command; some other command; $SomeEnvVar; $YS_MAIN_FILE",
           "AppID":"solver123"
       },
       "Resource":{
           "MinCores":8,
           "MaxCores":128,
           "Memory":512
       },
       "EnvVars":{
           "SomeEnvVar":"some value, EnvVars are used for application commands. like $SomeEnvVar, $YS_MAIN_FILE etc.",
           "YS_MAIN_FILE":"some_file.sh"
       }
  }
}`

// JobOptions 作业命令的参数
type JobOptions struct {
	BaseOptions
	UserID         string
	AppID          string
	WithDelete     bool
	IsSystemFailed bool
	FileSyncState  string
	Resource       Resource
	Zones          []string
	Shared         bool
	Fixed          bool
}

type Resource struct {
	MinCores *int `json:"MinCores" binding:"required"` //期望的最小核数
	MaxCores *int `json:"MaxCores" binding:"required"` //期望的最大核数
	Memory   *int `json:"Memory"`                      //期望的内存数，单位为M，暂未起作用
}

func init() {
	RegisterCmd(NewJobCommand())
}

// NewJobCommand 创建作业命令
func NewJobCommand() *cobra.Command {
	o := JobOptions{}
	cmd := &cobra.Command{
		Use:   "job",
		Short: "作业管理",
		Long:  "作业管理，用于提交作业，查询作业，删除作业，终止作业，重传作业结果文件",
		RunE:  helpRun,
	}
	cmd.AddCommand(
		newJobSubmitCmd(o),
		newJobGetCmd(o),
		newJobListCmd(o),
		newJobDeleteCmd(o),
		newJobTerminateCmd(o),
		newJobRetransmitCmd(o),
		newJobExampleCmd(o),
		newJobUpdateCmd(o),
		newJobPreSchedule(o),
		newParamsExampleCmd(o),
	)
	return cmd
}

func newJobSubmitCmd(o JobOptions) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "submit",
		Short: "提交作业",
		Long:  "提交作业, 指定一个.json文件作为参数文件, 可使用ysadmin job example生成参考参数文件",
		Example: `- 提交作业, 参数文件为job.json
  - ysadmin job submit -F job.json
- 提交作业, 参数文件为job.json, command文件为script.sh
  - ysadmin job submit -F job.json --sh script.sh`,
	}

	var shFile string
	cmd.Flags().StringVarP(&o.JsonFile, "file", "F", "", "JSON 文件路径")
	cmd.MarkFlagRequired("file")
	cmd.Flags().StringVarP(&shFile, "sh", "", "", "command文件路径,若指定则会覆盖json文件中的command字段")
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		data, err := checkReqJsonFileAndRead(o.JsonFile)
		if err != nil {
			return err
		}
		fmt.Printf("data=%v", string(data))
		req := new(jobcreate.Request)
		err = jsoniter.Unmarshal(data, req)
		if err != nil {
			return err
		}
		if shFile != "" {
			script, err := os.ReadFile(shFile)
			if err != nil {
				return fmt.Errorf("read script file failed: %w", err)
			}

			req.Params.Application.Command = string(script)
		}
		fmt.Printf("jobcreate.Request.Name=%v\n", req.Name)
		fmt.Printf("jobcreate.Request.Comment=%v\n", req.Comment)
		fmt.Printf("jobcreate.Request.PreScheduleID=%v\n", req.PreScheduleID)
		res, err := GetYsClient().Job.AdminJobCreate(
			GetYsClient().Job.AdminJobCreate.Name(req.Name),
			GetYsClient().Job.AdminJobCreate.Params(req.Params),
			GetYsClient().Job.AdminJobCreate.Timeout(req.Timeout),
			GetYsClient().Job.AdminJobCreate.Zone(req.Zone),
			GetYsClient().Job.AdminJobCreate.Comment(req.Comment),
			GetYsClient().Job.AdminJobCreate.Queue(req.Queue),
			GetYsClient().Job.AdminJobCreate.ChargeParams(req.ChargeParam),
			GetYsClient().Job.AdminJobCreate.NoRound(req.NoRound),
			GetYsClient().Job.AdminJobCreate.AllocType(req.AllocType),
			GetYsClient().Job.AdminJobCreate.PreScheduleID(req.PreScheduleID),
			GetYsClient().Job.AdminJobCreate.PayBy(req.PayBy),
		)
		PrintResp(res, err, "Admin Add Job")
		return nil
	}
	return cmd
}

func newJobGetCmd(o JobOptions) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "get",
		Short: "查询作业",
		Long:  "查询作业",
		Example: `- 查询作业ID为 52hgnbSTGWd 的作业
  - ysadmin job get -I 52hgnbSTGWd`,
	}
	cmd.Flags().StringVarP(&o.Id, "id", "I", "", "作业ID")
	cmd.MarkFlagRequired("id")
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		// 获取作业详情
		jobRes, jobErr := GetYsClient().Job.AdminJobGet(
			GetYsClient().Job.AdminJobGet.JobId(o.Id),
		)
		// 获取作业CPU使用情况
		cpuRes, cpuErr := GetYsClient().Job.AdminJobCpuUsage(
			GetYsClient().Job.AdminJobCpuUsage.JobId(o.Id),
		)
		// 打印作业详情结果
		PrintResp(jobRes, jobErr, "Get Job "+o.Id)
		// 如果有CPU使用情况的错误，则打印错误信息
		if cpuErr != nil {
			PrintResp(nil, cpuErr, "Get Job CPU Usage "+o.Id)
		} else {
			// 打印作业CPU使用情况结果
			PrintResp(cpuRes, nil, "Get Job CPU Usage "+o.Id)
		}
		return nil
	}
	return cmd
}

func newJobListCmd(o JobOptions) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "list",
		Short: "查询作业列表",
		Long:  "查询作业列表, 可以指定作业状态, 区域, 偏移量, 限制条数, 用户ID, 默认偏移量0, 限制条数1000",
		Example: `- 查询所有作业(默认1000条)
  - ysadmin job list
- 查询分区为az-zhigu的所有运行中的作业, 偏移量为0, 限制条数为1000
  - ysadmin job list -S running -Z az-zhigu -O 0 -L 1000
- 查询用户ID为4TiSsZonTa3的所有作业
  - ysadmin job list -U 4TiSsZonTa3`,
	}
	cmd.Flags().StringVarP(&o.State, "state", "S", "", "作业状态, 可选值: Initiated, InitiallySuspended, Pending, Running, Suspending, Suspended, Terminating, Terminated, Completed, Failed")
	cmd.Flags().StringVarP(&o.Zone, "zone", "Z", "", "指定区域, 例如az-jinan, az-zhigu。可以通过ysadmin zone list查看所有的区域信息")
	cmd.Flags().Int64VarP(&o.Offset, "offset", "O", 0, "指定的开始偏移量, 从0开始")
	cmd.Flags().Int64VarP(&o.Limit, "limit", "L", 1000, "指定的条数上限, 最多1000")
	cmd.Flags().StringVarP(&o.UserID, "user_id", "U", "", "用户ID, 例如4TiSsZonTa3")
	cmd.Flags().StringVarP(&o.AppID, "app_id", "A", "", "应用ID")
	cmd.Flags().BoolVarP(&o.WithDelete, "with_delete", "", true, "是否包含已删除的作业, 默认包含")
	cmd.Flags().BoolVarP(&o.IsSystemFailed, "is_system_failed", "", false, "仅查询系统失败的作业")
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		res, err := GetYsClient().Job.AdminJobList(
			GetYsClient().Job.AdminJobList.JobState(o.State),
			GetYsClient().Job.AdminJobList.Zone(o.Zone),
			GetYsClient().Job.AdminJobList.PageOffset(o.Offset),
			GetYsClient().Job.AdminJobList.PageSize(o.Limit),
			GetYsClient().Job.AdminJobList.UserID(o.UserID),
			GetYsClient().Job.AdminJobList.AppID(o.AppID),
			GetYsClient().Job.AdminJobList.WithDelete(o.WithDelete),
			GetYsClient().Job.AdminJobList.IsSystemFailed(o.IsSystemFailed),
		)
		PrintResp(res, err, "List Job")
		return nil
	}
	return cmd
}

func newJobDeleteCmd(o JobOptions) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "delete",
		Short: "删除作业",
		Long:  "删除作业, 只可删除状态在[Terminated, Completed, Failed]且回传状态在[Paused, Completed, Failed]的作业",
		Example: `- 删除作业ID为 52hgnbSTGWd 的作业
  - ysadmin job delete -I 52hgnbSTGWd`,
	}
	cmd.Flags().StringVarP(&o.Id, "id", "I", "", "作业ID")
	cmd.MarkFlagRequired("id")
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		res, err := GetYsClient().Job.AdminJobDelete(
			GetYsClient().Job.AdminJobDelete.JobId(o.Id),
		)
		PrintResp(res, err, "Delete Job "+o.Id)
		return nil
	}
	return cmd
}

func newJobTerminateCmd(o JobOptions) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "terminate",
		Short: "终止作业",
		Long:  "终止作业",
		Example: `- 终止作业ID为 52hgnbSTGWd 的作业
  - ysadmin job terminate -I 52hgnbSTGWd`,
	}
	cmd.Flags().StringVarP(&o.Id, "id", "I", "", "作业ID")
	cmd.MarkFlagRequired("id")
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		res, err := GetYsClient().Job.AdminJobTerminate(
			GetYsClient().Job.AdminJobTerminate.JobId(o.Id),
		)
		PrintResp(res, err, "Terminate Job "+o.Id)
		return nil
	}
	return cmd
}

func newJobRetransmitCmd(o JobOptions) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "retransmit",
		Short: "作业结果文件重传",
		Long:  "作业结果文件重传",
		Example: `- 重新传输ID为 52hgnbSTGWd 的作业结果文件(只有传输结束或最后一轮传输时间超过30min的Job才能重传)
  - ysadmin job retransmit -I 52hgnbSTGWd`,
	}
	cmd.Flags().StringVarP(&o.Id, "id", "I", "", "作业ID")
	cmd.MarkFlagRequired("id")
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		res, err := GetYsClient().Job.AdminJobRetransmit(
			GetYsClient().Job.AdminJobRetransmit.JobId(o.Id),
		)
		PrintResp(res, err, "Retransmit Job "+o.Id)
		return nil
	}
	return cmd
}

func newJobUpdateCmd(o JobOptions) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "update",
		Short: "更新作业",
		Long:  "更新作业",
		Example: `- 更新作业ID为 52hgnbSTGWd 的作业
          - ysadmin update --file-sync-state=Failed -I 52hgnbSTGWd`,
	}
	cmd.Flags().StringVarP(&o.Id, "id", "I", "", "作业ID")
	cmd.MarkFlagRequired("id")
	cmd.Flags().StringVar(&o.FileSyncState, "file-sync-state", "", "文件同步状态, 可选值: Failed")
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if o.FileSyncState == "" {
			return fmt.Errorf("file-sync-state is required")
		}
		if strings.ToLower(o.FileSyncState) != "failed" {
			return fmt.Errorf("file-sync-state only supports Failed")
		}
		res, err := GetYsClient().Job.AdminJobUpdate(
			GetYsClient().Job.AdminJobUpdate.JobId(o.Id),
			GetYsClient().Job.AdminJobUpdate.FileSyncState(o.FileSyncState),
		)
		PrintResp(res, err, "Update Job File Sync State "+o.Id)
		return nil
	}
	return cmd
}

func newJobExampleCmd(o JobOptions) *cobra.Command {
	var cmd = &cobra.Command{
		Use:     "submit_example",
		Short:   "创建作业示例文件",
		Long:    "创建作业示例文件",
		Example: `ysadmin job example`,
	}
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if checkFileExist(submitJobExampleFile) {
			fmt.Println("示例文件:", submitJobExampleFile)
			return nil
		}
		f, err := os.Create(submitJobExampleFile)
		if err != nil {
			return err
		}
		defer f.Close()
		if _, err = f.WriteString(submitJobExampleFileContent); err != nil {
			return err
		}
		fmt.Println("示例文件:", submitJobExampleFile)
		return nil
	}
	return cmd
}

func newJobPreSchedule(o JobOptions) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "preschedule",
		Short: "作业预调度",
		Long:  "作业预调度，提交作业预调度信息",
		Example: `- 开启作业预调度, 参数文件为Params.json
- ysadmin job preschedule -F Params.json -Z az-jinan -S 0 -X 0`,
	}
	cmd.Flags().StringVarP(&o.JsonFile, "file", "F", "", "JSON 文件路径")
	cmd.MarkFlagRequired("file")
	cmd.Flags().StringSliceVarP(&o.Zones, "zones", "Z", nil, "指定区域列表，多个区域用逗号分隔，例如 az-jinan,az-zhigu")
	cmd.Flags().BoolVarP(&o.Shared, "shared", "S", false, "是否共享节点")
	cmd.Flags().BoolVarP(&o.Fixed, "fixed", "X", false, "是否固定分区")
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		data, err := checkReqJsonFileAndRead(o.JsonFile)
		if err != nil {
			return err
		}
		req := new(jobpreschedule.Request)
		err = jsoniter.Unmarshal(data, req)
		if err != nil {
			return err
		}
		res, err := GetYsClient().Job.JobPreSchedule(
			GetYsClient().Job.JobPreSchedule.Params(req.Params),
			GetYsClient().Job.JobPreSchedule.Zones(req.Zones),
			GetYsClient().Job.JobPreSchedule.Fixed(req.Fixed),
			GetYsClient().Job.JobPreSchedule.Shared(req.Shared),
		)
		PrintResp(res, err, "JobPreSchedule success!")
		return nil
	}
	return cmd
}

func newParamsExampleCmd(o JobOptions) *cobra.Command {
	var cmd = &cobra.Command{
		Use:     "preschedule_example",
		Short:   "预调度参数文件示例 ",
		Long:    "预调度参数文件示例 ",
		Example: `ysadmin job preschedule example`,
	}
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if checkFileExist(scheduleJobExampleFile) {
			fmt.Println("示例文件:", scheduleJobExampleFile)
			return nil
		}
		f, err := os.Create(scheduleJobExampleFile)
		if err != nil {
			return err
		}
		defer f.Close()
		if _, err = f.WriteString(scheduleJobExampleFileContent); err != nil {
			return err
		}
		fmt.Println("示例文件:", scheduleJobExampleFile)
		return nil
	}
	return cmd
}
