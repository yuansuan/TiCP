package impl

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
)

// YsIDOptions ID命令的参数
type YsIDOptions struct{}

func init() {
	RegisterCmd(NewYsIDCommand())
}

// NewYsIDCommand 创建ID命令
func NewYsIDCommand() *cobra.Command {
	o := YsIDOptions{}
	cmd := &cobra.Command{
		Use:   "ysid",
		Short: "❄️  ID 生成 | 解码 | 编码",
		Long:  "❄️  ID 生成 | 解码 | 编码",
		RunE:  helpRun,
	}

	cmd.AddCommand(
		newYsIDCreateCmd(o),
		newYsIDDecodeCmd(o),
		newYsIDEncodeCmd(o),
	)

	return cmd
}

func newYsIDCreateCmd(o YsIDOptions) *cobra.Command {
	var cmd = &cobra.Command{
		Use:     "create",
		Short:   "生成ID",
		Long:    "生成ID",
		Args:    cobra.NoArgs,
		Example: `ysadmin ysid create`,
		Run: func(command *cobra.Command, args []string) {
			o.create()
		},
	}
	return cmd
}

func newYsIDDecodeCmd(o YsIDOptions) *cobra.Command {
	var cmd = &cobra.Command{
		Use:     "decode [id]",
		Short:   "解码ID",
		Long:    "解码ID, 将YSID解码成数字",
		Args:    cobra.ExactArgs(1),
		Example: `ysadmin ysid decode 52jm6EJC335`,
		Run: func(command *cobra.Command, args []string) {
			o.decode(args[0])
		},
	}
	return cmd
}

func newYsIDEncodeCmd(o YsIDOptions) *cobra.Command {
	var cmd = &cobra.Command{
		Use:     "encode [num]",
		Short:   "编码ID",
		Long:    "编码ID, 将数字编码成YSID",
		Args:    cobra.ExactArgs(1),
		Example: `ysadmin ysid encode 1732993997258887168`,
		Run: func(command *cobra.Command, args []string) {
			o.encode(args[0])
		},
	}
	return cmd
}

func (o *YsIDOptions) create() {
	IDGen, err := snowflake.NewNode(1)
	if err != nil {
		fmt.Println("GenSnowflakeFail: ", err.Error())
		return
	}
	id := IDGen.Generate().String()
	fmt.Println("SnowflakeId: ", id)
}

func (o *YsIDOptions) decode(id string) {
	num, err := snowflake.ParseString(id)
	if err != nil {
		fmt.Println("DecodeError: ", err.Error())
		return
	}
	fmt.Println(num.Int64())
}

func (o *YsIDOptions) encode(numStr string) {
	num, err := strconv.Atoi(numStr)
	if err != nil {
		fmt.Println("InputNotNum: ", err.Error())
		return
	}
	id := snowflake.ID(num).String()
	fmt.Println(id)
}
