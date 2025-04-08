package start

import (
	"fmt"
	"github.com/yuansuan/ticp/common/openapi-go/utils/payby"
	"github.com/yuansuan/ticp/common/project-root-api/cloud_app/v1/session"
	v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
	"time"
)

type Option func(req *session.ApiPostRequest) error

func (api API) HardwareId(hardwareId string) Option {
	return func(req *session.ApiPostRequest) error {
		req.HardwareId = &hardwareId
		return nil
	}
}

func (api API) SoftwareId(softwareId string) Option {
	return func(req *session.ApiPostRequest) error {
		req.SoftwareId = &softwareId
		return nil
	}
}

func (api API) MountPaths(mountPaths map[string]string) Option {
	return func(req *session.ApiPostRequest) error {
		req.MountPaths = &mountPaths
		return nil
	}
}

func (api API) ChargeParams(chargeParams v20230530.ChargeParams) Option {
	return func(req *session.ApiPostRequest) error {
		req.ChargeParams = &chargeParams
		return nil
	}
}

func (api API) PayBy(payBy string) Option {
	return func(req *session.ApiPostRequest) error {
		req.PayBy = &payBy
		return nil
	}
}

func (api API) PayByParams(payByAccessKeyID, payByAccessSecret string) Option {
	return func(req *session.ApiPostRequest) error {
		// 以 用户自定义payBy 优先
		if req.PayBy != nil {
			return nil
		}

		resourceTag := fmt.Sprintf("%s_%s", *req.HardwareId, *req.SoftwareId)
		timestamp := time.Now().UTC().UnixMilli()
		payBy, _ := payby.NewPayBy(payByAccessKeyID, payByAccessSecret, resourceTag, timestamp)
		token := payBy.Token()
		req.PayBy = &token
		return nil
	}
}
