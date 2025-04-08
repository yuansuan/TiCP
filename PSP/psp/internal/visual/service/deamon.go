package service

import (
	"github.com/pkg/errors"
	"github.com/yuansuan/ticp/common/go-kit/logging"

	mainconfig "github.com/yuansuan/ticp/PSP/psp/cmd/config"
	"github.com/yuansuan/ticp/PSP/psp/internal/visual/config"
	"github.com/yuansuan/ticp/PSP/psp/internal/visual/service/impl"
)

func InitDaemon() {
	logger := logging.Default()

	// 根据全局配置启用
	if !mainconfig.Custom.Main.EnableVisual {
		logger.Infof("sync session data routine has disabled by global config")
		return
	}

	// 根据配置启用
	syncData := config.GetConfig().SyncData
	if !syncData.Enable {
		logger.Infof("sync session data routine has disabled")
		return
	}

	sessionScanner, err := impl.NewSessionScanner()
	if err != nil {
		panic(errors.Wrap(err, "failed to init session scanner"))
	}

	go sessionScanner.SyncSessionDataRoutine()
	go sessionScanner.SyncSessionStatusRoutine()

	sessionNotification, err := impl.NewSessionNotification()
	if err != nil {
		logger.Errorf("failed to init session notification")
	}

	go sessionNotification.SessionCheckTicker()
}
