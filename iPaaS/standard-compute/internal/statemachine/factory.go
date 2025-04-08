package statemachine

import (
	"context"
	"fmt"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/tools/singularity/registry"

	"xorm.io/xorm"

	"github.com/yuansuan/ticp/iPaaS/standard-compute/config"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/backend"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/backend/job"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/dao"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/log"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/statemachine/channel"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/storage"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/pkg/with"
)

type Factory struct {
	ch      channel.Channel
	machine *StateMachine
	dao     *dao.Dao
	brk     *Breaker
	db      *xorm.Engine
}

func (f *Factory) Start(ctx context.Context) error {
	go func() {
		f.ConsumeJobs(ctx)
	}()

	<-ctx.Done()

	return nil
}

func (f *Factory) Name() string {
	return "job-consumer"
}

// RecoveryJobs 恢复未完成的所有任务
func (f *Factory) RecoveryJobs(ctx context.Context) error {
	jobs, err := f.dao.FindUnCompletedJobs(ctx)
	if err != nil {
		return err
	}

	for _, v := range jobs {
		go f.start(v.Id)
	}
	return nil
}

// ConsumeJobs 消费从 submit 提交过来的任务
func (f *Factory) ConsumeJobs(ctx context.Context) {
	for id := range f.ch.RecvMessage(ctx) {
		go f.start(id)
	}
}

func (f *Factory) ProduceJob(ctx context.Context, id int64) {
	f.ch.SendMessage(ctx, id)
}

// start 启动任务的状态机
func (f *Factory) start(jobID int64) {
	ctx := f.brk.Create(context.Background(), jobID)
	// for `with` package
	ctx = context.WithValue(ctx, with.OrmKey, f.db)

	j, err := f.dao.GetJobWithError(ctx, jobID)
	if err != nil {
		log.Warnw("unable to fetch the job", "sc_job", jobID, "error", err)
		return
	}

	jx, err := job.NewJob(j)
	if err != nil {
		log.Fatalf("create jobs failed: %v, please check db.job table %v", err, j.Id)
	}

	// 如果作业已取消则直接取消上下文对象
	if j.ControlBitTerminate {
		f.brk.Break(jobID)
	}

	go f.machine.Start(ctx, jx)
}

func (f *Factory) Cancel(jobId int64) {
	f.brk.Break(jobId)
}

func NewFactory(conf *config.Config, db *xorm.Engine, jobScheduler backend.Provider) (*Factory, error) {
	jobChannel := channel.NewChannel(conf.StateMachine)

	registryClient, err := registry.NewClient(conf.Singularity)
	if err != nil {
		return nil, fmt.Errorf("new registry client failed, %w", err)
	}

	storageManage, err := storage.NewManager(conf)
	if err != nil {
		return nil, fmt.Errorf("new storage manager, %w", err)
	}

	stateMachine := NewStateMachine(conf, jobScheduler, registryClient, storageManage, dao.Default, db)

	stateMachineBreaker, err := NewBreaker()
	if err != nil {
		return nil, fmt.Errorf("new statemachine breaker failed, %w", err)
	}

	return &Factory{
		machine: stateMachine,
		ch:      jobChannel,
		dao:     dao.Default,
		brk:     stateMachineBreaker,
		db:      db,
	}, nil
}
