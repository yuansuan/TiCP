package markdown

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/module/wx"
)

type WXMarkdownMessage struct {
	MsgType  string   `json:"msgtype"`
	Markdown Markdown `json:"markdown"`
}

type Markdown struct {
	Content string `json:"content"`
}

func NewWXMarkdownMessage() *WXMarkdownMessage {
	return &WXMarkdownMessage{
		MsgType: wx.MarkdownType,
	}
}

func (m *WXMarkdownMessage) Send(client *resty.Client, url string) (*resty.Response, error) {
	return wx.Send(client, url, m)
}

func (m *WXMarkdownMessage) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}

func (wx *WXMarkdownMessage) Content(content SimpleMarkdownFormat) {
	wx.Markdown.Content = content.String()
}

// SimpleMarkdownFormat is a simple markdown format struct
/*
	# Title
	## Subtitle
	Content
	> Item1
	> Item2
	> Item3
*/
type SimpleMarkdownFormat struct {
	Title    string
	Subtitle string
	Content  string
	Items    []string
}

// String returns the markdown format string
func (sm SimpleMarkdownFormat) String() string {
	var sb strings.Builder

	if sm.Title != "" {
		sb.WriteString(fmt.Sprintf("# %s\n\n", sm.Title))
	}

	if sm.Subtitle != "" {
		sb.WriteString(fmt.Sprintf("## %s\n\n", sm.Subtitle))
	}

	if sm.Content != "" {
		sb.WriteString(fmt.Sprintf("%s\n\n", sm.Content))
	}

	if len(sm.Items) > 0 {
		for _, item := range sm.Items {
			sb.WriteString(fmt.Sprintf("> %s\n", item))
		}
	}

	return sb.String()
}
