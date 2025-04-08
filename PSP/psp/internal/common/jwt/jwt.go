package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"google.golang.org/grpc/status"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/cache"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/gateway/config"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/strutil"
)

const (
	AccessToken           = "access_token"
	RefreshToken          = "refresh_token"
	WhiteListToken        = "WHITE_LIST_TOKEN"
	DefaultWhiteListValue = "UNKNOWN"
)

var (
	mySecret          = []byte("1234yskj")
	ErrorInvalidToken = status.Error(errcode.ErrRBACTokenInvalid, "verify Token Failed")
)

type MyClaim struct {
	UserID   int64  `json:"user_id"`
	UserName string `json:"user_name"`
	jwt.RegisteredClaims
}

// SetToken 设置token，放入缓存
func SetToken(userID int64, userName string, ctx *gin.Context) (aToken, rToken string, err error) {
	aToken, rToken, err = GenToken(userID, userName)

	if err != nil {
		return
	}

	oldToken, _ := ctx.Cookie(AccessToken)

	// 如果存在旧的token则删除缓存
	// 需要异步等待一会儿再做删除，否则并发访问的情况下，由于旧token从白名单中消失，接口会立刻返回401
	if !strutil.IsEmpty(oldToken) {
		go func() {
			time.Sleep(10 * time.Second)
			CleanWhiteListByToken(oldToken)
		}()
	}
	// 返回cookie
	logging.Default().Infof("set token, accessToken:[%v], refreshToken:[%v]", aToken, rToken)
	// 设置cookie
	SetCookie(ctx, aToken, rToken)

	// 存入redis白名单
	ip := ctx.ClientIP()
	if strutil.IsEmpty(ip) {
		ip = DefaultWhiteListValue
	}

	err = saveWhiteList(aToken, ip)

	return
}

func SetCookie(ctx *gin.Context, aToken, rToken string) {
	conf := config.GetConfig()
	RTokenExpiredDuration := time.Duration(conf.TokenExpire) * time.Minute
	ctx.SetCookie(AccessToken, aToken, int(RTokenExpiredDuration.Seconds()), "/", "", false, false)
	ctx.SetCookie(RefreshToken, rToken, int(RTokenExpiredDuration.Seconds()), "/", "", false, false)
}

// CleanWhiteListByJti 根据jti删除白名单 jti即jwt token id
func CleanWhiteListByJti(jti string, userName string) {
	logger := logging.Default()
	if strutil.IsEmpty(jti) || strutil.IsEmpty(userName) {
		logger.Warn("invoke CleanWhiteListByJti but jti or userName is empty")
	}

	oldKey := GetWhiteListKey(userName, jti)

	redisClient := boot.Middleware.DefaultRedis()
	if redisClient == nil || redisClient.Del(oldKey).Err() != nil {
		// 如果redis服务不可用，降级使用go-cache
		cache.Cache.Delete(oldKey)
	}
}

// CleanWhiteListByToken 根据token清理白名单
func CleanWhiteListByToken(aToken string) {
	if strutil.IsEmpty(aToken) {
		logger := logging.Default()
		logger.Warn("invoke CleanWhiteListByToken but aToken is empty")
	}

	myClaim, _ := VerifyToken(aToken)
	if myClaim == nil {
		return
	}

	oldKey := GetWhiteListKey(myClaim.UserName, myClaim.RegisteredClaims.ID)

	redisClient := boot.Middleware.DefaultRedis()
	if redisClient == nil || redisClient.Del(oldKey).Err() != nil {
		// 如果redis服务不可用，降级使用go-cache
		cache.Cache.Delete(oldKey)
	}
}

// CleanWhiteListByUserName 根据用户名清理白名单
func CleanWhiteListByUserName(userName string) {
	logger := logging.Default()
	if strutil.IsEmpty(userName) {
		logger.Warn("invoke CleanWhiteListByUserName but username is empty")
	}

	redisClient := boot.Middleware.DefaultRedis()

	prefixKey := GetWhiteListKey(userName, "")

	var redisSuccess bool

	if redisClient != nil {
		result := redisClient.Scan(0, fmt.Sprintf("%s*", prefixKey), 0)
		if result.Err() == nil {
			redisSuccess = true
		}
		iterator := result.Iterator()

		for iterator.Next() {
			redisClient.Del(iterator.Val())
		}
	}

	if !redisSuccess {
		// 如果redis服务不可用，降级使用go-cache
		for key := range cache.Cache.Items() {
			if strings.HasPrefix(key, prefixKey) {
				cache.Cache.Delete(key)
			}
		}
	}

}

