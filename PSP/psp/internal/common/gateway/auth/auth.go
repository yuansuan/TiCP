package auth

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/PSP/psp/internal/common"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/cache"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/gateway/config"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/gateway/openapicert"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/jwt"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/ginutil"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/strutil"
)

const (
	HeaderWsKey      = "Upgrade"
	Websocket        = "Websocket"
	RefreshedFlagKey = "RefreshedFlagKey"
)

type RefreshCache struct {
	newAccessToken  string
	newRefreshToken string
}

var mutex sync.Mutex

// BasicAuth 权限校验
func BasicAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		logger := logging.GetLogger(ctx)

		// 说明是openapi调用
		if strutil.IsNotEmpty(ctx.Request.Header.Get(common.HttpHeaderOpenapiCertificate)) {
			if openapicert.CertCheck(ctx) {
				ctx.Next()
				return
			} else {
				logger.Info("openapi certificate not exist, abort")
				ctx.AbortWithStatus(http.StatusUnauthorized)
				return
			}
		}

		// ws不拦截
		if ctx.Request.Header.Get(HeaderWsKey) == Websocket {
			ctx.Next()
			return
		}

		url := ctx.Request.URL.Path
		method := ctx.Request.Method
		logger.Infof("start basic auth, url:%s", url)

		// 请求url在白名单中则放行
		if CheckWhiteUrl(url, method) {
			logger.Infof("request url in white url list,pass")
			ctx.Next()
			return
		}

		// token是否存在
		accessToken, _ := ctx.Cookie(jwt.AccessToken)
		refreshToken, _ := ctx.Cookie(jwt.RefreshToken)
		if accessToken == "" || refreshToken == "" {
			logger.Info("token not exist, abort")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// at是否在redis白名单中
		if !jwt.CheckWhiteList(accessToken) {
			logger.Infof("access token:[%v] not in white token list, abort", accessToken)
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// at是否失效, 没失效直接放行
		myClaim, err := jwt.VerifyToken(accessToken)
		// 说明已at已失效，rt未失效
		if err != nil {
			// 判断此accessToken是否已经刷新,未刷新则加锁准备刷新
			key := fmt.Sprintf("%s:%s", RefreshedFlagKey, accessToken)
			refreshCache, ok := cache.Cache.Get(key)
			if !ok {
				func() {
					mutex.Lock()
					defer mutex.Unlock()
					// 双重检查
					_, ok = cache.Cache.Get(key)

					if !ok {
						// 尝试刷新token
						newAccessToken, newRefreshToken, err := jwt.RefreshJWTToken(accessToken, refreshToken, ctx)
						if err != nil {
							logger.Info("all token expired, abort")
							ctx.AbortWithStatus(http.StatusUnauthorized)
							return
						}
						cache.Cache.Set(key, RefreshCache{
							newAccessToken:  newAccessToken,
							newRefreshToken: newRefreshToken,
						}, time.Duration(30)*time.Second)
						logger.Infof("refresh token success, at:%s, rt:%s", newAccessToken, newRefreshToken)
					}
				}()
			} else {
				token := refreshCache.(RefreshCache)
				jwt.SetCookie(ctx, token.newAccessToken, token.newRefreshToken)
			}
		}

		ginutil.SetUser(ctx, myClaim.UserID, myClaim.UserName)
		v4, _ := uuid.NewV4()
		ginutil.SetTraceID(ctx, v4.String())
		ctx.Next()
	}
}

// CheckWhiteUrl 检查url是否在url白名单配置中
func CheckWhiteUrl(url, method string) bool {
	whiteUrlMap := config.GetConfig().WhiteUrlMap

	return whiteUrlMap != nil && (whiteUrlMap[fmt.Sprintf("%v::%v", url, method)] || whiteUrlMap[url])
}
