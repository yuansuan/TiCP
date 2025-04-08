package middleware

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging/trace"

	iamclient "github.com/yuansuan/ticp/common/project-root-iam/iam-client"

	"github.com/yuansuan/ticp/iPaaS/standard-compute/config"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/response"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/state"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/pkg/errorcode"
)

func SignatureValidator(c *gin.Context) {
	logger := trace.GetLogger(c)
	conf, err := getConfig(c)
	if err = response.InternalErrorIfError(c, err, errorcode.InternalServerError); err != nil {
		logger.Errorf("get config from gin ctx failed, %v", err)
		return
	}

	iamclient.SignatureValidateMiddleware(createIamConfig(conf))(c)
}

func createIamConfig(conf *config.Config) iamclient.IamConfig {
	iamConfig := iamclient.IamConfig{
		Endpoint:  conf.Iam.Endpoint,
		AppKey:    conf.Iam.AppKey,
		AppSecret: conf.Iam.AppSecret,
	}
	if conf.Iam.Proxy != "" {
		iamConfig.Proxy = conf.Iam.Proxy
	}

	return iamConfig
}

func getConfig(c *gin.Context) (*config.Config, error) {
	stateI, exist := c.Get("state")
	if !exist {
		return nil, errors.New("ginCtx['state'] not exist")
	}

	s, ok := stateI.(*state.State)
	if !ok {
		return nil, errors.New("ginCtx['state'] cannot convert to *state.State")
	}

	return s.Conf, nil
}
