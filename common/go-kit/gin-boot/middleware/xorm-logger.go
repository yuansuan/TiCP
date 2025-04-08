package middleware

import (
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"xorm.io/xorm/log"
)

type XormLogger struct {
	*logging.Logger

	showSQL bool
}

func NewXormLogger(l *logging.Logger, showSQL bool) *XormLogger {
	return &XormLogger{
		Logger:  l,
		showSQL: showSQL,
	}
}

// ShowSQL ...
func (w *XormLogger) ShowSQL(show ...bool) {
	if len(show) == 0 {
		w.showSQL = true
	} else {
		w.showSQL = show[0]
	}
}

// IsShowSQL ...
func (w *XormLogger) IsShowSQL() bool {
	return w.showSQL
}

// IsShowSQL ...
func (w *XormLogger) Level() log.LogLevel {
	return log.DEFAULT_LOG_LEVEL
}

// SetLevel ...
func (w *XormLogger) SetLevel(l log.LogLevel) {

}
