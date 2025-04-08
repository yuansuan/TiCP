package iamserver

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common/consts"

	bootHttp "github.com/yuansuan/ticp/common/go-kit/gin-boot/http"
	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	gindump "github.com/tpkeeper/gin-dump"
	"github.com/yuansuan/ticp/common/openapi-go/credential"
	"github.com/yuansuan/ticp/common/openapi-go/utils/signer"

	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common/leader"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/config"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/pkg/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/pkg/common/snowflake"
	v1 "github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/service/v1"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/store"

	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/store/dao"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/store/mysql"

	"github.com/mattn/go-isatty"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/pkg/log"

	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/controller/v1/policy"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/controller/v1/role"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/controller/v1/secret"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/controller/v1/user"
)

// discard gin-boot middlewares
// make apiserver faster
func clearMiddlewares(g *bootHttp.Driver) {
	g.Handlers = g.Handlers[:0]
}

func DumpFormatter(dumpStr string) {
	logging.Default().Info(dumpStr)
}

func InitRouter(g *bootHttp.Driver) {
	storeIns, err := mysql.GetMySQLFactoryOr()
	if err != nil {
		panic(err)
	}
	cleanAuditLog(storeIns)
	clearMiddlewares(g)
	// 下面都是 OpenAPI，供普通SSO 用户使用
	g.Use(requestIdMiddleware)
	router := g.Group("/iam/v1")
	g.GET("/healthz", healthOk)
	svr := v1.NewService(storeIns)

	secertCtl := secret.NewSecretController(storeIns)
	roleCtl := role.NewRoleController(storeIns)
	policyCtl := policy.NewPolicyController(storeIns)
	{
		// use less middleware to make console clean
		router.Use(gindump.DumpWithOptions(true, true, true, false, false, DumpFormatter))
		router.Use(gin.Recovery())
		router.Use(userAuth(storeIns, false))

		policyv1 := router.Group("/policies")
		{
			policyv1.POST("", policyCtl.AddPolicy)
			policyv1.GET(":policyName", policyCtl.GetPolicy)
			policyv1.GET("", policyCtl.ListPolicy)
			policyv1.DELETE(":policyName", policyCtl.DeletePolicy)
			policyv1.DELETE("", policyCtl.InvalidPolicyName)
			policyv1.PUT(":policyName", policyCtl.UpdatePolicy)
			policyv1.PUT("", policyCtl.InvalidPolicyName)
		}

		secretv1 := router.Group("/secrets")
		{
			secretv1.POST("", secertCtl.CreateSecret)
			secretv1.GET(":accessKeyId", secertCtl.GetSecret)
			secretv1.GET("", secertCtl.ListByParentUserID)
			secretv1.DELETE(":accessKeyId", secertCtl.DeleteSecretByParentUser)
			secretv1.DELETE("", secertCtl.InvalidAppKey)
		}

		rolev1 := router.Group("/roles")
		{
			rolev1.GET(":roleName", roleCtl.GetRole)
			rolev1.GET("", roleCtl.ListRole)
			rolev1.POST("", roleCtl.AddRole)
			rolev1.DELETE(":roleName", roleCtl.DeleteRole)
			rolev1.DELETE("", roleCtl.InvalidRoleName)
			rolev1.PUT(":roleName", roleCtl.UpdateRole)
			rolev1.PUT("", roleCtl.InvalidRoleName)

			rolev1.PATCH(":roleName/policies/:policyName", roleCtl.PatchPolicy)
			rolev1.PATCH(":roleName/policies", roleCtl.InvalidPolicyName)
			rolev1.DELETE(":roleName/policies/:policyName", roleCtl.DetachPolicy)
			rolev1.DELETE(":roleName/policies", roleCtl.InvalidPolicyName)
		}

		accountv1 := router.Group("/api/account")
		{
			accountv1.POST("", secertCtl.AddAccount)
		}

		router.POST("/AssumeRole", svr.AssumeRole)
		router.POST("/IsAllow", svr.IsAllow)
		router.GET("/IsYSProductAccount", secertCtl.IsYSProductAccount)
	}
	{
		// 账号密码换AK接口
		g.Group("/iam/v1/api/account").POST("/exchange", secertCtl.ExchangeCredentials)
	}
	internalRouter := g.Group("/iam/internal")
	{
		internalRouter.Use(gin.Recovery())
		internalRouter.Use(userAuth(storeIns, false))
		internalRouter.GET("/secrets", secertCtl.ListSecrets)
		internalRouter.POST("/secrets", secertCtl.InternalCreateSecret)
		internalRouter.GET("/secrets/user/:userId", secertCtl.InternalListSecret)
	}
	tempRouter := g.Group("/iam/temp")
	{
		tempRouter.Use(gin.Recovery())

		tempRouter.GET("/users/:phone", secertCtl.GetAKByPhone)
	}
	// admin接口, 允许增删改查  任何用户的策略，任何用户的角色，任何用户的AK
	adminRouter := g.Group("/iam/admin")
	{
		adminRouter.Use(gindump.Dump())
		adminRouter.Use(gin.Recovery())
		adminRouter.Use(userAuth(storeIns, true))

		adminRouter.POST("/policies", policyCtl.AdminAddPolicy)
		adminRouter.GET("/policies/:userId/:policyName", policyCtl.AdminGetPolicy)
		adminRouter.GET("/policies/:userId", policyCtl.AdminListPolicy)
		adminRouter.DELETE("/policies/:userId/:policyName", policyCtl.AdminDeletePolicy)
		adminRouter.PUT("/policies/:userId/:policyName", policyCtl.AdminUpdatePolicy)

		adminRouter.POST("/secrets", secertCtl.AdminCreateSecret)
		adminRouter.GET("/secrets/:accessKeyId", secertCtl.AdminGetSecret)
		adminRouter.GET("/secrets/user/:userId", secertCtl.AdminListSecret)
		adminRouter.GET("/secrets", secertCtl.AdminListSecrets)
		adminRouter.DELETE("/secrets/:accessKeyId", secertCtl.AdminDeleteSecret)
		adminRouter.PUT("/secrets/:accessKeyId/tag", secertCtl.AdminUpdateTag)

		adminRouter.POST("/roles", roleCtl.AdminAddRole)
		adminRouter.GET("/roles/:userId/:roleName", roleCtl.AdminGetRole)
		adminRouter.GET("/roles/:userId", roleCtl.AdminListRole)
		adminRouter.DELETE("/roles/:userId/:roleName", roleCtl.AdminDeleteRole)
		adminRouter.PUT("/roles/:userId/:roleName", roleCtl.AdminUpdateRole)
		adminRouter.PATCH("/roles/:roleName", roleCtl.AdminPatchPolicy)
		adminRouter.POST("/roles/:roleName", roleCtl.AdminDetachPolicy)

		adminRouter.GET("/users/:userId", user.AdminGetUser)
		adminRouter.GET("/users", user.AdminListUser)
		adminRouter.POST("/users", user.AdminAddUser)
		adminRouter.PUT("/users/:userId", user.AdminUpdateUser)
	}

	// paas api server， 可以单独移到一个微服务 TODO
	apiServer := g.Group("/", Logger())
	{
		apiSvr := NewApiServer()
		g.Any("/api/zones", apiSvr.ToJob)
		apiServer.Use(userAuth(storeIns, false))
		apiServer.Any("/api/jobs", apiSvr.ToJob)
		apiServer.Any("/api/jobs/*path", apiSvr.ToJob)
		apiServer.Any("/admin/jobs", apiSvr.ToJob)
		apiServer.Any("/admin/jobs/*path", apiSvr.ToJob)
		apiServer.Any("/api/apps", apiSvr.ToJob)
		apiServer.Any("/api/apps/*path", apiSvr.ToJob)
		apiServer.Any("/admin/apps", apiSvr.ToJob)
		apiServer.Any("/admin/apps/*path", apiSvr.ToJob)
		apiServer.Any("/system/jobs", apiSvr.ToJob)
		apiServer.Any("/system/jobs/*path", apiSvr.ToJob)

		apiServer.Any("/api/sessions", apiSvr.ToCloudApp)
		apiServer.Any("/api/sessions/*path", apiSvr.ToCloudApp)
		apiServer.Any("/api/hardwares", apiSvr.ToCloudApp)
		apiServer.Any("/api/hardwares/*path", apiSvr.ToCloudApp)
		apiServer.Any("/api/softwares", apiSvr.ToCloudApp)
		apiServer.Any("/api/softwares/*path", apiSvr.ToCloudApp)
		apiServer.Any("/admin/remoteapps", apiSvr.ToCloudApp)
		apiServer.Any("/admin/remoteapps/*path", apiSvr.ToCloudApp)
		apiServer.Any("/admin/sessions", apiSvr.ToCloudApp)
		apiServer.Any("/admin/sessions/*path", apiSvr.ToCloudApp)
		apiServer.Any("/admin/hardwares", apiSvr.ToCloudApp)
		apiServer.Any("/admin/hardwares/*path", apiSvr.ToCloudApp)
		apiServer.Any("/admin/softwares", apiSvr.ToCloudApp)
		apiServer.Any("/admin/softwares/*path", apiSvr.ToCloudApp)

		apiServer.Any("/admin/licenseManagers", apiSvr.ToLicManage)
		apiServer.Any("/admin/licenseManagers/*path", apiSvr.ToLicManage)
		apiServer.Any("/admin/licenses", apiSvr.ToLicManage)
		apiServer.Any("/admin/licenses/*path", apiSvr.ToLicManage)
		apiServer.Any("/admin/moduleConfigs", apiSvr.ToLicManage)
		apiServer.Any("/admin/moduleConfigs/*path", apiSvr.ToLicManage)

		apiServer.Any("/api/accounts", apiSvr.ToAccBill)
		apiServer.Any("/api/accounts/*path", apiSvr.ToAccBill)
		apiServer.Any("/internal/accounts", apiSvr.ToAccBill)
		apiServer.Any("/internal/accounts/*path", apiSvr.ToAccBill)
		apiServer.Any("/internal/users", apiSvr.ToAccBill)
		apiServer.Any("/internal/users/*path", apiSvr.ToAccBill)
		apiServer.Any("/internal/cashvouchers", apiSvr.ToAccBill)
		apiServer.Any("/internal/cashvouchers/*path", apiSvr.ToAccBill)
		apiServer.Any("/internal/accountcashvouchers", apiSvr.ToAccBill)
		apiServer.Any("/internal/accountcashvouchers/*path", apiSvr.ToAccBill)

		apiServer.Any("/internal/merchandises", apiSvr.ToMerchandise)
		apiServer.Any("/internal/merchandises/*path", apiSvr.ToMerchandise)
		apiServer.Any("/internal/specialprices", apiSvr.ToMerchandise)
		apiServer.Any("/internal/specialprices/*path", apiSvr.ToMerchandise)
		apiServer.Any("/internal/orders", apiSvr.ToMerchandise)
		apiServer.Any("/internal/orders/*path", apiSvr.ToMerchandise)
	}
}

