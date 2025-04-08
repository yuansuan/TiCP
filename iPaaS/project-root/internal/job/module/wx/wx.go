package wx

//go:generate mockgen -self_package=github.com/yuansuan/ticp/iPaaS/project-root/internal/job/module/wx -destination mock_sender.go -package wx github.com/yuansuan/ticp/iPaaS/project-root/internal/job/module/wx Sender
import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

const (
	TextType         = "text"
	MarkdownType     = "markdown"
	ImageType        = "image"
	NewsType         = "news"
	FileType         = "file"
	VoiceType        = "voice"
	TemplateCardType = "template_card"
)

const (
	TextNoticeCard = "text_notice"
	NewsNoticeCard = "news_notice"
)

const (
	WarningColor = "warning"
	CommentColor = "comment"
	InfoColor    = "info"
	RedColor     = "red"
)

// colorMsg warning, comment, info, red
func colorMsg(msg string, color string) string {
	switch color {
	case WarningColor:
		return fmt.Sprintf("<font color=\"warning\">%s</font>", msg)
	case CommentColor:
		return fmt.Sprintf("<font color=\"comment\">%s</font>", msg)
	case InfoColor:
		return fmt.Sprintf("<font color=\"info\">%s</font>", msg)
	case RedColor:
		return fmt.Sprintf("<font color=\"red\">%s</font>", msg)
	default:
		return msg
	}
}

type WXMessage interface {
	Send(client *resty.Client, webhookURL string) (*resty.Response, error)
	ToJSON() ([]byte, error)
}

func Send(client *resty.Client, url string, msg WXMessage) (*resty.Response, error) {
	if url == "" {
		return nil, fmt.Errorf("url is empty")
	}

	jsonData, err := msg.ToJSON()
	if err != nil {
		return nil, err
	}
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(jsonData).
		Post(url)
	return resp, err
}

type Sender interface {
	Send(msg WXMessage) (*resty.Response, error)
}

type HTTPSender struct {
	client *resty.Client
	url    string
}

func NewHTTPSender(url string) *HTTPSender {
	client := resty.New()
	return &HTTPSender{
		client: client,
		url:    url,
	}
}

// Send 发送消息
func (s *HTTPSender) Send(msg WXMessage) (*resty.Response, error) {
	return msg.Send(s.client, s.url)
}
