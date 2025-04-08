package util

import (
	"context"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	hpc "github.com/yuansuan/ticp/common/project-root-api/hpc/v1/job"
	adminjobcreate "github.com/yuansuan/ticp/common/project-root-api/job/v1/admin/jobcreate"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/jobcreate"
	v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/consts"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao/models"
	"go.uber.org/zap"
)

func TestModelToOpenApiJob(t *testing.T) {
	type args struct {
		job *models.Job
	}
	now := time.Now()
	tests := []struct {
		name string
		args args
		want *v20230530.JobInfo
	}{
		{
			name: "normal job",
			args: args{
				job: &models.Job{
					ID:                          0,
					Name:                        "TestJob",
					Comment:                     "Test comment",
					UserID:                      123,
					JobSource:                   "Source",
					State:                       1,
					SubState:                    0,
					StateReason:                 "Reason",
					ExitCode:                    "ExitCode",
					FileSyncState:               "Completed",
					Params:                      "Params",
					UserZone:                    "Zone",
					Timeout:                     60,
					FileClassifier:              "Classifier",
					ResourceUsageCpus:           2,
					ResourceUsageMemory:         2048,
					CustomStateRuleKeyStatement: "RuleKey",
					CustomStateRuleResultState:  "ResultState",
					HPCJobID:                    "HPCJobID",
					Zone:                        "JobZone",
					ResourceAssignCpus:          2,
					ResourceAssignMemory:        2048,
					Command:                     "Command",
					WorkDir:                     "/path/to/workdir",
					OriginJobID:                 "OriginJobID",
					Queue:                       "Queue",
					Priority:                    1,
					ExecHosts:                   "ExecHosts",
					ExecHostNum:                 1,
					ExecutionDuration:           120,
					InputType:                   "InputType",
					InputDir:                    "/path/to/input",
					Destination:                 "Destination",
					OutputType:                  "OutputType",
					OutputDir:                   "/path/to/output",
					NoNeededPaths:               "NoNeededPaths",
					NeededPaths:                 "NeededPaths",
					FileInputStorageZone:        "InputStorageZone",
					FileOutputStorageZone:       "OutputStorageZone",
					DownloadFileSizeTotal:       1024,
					DownloadFileSizeCurrent:     512,
					UploadFileSizeTotal:         2048,
					UploadFileSizeCurrent:       1024,
					AppID:                       456,
					AppName:                     "TestApp",
					UserCancel:                  0,
					IsFileReady:                 0,
					DownloadFinished:            0,
					IsDeleted:                   0,
					UploadTime:                  now,
					PendingTime:                 now,
					DownloadTime:                now,
					RunningTime:                 now,
					TerminatingTime:             now,
					TransmittingTime:            now,
					SuspendingTime:              now,
					SuspendedTime:               now,
					SubmitTime:                  now,
					EndTime:                     now,
					CreateTime:                  now,
					UpdateTime:                  now,
				},
			},
			want: &v20230530.JobInfo{
				ID:            "",
				Name:          "TestJob",
				JobState:      "Initiated",
				StateReason:   "Reason",
				FileSyncState: "Completed",
				AllocResource: &v20230530.AllocResource{
					Cores:  2,
					Memory: 2048,
				},
				Zone:             "JobZone",
				Workdir:          "/path/to/workdir",
				Parameters:       "Params",
				PendingTime:      ModelTimeToString(now),
				RunningTime:      ModelTimeToString(now),
				TerminatingTime:  ModelTimeToString(now),
				TransmittingTime: ModelTimeToString(now),
				SuspendingTime:   ModelTimeToString(now),
				SuspendedTime:    ModelTimeToString(now),
				EndTime:          ModelTimeToString(now),
				CreateTime:       ModelTimeToString(now),
				UpdateTime:       ModelTimeToString(now),
				DownloadProgress: &v20230530.DownloadProgress{
					Progress: &v20230530.Progress{
						TotalSize: 1024,
						Progress:  50,
					},
				},
				UploadProgress: &v20230530.UploadProgress{
					Progress: &v20230530.Progress{
						TotalSize: 2048,
						Progress:  50,
					},
				},
				ExecutionDuration: 120,
				ExitCode:          "ExitCode",
				StdoutPath:        "/path/to/output",
				StderrPath:        "",
			},
		},
		{
			name: "job nil",
			args: args{
				job: nil,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := ModelToOpenAPIJob(tt.args.job)
			assert.Equal(t, tt.want, resp)
		})
	}
}

