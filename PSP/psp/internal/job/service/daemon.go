package service

import (
	"github.com/pkg/errors"

	"github.com/yuansuan/ticp/PSP/psp/internal/job/service/dataloader"
)

// InitDaemon ...
func InitDaemon() {
	jobLoader, err := dataloader.NewJobLoader()
	if err != nil {
		panic(errors.Wrap(err, "failed to init daemon"))
	}

	jobLoader.JobLoaderStart()
}
