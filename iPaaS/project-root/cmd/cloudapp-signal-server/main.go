package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp-signal-server/config"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp-signal-server/handler"
)

const DefaultConfigFilename = "config/config.yaml"

func main() {
	cfg, err := config.Read(DefaultConfigFilename)
	if err != nil {
		logging.Default().Panicf("read config failed: %s", err)
	}

	h, err := handler.New(cfg)
	if err != nil {
		logging.Default().Panicf("init handler failed: %s", err)
	}

	server := http.Server{
		Addr:    cfg.Server.Addr,
		Handler: h.ExportHttp(),
	}

	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGILL, syscall.SIGTERM, syscall.SIGKILL)

		logging.Default().Infof("received signal %q from system, shutdown now ...", <-signals)

		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			if err := server.Shutdown(ctx); err != nil {
				logging.Default().Panicf("shutdown server failed: %s", err)
			}
		}()

		signal.Stop(signals)
		close(signals)
	}()

	logging.Default().Infof("server has started on %q", cfg.Server.Addr)
	if err = server.ListenAndServe(); err != nil {
		logging.Default().Errorf("http listen failed: %s", err)
	}
}
