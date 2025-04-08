package residual

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	api "github.com/yuansuan/ticp/common/project-root-api/common"
	schema "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/stat"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/config"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/consts"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao/models"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/module/storage"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/service/scheduler/residual/parser"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/util"
)

const (
	// ChunkThreshold 分片读取阈值
	ChunkThreshold = 10 * 1024 * 1024 // 10MB
	// ChunkSize 分片读取大小
	ChunkSize = 10 * 1024 * 1024 // 10MB
	// DefaultResidualMaxFileSize 默认残差解析读取文件大小上限
	DefaultResidualMaxFileSize = 1024 * 1024 * 1024 // 1GB
)

var (
	ErrResidualStat        = errors.New("stat file error")
	ErrResidualEmpty       = errors.New("file size is 0")
	ErrResidualTooLarge    = errors.New("file size is too large")
	ErrResidualAPIInternal = errors.New("api internal error")
	ErrResidualReadAt      = errors.New("readAt error")
	ErrResidualParse       = errors.New("parse error")
	ErrResidualMarshal     = errors.New("marshal error")
	ErrResidualParser      = errors.New("unsupported parser")
	ErrResidualJobInfo     = errors.New("job info error")
)

type ResidualPlugin struct {
	snowflake.IDGen
	residualDao dao.ResidualDao
}

func NewResidualPlugin(idgen snowflake.IDGen, residualDao dao.ResidualDao) *ResidualPlugin {
	return &ResidualPlugin{
		IDGen:       idgen,
		residualDao: residualDao,
	}
}

func (p *ResidualPlugin) Name() string {
	return "residual"
}

// Insert 插入作业残差图
func (p *ResidualPlugin) Insert(ctx context.Context, app *models.Application, job *models.Job) {
	logger := logging.GetLogger(ctx).With("func", "residual.Insert", "jobID", job.ID)

	// 残差图开关
	if !app.ResidualEnable {
		return
	}

	logger.Info("app: " + job.AppID.String() + ", residual enable")
	logger.Info("insert residual: start...")

	// 残差解析文件 默认值
	residualReg := app.ResidualLogRegexp
	if residualReg == "" {
		residualReg = consts.DefaultResidualReg
	}
	residualParser := app.ResidualLogParser
	if residualParser == "" {
		logger.Warnf("insert residual: parser is empty")
		return
	}

	id, err := p.GenID(ctx)
	if err != nil {
		logger.Errorf("generate a snowflake id fail: %v", err)
		return
	}

	now := time.Now()
	residualModel := &models.Residual{
		ID:                id,
		JobID:             job.ID,
		Content:           "",
		Finished:          false,
		ResidualLogRegexp: residualReg,
		ResidualLogParser: residualParser,
		CreateTime:        now,
		UpdateTime:        now,
	}

	err = p.residualDao.InsertResidual(ctx, residualModel)
	if err != nil {
		logger.Warnf("insert residual: insert db err: %v", err)
		return
	}
}

func checkFinishError(err error) bool {
	switch err {
	case ErrResidualAPIInternal:
		return false
	default:
		return true
	}
}

const maxRetries = 3 // todo: 配置

func statWithRetry(ctx context.Context, clientParams storage.ClientParams, path string) (*stat.Response, error) {
	logger := logging.GetLogger(ctx)
	var file *stat.Response
	var err error
	for i := 0; i < maxRetries; i++ {
		if i > 0 {
			time.Sleep(1 * time.Second)
		}

		file, err = storage.Client().Stat(clientParams, path)
		if err != nil {
			logger.Infof("residual: stat attempt %d failed: %v", i+1, err)
			continue
		}

		if file.Data.File.Size == 0 {
			logger.Infof("residual: stat attempt %d file size is 0", i+1)
			continue
		}

		// success
		break
	}

	if err != nil {
		if file != nil && file.ErrorCode == api.InternalServerErrorCode {
			return nil, errors.WithMessage(ErrResidualAPIInternal, err.Error())
		}

		if file != nil && file.ErrorCode == api.PathNotFound {
			return nil, common.ErrPathNotFound
		}

		return nil, errors.WithMessage(ErrResidualStat, err.Error())
	}

	if file.Data.File.Size == 0 {
		logger.Infof("residual: file size is 0")
		return nil, ErrResidualEmpty
	}

	return file, nil
}