func TestModelToAdminOpenApiJob(t *testing.T) {
	type args struct {
		job *models.Job
	}
	now := time.Now()
	tests := []struct {
		name string
		args args
		want *v20230530.AdminJobInfo
	}{
		{
			name: "admin job",
			args: args{
				job: &models.Job{
					ID:                          0,
					Name:                        "TestJob",
					Comment:                     "Test comment",
					UserID:                      123,
					JobSource:                   "Source",
					State:                       1,
					SubState:                    0,
					StateReason:                 "Reason",
					ExitCode:                    "ExitCode",
					Params:                      "Params",
					UserZone:                    "Zone",
					Timeout:                     60,
					FileClassifier:              "Classifier",
					ResourceUsageCpus:           2,
					ResourceUsageMemory:         2048,
					CustomStateRuleKeyStatement: "RuleKey",
					CustomStateRuleResultState:  "ResultState",
					HPCJobID:                    "HPCJobID",
					Zone:                        "JobZone",
					ResourceAssignCpus:          2,
					ResourceAssignMemory:        2048,
					Command:                     "Command",
					WorkDir:                     "/path/to/workdir",
					OriginJobID:                 "OriginJobID",
					Queue:                       "Queue",
					Priority:                    1,
					ExecHosts:                   "ExecHosts",
					ExecHostNum:                 1,
					ExecutionDuration:           120,
					InputType:                   "InputType",
					InputDir:                    "/path/to/input",
					Destination:                 "Destination",
					OutputType:                  "OutputType",
					OutputDir:                   "/path/to/output",
					NoNeededPaths:               "NoNeededPaths",
					NeededPaths:                 "NeededPaths",
					FileInputStorageZone:        "InputStorageZone",
					FileOutputStorageZone:       "OutputStorageZone",
					DownloadFileSizeTotal:       1024,
					DownloadFileSizeCurrent:     512,
					UploadFileSizeTotal:         2048,
					UploadFileSizeCurrent:       1024,
					AppID:                       456,
					AppName:                     "TestApp",
					UserCancel:                  0,
					IsFileReady:                 0,
					DownloadFinished:            0,
					IsDeleted:                   0,
					UploadTime:                  now,
					PendingTime:                 now,
					DownloadTime:                now,
					RunningTime:                 now,
					TerminatingTime:             now,
					TransmittingTime:            now,
					SuspendingTime:              now,
					SuspendedTime:               now,
					SubmitTime:                  now,
					EndTime:                     now,
					CreateTime:                  now,
					UpdateTime:                  now,
				},
			},
			want: &v20230530.AdminJobInfo{
				JobInfo: v20230530.JobInfo{
					ID:          "",
					Name:        "TestJob",
					JobState:    "Initiated",
					StateReason: "Reason",
					AllocResource: &v20230530.AllocResource{
						Cores:  2,
						Memory: 2048,
					},
					ExecHostNum:      1,
					Zone:             "JobZone",
					Workdir:          "/path/to/workdir",
					Parameters:       "Params",
					PendingTime:      ModelTimeToString(now),
					RunningTime:      ModelTimeToString(now),
					TerminatingTime:  ModelTimeToString(now),
					TransmittingTime: ModelTimeToString(now),
					SuspendingTime:   ModelTimeToString(now),
					SuspendedTime:    ModelTimeToString(now),
					EndTime:          ModelTimeToString(now),
					CreateTime:       ModelTimeToString(now),
					UpdateTime:       ModelTimeToString(now),
					DownloadProgress: &v20230530.DownloadProgress{
						Progress: &v20230530.Progress{
							TotalSize: 1024,
							Progress:  50,
						},
					},
					UploadProgress: &v20230530.UploadProgress{
						Progress: &v20230530.Progress{
							TotalSize: 2048,
							Progress:  50,
						},
					},
					ExecutionDuration: 120,
					ExitCode:          "ExitCode",
					StdoutPath:        "/path/to/output",
					StderrPath:        "",
				},
				Queue:       "Queue",
				Priority:    1,
				OriginJobID: "OriginJobID",
				ExecHosts:   "ExecHosts",
				SubmitTime:  ModelTimeToString(now),
				UserID:      "38",
				HPCJobID:    "HPCJobID",
			},
		},
		{
			name: "admin job nil",
			args: args{
				job: nil,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := ModelToAdminOpenAPIJob(tt.args.job)
			assert.Equal(t, tt.want, resp)
		})
	}
}

func TestModelTimeToString(t *testing.T) {
	type args struct {
		timeInput time.Time
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "time to string",
			args: args{
				timeInput: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			want: "2021-01-01T08:00:00+08:00",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timeString := ModelTimeToString(tt.args.timeInput)
			assert.Equal(t, tt.want, timeString)
		})
	}
}

func Test_calPercent(t *testing.T) {
	type args struct {
		numerator   int64
		denominator int64
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "calculate percent",
			args: args{
				numerator:   50,
				denominator: 200,
			},
			want: 25,
		},
		{
			name: "calculate percent, denominator is zero",
			args: args{
				numerator:   50,
				denominator: 0,
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := calPercent(tt.args.numerator, tt.args.denominator)
			assert.Equal(t, tt.want, resp)
		})
	}
}

