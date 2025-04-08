package util

import (
	"fmt"
	"testing"
	"time"

	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/common"
)

func TestJWTGenerate(t *testing.T) {
	token, err := JWTGenerate(common.JWTSignup, "124", time.Now().Add(time.Hour*48))
	fmt.Println(token)
	fmt.Println(err)
}

var tokenString = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1Njk5MDk3NjMsImlhdCI6MTU2OTczNjk2MywibmJmIjoxNTY5NzE1MjAwLCJzdWIiOiJzaWdudXAtMTIzNCJ9.V61rbiGWyx7il8JwLqFTyZgqBcKNfMi6aZWXSsi-mLE"

func TestJWTParse(t *testing.T) {
	//tokenString, err := JWTGenerate(common.JWTSignup, "124", time.Now().Add(time.Hour*48))
	claims, err := JWTParse(tokenString)
	fmt.Println(claims)
	fmt.Println(err)
}

func TestJWTGetSubject(t *testing.T) {
	//tokenString, err := JWTGenerate(common.JWTSignup, "124", time.Now().Add(time.Minute))
	sub, err := JWTGetSubject(common.JWTSignup, tokenString)
	fmt.Println(sub)
	fmt.Println(err)
}
