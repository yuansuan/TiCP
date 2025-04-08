package text

import (
	"testing"

	"github.com/go-resty/resty/v2"
)

func TestWXTextMessage_Send(t *testing.T) {
	type args struct {
		URL string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "normal",
			args: args{
				// URL: "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=f8848b25-709d-41af-bbfb-6c3d1651e533",
				URL: "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=xxxxxxxx",
			},
			wantErr: false,
		},
	}
	client := resty.New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wxmsg := NewWXTextMessage("just a test message, don't worry")
			if resp, err := wxmsg.Send(client, tt.args.URL); (err != nil) != tt.wantErr {
				t.Errorf("WXTextMessage.Send() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				t.Log(resp)
			}
		})
	}
}
