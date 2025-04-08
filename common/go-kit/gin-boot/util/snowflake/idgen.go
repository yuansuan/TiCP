package snowflake

//go:generate mockgen -destination mock_id.go -package snowflake yuansuan.cn/project-root/pkg/common/idgen/snowflake IDGen
import (
	"context"
)

// IDGen ID生成器
type IDGen interface {
	GenID(ctx context.Context) (ID, error)
}
