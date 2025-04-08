package singularity

import (
	"fmt"
	"github.com/k0kubun/go-ansi"
	"github.com/olekukonko/tablewriter"
	"github.com/schollz/progressbar/v3"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/tools/singularity/registry"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/tools/singularity/registry/image"
	"os"
	"time"
)

// Clog 将日志打印到终端上
func Clog(head, msg string, _ map[string]interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, "%s: %s\n", head, msg)
}

// LogToConsole 将日志打印到终端上
func LogToConsole(_ registry.Kind, head, msg string, ctx map[string]interface{}) {
	Clog(head, msg, ctx)
}

// GetProgressBarOpts 返回进度条参数
func GetProgressBarOpts(desc string) []progressbar.Option {
	return []progressbar.Option{
		progressbar.OptionSetWriter(ansi.NewAnsiStderr()),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetDescription(desc),
		progressbar.OptionClearOnFinish(),
		progressbar.OptionUseANSICodes(true),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[dark_gray]=[reset]",
			SaucerHead:    "[dark_gray]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	}
}

// PrintDefaultedLocators 打印内容
func PrintDefaultedLocators(locators []*image.DefaultedLocator) {
	table := tablewriter.NewWriter(os.Stderr)
	table.SetHeader([]string{"DEFAULT", "NAME", "TAG", "HASH", "CREATED", "DANGLING"})
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("\t")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t\t")
	table.SetNoWhiteSpace(true)
	for _, locator := range locators {
		var createdTime string
		if locator.Created() != nil {
			createdTime = locator.Created().In(time.Local).Format("2006/01/02 15:04:05")
		}

		row := []string{
			DefaultString(locator.IsDefaultTag() && locator.IsDefaultedHash(), "*", ""),
			locator.Name(),
			DefaultString(locator.Tag() == "", "<none>", locator.Tag()),
			DefaultString(locator.Hash() == "", "<none>", locator.Hash()),
			DefaultString(createdTime == "", "<none>", createdTime),
			DefaultString(locator.Hash() == "" || !locator.IsDefaultedHash(), "*", ""),
		}

		if locator.Hash() == "" || !locator.IsDefaultedHash() {
			table.Rich(row, []tablewriter.Colors{
				{},
				{tablewriter.FgHiBlackColor},
				{tablewriter.FgHiBlackColor},
				{tablewriter.FgHiBlackColor},
				{tablewriter.FgHiBlackColor},
				{tablewriter.FgHiBlackColor},
				{tablewriter.FgHiBlackColor},
				{},
			})
		} else {
			table.Append(row)
		}
	}
	table.Render()
}

// DefaultString 字符串的三目运算符
func DefaultString(b bool, t string, f string) string {
	if b {
		return t
	}
	return f
}
