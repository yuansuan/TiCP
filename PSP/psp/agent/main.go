package main

import (
	"fmt"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/natefinch/lumberjack"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	stdlog "log"
	"net/http"

	"github.com/yuansuan/ticp/PSP/psp/internal/agent/collector"
	"github.com/yuansuan/ticp/PSP/psp/internal/agent/util"
)

func main() {
	var logPath string
	var port string
	var maxSize string
	var maxNum string
	cmd := &cobra.Command{
		Run: func(cmd *cobra.Command, args []string) {
			port = cmd.Flag("port").Value.String()
			logPath = cmd.Flag("log_path").Value.String()
			maxSize = cmd.Flag("log_max_size").Value.String()
			maxNum = cmd.Flag("log_max_num").Value.String()
		},
	}

	// 添加 flag
	cmd.Flags().Int("port", 9100, "Port number")
	cmd.Flags().String("log_path", "", "Log directory")
	cmd.Flags().String("log_max_size", "50", "Log maxSize")
	cmd.Flags().String("log_max_num", "20", "Log maxBackups")

	// 执行根命令
	if err := cmd.Execute(); err != nil {
		stdlog.Println(err)
		return
	}

	//判断文件路径是否合规
	result := util.CheckFilePath(logPath)
	if !result {
		return
	}

	size, num := util.GetLogFileArg(maxSize, maxNum)
	logFile := &lumberjack.Logger{
		Filename:   logPath, // 日志文件路径
		MaxSize:    size,    // 每个日志文件的最大大小（MB）
		MaxBackups: num,     // 保留的旧日志文件数量
		Compress:   false,   // 是否压缩旧日志文件
	}

	defer logFile.Close()

	// 设置日志输出目标为文件
	logger := log.NewLogfmtLogger(logFile)
	//设置日志时间
	logger = log.With(logger, "time", log.DefaultTimestamp)

	level.Info(logger).Log("msg", fmt.Sprintf("Service port :%s", port))

	handler, err := NewHandler(logger)
	if err != nil {
		level.Error(logger).Log("Error occur when start server ", err)
		return
	}

	http.Handle("/metrics", handler)
	if err = http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
		level.Error(logger).Log("Error occur when start server ", err)
		return
	}
}

func NewHandler(logger log.Logger) (http.Handler, error) {
	nc, err := collector.NewCollector(logger)
	if err != nil {
		return nil, fmt.Errorf("couldn't create collector: %s", err)
	}

	registry := prometheus.NewRegistry()
	if err := registry.Register(nc); err != nil {
		return nil, fmt.Errorf("couldn't register node collector: %s", err)
	}

	handler := promhttp.HandlerFor(
		prometheus.Gatherers{registry},
		promhttp.HandlerOpts{
			ErrorLog:      stdlog.New(log.NewStdlibAdapter(level.Error(logger)), "", 0),
			ErrorHandling: promhttp.ContinueOnError,
		})

	return handler, nil
}
