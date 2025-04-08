package dao

import (
	"context"

	"github.com/pkg/errors"
	"xorm.io/xorm"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/common/project-root-api/hpc/jobstate"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/config"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/dao/models"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/pkg/with"
)

var (
	// ErrJobNotFound 数据库中找不到指定的作业
	ErrJobNotFound = errors.New("dao: job not found")
)

var Default *Dao

func init() {
	node, _ := snowflake.NewNode(0)
	Default = &Dao{idGenerator: node}
}

type Dao struct {
	idGenerator *snowflake.Node
}

func NewDao(cfg config.Snowflake) (*Dao, error) {
	node, err := snowflake.NewNode(cfg.Node)
	if err != nil {
		return nil, err
	}

	return &Dao{idGenerator: node}, nil
}

func (d *Dao) InsertJobWithGenerateID(ctx context.Context, data *models.Job) (int64, error) {
	data.Id = d.idGenerator.Generate().Int64()
	return data.Id, with.DefaultSession(ctx, func(db *xorm.Session) error {
		_, err := db.Insert(data)
		return errors.Wrap(err, "dao")
	})
}

// GetJobWithError 从数据库中获取最新的作业信息
func (*Dao) GetJobWithError(ctx context.Context, id int64) (*models.Job, error) {
	job := &models.Job{}
	err := with.DefaultSession(ctx, func(db *xorm.Session) error {
		exists, err := db.ID(id).Get(job)
		if err != nil {
			return errors.Wrap(err, "dao")
		} else if !exists {
			return ErrJobNotFound
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return job, nil
}

func (d *Dao) GetJob(ctx context.Context, id int64) (exists bool, data *models.Job, err error) {
	data = &models.Job{}
	err = with.DefaultSession(ctx, func(db *xorm.Session) error {
		exists, err = db.ID(id).Get(data)
		return errors.Wrap(err, "dao")
	})
	return
}

func (d *Dao) GetJobByIdempotentId(ctx context.Context, idempotentId string) (exists bool, data *models.Job, err error) {
	data = &models.Job{}
	err = with.DefaultSession(ctx, func(db *xorm.Session) error {
		exists, err = db.Where("idempotent_id = ?", idempotentId).Get(data)
		return errors.Wrap(err, "dao")
	})
	return
}

type GetJobsArg struct {
	IDs        []int64
	State      string
	PageOffset int
	PageSize   int
}

func (d *Dao) GetJobs(ctx context.Context, arg GetJobsArg) ([]*models.Job, error) {
	var jobs []*models.Job
	err := with.DefaultSession(ctx, func(db *xorm.Session) error {
		db = db.Table("sc_job")
		if len(arg.IDs) > 0 {
			db = db.In("id", arg.IDs)
		}

		if arg.State != "" {
			db = db.Where("state = ?", arg.State)
		}

		db.Limit(arg.PageSize, arg.PageOffset)

		e := db.Find(&jobs)
		return e
	})
	if err != nil {
		return nil, errors.Wrap(err, "dao")
	}

	return jobs, nil
}

func (d *Dao) GetAllJobsByState(ctx context.Context, state jobstate.State) ([]*models.Job, error) {
	jobs := make([]*models.Job, 0)

	err := with.DefaultSession(ctx, func(db *xorm.Session) error {
		return db.Where("state = ?", state).
			Find(&jobs)
	})
	if err != nil {
		return nil, errors.Wrap(err, "dao")
	}

	return jobs, nil
}

func (d *Dao) DeleteJob(ctx context.Context, id int64) error {
	err := with.DefaultSession(ctx, func(db *xorm.Session) error {
		_, e := db.ID(id).Delete(new(models.Job))
		return e
	})
	return errors.Wrap(err, "dao")
}

func (d *Dao) TerminateJob(ctx context.Context, id int64) error {
	data := &models.Job{ControlBitTerminate: true}
	return with.DefaultSession(ctx, func(db *xorm.Session) error {
		_, err := db.ID(id).Cols("control_bit_terminate").Update(data)
		return errors.Wrap(err, "dao")
	})
}

func (d *Dao) FailedJob(ctx context.Context, id int64, msg string) error {
	return with.DefaultSession(ctx, func(db *xorm.Session) error {
		data := &models.Job{}
		data.State = jobstate.Failed
		data.StateReason = msg
		_, err := db.ID(id).Cols("state", "state_reason").Update(data)
		return errors.Wrap(err, "dao")
	})
}

func (d *Dao) UpdateJob(ctx context.Context, data *models.Job) error {
	return with.DefaultSession(ctx, func(db *xorm.Session) error {
		_, err := db.ID(data.Id).MustCols("sub_state").Update(data)
		return errors.Wrap(err, "dao")
	})
}

type UpdateJobFileProgressArgs struct {
	DownloadCurrentSize *int64
	DownloadTotalSize   *int64
}

func (d *Dao) UpdateJobFileProgress(ctx context.Context, id int64, arg UpdateJobFileProgressArgs) error {
	return with.DefaultSession(ctx, func(db *xorm.Session) error {
		update := make(map[string]interface{})
		if arg.DownloadCurrentSize != nil {
			update["download_current_size"] = *arg.DownloadCurrentSize
		}

		if arg.DownloadTotalSize != nil {
			update["download_total_size"] = *arg.DownloadTotalSize
		}

		_, err := db.ID(id).Table(new(models.Job)).Update(update)
		return errors.Wrap(err, "dao")
	})
}

func (d *Dao) FindUnCompletedJobs(ctx context.Context) (jobs []*models.Job, err error) {
	jobs = []*models.Job{}

	err = with.DefaultSession(ctx, func(db *xorm.Session) error {
		err := db.Where("state != ? and state != ? and state != ?", jobstate.Completed, jobstate.Failed, jobstate.Canceled).Find(&jobs)
		return errors.Wrap(err, "dao")
	})
	return
}

func (*Dao) GetJobControlBitTerminateById(ctx context.Context, id int64) (exists bool, controlBitTerminate bool, err error) {
	err = with.DefaultSession(ctx, func(db *xorm.Session) error {
		data := &models.Job{}
		exists, err = db.ID(id).Cols("control_bit_terminate").Get(data)
		controlBitTerminate = data.ControlBitTerminate

		return errors.Wrap(err, "dao")
	})
	return
}
