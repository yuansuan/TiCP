package monitorchart

import (
	"context"
	"io"
	"net/http"
	"regexp"

	"github.com/pkg/errors"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/common/project-root-api/common"
	schema "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/config"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/consts"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao/models"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/module/storage"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/service/scheduler/monitorchart/parser"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/common/with"
	"xorm.io/xorm"
)

const (
	// ChunkThreshold 分片读取阈值
	ChunkThreshold = 10 * 1024 * 1024 // 10MB // TODO: 重复代码可封装
	// ChunkSize 分片读取大小
	ChunkSize = 10 * 1024 * 1024 // 10MB // TODO: 重复代码可封装
	// DefaultMonitorChartMaxFileSize 默认监控图表解析读取文件大小上限
	DefaultMonitorChartMaxFileSize = 1024 * 1024 * 1024 // 1GB
)

var (
	ErrMonitorChartParser      = errors.New("unsupported parser")
	ErrMonitorChartJobInfo     = errors.New("job info error")
	ErrMonitorChartMarshal     = errors.New("marshal error")
	ErrMonitorChartAPIInternal = errors.New("api internal error") // TODO: 重复代码可封装
	ErrMonitorChartLs          = errors.New("ls files error")     // TODO: 重复代码可封装
	ErrMonitorChartReadAt      = errors.New("readAt error")       // TODO: 重复代码可封装
)

type MonitorChartPlugin struct{}

func NewMonitorChartPlugin() *MonitorChartPlugin {
	return &MonitorChartPlugin{}
}

func (p *MonitorChartPlugin) Name() string {
	return "monitorChart"
}

// Insert 插入作业监控图表
func (p *MonitorChartPlugin) Insert(ctx context.Context, app *models.Application, job *models.Job) {
	// func Insert(ctx context.Context, app *models.Application, job *models.Job) {
	logger := logging.GetLogger(ctx).With("func", "monitorchart.Insert", "jobID", job.ID)
	if !app.MonitorChartEnable {
		return
	}

	logger.Info("app: " + job.AppID.String() + ", monitorChart enable")
	logger.Info("insert monitorChart: start...")

	monitorChartReg := app.MonitorChartRegexp
	if monitorChartReg == "" {
		monitorChartReg = consts.DefaultMonitorChartRegexp
	}
	monitorChartParser := app.MonitorChartParser
	if monitorChartParser == "" {
		logger.Warnf("insert monitorChart: parser is empty")
		return
	}

	monitorChartModel := &models.MonitorChart{
		JobID:              job.ID,
		Content:            "",
		Finished:           false,
		MonitorChartRegexp: monitorChartReg,
		MonitorChartParser: monitorChartParser,
		FailedReason:       "",
	}

	err := with.DefaultSession(ctx, func(db *xorm.Session) error {
		_, err := db.Insert(monitorChartModel)
		return err
	})
	if err != nil {
		logger.Warnf("insert monitorChart: insert db err: %v", err)
		return
	}

}

// TODO: 重复代码可封装
func checkFinishError(err error) bool {
	switch err {
	case ErrMonitorChartAPIInternal:
		return false
	default:
		return true
	}
}

func readAndParseMonitorChart(ctx context.Context, clientParams storage.ClientParams, path, reg string, p parser.Parser) ([]*parser.Result, error) {
	logger := logging.GetLogger(ctx)
	results := []*parser.Result{}

	var re = regexp.MustCompile(reg)
	lsParams := storage.LsParams{
		Offset: 0,
		Lspath: path,
	}
	files, err := lsAll(ctx, clientParams, lsParams, re)
	if err != nil {
		return nil, err
	}

	monitorChartMaxFileSize := config.GetConfig().MonitorChartMaxFileSize
	if monitorChartMaxFileSize <= 0 {
		monitorChartMaxFileSize = DefaultMonitorChartMaxFileSize
	}

	for _, file := range files {
		pathWithName := path + file.Name
		size := file.Size

		if size > monitorChartMaxFileSize {
			logger.Errorf("monitorChart: file size is too large: %v", size)
			continue
		}

		// Read monitorchart file
		bodyReader, err := startBackgroundChunkedRead(ctx, clientParams, pathWithName, size, ChunkSize)
		if err != nil {
			logger.Warnf("monitorChart: startBackgroundChunkedRead err: %v", err)
			continue
		}

		// Parse monitorchart file
		resultMap, err := p.Parse(pathWithName, bodyReader)
		if err != nil {
			logger.Warnf("monitorChart: parse err: %v", err)
			continue
		}

		// merge resultMap to results
		results = mergeResultMap(results, resultMap)
	}

	return results, nil
}

func lsAll(ctx context.Context, clientParams storage.ClientParams, lsParams storage.LsParams, re *regexp.Regexp) ([]*schema.FileInfo, error) {
	logger := logging.GetLogger(ctx)
	regFile := make([]*schema.FileInfo, 0)
	// TODO: 重复代码可封装
	for {
		if lsParams.Offset == -1 { // -1表示已经最后一页
			break
		}

		resp, err := storage.Client().LsWithPage(clientParams, lsParams)
		if err != nil {
			logger.Warnf("ls error: %v", err)
			if resp.ErrorCode == common.InternalServerErrorCode {
				return nil, errors.WithMessage(ErrMonitorChartAPIInternal, err.Error())
			}
			return nil, errors.WithMessage(ErrMonitorChartLs, err.Error())
		}

		files := resp.Data.Files
		for _, file := range files {
			if !file.IsDir && file.Size != 0 && re.MatchString(file.Name) {
				regFile = append(regFile, file)
			}
		}
		lsParams.Offset = resp.Data.NextMarker
	}

	return regFile, nil
}

// TODO: 重复代码可封装
func startBackgroundChunkedRead(ctx context.Context, clientParams storage.ClientParams, path string, size, chunkSize int64) (io.Reader, error) {
	logger := logging.GetLogger(ctx)
	pr, pw := io.Pipe()

	go func() {
		defer pw.Close()
		// defer pr.Close()

		offset := int64(0)
		for offset < size {
			length := chunkSize
			if offset+length > size {
				length = size - offset
			}

			logger.Info("monitorChart: readAt offset: ", offset, ", length: ", length, ", path: ", path)

			readAtParams := storage.ReadAtParams{
				Readpath: path,
				Length:   length,
				Offset:   offset,
				Resolver: func(resp *http.Response) error {
					if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusPartialContent {
						pw.CloseWithError(errors.Errorf("ReadAt: status code: %v", resp.StatusCode))
						if resp.StatusCode == http.StatusInternalServerError {
							return ErrMonitorChartAPIInternal
						}
						return ErrMonitorChartReadAt
					}

					_, err := io.Copy(pw, resp.Body)
					if err != nil {
						logger.Warnf("monitorChart: chunked read: %v", err)
						pw.CloseWithError(err)
						return err
					}

					return nil
				},
			}
			_, err := storage.Client().ReadAt(clientParams, readAtParams)
			if err != nil {
				logger.Warnf("monitorChart: chunked read: %v", err)
				pw.CloseWithError(err)
				return
			}

			offset += length
		}
	}()

	return pr, nil
}
