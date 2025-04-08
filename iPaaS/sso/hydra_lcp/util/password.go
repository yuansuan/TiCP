package util

import (
	"bytes"
	"crypto/rand"
	"math/big"
	"regexp"
	"time"

	mathRand "math/rand"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/status"

	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/consts"
)

const (
	constMobileRegexp = "^1[0-9]{10}$"
	constLowerRegexp  = "[a-z]"
	constUpperRegexp  = "[A-Z]"
	constNumberRegexp = "[0-9]"
	// supported special characters:
	// ` ~ ! @ # $ % ^ & * _ + - = ( ) [ ] { } < > \ | ; : ' " , . / ?
	constSpecialCharRegexp = "[`" + `~!@#$%^&*_+\-=()\[\]{}<>\\|;:'",./?]`

	constMinPasswdLength      = 8
	constMaxPasswdLength      = 16
	constMinPasswdCombination = 2
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

func GenerateValidPassword() string {
	src := mathRand.NewSource(time.Now().UnixNano())
	rand := mathRand.New(src)
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_#?!@$%^&*-"
	length := rand.Intn(9) + 8 // Random length between 8 and 16
	password := make([]byte, length)

	for i := 0; i < length; i++ {
		password[i] = chars[rand.Intn(len(chars))]
	}

	return string(password)
}

var (
	// _PhoneRegexp 手机号简易正则匹配
	_PhoneRegexp = regexp.MustCompile("^1\\d{10}$")
)

// IsPhoneValid returns whether the phone is valid
func IsPhoneValid(phone string) bool {
	return _PhoneRegexp.MatchString(phone)
}

// RandomString ...
func RandomString(len int) string {
	var randString string
	var str = "abcdefghijkmnpqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	b := bytes.NewBufferString(str)
	length := b.Len()
	bigInt := big.NewInt(int64(length))
	for i := 0; i < len; i++ {
		randomInt, _ := rand.Int(rand.Reader, bigInt)
		randString += string(str[randomInt.Int64()])
	}
	return randString
}
