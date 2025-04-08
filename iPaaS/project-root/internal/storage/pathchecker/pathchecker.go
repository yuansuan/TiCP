package pathchecker

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	commoncode "github.com/yuansuan/ticp/common/project-root-api/common"
	iam_client "github.com/yuansuan/ticp/common/project-root-iam/iam-client"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/config"
)

const (
	cloudStorageYSProductName = "cs"

	// SystemURLPrefix 系统api前缀
	SystemURLPrefix = "/system"
	// AdminURLPrefix 管理api前缀
	AdminURLPrefix = "/admin"
)

// PathAccessChecker 路径访问权限检查
type PathAccessChecker interface {
	CheckPathAccess(accessKey, currentUserID, path string, logger *logging.Logger) (bool, error)
	CheckPathAccessAndHandleError(accessKey string, userID string, path string, logger *logging.Logger, ctx *gin.Context) bool
	GetUserIDAndAKAndHandleError(ctx *gin.Context, prefix string) (string, string, bool, error)
}

// PathAccessCheckerImpl 路径访问权限检查
type PathAccessCheckerImpl struct {
	AuthEnabled       bool
	IamClient         *iam_client.IamClient
	UserIDKey         string
	UserAppKeyInQuery string
}

// CheckPathAccess 检查用户是否有访问路径的权限
func (p *PathAccessCheckerImpl) CheckPathAccess(accessKey, currentUserID, path string, logger *logging.Logger) (bool, error) {
	parts := strings.Split(path, "/")
	userID := parts[1]
	if currentUserID == userID {
		return true, nil
	}

	logger.Infof("check path access, accessKey: %s, path: %s", accessKey, path)
	allowDefault, err := p.IamClient.IsAllowDefault(accessKey, path, cloudStorageYSProductName)
	if err != nil {
		logger.Warnf("call iamClient.IsAllowDefault error, err: %v", err)
		return false, err
	}
	if !allowDefault.Allow {
		logger.Infof("user %s is not allow to access path %s, msg: %s", accessKey, path, allowDefault.Message)
	}

	return allowDefault.Allow, nil
}

// CheckPathAccessAndHandleError 检查用户是否有访问路径的权限，如果没有权限，返回错误信息
func (p *PathAccessCheckerImpl) CheckPathAccessAndHandleError(accessKey string, userID string, path string, logger *logging.Logger, ctx *gin.Context) bool {

	if !p.AuthEnabled {
		return true
	}

	flag, err := p.CheckPathAccess(accessKey, userID, path, logger)
	if err != nil {
		if res, ok := err.(iam_client.ErrorResponse); ok {
			switch res.Status {
			case http.StatusForbidden:
				common.ErrorResp(ctx, http.StatusForbidden, res.Code, res.Message)
			case http.StatusNotFound:
				common.ErrorResp(ctx, http.StatusNotFound, res.Code, res.Message)
			case http.StatusBadRequest:
				common.ErrorResp(ctx, http.StatusBadRequest, res.Code, res.Message)
			default:
				common.ErrorResp(ctx, http.StatusInternalServerError, res.Code, res.Message)
			}
			logger.Infof("call iamClient.IsAllowDefault error, err: %v", err)
			return false
		} else {
			common.InternalServerError(ctx, err.Error())
			logger.Errorf("call iamClient.IsAllowDefault error, err: %v", err)
			return false
		}
	}
	if !flag {
		msg := fmt.Sprintf("user has no access to the path, path: %s, accessKey: %s", path, accessKey)
		logger.Info(msg)
		common.ErrorResp(ctx, http.StatusForbidden, commoncode.AccessDeniedErrorCode, msg)
		return false
	}

	return true
}

// GetUserIDAndAKAndHandleError 获取用户id和ak
func (p *PathAccessCheckerImpl) GetUserIDAndAKAndHandleError(ctx *gin.Context, prefix string) (string, string, bool, error) {
	var userID, accessKey string

	if prefix != SystemURLPrefix && prefix != AdminURLPrefix {
		msg := fmt.Sprintf("prefix must be %s or %s", SystemURLPrefix, AdminURLPrefix)
		logging.GetLogger(ctx).Info(msg)
	}

	prefixMsg := strings.TrimPrefix(prefix, "/")

	flag := false
	if strings.HasPrefix(ctx.Request.URL.String(), prefix) {
		flag = true
	}

	if p.AuthEnabled {
		userID = ctx.Request.Header.Get(p.UserIDKey)
		accessKey = ctx.Request.Header.Get(p.UserAppKeyInQuery)
		if flag && accessKey != config.GetConfig().AccessKeyId {
			msg := fmt.Sprintf("call %s api must use %s access key,current accessKey: %s", prefixMsg, prefixMsg, accessKey)
			logging.GetLogger(ctx).Info(msg)
			common.ErrorResp(ctx, http.StatusForbidden, commoncode.AccessDeniedErrorCode, msg)
			return "", "", flag, errors.New(msg)
		}
	} else {
		userID = config.GetConfig().YsId
		accessKey = config.GetConfig().AccessKeyId
	}

	return userID, accessKey, flag, nil
}

// GetAuthInfo 获取信息
func (p *PathAccessCheckerImpl) GetAuthInfo() (authEnabled bool, userIDKey, userAppKeyInQuery string) {
	return p.AuthEnabled, p.UserIDKey, p.UserAppKeyInQuery
}

// GetIamClient 获取iam客户端
func (p *PathAccessCheckerImpl) GetIamClient() *iam_client.IamClient {
	return p.IamClient
}