const (
	expireTime = 15 * time.Minute
)

func requestIdMiddleware(c *gin.Context) {
	reqId := common.SetRequestID(c)
	reqIdKey := common.RequestIDKey
	c.Request.Header.Set(reqIdKey, reqId)
	c.Writer.Header().Set(reqIdKey, reqId)
	c.Next()
}

// isAdmin 表示iam 的admin
func userAuth(store store.Factory, isAdmin bool) func(c *gin.Context) {
	return func(c *gin.Context) {
		qs := c.Request.URL.Query()
		akId := qs.Get("AccessKeyId")
		if akId == "" {
			akId = qs.Get("AppKey")
		}
		if akId == "" {
			logging.Default().Infof("EmptyAccessKeyId, RequestId: %s", common.GetRequestID(c))
			common.ErrorRespWithAbort(c, http.StatusUnauthorized, "InvalidAppKey", "empty access key id")
			return
		}
		signature := qs.Get("Signature")
		timestamp, err := parseTimestamp(qs.Get("Timestamp"))
		if err != nil {
			logging.Default().Infof("InvalidParameters.Timestamp, RequestId: %s", common.GetRequestID(c))
			common.ErrorRespWithAbort(c, http.StatusUnauthorized, "InvalidTime", "time param error")
			return
		}
		now := time.Now()
		if now.Sub(timestamp) > expireTime {
			errMsg := fmt.Sprintf("time expired. server current time: %d, your timestamp: %d",
				now.Unix(), timestamp.Unix())
			logging.Default().Infof("ExpiredParameters.Timestamp, RequestId: %s, Error Info: %s",
				common.GetRequestID(c), errMsg)
			common.ErrorRespWithAbort(c, http.StatusUnauthorized, "InvalidTime", errMsg)
			return
		}
		secret := &dao.Secret{
			AccessKeyId:     config.GetConfig().AdminAKID,
			AccessKeySecret: config.GetConfig().AdminAKSECRET,
			ParentUser:      snowflake.ID(0).String(),
			Tag:             common.IamAdminTag,
		}
		if !isAdmin {
			secret, err = store.Secrets().Get(context.Background(), akId)
			if err != nil {
				if errors.Is(err, common.ErrRecordNotFound) {
					logging.Default().Infof("NotFoundAccessKeyId, RequestId: %s", common.GetRequestID(c))
					common.ErrorRespWithAbort(c, http.StatusUnauthorized, "InvalidAppKey",
						fmt.Sprintf("%s not found", akId))
					return
				}
				common.InternalServerError(c, "")
				return
			}
		}
		cred := credential.NewCredential(secret.AccessKeyId, secret.AccessKeySecret)
		sig, err := signer.NewSigner(cred)
		if err != nil {
			logging.Default().Errorf("NewSignerError, RequestId: %s, Error: %v", common.GetRequestID(c), err)
			common.InternalServerError(c, "")
			return
		}
		sigResult, err := sig.SignHttp(c.Request)
		if err != nil {
			logging.Default().Infof("SignHttpError, RequestId: %s, Error: %v", common.GetRequestID(c), err)
			common.ErrorResp(c, http.StatusBadRequest, "InvalidRequest", "SignHttpFailed")
			return
		}
		if sigResult.Signature != signature {
			logging.Default().Infof("SigMatchFail, RequestId: %s", common.GetRequestID(c))
			common.ErrorRespWithAbort(c, http.StatusUnauthorized, "InvalidSignature", "SigMatchFail")
			return
		}

		// check if allow to access open api admin
		allow := isAdminAllow(c, secret.Tag)
		if !allow {
			common.ErrorRespWithAbort(c, http.StatusForbidden, "Forbidden", "NotAllow")
			return
		}

		common.SetUserInfo(c,
			&common.UserInfo{
				UserID: snowflake.MustParseString(secret.ParentUser),
				Tag:    secret.Tag,
				IsTmp:  secret.SessionToken != "",
			})
		c.Request.Header.Set(common.UserInfoKey, secret.ParentUser)
		c.Request.Header.Set(common.UserAccessKeyId, secret.AccessKeyId)
		c.Next()
		logging.Default().Infof("ReuqestId: %s, Url: %s, Status: %d",
			common.GetRequestID(c), c.Request.URL.String(), c.Writer.Status())
	}
}

