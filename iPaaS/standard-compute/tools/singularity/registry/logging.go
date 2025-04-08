package registry

type Kind int

const (
	// KindDebug 调试信息
	KindDebug Kind = iota // 0
	// KindInfo 普通信息
	KindInfo // 1
)

// Logger 是一个专用的日志记录器
type Logger func(kind Kind, head, msg string, ctx map[string]interface{})

// NewDiscardLogger 创建一个什么都不做的日志记录器
func NewDiscardLogger() Logger {
	return func(kind Kind, head, msg string, ctx map[string]interface{}) {}
}
