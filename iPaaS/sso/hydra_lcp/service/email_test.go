//go:build darwin
// +build darwin

package service

import (
	"context"
	"fmt"
	"testing"
)

var s = EmailService{
	host:     "smtp.partner.outlook.cn:587",
	user:     "internal@lambdacal.com",
	password: "1234password",
}

var email = "test1@yuansuan.cn"
var pwd = "helloworld"

func TestEmailService_Send(t *testing.T) {
	ctx := context.TODO()
	err := s.Send(ctx, "test", "test", "fwchen@yuansuan.cn")
	fmt.Println(err)
}

func TestEmailService_Signup(t *testing.T) {
	userID := int64(1234)
	email := "fwchen@yuansuan.cn"
	pwd := "2"

	ctx := context.TODO()
	err := s.Signup(ctx, userID, email, pwd, "localhost:8900")
	fmt.Println(err)
}

var tokenString = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1Njk5MDk5NzUsImlhdCI6MTU2OTczNzE3NSwibmJmIjoxNTY5NzE1MjAwLCJzdWIiOiJzaWdudXAtMTIzNCJ9.Fiote10yezlgKlaESepbQ1176Hte3xqHFiJ999LTQvI"

func TestEmailService_Activate(t *testing.T) {
	userID := int64(1234)

	ctx := context.TODO()
	err := s.Activate(ctx, tokenString, userID)
	fmt.Println(err)
}
