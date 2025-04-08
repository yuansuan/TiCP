package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/common/go-kit/logging/trace"
	"go.uber.org/zap"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/api/util"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp_dirnfs/model"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common/hashid"

	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp_dirnfs/module/samba"
)

const (
	usernameKey = "username"
)

var options struct {
	Addr    string
	Base    string
	Config  string
	Prefix  string
	UserCfg string

	samba *samba.Samba
}

func init() {
	var err error

	logging.SetLevel(int(zap.DebugLevel))
	if _, err = logging.SetDefault(); err != nil {
		panic(err)
	}

	if options.Addr = os.Getenv("SERVICE_ADDR"); len(options.Addr) == 0 {
		logging.Default().Fatalf("invalid environment for listen address")
	}

	if options.Base = os.Getenv("SHARE_BASE"); len(options.Base) == 0 {
		logging.Default().Fatalf("invalid environment for base dir")
	}

	if options.Config = os.Getenv("SAMBA_CONFIG"); len(options.Config) == 0 {
		logging.Default().Fatalf("invalid environment for samba config")
	}

	if options.UserCfg = os.Getenv("USER_CONFIG"); len(options.UserCfg) == 0 {
		logging.Default().Fatalf("invalid environment for user config")
	}

	options.samba, err = samba.New(options.Config, options.UserCfg)
	if err != nil {
		logging.Default().Fatalf("startup samba daemon failed: %s", err)
	}

	go func() {
		if err := options.samba.Start(context.Background()); err != nil {
			logging.Default().Fatalf("unable to start samba: %s", err)
		}
	}()

	options.Prefix = "YS"
}

func addShareEntry(userID snowflake.ID, username, password, subPath string, excludeUserID bool) error {
	logging.Default().Infof("add share entry: userID: %s, username: %s, password: %s, subPath: %s", userID, username, password, subPath)

	var home string
	if excludeUserID {
		home = filepath.Join(options.Base, subPath)
	} else {
		home = filepath.Join(options.Base, userID.String(), subPath)
	}

	if stat, err := os.Stat(home); err != nil {
		if os.IsNotExist(err) {
			logging.Default().Infof("going to make dir: %s", home)
			if err = os.MkdirAll(home, 0755); err != nil {
				return errors.New("invalid project id or permission denied")
			}
		}

		if err != nil {
			return errors.Wrap(err, "unknown stat error")
		}
	} else if !stat.IsDir() {
		return errors.New("invalid project id or UserHome")
	}

	ctx := context.Background()
	if err := options.samba.AddUser(ctx, username, password, home, true); err == nil {
		err = options.samba.Reload(ctx)
		if err != nil {
			logging.Default().Warnf("ReloadSambaFail, username: %s, Password: %s, home: %s, error: %s",
				username, password, home, err.Error())
			return err
		}
	} else {
		logging.Default().Warnf("AddUserFail, username: %s, Password: %s, home: %s, error: %s",
			username, password, home, err.Error())
		return err
	}
	return nil
}

func addUser(c *gin.Context) {
	log := trace.GetLogger(c)

	username := c.Param(usernameKey)

	userId, err := hashid.Decode(username)
	if err = renderErrorResp(c, http.StatusUnauthorized, err, fmt.Sprintf("decode username %s failed", username)); err != nil {
		log.Errorf("decode username %s failed, %v", username, err)
		return
	}

	req := new(model.AddUserRequest)
	err = c.ShouldBindJSON(req)
	if err = renderErrorResp(c, http.StatusBadRequest, err, "parse add user request body failed"); err != nil {
		logging.Default().Errorf("parse add user request body failed, %v", err)
		return
	}

	username = options.Prefix + username
	logging.Default().Infof("[add user] received request from %q with username is %q, password is %q", c.Request.RemoteAddr, username, req.Password)

	subPath, err := hashid.DecodeStr(req.SubPath)
	if err = renderErrorResp(c, http.StatusBadRequest, err, fmt.Sprintf("decode sub path failed, Username: %s, RawSubPath: %s", username, req.SubPath)); err != nil {
		logging.Default().Errorf("decode sub path failed, Error: %v, Username: %s, RawSubPath: %s", err, username, req.SubPath)
		return
	}

	err = addShareEntry(userId, username, req.Password, subPath, req.ExcludeUserID)
	if err = renderErrorResp(c, http.StatusInternalServerError, err, "add share entry failed"); err != nil {
		logging.Default().Errorf("add share entry for %s failed, %v", username, err)
		return
	}

	logging.Default().Infof("sharing data to %s from %s succeed", username, c.Request.RemoteAddr)
	renderJson(c)
}

func delUser(c *gin.Context) {
	username := c.Param(usernameKey)

	_, err := hashid.Decode(username)
	if err = renderErrorResp(c, http.StatusUnauthorized, err, fmt.Sprintf("decode username %s failed", username)); err != nil {
		logging.Default().Errorf("decode username %s failed, %v", username, err)
		return
	}

	username = options.Prefix + username
	logging.Default().Infof("[delete user] received request from %s with username is %s", c.Request.RemoteAddr, username)

	err = options.samba.DelUser(c, username)
	if err = renderErrorResp(c, http.StatusInternalServerError, err, fmt.Sprintf("delete user failed")); err != nil {
		logging.Default().Errorf("delete user %s failed, %v", username, err)
		return
	}

	logging.Default().Infof("delete share data user %s from %s success", username, c.Request.RemoteAddr)

	renderJson(c)
}

func main() {
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Set(util.RequestIdKeyInHeader, uuid.NewString())
	})

	r.POST("/users/:username", addUser)
	r.DELETE("/users/:username", delUser)

	if err := r.Run(options.Addr); err != nil {
		logging.Default().Fatalf("run http server failed, %v", err)
	}
}

func renderErrorResp(c *gin.Context, code int, err error, errMsg string) error {
	resp := model.BaseResponse{
		RequestId: getRequestId(c),
	}

	if err != nil {
		resp.ErrorMessage = errMsg
		c.JSON(code, resp)
	}
	return err
}

func renderJson(c *gin.Context) {
	c.JSON(http.StatusOK, model.BaseResponse{
		RequestId: getRequestId(c),
	})
}

func getRequestId(c *gin.Context) string {
	requestIdV, exist := c.Get(util.RequestIdKeyInHeader)
	if !exist {
		return ""
	}

	requestId, ok := requestIdV.(string)
	if !ok {
		return ""
	}

	return requestId
}
