package service

import (
	"context"

	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/PSP/psp/internal/app/service/impl"
)

func InitDaemon() {
	ctx := context.Background()
	logger := logging.GetLogger(ctx)

	appLoader, err := impl.NewAppLoader()
	if err != nil {
		logger.Errorf("failed to init app loader")
	}

	go appLoader.CheckInternalAppTemplate(ctx)
}
