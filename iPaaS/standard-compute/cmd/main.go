package main

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/daemon"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/shutdown"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/version"
)

type globalCmd struct {
	cmd *cobra.Command

	configFlag configFlag
	logFlag    logFlag
}

type configFlag struct {
	path string
}

type logFlag struct {
	logLevel     string
	useConsole   bool
	releaseLevel string
	path         string
	maxSize      int // MB
	maxAge       int // Day
	maxBackups   int
}

func main() {
	gc := &globalCmd{
		cmd: &cobra.Command{
			Use:   "standard-compute",
			Short: "interfaces of HPC abilities",
		},
		logFlag: logFlag{},
	}
	gc.cmd.RunE = gc.run
	initFlag(gc)

	gc.cmd.AddCommand(version.Cmd())

	if err := gc.cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initFlag(gc *globalCmd) {
	gc.cmd.Flags().StringVarP(&gc.configFlag.path, "config", "c", "", "config path")
	gc.cmd.Flags().StringVarP(&gc.logFlag.logLevel, "log-level", "l", "", "log level [info | debug]")
	gc.cmd.Flags().BoolVar(&gc.logFlag.useConsole, "use-console", false, "log use console")
	gc.cmd.Flags().StringVar(&gc.logFlag.releaseLevel, "release-level", "", "release level [development | production]")
	gc.cmd.Flags().StringVar(&gc.logFlag.path, "log-path", "", "log path")
	gc.cmd.Flags().IntVar(&gc.logFlag.maxSize, "log-max-size", 0, "log file max size [MB]")
	gc.cmd.Flags().IntVar(&gc.logFlag.maxAge, "log-max-age", 0, "log max save age [DAY]")
	gc.cmd.Flags().IntVar(&gc.logFlag.maxBackups, "log-max-backups", 0, "log max backups")
}

func (gc *globalCmd) run(_ *cobra.Command, _ []string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	d := daemon.New(ctx,
		daemon.WithConfigPath(gc.configFlag.path),
		daemon.WithLogLevel(gc.logFlag.logLevel),
		daemon.WithReleaseLevel(gc.logFlag.releaseLevel),
		daemon.WithLogPath(gc.logFlag.path),
		daemon.WithLogUseConsole(gc.logFlag.useConsole),
		daemon.WithLogMaxAge(gc.logFlag.maxAge),
		daemon.WithLogMaxBackups(gc.logFlag.maxBackups),
		daemon.WithLogMaxSize(gc.logFlag.maxSize))

	if err := d.Init(); err != nil {
		return fmt.Errorf("init daemon failed, %w", err)
	}

	shutdown.NewListener().WaitWithCancel(cancel)
	d.Wait()

	return nil
}
