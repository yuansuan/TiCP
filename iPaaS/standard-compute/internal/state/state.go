package state

import (
	"xorm.io/xorm"

	"github.com/yuansuan/ticp/iPaaS/standard-compute/config"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/backend"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/statemachine"
)

type State struct {
	Conf         *config.Config
	DB           *xorm.Engine
	JobScheduler backend.Provider
	Factory      *statemachine.Factory
}
