package state

import (
	"fmt"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/config"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/module/cloud"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/module/cloud/aggregator"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/module/iam"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/module/openapi"
)

type State struct {
	Cloud         cloud.Aggregator
	OpenAPIClient *openapi.Client
	IamClient     *iam.Client
}

func New() (*State, error) {
	s := &State{}

	var err error
	if err = s.init(); err != nil {
		return nil, fmt.Errorf("init failed, %w", err)
	}

	return s, nil
}

func (s *State) init() error {
	c, err := aggregator.New()
	if err != nil {
		return fmt.Errorf("init cloud failed, %w", err)
	}
	s.Cloud = c

	cli, err := openapi.NewClient(config.GetConfig().OpenAPI)
	if err != nil {
		return fmt.Errorf("new openapi client failed, %w", err)
	}
	s.OpenAPIClient = cli

	iamCli, err := iam.NewClient(config.GetConfig().OpenAPI)
	if err != nil {
		return fmt.Errorf("new iam client failed, %w", err)
	}
	s.IamClient = iamCli

	return nil
}

func (s *State) NewCloudStorageClientAfterAssumeRole(zone string, assumeUserId snowflake.ID) (*openapi.Client, error) {
	cloudStorageEndpoint, err := s.OpenAPIClient.GetCloudStorageEndpointByZone(zone)
	if err != nil {
		return nil, fmt.Errorf("get cloud storage endpoint failed, %w", err)
	}

	assumedCredential, err := s.IamClient.AssumeRole(assumeUserId)
	if err != nil {
		return nil, fmt.Errorf("assume role failed, %w", err)
	}

	cloudStorageCli, err := openapi.NewClient(config.OpenAPI{
		AccessKeyId:     assumedCredential.AccessKeyId,
		AccessKeySecret: assumedCredential.AccessKeySecret,
		Endpoint:        cloudStorageEndpoint,
	})
	if err != nil {
		return nil, fmt.Errorf("new cloud storage endpoint failed, %w", err)
	}

	return cloudStorageCli, nil
}