func TestAssembleHPCJobRequest(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	logger, err := logging.NewLogger()
	if err != nil {
		t.Error(err)
	}

	ctx.Set(logging.LoggerName, logger)

	type args struct {
		ctx        context.Context
		logger     *zap.SugaredLogger
		job        *models.Job
		appImage   string
		envVars    map[string]string
		noTransfer bool
		localImage bool
	}
	now := time.Now()
	job := &models.Job{
		ID:                          10086, // 对应:3ZU
		Name:                        "TestJob",
		Comment:                     "Test job comment",
		UserID:                      123,
		JobSource:                   "TestSource",
		State:                       1,
		SubState:                    2,
		StateReason:                 "Test reason",
		ExitCode:                    "ExitCode",
		Params:                      "Test params",
		UserZone:                    "TestUserZone",
		Timeout:                     60,
		FileClassifier:              "TestClassifier",
		ResourceUsageCpus:           4,
		ResourceUsageMemory:         8192,
		CustomStateRuleKeyStatement: "TestKeyStatement",
		CustomStateRuleResultState:  "TestResultState",
		HPCJobID:                    "HPCJobID",
		Zone:                        "TestZone",
		ResourceAssignCpus:          8,
		ResourceAssignMemory:        16384,
		Command:                     "TestCommand",
		WorkDir:                     "/path/to/workdir",
		OriginJobID:                 "OriginJobID",
		Queue:                       "TestQueue",
		Priority:                    1,
		ExecHosts:                   "TestExecHosts",
		ExecHostNum:                 2,
		ExecutionDuration:           3600,
		InputType:                   "TestInputType",
		InputDir:                    "/path/to/input",
		Destination:                 "/path/to/destination",
		OutputType:                  "TestOutputType",
		OutputDir:                   "/path/to/output",
		NoNeededPaths:               "TestNoNeededPaths",
		NeededPaths:                 "TestNeededPaths",
		FileInputStorageZone:        "TestInputStorageZone",
		FileOutputStorageZone:       "TestOutputStorageZone",
		DownloadFileSizeTotal:       1024,
		DownloadFileSizeCurrent:     512,
		UploadFileSizeTotal:         2048,
		UploadFileSizeCurrent:       1024,
		AppID:                       456,
		AppName:                     "TestApp",
		UserCancel:                  0,
		IsFileReady:                 0,
		DownloadFinished:            0,
		IsDeleted:                   0,
		UploadTime:                  now,
		DownloadTime:                now,
		PendingTime:                 now,
		RunningTime:                 now,
		TerminatingTime:             now,
		TransmittingTime:            now,
		SuspendingTime:              now,
		SuspendedTime:               now,
		SubmitTime:                  now,
		EndTime:                     now,
		CreateTime:                  now,
		UpdateTime:                  now,
	}
	coresPerNode := 4
	tests := []struct {
		name string
		args args
		want hpc.SystemPostRequest
	}{
		{
			name: "make hpc job request",
			args: args{
				ctx:        ctx,
				logger:     logger,
				job:        job,
				appImage:   "test/image",
				envVars:    map[string]string{"key1": "value1", "key2": "value2"},
				noTransfer: false,
				localImage: false,
			},
			want: hpc.SystemPostRequest{
				IdempotentID: "3ZU",
				Application:  "image:test/image",
				Environment:  map[string]string{"key1": "value1", "key2": "value2"},
				Command:      "TestCommand",
				Override: v20230530.JobInHPCOverride{
					Enable:  true,
					WorkDir: "/path/to/workdir",
				},
				Queue: "TestQueue",
				Resource: v20230530.JobInHPCResource{
					Cores:        4,
					CoresPerNode: &coresPerNode,
				},
				Inputs: []v20230530.JobInHPCInputStorage{{
					Src:  "/path/to/input",
					Dst:  "",
					Type: "TestInputType",
				}},
				Output: &v20230530.JobInHPCOutputStorage{
					Dst:           "/path/to/output",
					Type:          "TestOutputType",
					NoNeededPaths: "TestNoNeededPaths",
					NeededPaths:   "TestNeededPaths",
				},
				CustomStateRule: &v20230530.JobInHPCCustomStateRule{
					KeyStatement: "TestKeyStatement",
					ResultState:  "TestResultState",
				},
			},
		},
		{
			name: "make hpc job request,noTransfer is true",
			args: args{
				ctx:        ctx,
				logger:     logger,
				job:        job,
				appImage:   "test/image",
				envVars:    map[string]string{"key1": "value1", "key2": "value2"},
				noTransfer: true,
				localImage: false,
			},
			want: hpc.SystemPostRequest{
				IdempotentID: "3ZU",
				Application:  "image:test/image",
				Environment:  map[string]string{"key1": "value1", "key2": "value2"},
				Command:      "TestCommand",
				Override: v20230530.JobInHPCOverride{
					Enable:  true,
					WorkDir: "/path/to/workdir",
				},
				Queue: "TestQueue",
				Resource: v20230530.JobInHPCResource{
					Cores:        4,
					CoresPerNode: &coresPerNode,
				},
				Inputs: nil,
				Output: nil,
				CustomStateRule: &v20230530.JobInHPCCustomStateRule{
					KeyStatement: "TestKeyStatement",
					ResultState:  "TestResultState",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := AssembleHPCJobRequest(tt.args.ctx, tt.args.logger, tt.args.job, tt.args.appImage, tt.args.envVars, nil, tt.args.noTransfer, tt.args.localImage, int(tt.args.job.ResourceUsageCpus))
			assert.Equal(t, tt.want, resp)
		})
	}
}

func TestConvertJobModel(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	logger, err := logging.NewLogger()
	if err != nil {
		t.Error(err)
	}

	ctx.Set(logging.LoggerName, logger)

	now := time.Now()

	type args struct {
		ctx        context.Context
		logger     *zap.SugaredLogger
		req        *jobcreate.Request
		userID     snowflake.ID
		jobID      snowflake.ID
		appInfo    *models.Application
		inputZone  consts.Zone
		outputZone consts.Zone
	}
	core := 2
	memory := 256
	tests := []struct {
		name    string
		args    args
		want    *models.Job
		wantErr bool
	}{
		{
			name: "normal",
			args: args{
				ctx:    ctx,
				logger: logger,
				req: &jobcreate.Request{
					Name: "TestJob",
					Params: jobcreate.Params{
						Application: jobcreate.Application{
							Command: "Command",
							AppID:   "456",
						},
						Resource: &jobcreate.Resource{
							Cores:  &core,
							Memory: &memory,
						},
						EnvVars: map[string]string{},
						Input: &jobcreate.Input{
							Type:        "CloudStorage",
							Source:      "https://test/10086/path/to/input",
							Destination: "Destination",
						},
						Output: &jobcreate.Output{
							Type:          "CloudStorage",
							Address:       "https://test/10086/path/to/output",
							NoNeededPaths: "NoNeededPaths",
							NeededPaths:   "NeededPaths",
						},
						TmpWorkdir:        false,
						SubmitWithSuspend: false,
						CustomStateRule: &jobcreate.CustomStateRule{
							KeyStatement: "RuleKey",
							ResultState:  "ResultState",
						},
					},
					Timeout: 60,
					Zone:    "JobZone",
					Comment: "Test comment",
				},
				userID: 123,
				jobID:  10086, // 对应:3ZU
				appInfo: &models.Application{
					ID:                snowflake.ID(123456),
					Name:              "TestApp",
					Type:              "",
					Version:           "",
					AppParamsVersion:  0,
					Image:             "",
					Endpoint:          "",
					Command:           "",
					PublishStatus:     "",
					Description:       "",
					IconUrl:           "",
					CoresMaxLimit:     0,
					CoresPlaceholder:  "",
					FileFilterRule:    "",
					ResidualEnable:    false,
					ResidualLogRegexp: "",
					ResidualLogParser: "",
					LicManagerId:      0,
					SnapshotEnable:    false,
					BinPath:           "",
				},
				inputZone:  "InputStorageZone",
				outputZone: "OutputStorageZone",
			},
			want: &models.Job{
				ID:                          10086,
				Name:                        "TestJob",
				Comment:                     "Test comment",
				UserID:                      123,
				JobSource:                   "",
				State:                       2,
				SubState:                    200,
				StateReason:                 "User Submit",
				ExitCode:                    "",
				Params:                      "Params",
				UserZone:                    "JobZone",
				Timeout:                     60,
				FileClassifier:              "",
				ResourceUsageCpus:           2,
				ResourceUsageMemory:         256,
				CustomStateRuleKeyStatement: "RuleKey",
				CustomStateRuleResultState:  "ResultState",
				HPCJobID:                    "",
				Zone:                        "JobZone",
				ResourceAssignCpus:          2,
				ResourceAssignMemory:        256,
				Command:                     "Command",
				WorkDir:                     "Destination/",
				OriginJobID:                 "",
				Queue:                       "",
				Priority:                    0,
				ExecHosts:                   "",
				ExecHostNum:                 0,
				ExecutionDuration:           0,
				InputType:                   "CloudStorage",
				InputDir:                    "https://test/10086/path/to/input",
				Destination:                 "Destination",
				OutputType:                  "CloudStorage",
				OutputDir:                   "https://test/10086/path/to/output/",
				NoNeededPaths:               "NoNeededPaths",
				NeededPaths:                 "NeededPaths",
				FileInputStorageZone:        "InputStorageZone",
				FileOutputStorageZone:       "OutputStorageZone",
				DownloadFileSizeTotal:       0,
				DownloadFileSizeCurrent:     0,
				UploadFileSizeTotal:         0,
				UploadFileSizeCurrent:       0,
				AppID:                       123456,
				AppName:                     "TestApp",
				UserCancel:                  0,
				IsFileReady:                 0,
				DownloadFinished:            0,
				IsDeleted:                   0,
				UploadTime:                  InvalidTime,
				PendingTime:                 InvalidTime,
				DownloadTime:                InvalidTime,
				RunningTime:                 InvalidTime,
				TerminatingTime:             InvalidTime,
				TransmittingTime:            InvalidTime,
				SuspendingTime:              InvalidTime,
				SuspendedTime:               InvalidTime,
				SubmitTime:                  InvalidTime,
				EndTime:                     InvalidTime,
				CreateTime:                  now,
				UpdateTime:                  now,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := ConvertJobModel(tt.args.ctx, tt.args.logger, tt.args.req, tt.args.userID, tt.args.jobID, tt.args.appInfo, tt.args.inputZone, tt.args.outputZone, nil, nil)

			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertJobModel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				t.Log(err)
				return
			}

			assert.Equal(t, tt.want.ID, resp.ID)
			assert.Equal(t, tt.want.Name, resp.Name)
			assert.Equal(t, tt.want.Comment, resp.Comment)
			assert.Equal(t, tt.want.UserID, resp.UserID)
			assert.Equal(t, tt.want.JobSource, resp.JobSource)
			assert.Equal(t, tt.want.State, resp.State)
			assert.Equal(t, tt.want.SubState, resp.SubState)
			assert.Contains(t, resp.StateReason, tt.want.StateReason)
			assert.Equal(t, tt.want.ExitCode, resp.ExitCode)
			assert.NotEmpty(t, resp.Params)
			assert.Equal(t, tt.want.UserZone, resp.UserZone)
			assert.Equal(t, tt.want.Timeout, resp.Timeout)
			assert.Equal(t, tt.want.FileClassifier, resp.FileClassifier)
			assert.Equal(t, tt.want.ResourceUsageCpus, resp.ResourceUsageCpus)
			assert.Equal(t, tt.want.ResourceUsageMemory, resp.ResourceUsageMemory)
			assert.Equal(t, tt.want.CustomStateRuleKeyStatement, resp.CustomStateRuleKeyStatement)
			assert.Equal(t, tt.want.CustomStateRuleResultState, resp.CustomStateRuleResultState)
			assert.Equal(t, tt.want.HPCJobID, resp.HPCJobID)
			assert.Equal(t, tt.want.Zone, resp.Zone)
			assert.Equal(t, tt.want.ResourceAssignCpus, resp.ResourceAssignCpus)
			assert.Equal(t, tt.want.ResourceAssignMemory, resp.ResourceAssignMemory)
			assert.Equal(t, tt.want.Command, resp.Command)
			assert.Equal(t, tt.want.WorkDir, resp.WorkDir)
			assert.Equal(t, tt.want.OriginJobID, resp.OriginJobID)
			assert.Equal(t, tt.want.Queue, resp.Queue)
			assert.Equal(t, tt.want.Priority, resp.Priority)
			assert.Equal(t, tt.want.ExecHosts, resp.ExecHosts)
			assert.Equal(t, tt.want.ExecHostNum, resp.ExecHostNum)
			assert.Equal(t, tt.want.ExecutionDuration, resp.ExecutionDuration)
			assert.Equal(t, tt.want.InputType, resp.InputType)
			assert.Equal(t, tt.want.InputDir, resp.InputDir)
			assert.Equal(t, tt.want.Destination, resp.Destination)
			assert.Equal(t, tt.want.OutputType, resp.OutputType)
			assert.Equal(t, tt.want.OutputDir, resp.OutputDir)
			assert.Equal(t, tt.want.NoNeededPaths, resp.NoNeededPaths)
			assert.Equal(t, tt.want.NeededPaths, resp.NeededPaths)
			assert.Equal(t, tt.want.FileInputStorageZone, resp.FileInputStorageZone)
			assert.Equal(t, tt.want.FileOutputStorageZone, resp.FileOutputStorageZone)
			assert.Equal(t, tt.want.DownloadFileSizeTotal, resp.DownloadFileSizeTotal)
			assert.Equal(t, tt.want.DownloadFileSizeCurrent, resp.DownloadFileSizeCurrent)
			assert.Equal(t, tt.want.UploadFileSizeTotal, resp.UploadFileSizeTotal)
			assert.Equal(t, tt.want.UploadFileSizeCurrent, resp.UploadFileSizeCurrent)
			assert.Equal(t, tt.want.AppID, resp.AppID)
			assert.Equal(t, tt.want.AppName, resp.AppName)
			assert.Equal(t, tt.want.UserCancel, resp.UserCancel)
			assert.Equal(t, tt.want.IsFileReady, resp.IsFileReady)
			assert.Equal(t, tt.want.DownloadFinished, resp.DownloadFinished)
			assert.Equal(t, tt.want.IsDeleted, resp.IsDeleted)
			assert.Equal(t, tt.want.UploadTime, resp.UploadTime)
			assert.NotEqual(t, resp.PendingTime, InvalidTime)
			assert.Equal(t, tt.want.DownloadTime, resp.DownloadTime)
			assert.Equal(t, tt.want.RunningTime, resp.RunningTime)
			assert.Equal(t, tt.want.TerminatingTime, resp.TerminatingTime)
			assert.NotEqual(t, resp.CreateTime, InvalidTime)
			assert.NotEqual(t, resp.UpdateTime, InvalidTime)

		})
	}
}

