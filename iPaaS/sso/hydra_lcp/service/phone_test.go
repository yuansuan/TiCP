//go:build darwin
// +build darwin

package service

import (
	"context"
	"fmt"
	"testing"
)

var phoneSrv = NewPhoneSrv()

func TestPhoneService_Send(t *testing.T) {
	ctx := context.TODO()
	phone := "15910957558"
	content := "恭喜你，支付成功！"
	sign := "远算云"

	err := phoneSrv.sendSmsByMonyun(ctx, phone, content, sign)

	fmt.Println(err)
}

func TestPhoneService_SendVerificationCode(t *testing.T) {
	ctx := context.TODO()
	phone := "15910957558"
	sign := "智算云"
	err := phoneSrv.SendVerificationCode(ctx, phone, sign)

	fmt.Println(err)
}

func TestPhoneService_VerifyCode(t *testing.T) {
	ctx := context.TODO()

	code := "826590"
	phone := "15910957558"

	err := phoneSrv.VerifyCode(ctx, phone, code)
	fmt.Println(err)
}
