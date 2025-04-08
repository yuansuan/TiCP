package impl

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	openapi "github.com/yuansuan/ticp/common/openapi-go"
)

// HpcOptions 直连超算命令的参数
type HpcOptions struct {
	Zone    string
	Timeout int
}

func init() {
	RegisterCmd(NewHpcCmdCommand())
}

// NewHpcCmdCommand 创建直连超算命令
func NewHpcCmdCommand() *cobra.Command {
	o := HpcOptions{}
	cmd := &cobra.Command{
		Use:   "hpc",
		Short: "超算中心管理",
		Long:  "超算中心管理, 直连超算, 可查看HPC剩余资源以及在HPC集群上执行shell命令",
		RunE:  helpRun,
	}

	cmd.AddCommand(
		newHpcCmdCmd(o),
		newFreeResourceCmd(o),
	)

	return cmd
}

func newFreeResourceCmd(o HpcOptions) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "freeresource",
		Short: "获取HPC剩余资源",
		Long:  "获取HPC剩余资源",
		Example: `- 获取HPC剩余资源, 区域为az-zhigu
  - ysadmin hpc freeresource -Z az-zhigu`,
	}

	cmd.Flags().StringVarP(&o.Zone, "zone", "Z", "", "区域")
	cmd.MarkFlagRequired("zone")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		o.freeResource()
		return nil
	}

	return cmd
}

func newHpcCmdCmd(o HpcOptions) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "cmd [cmd]",
		Short: "在HPC集群上执行shell命令",
		Long:  "在HPC集群上执行shell命令",
		Example: `- 在HPC集群上执行ls -l命令, 区域为az-zhigu
  - ysadmin hpc cmd "ls -l" -Z az-zhigu
- 在HPC集群上执行ls -l命令, 超时时间为20秒, 区域为az-zhigu
  - ysadmin hpc cmd "ls -l" -Z az-zhigu -T 20`,
		Args: cobra.ExactArgs(1),
	}

	cmd.Flags().StringVarP(&o.Zone, "zone", "Z", "", "区域")
	cmd.MarkFlagRequired("zone")
	cmd.Flags().IntVarP(&o.Timeout, "timeout", "T", 10, "执行命令的超时时间")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		o.execCmd(args[0])
		return nil
	}

	return cmd
}

func (o *HpcOptions) execCmd(cmd string) {
	c := NewHpcClient(o.Zone)
	res, err := c.HPC.Command.System.Execute(
		c.HPC.Command.System.Execute.Command(cmd),
		c.HPC.Command.System.Execute.Timeout(o.Timeout),
	)
	PrintResp(res, err, "Exec Hpc Cmd")
}

func (o *HpcOptions) freeResource() {
	c := NewHpcClient(o.Zone)
	res, err := c.HPC.Resource.System.Get()
	PrintResp(res, err, "Hpc FreeResource")
}

// NewHpcClient 创建直连超算的Client
func NewHpcClient(zone string) *openapi.Client {
	res, err := GetYsClient().Job.ZoneList()
	if err != nil {
		fmt.Printf("List Zones Fail, Error: %s\n", err.Error())
		os.Exit(1)
	}
	if _, ok := res.Data.Zones[zone]; !ok {
		fmt.Printf("Zone %s not Found\n", zone)
		os.Exit(1)
	}
	GetYsClient().SetBaseUrl(res.Data.Zones[zone].HPCEndpoint)
	return GetYsClient()
}
