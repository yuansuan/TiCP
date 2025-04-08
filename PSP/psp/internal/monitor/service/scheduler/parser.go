package scheduler

import (
	"context"

	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi"
)

const (
	// Slurm PbsPro 调度器类型
	Slurm  = "slurm"
	PbsPro = "pbspro"
)

// CreateParser 根据数据类型返回对应的解析器实例
func CreateParser(dataType string) Parser {
	switch dataType {
	case Slurm:
		return &SlurmParser{}
	case PbsPro:
		return &PbsProParser{}
	default:
		logging.Default().Warnf("not support scheduler type [%s]", dataType)
	}
	return nil
}

// Parser 解析器接口
type Parser interface {
	Parse(ctx context.Context, openapi *openapi.OpenAPI) ([]map[string]string, error)
	GetCommand(nodeName, operation string) (command string)
	GetStateInfoMap(ctx context.Context, openapi *openapi.OpenAPI) (map[string]string, error)
	GetClusterName(ctx context.Context, openapi *openapi.OpenAPI) (string, error)
}
