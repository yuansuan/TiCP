package taskgroup

import (
	"context"
)

type task interface {
	Name() string
	Start(ctx context.Context) error
}
