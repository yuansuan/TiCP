package errors

import (
	"github.com/pkg/errors"
)

var (
	ErrSnowflakeGeneration = errors.New("SnowflakeID.Generation")
)