func TestConvertAdminJobModel(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	logger, err := logging.NewLogger()
	if err != nil {
		t.Error(err)
	}

	ctx.Set(logging.LoggerName, logger)

	now := time.Now()

	type args struct {
		ctx        context.Context
		logger     *zap.SugaredLogger
		req        *adminjobcreate.Request
		userID     snowflake.ID
		jobID      snowflake.ID
		appInfo    *models.Application
		inputZone  consts.Zone
		outputZone consts.Zone
		queue      string
	}
	core := 2
	memory := 256
	tests := []struct {
		name    string
		args    args
		want    *models.Job
		wantErr bool
	}{
		{
			name: "admin normal",
			args: args{
				ctx:    ctx,
				logger: logger,
				req: &adminjobcreate.Request{
					Request: jobcreate.Request{
						Name: "TestJob",
						Params: jobcreate.Params{
							Application: jobcreate.Application{
								Command: "Command",
								AppID:   "456",
							},
							Resource: &jobcreate.Resource{
								Cores:  &core,
								Memory: &memory,
							},
							EnvVars: map[string]string{},
							Input: &jobcreate.Input{
								Type:        "CloudStorage",
								Source:      "https://test/10086/path/to/input",
								Destination: "Destination",
							},
							Output: &jobcreate.Output{
								Type:          "CloudStorage",
								Address:       "https://test/10086/path/to/output",
								NoNeededPaths: "NoNeededPaths",
								NeededPaths:   "NeededPaths",
							},
							TmpWorkdir:        false,
							SubmitWithSuspend: false,
							CustomStateRule: &jobcreate.CustomStateRule{
								KeyStatement: "RuleKey",
								ResultState:  "ResultState",
							},
						},
						Timeout: 60,
						Zone:    "JobZone",
						Comment: "Test comment",
					},
					Queue: "QUEUE",
				},
				userID: 123,
				jobID:  10086, // 对应:3ZU
				appInfo: &models.Application{
					ID:                snowflake.ID(123456),
					Name:              "TestApp",
					Type:              "",
					Version:           "",
					AppParamsVersion:  0,
					Image:             "",
					Endpoint:          "",
					Command:           "",
					PublishStatus:     "",
					Description:       "",
					IconUrl:           "",
					CoresMaxLimit:     0,
					CoresPlaceholder:  "",
					FileFilterRule:    "",
					ResidualEnable:    false,
					ResidualLogRegexp: "",
					ResidualLogParser: "",
					LicManagerId:      0,
					SnapshotEnable:    false,
					BinPath:           "",
				},
				inputZone:  "InputStorageZone",
				outputZone: "OutputStorageZone",
				queue:      "QUEUE",
			},
			want: &models.Job{
				ID:                          10086,
				Name:                        "TestJob",
				Comment:                     "Test comment",
				UserID:                      123,
				JobSource:                   "",
				State:                       2,
				SubState:                    200,
				StateReason:                 "User Submit",
				ExitCode:                    "",
				Params:                      "Params",
				UserZone:                    "JobZone",
				Timeout:                     60,
				FileClassifier:              "",
				ResourceUsageCpus:           2,
				ResourceUsageMemory:         256,
				CustomStateRuleKeyStatement: "RuleKey",
				CustomStateRuleResultState:  "ResultState",
				HPCJobID:                    "",
				Zone:                        "JobZone",
				ResourceAssignCpus:          2,
				ResourceAssignMemory:        256,
				Command:                     "Command",
				WorkDir:                     "Destination/",
				OriginJobID:                 "",
				Queue:                       "QUEUE",
				Priority:                    0,
				ExecHosts:                   "",
				ExecHostNum:                 0,
				ExecutionDuration:           0,
				InputType:                   "CloudStorage",
				InputDir:                    "https://test/10086/path/to/input",
				Destination:                 "Destination",
				OutputType:                  "CloudStorage",
				OutputDir:                   "https://test/10086/path/to/output/",
				NoNeededPaths:               "NoNeededPaths",
				NeededPaths:                 "NeededPaths",
				FileInputStorageZone:        "InputStorageZone",
				FileOutputStorageZone:       "OutputStorageZone",
				DownloadFileSizeTotal:       0,
				DownloadFileSizeCurrent:     0,
				UploadFileSizeTotal:         0,
				UploadFileSizeCurrent:       0,
				AppID:                       123456,
				AppName:                     "TestApp",
				UserCancel:                  0,
				IsFileReady:                 0,
				DownloadFinished:            0,
				IsDeleted:                   0,
				UploadTime:                  InvalidTime,
				PendingTime:                 InvalidTime,
				DownloadTime:                InvalidTime,
				RunningTime:                 InvalidTime,
				TerminatingTime:             InvalidTime,
				TransmittingTime:            InvalidTime,
				SuspendingTime:              InvalidTime,
				SuspendedTime:               InvalidTime,
				SubmitTime:                  InvalidTime,
				EndTime:                     InvalidTime,
				CreateTime:                  now,
				UpdateTime:                  now,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := ConvertAdminJobModel(tt.args.ctx, tt.args.logger, tt.args.req, tt.args.userID, tt.args.jobID, tt.args.appInfo, tt.args.inputZone, tt.args.outputZone, tt.args.queue, nil)

			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertAdminJobModel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				t.Log(err)
				return
			}

			assert.Equal(t, tt.want.ID, resp.ID)
			assert.Equal(t, tt.want.Name, resp.Name)
			assert.Equal(t, tt.want.Comment, resp.Comment)
			assert.Equal(t, tt.want.UserID, resp.UserID)
			assert.Equal(t, tt.want.JobSource, resp.JobSource)
			assert.Equal(t, tt.want.State, resp.State)
			assert.Equal(t, tt.want.SubState, resp.SubState)
			assert.Contains(t, resp.StateReason, tt.want.StateReason)
			assert.Equal(t, tt.want.ExitCode, resp.ExitCode)
			assert.NotEmpty(t, resp.Params)
			assert.Equal(t, tt.want.UserZone, resp.UserZone)
			assert.Equal(t, tt.want.Timeout, resp.Timeout)
			assert.Equal(t, tt.want.FileClassifier, resp.FileClassifier)
			assert.Equal(t, tt.want.ResourceUsageCpus, resp.ResourceUsageCpus)
			assert.Equal(t, tt.want.ResourceUsageMemory, resp.ResourceUsageMemory)
			assert.Equal(t, tt.want.CustomStateRuleKeyStatement, resp.CustomStateRuleKeyStatement)
			assert.Equal(t, tt.want.CustomStateRuleResultState, resp.CustomStateRuleResultState)
			assert.Equal(t, tt.want.HPCJobID, resp.HPCJobID)
			assert.Equal(t, tt.want.Zone, resp.Zone)
			assert.Equal(t, tt.want.ResourceAssignCpus, resp.ResourceAssignCpus)
			assert.Equal(t, tt.want.ResourceAssignMemory, resp.ResourceAssignMemory)
			assert.Equal(t, tt.want.Command, resp.Command)
			assert.Equal(t, tt.want.WorkDir, resp.WorkDir)
			assert.Equal(t, tt.want.OriginJobID, resp.OriginJobID)
			assert.Equal(t, tt.want.Queue, resp.Queue)
			assert.Equal(t, tt.want.Priority, resp.Priority)
			assert.Equal(t, tt.want.ExecHosts, resp.ExecHosts)
			assert.Equal(t, tt.want.ExecHostNum, resp.ExecHostNum)
			assert.Equal(t, tt.want.ExecutionDuration, resp.ExecutionDuration)
			assert.Equal(t, tt.want.InputType, resp.InputType)
			assert.Equal(t, tt.want.InputDir, resp.InputDir)
			assert.Equal(t, tt.want.Destination, resp.Destination)
			assert.Equal(t, tt.want.OutputType, resp.OutputType)
			assert.Equal(t, tt.want.OutputDir, resp.OutputDir)
			assert.Equal(t, tt.want.NoNeededPaths, resp.NoNeededPaths)
			assert.Equal(t, tt.want.NeededPaths, resp.NeededPaths)
			assert.Equal(t, tt.want.FileInputStorageZone, resp.FileInputStorageZone)
			assert.Equal(t, tt.want.FileOutputStorageZone, resp.FileOutputStorageZone)
			assert.Equal(t, tt.want.DownloadFileSizeTotal, resp.DownloadFileSizeTotal)
			assert.Equal(t, tt.want.DownloadFileSizeCurrent, resp.DownloadFileSizeCurrent)
			assert.Equal(t, tt.want.UploadFileSizeTotal, resp.UploadFileSizeTotal)
			assert.Equal(t, tt.want.UploadFileSizeCurrent, resp.UploadFileSizeCurrent)
			assert.Equal(t, tt.want.AppID, resp.AppID)
			assert.Equal(t, tt.want.AppName, resp.AppName)
			assert.Equal(t, tt.want.UserCancel, resp.UserCancel)
			assert.Equal(t, tt.want.IsFileReady, resp.IsFileReady)
			assert.Equal(t, tt.want.DownloadFinished, resp.DownloadFinished)
			assert.Equal(t, tt.want.IsDeleted, resp.IsDeleted)
			assert.Equal(t, tt.want.UploadTime, resp.UploadTime)
			assert.NotEqual(t, resp.PendingTime, InvalidTime)
			assert.Equal(t, tt.want.DownloadTime, resp.DownloadTime)
			assert.Equal(t, tt.want.RunningTime, resp.RunningTime)
			assert.Equal(t, tt.want.TerminatingTime, resp.TerminatingTime)
			assert.NotEqual(t, resp.CreateTime, InvalidTime)
			assert.NotEqual(t, resp.UpdateTime, InvalidTime)
		})
	}
}

