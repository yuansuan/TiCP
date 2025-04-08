package daemon

import (
	"github.com/pkg/errors"
)

// InitDaemon ...
func InitDaemon() {
	check, err := NewProjectCheck()
	if err != nil {
		panic(errors.Wrap(err, "failed to init daemon"))
	}

	check.ProjectCheckStart()
}
