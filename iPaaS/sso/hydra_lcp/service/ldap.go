package service

import (
	"context"
	"fmt"
	"github.com/go-ldap/ldap/v3"
	"os"
	"strconv"

	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/consts"
	"google.golang.org/grpc/status"

	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/common"
)

// LdapService LdapService
type LdapService struct {
	dsn string
}

// NewLdapService NewLdapService
func NewLdapService(dsn string) *LdapService {
	return &LdapService{
		dsn: dsn,
	}
}

// VerifyPassword VerifyPassword
func (s *LdapService) VerifyPassword(ctx context.Context, cn string, pwd string) (int64, error) {
	dn := fmt.Sprintf("uid=%v,"+common.LdapBaseDN, cn)
	filter := fmt.Sprintf("(&(uid=%v))", cn)

	l, err := ldap.Dial("tcp", s.dsn)
	if err != nil {
		return 0, status.Error(consts.ErrHydraLcpLdapError, err.Error())
	}
	defer l.Close()

	err = l.Bind(dn, pwd)
	if err != nil {
		return 0, status.Error(consts.ErrHydraLcpLoginFailed, err.Error())
	}

	searchRequest := ldap.NewSearchRequest(
		common.LdapBaseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0, 0, false,
		filter,
		[]string{common.LdapID},
		nil)
	sr, err := l.Search(searchRequest)
	if err != nil {
		return 0, status.Error(consts.ErrHydraLcpLdapError, err.Error())
	}

	if len(sr.Entries) != 1 {
		return 0, status.Error(consts.ErrHydraLcpLdapError, "User does not exist or too many entries returned")
	}

	var sid string

	for _, attr := range sr.Entries[0].Attributes {
		if attr.Name == common.LdapID {
			sid = attr.Values[0]
		}
	}

	uid, err := strconv.ParseInt(sid, 10, 64)
	if err != nil {
		return 0, status.Errorf(consts.ErrHydraLcpLdapError, "uid should be int64, %v", err)
	}
	return uid, nil
}

// ResetPassword 重置密码
func (s *LdapService) ResetPassword(ctx context.Context, cn string, oldPwd string, newPwd string) error {
	l, err := ldap.Dial("tcp", s.dsn)
	if err != nil {
		return status.Error(consts.ErrHydraLcpLdapError, err.Error())
	}
	defer l.Close()

	dn := fmt.Sprintf("uid=%v,"+common.LdapBaseDN, cn)
	err = l.Bind(dn, oldPwd)
	if err != nil {
		return status.Error(consts.ErrHydraLcpLdapError, err.Error())
	}

	passwordModifyRequest := ldap.NewPasswordModifyRequest("", oldPwd, newPwd)
	_, err = l.PasswordModify(passwordModifyRequest)

	if err != nil {
		return status.Error(consts.ErrHydraLcpLdapError, err.Error())
	}
	return nil
}

// ModifyPwdByAdmin ...
func (s *LdapService) ModifyPwdByAdmin(cn string, pwd string) error {

	adminCn := os.Getenv("LDAP_ADMIN_CN")
	adminPwd := os.Getenv("LDAP_ADMIN_PASSWORD")

	if adminCn == "" || adminPwd == "" {
		return status.Error(consts.ErrHydraLcpLdapError, "ldap administrator cn or password is not be seted")
	}

	l, err := ldap.Dial("tcp", s.dsn)
	if err != nil {
		return status.Error(consts.ErrHydraLcpLdapError, err.Error())
	}
	defer l.Close()

	adminDn := fmt.Sprintf("cn=%v,"+common.LdapBaseDN, adminCn)
	err = l.Bind(adminDn, adminPwd)
	if err != nil {
		return status.Error(consts.ErrHydraLcpLdapError, err.Error())
	}

	userDn := fmt.Sprintf("uid=%v,"+common.LdapBaseDN, cn)
	passwordModifyRequest := ldap.NewPasswordModifyRequest(userDn, "", pwd)
	_, err = l.PasswordModify(passwordModifyRequest)

	if err != nil {
		status.Error(consts.ErrHydraLcpLdapError, err.Error())
	}
	return nil
}

// SearchUser ...
func (s *LdapService) SearchUser(cn string) (bool, error) {

	l, err := ldap.Dial("tcp", s.dsn)
	if err != nil {
		return false, status.Error(consts.ErrHydraLcpLdapError, err.Error())
	}
	defer l.Close()

	searchRequest := ldap.NewSearchRequest(
		common.LdapBaseDN, // The base dn to search
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(uid=%v))", cn),
		[]string{"dn", "uid"}, // A list attributes to retrieve
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		return false, status.Error(consts.ErrHydraLcpLdapError, err.Error())
	}

	exist := true
	if len(sr.Entries) == 0 {
		exist = false
	}

	return exist, nil
}
