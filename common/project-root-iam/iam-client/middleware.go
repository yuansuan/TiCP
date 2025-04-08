package iam_client

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/common/openapi-go/credential"
	"github.com/yuansuan/ticp/common/openapi-go/utils/signer"
	v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
	iam_api "github.com/yuansuan/ticp/common/project-root-iam/iam-api"
)

const (
	expireTime   = 5 * time.Minute
	userIDKey    = "x-ys-user-id"
	requestIDKey = "x-ys-request-id"

	requestIDKeyInCtx    = "RequestId"
	userAppKeyInQuery    = "AccessKeyId"
	userAppKeyInQueryBak = "AppKey"
	signatureKey         = "Signature"
	timestampKey         = "Timestamp"
)

type IamConfig struct {
	Endpoint  string
	AppKey    string
	AppSecret string
	Proxy     string
}

// middleware valid url start with "/system" and "/admin"
func ValidUserIDMiddleware(userID string) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		isAdminUrl := strings.HasPrefix(ctx.Request.URL.Path, "/admin") || strings.HasPrefix(ctx.Request.URL.Path, "/system")
		if isAdminUrl {
			// get userID from header
			userIDFromHeader := ctx.GetHeader(userIDKey)
			if userIDFromHeader != userID {
				logging.Default().Infof("user id is not valid, userIDFromHeader: %s, userID: %s", userIDFromHeader, userID)
				ifError(ctx, http.StatusUnauthorized, "InvalidUserID", "user id is not valid")
				return
			}
		}
	}
}

func SignatureValidateMiddleware(iamConfig IamConfig) func(ctx *gin.Context) {
	return func(c *gin.Context) {
		logger := logging.GetLogger(c)
		logger.Debug("signature validate middleware called")
		if getRequestId(c) == "" {
			reqId := uuid.New().String()
			c.Set(requestIDKey, reqId)
			c.Set(requestIDKeyInCtx, reqId)
			c.Request.Header.Set(requestIDKey, reqId)
			c.Writer.Header().Set(requestIDKey, reqId)
		}
		var err error
		requestID := getRequestId(c)
		userAppKey := c.Query(userAppKeyInQuery)
		if userAppKey == "" {
			userAppKey = c.Query(userAppKeyInQueryBak)
			if userAppKey == "" {
				logging.Default().Infof("AppKey is empty in query params, requestID: %s", requestID)
				ifError(c, http.StatusUnauthorized, "InvalidAppKey", "AppKey is empty in query params")
				return
			}
		}

		signature := c.Query(signatureKey)
		if signature == "" {
			logging.Default().Infof("Signature is empty in query params requestID: %s", requestID)
			ifError(c, http.StatusUnauthorized, "InvalidSignature", "Signature is empty in query params")
			return
		}

		timestamp, err := parseTimestamp(c.Query(timestampKey))
		if err != nil {
			logging.Default().Infof("Timestamp is invalid in query params %s, requestID: %s", err.Error(), requestID)
			ifError(c, http.StatusUnauthorized, "InvalidTimestamp", "Timestamp is invalid in query params")
			return
		}

		now := time.Now()
		if now.Sub(timestamp) > expireTime {
			logging.Default().Infof("Signature expired, requestID: %s", requestID)
			ifError(c, http.StatusUnauthorized, "InvalidSignature", "Signature expired")
			return
		}

		client := NewClient(iamConfig.Endpoint, iamConfig.AppKey, iamConfig.AppSecret)
		if iamConfig.Proxy != "" {
			client.SetProxy(iamConfig.Proxy)
		}
		secret, err := client.GetSecret(&iam_api.GetSecretRequest{
			AccessKeyId: userAppKey,
		})
		if err != nil {
			logging.Default().Infof("get secret %s from iam server failed %s, requestID: %s", userAppKey, err.Error(), requestID)
			switchErr(c, err)
			return
		}

		if keyExpired(secret.Expire) {
			logging.Default().Infof("service appKey expired requestID: %s", requestID)
			ifError(c, http.StatusUnauthorized, "InvalidAppKey", "invalid appKey")
			return
		}

		signAgain, sigSourceStr, err := sign(secret.AccessKeyId, secret.AccessKeySecret, c.Request)
		if err != nil {
			logging.Default().Infof("sign failed %s, requestID: %s", err.Error(), requestID)
			ifError(c, http.StatusUnauthorized, "InvalidSignature", "invalid signature")
			return
		}

		if signature != signAgain {
			logging.Default().Infof("new signature %s calculated not equal with old signature in query params %s, requestID: %s, raw string: %s, head: %+v",
				signAgain, signature, requestID, sigSourceStr, c.Request.Header)
			ifError(c, http.StatusUnauthorized, "InvalidSignature", "invalid signature")
			return
		}

		c.Request.Header.Set(userIDKey, secret.YSId)
		c.Request.Header.Set(userAppKeyInQuery, userAppKey)
		c.Set(userIDKey, secret.YSId)
		c.Set(userAppKeyInQuery, userAppKey)
		c.Next()
	}
}

func sign(akId, akSecret string, req *http.Request) (string, string, error) {
	cred := credential.NewCredential(akId, akSecret)
	signer, err := signer.NewSigner(cred)
	if err != nil {
		return "", "", fmt.Errorf("new signer failed, %w", err)
	}

	fingerprint, err := signer.SignHttp(req)
	if err != nil {
		return "", "", fmt.Errorf("sign http failed, %w", err)
	}

	return fingerprint.Signature, fingerprint.SourceStr, nil
}

// keyExpired checks if a key has expired, if the value of user.SessionState.Expires is 0, it will be ignored.
func keyExpired(expires time.Time) bool {
	return time.Now().After(expires)
}

func parseTimestamp(ts string) (time.Time, error) {
	i, err := strconv.ParseInt(ts, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	tm := time.Unix(i, 0)
	return tm, nil
}

func getRequestId(c *gin.Context) string {
	requestIDI, exist := c.Get(requestIDKey)
	if !exist {
		return c.GetHeader(requestIDKey)
	}

	requestID, ok := requestIDI.(string)
	if !ok {
		return ""
	}

	return requestID
}

func ifError(c *gin.Context, status int, errorCode, errorMsg string) {
	c.AbortWithStatusJSON(status, v20230530.Response{
		ErrorCode: errorCode,
		ErrorMsg:  errorMsg,
		RequestID: getRequestId(c),
	})

}

func switchErr(c *gin.Context, err error) {
	var res ErrorResponse
	isErrRespType := errors.As(err, &res)
	if isErrRespType {
		switch res.Status {
		case http.StatusNotFound:
			ifError(c, http.StatusUnauthorized, "Unauthorized", "Unauthorized")
			return
		default:
			ifError(c, http.StatusInternalServerError, res.Code, res.Message)
			return
		}
	} else {
		ifError(c, http.StatusInternalServerError, "Internal Server Error", err.Error())
		return
	}
}
