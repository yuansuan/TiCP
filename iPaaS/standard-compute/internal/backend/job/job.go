package job

import (
	"encoding/json"
	"time"

	"github.com/yuansuan/ticp/common/go-kit/logging"
	v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"

	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/dao/models"
)

type State int

const (
	StatePending   = 1
	StateRunning   = 2
	StateCompleted = 3
)

type Job struct {
	*models.Job

	// extern the Job object
	EnvVars              []string
	Inputs               []*v20230530.JobInHPCInputStorage
	Output               *v20230530.JobInHPCOutputStorage
	CustomStateRule      *v20230530.JobInHPCCustomStateRule
	SchedulerSubmitFlags map[string]string

	BackendJobState  State
	PreparedFilePath string

	TraceLogger      *logging.Logger
	TransmittingTime time.Time
}

func NewJob(j *models.Job) (*Job, error) {
	newJob := &Job{
		Job: j,

		EnvVars:         []string{},
		Inputs:          []*v20230530.JobInHPCInputStorage{},
		Output:          new(v20230530.JobInHPCOutputStorage),
		CustomStateRule: new(v20230530.JobInHPCCustomStateRule),
	}

	if err := json.Unmarshal([]byte(newJob.Job.EnvVars), &newJob.EnvVars); err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(newJob.Job.Output), &newJob.Output); err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(newJob.Job.Inputs), &newJob.Inputs); err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(newJob.Job.CustomStateRule), &newJob.CustomStateRule); err != nil {
		return nil, err
	}

	if newJob.Job.SchedulerSubmitFlags != "" {
		if err := json.Unmarshal([]byte(newJob.Job.SchedulerSubmitFlags), &newJob.SchedulerSubmitFlags); err != nil {
			return nil, err
		}
	}

	return newJob, nil
}
