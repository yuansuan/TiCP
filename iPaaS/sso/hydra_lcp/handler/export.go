package handler

import (
	"os"

	grpc_boot "github.com/yuansuan/ticp/common/go-kit/gin-boot/grpc-boot"

	"github.com/yuansuan/ticp/iPaaS/sso/protos/platform/idgen"

	"github.com/yuansuan/ticp/iPaaS/sso/protos/platform/company"

	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/config"
	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/service"
	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/util"
)

// Handler Handler
type Handler struct {
	conf *config.CustomT

	*util.HydraConfig

	userSrv            *service.UserService
	emailSrv           *service.EmailService
	phoneSrv           *service.PhoneService
	ldapSrv            *service.LdapService
	offiacctBindingSrv *service.OffiaccountBindingService
	captchaSrv         *service.CaptchaService

	Idgen              idgen.IdGenClient            `grpc_client_inject:"idgen"`
	CompanyService     company.CompanyServiceClient `grpc_client_inject:"company"`
	CompanyUserService company.UserServiceClient    `grpc_client_inject:"company"`
}

// CreateHandler CreateHandler
func CreateHandler() *Handler {
	h := &Handler{}

	h.conf = &config.Custom
	h.userSrv = service.NewUserSrv()
	h.emailSrv = service.NewEmailSrv(h.conf.SMTP.Host, h.conf.SMTP.UserName, h.conf.SMTP.Password)
	h.phoneSrv = service.NewPhoneSrv()
	h.offiacctBindingSrv = service.NewOffiaccountBindingSrv()
	h.captchaSrv = service.NewCaptcha()

	h.HydraConfig = util.GetHydraConfig()

	// config ldap from env
	// if env LDAP_STARTUP exists and its value is "yes", use ldap. if not, can't use ldap
	vv, startup := os.LookupEnv("LDAP_STARTUP")
	if startup && vv == "yes" {
		dsn := os.Getenv("LDAP_DSN")
		h.ldapSrv = service.NewLdapService(dsn)
	}

	grpc_boot.InjectAllClient(h)

	return h
}