func parseTimestamp(ts string) (time.Time, error) {
	i, err := strconv.ParseInt(ts, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	tm := time.Unix(i, 0)
	return tm, nil
}

func healthOk(c *gin.Context) {
	c.String(200, "ok")
}

func isAdminAllow(c *gin.Context, tag string) bool {
	account := v1.IsYuansuanProductAccount(tag)
	// url start with /admin or /internal
	url := strings.HasPrefix(c.Request.URL.Path, "/admin") ||
		strings.HasPrefix(c.Request.URL.Path, "/internal")
	// if account is YS_yuansuan, allow /admin and /internal
	if account && url {
		return true
	}
	// if url not start with /admin or /internal, allow
	if !url {
		return true
	}
	return false
}

func Logger() gin.HandlerFunc {
	return LoggerWithConfig(GetLoggerConfig(nil, nil, nil))
}

func GetLoggerConfig(formatter gin.LogFormatter, output io.Writer,
	skipPaths []string) gin.LoggerConfig {
	return gin.LoggerConfig{
		Formatter: formatter,
		Output:    output,
		SkipPaths: skipPaths,
	}
}

func LoggerWithConfig(conf gin.LoggerConfig) gin.HandlerFunc {
	formatter := conf.Formatter
	if formatter == nil {
		formatter = defaultLogFormatter
	}
	out := conf.Output
	if out == nil {
		out = gin.DefaultWriter
	}
	notlogged := conf.SkipPaths
	isTerm := true
	if w, ok := out.(*os.File); !ok || os.Getenv("TERM") == "dumb" ||
		(!isatty.IsTerminal(w.Fd()) && !isatty.IsCygwinTerminal(w.Fd())) {
		isTerm = false
	}

	if isTerm {
		gin.ForceConsoleColor()
	}

	var skip map[string]struct{}
	if length := len(notlogged); length > 0 {
		skip = make(map[string]struct{}, length)

		for _, path := range notlogged {
			skip[path] = struct{}{}
		}
	}

	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Log only when path is not being skipped
		if _, ok := skip[path]; !ok {
			param := gin.LogFormatterParams{
				Request: c.Request,
				Keys:    c.Keys,
			}

			// Stop timer
			param.TimeStamp = time.Now()
			param.Latency = param.TimeStamp.Sub(start)
			param.ClientIP = c.ClientIP()
			param.Method = c.Request.Method
			param.StatusCode = c.Writer.Status()
			param.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()
			param.BodySize = c.Writer.Size()
			if raw != "" {
				path = path + "?" + raw
			}
			param.Path = path
			log.L(c).Info(formatter(param))
		}
	}
}

