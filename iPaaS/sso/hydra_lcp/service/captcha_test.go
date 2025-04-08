//go:build darwin
// +build darwin

package service

import (
	"fmt"
	"testing"
)

var captchaSrv = NewCaptcha()
var c = customStore{}

func TestCustomStore_Set(t *testing.T) {
	c.Set("1", []byte{1, 2, 3})
}

func TestCustomStore_Get(t *testing.T) {
	digits := c.Get("zj8kDKkFxLceCtBAY1yc", false)
	fmt.Println(digits)
}

func TestCaptchaService_CreateDigitCaptcha(t *testing.T) {
	i, d := captchaSrv.CreateDigitCaptcha()

	fmt.Println(">>>>", i)
	fmt.Println(">>>>", d)

}

func TestCaptchaService_ValidateCaptcha(t *testing.T) {
	i := "XppoIb2YRs7RSXwrHstd"
	d := "916794"
	err := captchaSrv.ValidateCaptcha(i, d)
	fmt.Println(err)
}
