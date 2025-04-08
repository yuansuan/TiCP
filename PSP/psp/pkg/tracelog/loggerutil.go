package tracelog

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/natefinch/lumberjack"
	"github.com/pkg/errors"
	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/PSP/psp/cmd/config"
	"github.com/yuansuan/ticp/PSP/psp/internal/common"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/ginutil"
)

const (
	Msg  = "msg"
	Time = "time"
)

var logFile *lumberjack.Logger
var logger log.Logger

func InitLogger() {
	//获取日志大小和日志文件数量
	loggerConfig := config.GetConfig().Logger

	//检查日志文件路径是否规范
	if err := checkFilePath(loggerConfig.LogDir); err != nil {
		logging.Default().Errorf("init logger failed err: %v", err)
		panic(errors.Wrap(err, "failed to init logger info"))
	}

	maxSize, maxBackups := getLogFileArgs(loggerConfig.MaxSize, loggerConfig.BackupCount)
	//获取日志文件路径
	logFile = &lumberjack.Logger{
		Filename:   loggerConfig.LogDir, // 日志文件路径
		MaxSize:    maxSize,             // 每个日志文件的最大大小（MB）
		MaxBackups: maxBackups,          // 保留的旧日志文件数量
		MaxAge:     loggerConfig.MaxAge, // 保留的旧日志文件天数
		Compress:   true,                // 是否压缩旧日志文件
	}

	// 设置日志输出目标为文件
	logger = log.NewLogfmtLogger(logFile)
	logger = log.With(logger, Time, log.DefaultTimestamp)
}

func Info(ctx context.Context, msg string) {
	level.Info(logger).Log(common.TraceId, getTraceID(ctx), Msg, msg)
}

func Error(ctx context.Context, msg string) {
	level.Error(logger).Log(common.TraceId, getTraceID(ctx), Msg, msg)
}

func CloseLogger() {
	logFile.Close()
}

func getLogFileArgs(maxSize, maxNum int) (outSize int, outNum int) {
	if maxSize == 0 {
		outSize = 50
	}
	if maxNum == 0 {
		outNum = 20
	}
	return outSize, outNum
}

func checkFilePath(path string) error {
	if path == "" {
		return errors.Errorf("custom config The log_dir parameter is not set")
	}

	isAbs := filepath.IsAbs(path)
	if !isAbs {
		return errors.Errorf("log_dir must be absolute path")
	}

	_, fileName := filepath.Split(path)
	if !strings.Contains(fileName, ".") {
		return errors.Errorf("log_dir must be a full path")
	}

	return nil
}

func getTraceID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	if ginContext, ok := ctx.(*gin.Context); ok {
		return ginutil.GetTraceID(ginContext)
	}

	return ""
}