func readAndParseResidual(ctx context.Context, clientParams storage.ClientParams, path string, p parser.Parser) (*parser.Result, error) {
	logger := logging.GetLogger(ctx)
	file, err := statWithRetry(ctx, clientParams, path)
	if err != nil {
		return nil, err
	}

	size := file.Data.File.Size
	logger.Infof("residual: file size: %v", size)

	residualMaxFileSize := config.GetConfig().ResidualMaxFileSize
	if residualMaxFileSize <= 0 {
		residualMaxFileSize = DefaultResidualMaxFileSize
	}

	var bodyReader io.Reader
	if size > residualMaxFileSize {
		logger.Errorf("residual: file size is too large: %v", size)
		return nil, errors.WithMessage(ErrResidualTooLarge, fmt.Sprintf("file size: %v", size))
	}

	bodyReader, err = startBackgroundChunkedRead(ctx, clientParams, path, size, ChunkSize)
	if err != nil {
		logger.Warnf("residual: chunked read: %v", err)
		return nil, err
	}

	// Parse residual file
	result, err := p.Parse(bodyReader)
	if err != nil {
		logger.Warnf("parse err: %v", err)
		return nil, errors.WithMessage(ErrResidualParse, err.Error())
	}

	return result, nil
}

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

			logger.Info("residual: readAt offset: ", offset, ", length: ", length, ", path: ", path)

			readAtParams := storage.ReadAtParams{
				Readpath: path,
				Length:   length,
				Offset:   offset,
				Resolver: func(resp *http.Response) error {
					if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusPartialContent {
						pw.CloseWithError(errors.Errorf("ReadAt: status code: %v", resp.StatusCode))
						if resp.StatusCode == http.StatusInternalServerError {
							return ErrResidualAPIInternal
						}
						return ErrResidualReadAt
					}

					_, err := io.Copy(pw, resp.Body)
					if err != nil {
						logger.Warnf("residual: chunked read: %v", err)
						pw.CloseWithError(err)
						return err
					}

					return nil
				},
			}
			_, err := storage.Client().ReadAt(clientParams, readAtParams)
			if err != nil {
				logger.Warnf("residual: chunked read: %v", err)
				pw.CloseWithError(err)
				return
			}

			offset += length
		}
	}()

	return pr, nil
}

type ResidualHandler interface {
	HandlerHpcResidual(ctx context.Context, job *models.Job, residual *models.Residual, zones schema.Zones) (*schema.Residual, error)
}

type Handler struct {
	Zones schema.Zones
}

func NewHandler(zones schema.Zones) *Handler {
	return &Handler{
		Zones: zones,
	}
}

func (h *Handler) HandlerHpcResidual(ctx context.Context, job *models.Job, residual *models.Residual, zones schema.Zones) (*schema.Residual, error) {
	logger := logging.GetLogger(ctx).With("func", "HandlerHpcResidual", "jobID", job.ID)
	// 残差解析文件 默认值
	residualReg := residual.ResidualLogRegexp
	if residualReg == "" {
		residualReg = consts.DefaultResidualReg
	}
	residualParser := residual.ResidualLogParser
	if residualParser == "" {
		logger.Warnf("residual: parser is empty")
		return nil, errors.WithMessage(ErrResidualParser, "parser is empty")
	}

	// residual analyse
	p := NewParser(residualParser)
	if p == nil {
		logger.Warnf("residual: unsupported parser: %v", residualParser)
		return nil, errors.WithMessage(ErrResidualParser, fmt.Sprintf("unsupported parser: %v", residualParser))
	}

	zone, ok := zones[job.Zone]
	if !ok {
		logger.Warnf("residual: zone not found: %v", job.Zone)
		return nil, errors.WithMessage(ErrResidualJobInfo, fmt.Sprintf("zone not found: %v", job.Zone))
	}

	if zone.HPCEndpoint == "" {
		logger.Warnf("residual: zone domain is empty")
		return nil, errors.WithMessage(ErrResidualJobInfo, "zone domain is empty")
	}

	clientParams := storage.ClientParams{
		Endpoint: zone.HPCEndpoint,
		Timeout:  0, // 0 for no timeout
		AdminAPI: true,
	}

	workdir := job.WorkDir
	if workdir == "" {
		logger.Warnf("residual: workdir is empty")
		return nil, errors.WithMessage(ErrResidualJobInfo, "workdir is empty")
	}

	workdir = strings.TrimPrefix(workdir, zone.HPCEndpoint)
	workdir = util.AddPrefixSlash(workdir)
	workdir = util.AddSuffixSlash(workdir)

	// storage readAt
	logger.Info("residual: readAt path: ", workdir+residualReg)
	logger.Info("residual: readAt client base url: ", clientParams.Endpoint)
	ctx = logging.AppendWith(ctx, "func", "readAndParseResidual", "endpoint", clientParams.Endpoint, "path", workdir+residualReg, "parser", residualParser)
	result, err := readAndParseResidual(ctx, clientParams, workdir+residualReg, p)
	if err != nil {
		return nil, err
	}

	// store result
	residualData := ConvertResidual(result)
	return residualData, nil
}