func Test_timeParse(t *testing.T) {
	type args struct {
		currentTime *time.Time
	}
	now := time.Now()
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		{
			name: "normal",
			args: args{
				currentTime: &now,
			},
			want: now,
		},
		{
			name: "nil",
			args: args{
				currentTime: nil,
			},
			want: InvalidTime,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := TimeParse(tt.args.currentTime)
			assert.Equal(t, tt.want, resp)
		})
	}
}

func TestHpcModelToYsJobModel(t *testing.T) {
	type args struct {
		job *v20230530.JobInHPC
	}
	now := time.Now()
	tests := []struct {
		name string
		args args
		want *models.Job
	}{
		{
			name: "normal",
			args: args{
				job: &v20230530.JobInHPC{
					ID:          "job123",
					Application: "app123",
					Environment: map[string]string{
						"key1": "value1",
						"key2": "value2",
					},
					Command: "run",
					Override: v20230530.JobInHPCOverride{
						Enable:  false,
						WorkDir: "",
					},
					Queue: "queue123",
					Resource: v20230530.JobInHPCResource{
						Cores: 20,
					},
					Inputs: []v20230530.JobInHPCInputStorage{
						{
							Src:  "/input",
							Dst:  "/dst",
							Type: "CloudStorage",
						},
					},
					Output: &v20230530.JobInHPCOutputStorage{
						Dst:           "/output",
						Type:          "CloudStorage",
						NoNeededPaths: "",
						NeededPaths:   "",
					},
					CustomStateRule: &v20230530.JobInHPCCustomStateRule{
						KeyStatement: "key statement",
						ResultState:  "result state",
					},
					SchedulerID:       "scheduler123",
					Status:            "running",
					FileSyncState:     "synced",
					StateReason:       "reason",
					PendingTime:       &now,
					RunningTime:       &now,
					CompletingTime:    &now,
					CompletedTime:     &now,
					AllocCores:        8,
					ExitCode:          "0",
					ExecutionDuration: 120,
					DownloadProgress: v20230530.JobInHPCProgress{
						Total:   200,
						Current: 50,
					},
					UploadProgress: v20230530.JobInHPCProgress{
						Total:   200,
						Current: 200,
					},
					Priority:     1,
					ExecHosts:    "host1,host2",
					ExecHostsNum: 2,
				},
			},
			want: &models.Job{
				ID:                          0,
				Name:                        "",
				Comment:                     "",
				UserID:                      0,
				JobSource:                   "",
				State:                       3,
				SubState:                    300,
				StateReason:                 "reason",
				ExitCode:                    "0",
				Params:                      "",
				UserZone:                    "",
				Timeout:                     0,
				FileClassifier:              "",
				ResourceUsageCpus:           20,
				ResourceUsageMemory:         0,
				CustomStateRuleKeyStatement: "",
				CustomStateRuleResultState:  "",
				HPCJobID:                    "",
				Zone:                        "",
				ResourceAssignCpus:          8,
				ResourceAssignMemory:        0,
				Command:                     "",
				WorkDir:                     "",
				OriginJobID:                 "scheduler123",
				Queue:                       "queue123",
				Priority:                    1,
				ExecHosts:                   "host1,host2",
				ExecHostNum:                 2,
				ExecutionDuration:           120,
				InputType:                   "CloudStorage",
				InputDir:                    "/input",
				Destination:                 "/dst",
				OutputType:                  "CloudStorage",
				OutputDir:                   "/output",
				NoNeededPaths:               "",
				NeededPaths:                 "",
				FileInputStorageZone:        "",
				FileOutputStorageZone:       "",
				DownloadFileSizeTotal:       200,
				DownloadFileSizeCurrent:     50,
				UploadFileSizeTotal:         200,
				UploadFileSizeCurrent:       200,
				AppID:                       0,
				AppName:                     "",
				UserCancel:                  0,
				IsFileReady:                 0,
				DownloadFinished:            0,
				IsDeleted:                   0,
				UploadTime:                  InvalidTime,
				DownloadTime:                InvalidTime,
				PendingTime:                 now,
				RunningTime:                 now,
				TerminatingTime:             InvalidTime,
				TransmittingTime:            now,
				SuspendingTime:              InvalidTime,
				SuspendedTime:               InvalidTime,
				SubmitTime:                  InvalidTime,
				EndTime:                     InvalidTime,
				CreateTime:                  InvalidTime,
				UpdateTime:                  InvalidTime,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HpcModelToYsJobModel(tt.args.job)
			assert.NotNil(t, got)

			assert.Equal(t, tt.want.State, got.State)
			assert.Equal(t, tt.want.SubState, got.SubState)
			assert.Contains(t, got.StateReason, tt.want.StateReason)
			assert.Equal(t, tt.want.ExitCode, got.ExitCode)
			assert.Equal(t, tt.want.ResourceUsageCpus, got.ResourceUsageCpus)
			assert.Equal(t, tt.want.ResourceAssignCpus, got.ResourceAssignCpus)
			assert.Equal(t, tt.want.OriginJobID, got.OriginJobID)
			assert.Equal(t, tt.want.Queue, got.Queue)
			assert.Equal(t, tt.want.Priority, got.Priority)
			assert.Equal(t, tt.want.ExecHosts, got.ExecHosts)
			assert.Equal(t, tt.want.ExecHostNum, got.ExecHostNum)
			assert.Equal(t, tt.want.ExecutionDuration, got.ExecutionDuration)
			assert.Equal(t, tt.want.InputType, got.InputType)
			assert.Equal(t, tt.want.InputDir, got.InputDir)
			assert.Equal(t, tt.want.Destination, got.Destination)
			assert.Equal(t, tt.want.OutputType, got.OutputType)
			assert.Equal(t, tt.want.OutputDir, got.OutputDir)
			assert.Equal(t, tt.want.DownloadFileSizeTotal, got.DownloadFileSizeTotal)
			assert.Equal(t, tt.want.DownloadFileSizeCurrent, got.DownloadFileSizeCurrent)
			assert.Equal(t, tt.want.UploadFileSizeTotal, got.UploadFileSizeTotal)
			assert.Equal(t, tt.want.UploadFileSizeCurrent, got.UploadFileSizeCurrent)
			assert.Equal(t, tt.want.RunningTime, got.RunningTime)
			assert.Equal(t, tt.want.TransmittingTime, got.TransmittingTime)
		})
	}
}
