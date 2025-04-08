package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/shirou/gopsutil/net"
	"github.com/spf13/cobra"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"go.uber.org/zap"
)

type Command struct {
	*cobra.Command
	VideoPath     string
	NetName       string
	Duration      time.Duration
	CheckInterval time.Duration
}

const (
	defaultDuration      = 60 * time.Second
	defaultCheckInterval = 1 * time.Second
)

func main() {
	gc := Command{
		Command: &cobra.Command{
			Use:   "benchmark-bandwidth",
			Short: "bandwidth benchmark tools",
			Long:  "bandwidth benchmark tools",
		},
	}

	gc.Command.RunE = gc.runE

	gc.Flags().StringVar(&gc.VideoPath, "video-path", "benchmark.mp4", "video path for benchmark")
	gc.Flags().StringVar(&gc.NetName, "net-name", "", "catch bandwidth net name")
	gc.Flags().DurationVar(&gc.Duration, "duration", defaultDuration, "net traffic collect duration e.g. [60s | 10m]")
	gc.Flags().DurationVar(&gc.CheckInterval, "check-interval", defaultCheckInterval, "net traffic check interval e.g. [1s | 1m]")

	if err := gc.Execute(); err != nil {
		logging.Default().Error(err)
		os.Exit(1)
	}
}

func (c *Command) runE(cmd *cobra.Command, args []string) error {
	logging.Default().Infof("benchmark video path [%s]", c.VideoPath)
	logging.Default().Infof("net name [%s]. It would listen first net iface if name is empty", c.NetName)
	logging.Default().Infof("duration [%s]", c.Duration.String())

	if _, err := os.Stat(c.VideoPath); os.IsNotExist(err) {
		return fmt.Errorf("benchmark video file not exist, %w", err)
	}

	ec := exec.Command("cmd", "/c", "start", c.VideoPath)
	err := ec.Run()
	if err != nil {
		return fmt.Errorf("cannot open benchmarch video file, %w", err)
	}

	startInfo, err := c.getNetStat()
	if err != nil {
		return fmt.Errorf("get net stat failed, %w", err)
	}

	logging.Default().Infof("start info: %s", startInfo.String())
	bytesStartSend := startInfo.BytesSent
	bytesStartRecv := startInfo.BytesRecv

	var bytesSendDiffPeek, bytesRecvDiffPeek uint64
	currentBytesSend, currentBytesRecv := bytesStartSend, bytesStartRecv
	// 另一个线程每隔一段时间计算一下上下行流量，统计出监测时间段内的最大上下行流量
	go func() {
		for {
			time.Sleep(c.CheckInterval)

			// 获得当前的流量统计
			currentNetStat, err := c.getNetStat()
			if err != nil {
				logging.Default().Errorf("get net stat failed, %v", err)
				continue
			}

			// 与上一次相减，得差值
			currentBytesSendDiff := currentNetStat.BytesSent - currentBytesSend
			currentBytesRecvDiff := currentNetStat.BytesRecv - currentBytesRecv
			// 更新current
			currentBytesSend = currentNetStat.BytesSent
			currentBytesRecv = currentNetStat.BytesRecv

			// 与peek值对比
			if currentBytesSendDiff > bytesSendDiffPeek {
				bytesSendDiffPeek = currentBytesSendDiff
			}
			if currentBytesRecvDiff > bytesRecvDiffPeek {
				bytesRecvDiffPeek = currentBytesRecvDiff
			}
		}
	}()

	time.Sleep(c.Duration)

	endInfo, err := c.getNetStat()
	if err != nil {
		return fmt.Errorf("get net stat failed, %w", err)
	}

	logging.Default().Infof("end info: %s", endInfo.String())
	bytesEndSend := endInfo.BytesSent
	bytesEndRecv := endInfo.BytesRecv

	bytesSendDiff := bytesEndSend - bytesStartSend
	bytesRecvDiff := bytesEndRecv - bytesStartRecv

	logging.Default().Infof("bytesSendDiff: %f MB", float64(bytesSendDiff)/1024/1024)
	logging.Default().Infof("bytesRecvDiff: %f MB", float64(bytesRecvDiff)/1024/1024)

	logging.Default().Infof("================================")
	res := Result{
		UploadAverageBandwidth:   fmt.Sprintf("%f MB/s", float64(bytesSendDiff)/1024/1024/c.Duration.Seconds()),
		UploadPeekBandwidth:      fmt.Sprintf("%f MB/s", float64(bytesSendDiffPeek)/1024/1024/c.CheckInterval.Seconds()),
		DownloadAverageBandwidth: fmt.Sprintf("%f MB/s", float64(bytesRecvDiff)/1024/1024/c.Duration.Seconds()),
		DownloadPeekBandwidth:    fmt.Sprintf("%f MB/s", float64(bytesRecvDiffPeek)/1024/1024/c.CheckInterval.Seconds()),
	}
	logging.Default().Infof("upload average bandwidth: %s", res.UploadAverageBandwidth)
	logging.Default().Infof("upload peek bandwidth: %s", res.UploadPeekBandwidth)

	logging.Default().Infof("download average bandwidth: %s", res.DownloadAverageBandwidth)
	logging.Default().Infof("download peek bandwidth: %s", res.DownloadPeekBandwidth)

	content, err := json.MarshalIndent(&res, "", " ")
	if err != nil {
		return err
	}

	if err = os.WriteFile(fmt.Sprintf("benchmark-bandwidth-result-%s.json", time.Now().Format("20060102-150405")), content, 0644); err != nil {
		return err
	}

	return nil
}

func (c *Command) getNetStat() (*net.IOCountersStat, error) {
	netStats, err := net.IOCounters(true)
	if err != nil {
		return nil, logErrAndReturn(err)
	}

	if len(netStats) == 0 {
		return nil, logErrAndReturn(fmt.Errorf("len of netStats is 0"))
	}

	if c.NetName == "" {
		c.NetName = netStats[0].Name
		logging.Default().Infof("net stat: %s", netStats[0].String())
		return &netStats[0], nil
	}

	for i := range netStats {
		if netStats[i].Name == c.NetName {
			logging.Default().Infof("net stat: %s", netStats[0].String())
			return &netStats[i], nil
		}
	}

	return nil, logErrAndReturn(fmt.Errorf("not found net where name = [%s]", c.NetName))
}

type Result struct {
	UploadAverageBandwidth string `json:"upload_average_bandwidth"`
	UploadPeekBandwidth    string `json:"upload_peek_bandwidth"`

	DownloadAverageBandwidth string `json:"download_average_bandwidth"`
	DownloadPeekBandwidth    string `json:"download_peek_bandwidth"`
}

func logErrAndReturn(err error) error {
	if err == nil {
		return nil
	}

	logging.Default().With(zap.AddCallerSkip(-1)).Error(err)
	return err
}
