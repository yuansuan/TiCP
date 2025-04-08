package impl

import (
	"github.com/spf13/cobra"
)

// ZoneOptions 区域命令的参数
type ZoneOptions struct{}

func init() {
	RegisterCmd(NewZoneCommand())
}

// NewZoneCommand 创建区域命令
func NewZoneCommand() *cobra.Command {
	o := ZoneOptions{}
	cmd := &cobra.Command{
		Use:   "zone",
		Short: "区域管理",
		Long:  "区域管理, 用于获取区域列表等",
		RunE:  helpRun,
	}

	cmd.AddCommand(
		newZoneListCmd(o),
	)
	return cmd
}

func newZoneListCmd(o ZoneOptions) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "list",
		Short: "获取所有区域的域名信息",
		Long:  "获取所有区域的域名信息, 包含HPCEndpoint、StorageEndpoint、CloudAppEnable信息",
		Run: func(command *cobra.Command, args []string) {
			o.listZones()
		},
	}
	return cmd
}

func (o *ZoneOptions) listZones() {
	res, err := GetYsClient().Job.ZoneList()
	PrintResp(res, err, "List Zone")
}