// GetWhiteListKey 获取白名单key
func GetWhiteListKey(userName, jti string) string {
	return fmt.Sprintf("PSP:%s:%s:%s", WhiteListToken, userName, jti)
}

// GenToken 颁发token access token 和 refresh token
func GenToken(userID int64, userName string) (aToken, rToken string, err error) {
	conf := config.GetConfig()
	RTokenExpiredDuration := time.Duration(conf.TokenExpire) * time.Minute
	ATokenExpiredDuration := time.Duration(conf.TokenExpire/2) * time.Minute
	v4, _ := uuid.NewV4()
	rc := jwt.RegisteredClaims{
		ExpiresAt: getJWTTime(ATokenExpiredDuration),
		ID:        v4.String(),
	}
	myClaim := MyClaim{
		UserID:           userID,
		UserName:         userName,
		RegisteredClaims: rc,
	}

	// refresh token 不需要保存任何用户信息
	standardClaims := rc
	value, _ := uuid.NewV4()
	standardClaims.ID = value.String()
	standardClaims.ExpiresAt = getJWTTime(RTokenExpiredDuration)

	aToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, myClaim).SignedString(mySecret)
	rToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, standardClaims).SignedString(mySecret)
	return
}

// VerifyToken 验证Token
func VerifyToken(tokenID string) (*MyClaim, error) {
	var myc = new(MyClaim)
	token, err := jwt.ParseWithClaims(tokenID, myc, keyFunc)
	if err != nil {
		return myc, err
	}
	if !token.Valid {
		err = ErrorInvalidToken
		return myc, err
	}

	return myc, nil
}

// RefreshJWTToken 通过 refresh token 刷新 atoken
func RefreshJWTToken(accessToken, refreshToken string, ctx *gin.Context) (newAToken, newRToken string, err error) {
	logging.Default().Infof("begin refresh token, oldAccessToken:[%v],refreshToken:[%v]", accessToken, refreshToken)
	// rToken 无效直接返回
	if _, err = jwt.Parse(refreshToken, keyFunc); err != nil {
		return
	}
	// 从旧access token 中解析出claims数据
	var claim MyClaim
	_, err = jwt.ParseWithClaims(accessToken, &claim, keyFunc)
	// 判断错误是不是因为access token 正常过期导致的
	v, _ := err.(*jwt.ValidationError)
	if v.Errors == jwt.ValidationErrorExpired {
		return SetToken(claim.UserID, claim.UserName, ctx)
	}
	return
}

// 存入缓存白名单
func saveWhiteList(aToken string, cacheValue string) error {
	myClaim, err := VerifyToken(aToken)

	if err != nil {
		return err
	}

	conf := config.GetConfig()
	RTokenExpiredDuration := time.Duration(conf.TokenExpire) * time.Minute

	redisClient := boot.Middleware.DefaultRedis()

	if redisClient == nil || redisClient.Set(GetWhiteListKey(myClaim.UserName, myClaim.RegisteredClaims.ID), cacheValue, RTokenExpiredDuration).Err() != nil {
		// 如果redis不可用，则降级为使用go-cache
		cache.Cache.Set(GetWhiteListKey(myClaim.UserName, myClaim.RegisteredClaims.ID), cacheValue, RTokenExpiredDuration)
	}
	return nil
}

// CheckWhiteList 检查access_token是否存在于白名单中
func CheckWhiteList(accessToken string) bool {
	redisClient := boot.Middleware.DefaultRedis()

	myClaim, _ := VerifyToken(accessToken)

	if myClaim == nil {
		return false
	}

	key := GetWhiteListKey(myClaim.UserName, myClaim.RegisteredClaims.ID)
	var value string
	var err error

	if redisClient != nil {
		value, err = redisClient.Get(key).Result()
	}

	if redisClient == nil || err != nil {
		v, ok := cache.Cache.Get(key)
		if ok {
			value = v.(string)
		}
	}

	return !strutil.IsEmpty(value)
}

func getJWTTime(t time.Duration) *jwt.NumericDate {
	return jwt.NewNumericDate(time.Now().Add(t))
}

func keyFunc(token *jwt.Token) (interface{}, error) {
	return mySecret, nil
}
