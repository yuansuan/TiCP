//go:build darwin
// +build darwin

package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/pkg/snowflake"
)

var ldapSrv = NewLdapService("10.0.1.155:389")

func TestLdapService_VerifyPassword(t *testing.T) {
	ctx := context.TODO()
	id, err := ldapSrv.VerifyPassword(ctx, "xzheng", "8c27b8daa86a")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("user id is %v\n", id)
		fmt.Printf("user base58encode id is %v\n", snowflake.ID(id).String())
	}
}
