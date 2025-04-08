package impl

import (
	"bytes"
	"fmt"
	"html"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// DocOptions doc命令的参数
type DocOptions struct{}

func init() {
	RegisterCmd(NewDocCommand())
}

// NewDocCommand 创建doc命令
func NewDocCommand() *cobra.Command {
	o := DocOptions{}
	cmd := &cobra.Command{
		Use:   "doc",
		Short: "📖 生成某个命令的文档",
		Long:  "📖 生成某个命令的文档, 通过解析cobra的结构体, 生成markdown格式的文档, 包含描述、用法、示例、子命令、flag等信息, 支持递归生成子命令的文档, 支持输出到文件",
		Args:  cobra.MinimumNArgs(1),
		Example: `ysadmin doc job
ysadmin doc job submit`,
	}

	cmd.Flags().BoolP("sub", "S", false, "是否递归生成子命令的文档")
	cmd.Flags().StringP("file", "F", "", "指定输出文件，默认输出到标准输出")

	cmd.Run = func(cmd *cobra.Command, args []string) {
		o.Run(cmd, args)
	}

	cmd.AddCommand(
		newDocSubCmd(),
	)
	return cmd
}

func newDocSubCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "sub",
		Short:   "用于生成子命令文档测试的命令",
		Long:    "用于生成子命令文档测试的命令",
		Example: `ysadmin doc sub`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("sub command")
		},
	}

	cmd.AddCommand(
		newDocSubSubCmd(),
	)

	// flag
	cmd.Flags().StringP("aaa", "a", "", "用于测试的flagA")
	cmd.Flags().IntP("bbb", "", 0, "用于测试的flagB")
	cmd.Flags().BoolP("ccc", "", false, "用于测试的flagC")

	return cmd
}

func newDocSubSubCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "sub",
		Short:   "用于生成子命令的子命令文档测试的命令",
		Long:    "用于生成子命令的子命令文档测试的命令",
		Example: `ysadmin doc sub`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("sub sub command")
		},
	}

	cmd.Flags().StringP("aaa", "a", "", "用于测试的flagA")
	cmd.Flags().IntP("bbb", "", 0, "用于测试的flagB")
	cmd.Flags().BoolP("ccc", "", false, "用于测试的flagC")

	return cmd
}

// Run 生成某个命令的文档
func (o *DocOptions) Run(command *cobra.Command, args []string) {
	// 获取rootCmd
	rootCmd := command.Root()

	// 获取命令
	cmd, _, err := rootCmd.Find(args)
	if err != nil {
		fmt.Println("Find command error: ", err)
		return
	}

	// 输出源
	out := rootCmd.OutOrStdout()
	if file, _ := command.Flags().GetString("file"); file != "" {
		fmt.Println("output to file: ", file)
		// 打开
		f, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println("Open file error: ", err)
			return
		}
		defer f.Close()
		out = f
	}

	topLevel := 1
	printCmd(cmd, out, topLevel)

	// 是否生成子命令的文档
	if command.Flag("sub").Changed {
		// 递归打印子命令
		printSubCmd(cmd, out, topLevel+1)
		return
	}
}

// 递归打印子命令
func printSubCmd(cmd *cobra.Command, out io.Writer, level int) {
	for _, c := range cmd.Commands() {
		printCmd(c, out, level)
		printSubCmd(c, out, level+1)
	}
}

func printCmd(cmd *cobra.Command, out io.Writer, level int) {
	// 打印cmd的描述
	printDesc(cmd, out, level)

	// 打印cmd的用法
	printUsage(cmd, out)

	// 打印cmd的示例（如果有）
	printExample(cmd, out)

	// 打印cmd的子命令（如果有）
	printSubCmds(cmd, out)

	// 打印cmd的flag（如果有）
	printFlags(cmd, out)
}

func printDesc(cmd *cobra.Command, out io.Writer, level int) {
	subtitle(out, level, cmd.Short)
	fmt.Fprintln(out, "")
	fmt.Fprintln(out, cmd.Long)
	fmt.Fprintln(out, "")
}

func printUsage(cmd *cobra.Command, out io.Writer) {
	fmt.Fprintln(out, "使用方法:")
	fmt.Fprintln(out, "")
	fmt.Fprintln(out, cmd.UseLine())
	fmt.Fprintln(out, "")
}

func printExample(cmd *cobra.Command, out io.Writer) {
	if len(cmd.Example) > 0 {
		fmt.Fprintln(out, "示例:")
		fmt.Fprintln(out, "")
		fmt.Fprintln(out, escapeMarkdown(cmd.Example))
		fmt.Fprintln(out, "")
	}
}

func printSubCmds(cmd *cobra.Command, out io.Writer) {
	if len(cmd.Commands()) > 0 {
		fmt.Fprintln(out, "可用子命令:")
		fmt.Fprintln(out, "")
		for _, c := range cmd.Commands() {
			fmt.Fprintln(out, "-", c.Name(), " ", c.Short)
		}
		fmt.Fprintln(out, "")
	}
}

