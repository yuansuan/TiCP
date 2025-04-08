package main

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp_agent/pkg/daemon"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp_agent/pkg/shutdown"
)

type globalCmd struct {
	cmd *cobra.Command

	serverAddr    string
	useConsole    bool
	logLevel      string
	logPath       string
	customEnvPath string
}

func main() {
	gc := &globalCmd{
		cmd: &cobra.Command{
			Use:   "ys-agent",
			Short: "ys-agent",
			Long:  "ys-agent",
		},
	}
	gc.cmd.RunE = gc.run

	gc.initFlags()

	if err := gc.cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func (gc *globalCmd) initFlags() {
	gc.cmd.Flags().StringVarP(&gc.serverAddr, "listen", "l", ":3390", "http server listen address")
	gc.cmd.Flags().BoolVar(&gc.useConsole, "use-console", false, "log in console or file")
	gc.cmd.Flags().StringVar(&gc.logLevel, "log-level", "", "log level [ info | debug ]")
	gc.cmd.Flags().StringVar(&gc.logPath, "log-path", "", "log path")
	gc.cmd.Flags().StringVar(&gc.customEnvPath, "custom-env", "", "custom env file path")
}

func (gc *globalCmd) run(_ *cobra.Command, _ []string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	d := daemon.New(ctx,
		daemon.WithUseConsole(gc.useConsole),
		daemon.WithLogLevel(gc.logLevel),
		daemon.WithHTTPServerAddress(gc.serverAddr),
		daemon.WithLogPath(gc.logPath),
		daemon.WithCustomEnvPath(gc.customEnvPath))
	if err := d.Init(); err != nil {
		return fmt.Errorf("init daemon failed, %w", err)
	}

	shutdown.NewListener().WaitWithCancel(cancel)
	d.Wait()

	return nil
}