// defaultLogFormatter is the default log format function Logger middleware uses.
var defaultLogFormatter = func(param gin.LogFormatterParams) string {
	var statusColor, methodColor, resetColor string
	if param.IsOutputColor() {
		statusColor = param.StatusCodeColor()
		methodColor = param.MethodColor()
		resetColor = param.ResetColor()
	}
	if param.Latency > time.Minute {
		// Truncate in a golang < 1.8 safe way
		param.Latency = param.Latency - param.Latency%time.Second
	}
	return fmt.Sprintf("%s%3d%s - [%s] \"%v %s%s%s %s\" %s",
		// param.TimeStamp.Format("2006/01/02 - 15:04:05"),
		statusColor, param.StatusCode, resetColor,
		param.ClientIP, param.Latency,
		methodColor, param.Method, resetColor,
		param.Path, param.ErrorMessage,
	)
}

func cleanAuditLog(storeIn store.Factory) {
	logging.Default().Infof("start clean audit log client, interval: %d day",
		config.GetConfig().CleanAuditLogInterval)
	c := cron.New(cron.WithSeconds())
	// Schedule the clean job to run every day at 22:00
	_, errAdded := c.AddFunc("0 0 22 * * *", func() {
		logging.Default().Info("start clean audit log")
		err := leader.WithLock(fmt.Sprintf("%s_clean_audit_log", consts.DefaultNamespace), func() error {
			// clean audit log
			dbErr := storeIn.PolicyAudits().CleanThreeMonthAgoData(context.Background())
			if dbErr != nil {
				logging.Default().Errorf("DB clean audit log failed, error: %+v", dbErr)
				return dbErr
			}
			logging.Default().Info("clean audit log success")
			return nil
		}, leader.SetLeaderKeyExpire(1*time.Hour))
		// 多个实例启动时，只有一个实例会抢到锁，其他都会失败， 忽略抢锁失败的错误
		if err != nil && err != leader.ErrFailed {
			logging.Default().Errorf("clean audit log failed, error: %+v", err)
		}
	})
	if errAdded != nil {
		logging.Default().Errorf("add clean audit log cron job failed, error: %+v", errAdded)
	}
	c.Start()
}
