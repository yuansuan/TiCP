package text

import (
	"encoding/json"

	"github.com/go-resty/resty/v2"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/module/wx"
)

type WXTextMessage struct {
	MsgType string `json:"msgtype"`
	Text    Text   `json:"text"`
}

// Text text message
type Text struct {
	Content             string   `json:"content"` // 不超过2048个字节，必须是utf8编码
	MentionedList       []string `json:"mentioned_list,omitempty"`
	MentionedMobileList []string `json:"mentioned_mobile_list,omitempty"`
}

func NewWXTextMessage(content string) *WXTextMessage {
	msg := &WXTextMessage{
		MsgType: wx.TextType,
	}
	msg.Content(content)
	return msg
}

func (m *WXTextMessage) Send(client *resty.Client, url string) (*resty.Response, error) {
	return wx.Send(client, url, m)
}

func (m *WXTextMessage) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}

func (wx *WXTextMessage) Content(content string) {
	wx.Text.Content = content
}

func (wx *WXTextMessage) MentionedList(mentionedList []string) {
	wx.Text.MentionedList = mentionedList
}

func (wx *WXTextMessage) MentionedMobileList(mentionedMobileList []string) {
	wx.Text.MentionedMobileList = mentionedMobileList
}