func printFlags(cmd *cobra.Command, out io.Writer) {
	// 打印cmd的flag（如果有）
	if len(cmd.Flags().FlagUsages()) > 0 {
		printFlagCommon(out)
		cmd.Flags().VisitAll(func(flag *pflag.Flag) {
			printFlag(out, flag)
		})
		fmt.Fprintln(out, "")
	}
}

func printFlagCommon(out io.Writer) {
	fmt.Fprintln(out, "可用Flags:")
	fmt.Fprintln(out, "")
	fmt.Fprintln(out, "| 命令参数 | 类型 | 说明 |")
	fmt.Fprintln(out, "| ---: | :---: | :--- |")
}

func printFlag(out io.Writer, flag *pflag.Flag) {
	if flag.Hidden {
		return
	}

	fmt.Fprintf(out, "| %s | %s | %s |\n", handleFlagName(flag), handleFlagType(flag), escapeMarkdown(handleFlagDescription(flag)))
}

func handleFlagName(flag *pflag.Flag) string {
	line := ""
	if flag.Shorthand != "" && flag.ShorthandDeprecated == "" {
		line = fmt.Sprintf("-%s, --%s", flag.Shorthand, flag.Name)
	} else {
		line = fmt.Sprintf("--%s", flag.Name)
	}
	return line
}

func handleFlagType(flag *pflag.Flag) string {
	line := ""
	varname, _ := pflag.UnquoteUsage(flag)
	if varname != "" {
		line += " " + varname
	}
	if flag.NoOptDefVal != "" {
		switch flag.Value.Type() {
		case "string":
			line += fmt.Sprintf("[=\"%s\"]", flag.NoOptDefVal)
		case "bool":
			if flag.NoOptDefVal != "true" {
				line += fmt.Sprintf("[=%s]", flag.NoOptDefVal)
			}
		case "count":
			if flag.NoOptDefVal != "+1" {
				line += fmt.Sprintf("[=%s]", flag.NoOptDefVal)
			}
		default:
			line += fmt.Sprintf("[=%s]", flag.NoOptDefVal)
		}
	}
	return line
}

func handleFlagDescription(flag *pflag.Flag) string {
	line := ""
	maxlen := 0
	_, usage := pflag.UnquoteUsage(flag)

	line += "  "
	if len(line) > maxlen {
		maxlen = len(line)
	}

	line += usage
	anno := flag.Annotations
	if required := anno[cobra.BashCompOneRequiredFlag]; required != nil && len(required) > 0 {
		line += fmt.Sprintf(" (必填)")
	}

	if !defaultIsZeroValue(flag) {
		if flag.Value.Type() == "string" {
			line += fmt.Sprintf(" (default %q)", flag.DefValue)
		} else {
			line += fmt.Sprintf(" (default %s)", flag.DefValue)
		}
	}
	if len(flag.Deprecated) != 0 {
		line += fmt.Sprintf(" (DEPRECATED: %s)", flag.Deprecated)
	}
	return line
}

func defaultIsZeroValue(f *pflag.Flag) bool {
	switch f.Value.Type() {
	case "bool":
		return f.DefValue == "false"
	case "duration":
		// Beginning in Go 1.7, duration zero values are "0s"
		return f.DefValue == "0" || f.DefValue == "0s"
	case "int", "int8", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64", "count", "float32", "float64":
		return f.DefValue == "0"
	case "string":
		return f.DefValue == ""
	case "ip", "ipMask", "ipNet":
		return f.DefValue == "<nil>"
	case "intSlice", "stringSlice", "stringArray":
		return f.DefValue == "[]"
	default:
		switch f.Value.String() {
		case "false":
			return true
		case "<nil>":
			return true
		case "":
			return true
		case "0":
			return true
		}
		return false
	}
}

func subtitle(out io.Writer, level int, text string) {
	fmt.Fprintf(out, "%s %s\n", strings.Repeat("#", level), text)
}

func escapeMarkdown(input string) string {
	escaped := input

	// 处理Markdown有意义的字符
	escaped = strings.ReplaceAll(escaped, "*", "\\*")
	escaped = strings.ReplaceAll(escaped, "_", "\\_")
	escaped = strings.ReplaceAll(escaped, "~", "\\~")
	escaped = strings.ReplaceAll(escaped, "`", "\\`")
	escaped = strings.ReplaceAll(escaped, "|", "\\|")
	escaped = strings.ReplaceAll(escaped, "[", "\\[")
	escaped = strings.ReplaceAll(escaped, "]", "\\]")
	escaped = strings.ReplaceAll(escaped, "(", "\\(")
	escaped = strings.ReplaceAll(escaped, ")", "\\)")
	escaped = strings.ReplaceAll(escaped, "{", "\\{")
	escaped = strings.ReplaceAll(escaped, "}", "\\}")

	return escaped
}

func escapeCodeBlocks(input string) string {
	var buffer bytes.Buffer
	lines := strings.Split(input, "\n")
	inCodeBlock := false

	for _, line := range lines {
		if strings.HasPrefix(line, "```") {
			inCodeBlock = !inCodeBlock
		}

		if inCodeBlock {
			buffer.WriteString(line)
		} else {
			buffer.WriteString(html.EscapeString(line))
		}

		buffer.WriteString("\n")
	}

	return buffer.String()
}
