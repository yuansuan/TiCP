package util

import (
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/status"
	"regexp"

	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/consts"
)

// PwdGenerate PwdGenerate
func PwdGenerate(pwd string) (pwdHash string, err error) {
	h, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return "", status.Error(consts.ErrHydraLcpPwdFailToGenerate, err.Error())
	}
	return string(h), nil
}

// PwdVerify PwdVerify
func PwdVerify(pwd string, pwdHash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(pwdHash), []byte(pwd))
	if err != nil {
		return status.Errorf(consts.ErrHydraLcpPwdNotMatch, "password not match, err: %v", err.Error())
	}
	return nil
}

// IsPasswordValid tells whether the password is valid:
// 8-16 characters, only accept lower case letters, upper case letters, numbers, and serverial special characters(_#?!@$%^&*-)
func IsPasswordValid(passwd string) bool {
	ok, _ := regexp.MatchString("[a-zA-Z0-9_#?!@$%^&*-]{8,16}$", passwd)
	return ok
}

var (
	// _PhoneRegexp 手机号简易正则匹配
	_PhoneRegexp = regexp.MustCompile("^1\\d{10}$")
)

// IsPhoneValid returns whether the phone is valid
func IsPhoneValid(phone string) bool {
	return _PhoneRegexp.MatchString(phone)
}
