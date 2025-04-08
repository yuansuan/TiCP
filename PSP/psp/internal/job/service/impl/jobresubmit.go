package impl

import (
	"context"
	"encoding/json"
	"path/filepath"
	"strings"

	"google.golang.org/grpc/status"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/storage"
	"github.com/yuansuan/ticp/PSP/psp/internal/job/consts"
	"github.com/yuansuan/ticp/PSP/psp/internal/job/dao/model"
	"github.com/yuansuan/ticp/PSP/psp/internal/job/dto"
	"github.com/yuansuan/ticp/PSP/psp/internal/job/util"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
)

func (s *jobServiceImpl) JobResubmit(ctx context.Context, req *dto.ResubmitRequest, loginUserId snowflake.ID, username string) (*dto.ResubmitResponse, error) {
	job, jobAttr, err := s.checkJobAndJobAttr(ctx, req)
	if err != nil {
		return nil, err
	}

	submitParam := &dto.JobSubmitParamInfo{}
	if err := json.Unmarshal([]byte(jobAttr.Value), submitParam); err != nil {
		return nil, err
	}
	param := submitParam.Param
	param.UserID = 0

	app, err := s.GetAndCheckApp(ctx, param.AppID)
	if err != nil {
		return nil, status.Errorf(errcode.ErrJobResubmitAppNotExistOrNotPublish, "compute app: [%v] not exist or not publish", param.AppID)
	}

	tempUploadDir, err := s.CreateJobTempDir(ctx, username, app.computeType)
	if err != nil {
		return nil, err
	}
	if param.WorkDir.IsTemp {
		param.WorkDir.Path = tempUploadDir
	}

	taskKey := ""
	if app.isCloud {
		err = s.uploadFiles(ctx, username, param.WorkDir.Path, &taskKey, submitParam.Dirs, submitParam.Files, param.WorkDir.IsTemp)
		if err != nil {
			return nil, err
		}
	} else {
		err = s.linkFiles(ctx, username, job.WorkDir, param.WorkDir.Path, submitParam.Dirs, submitParam.Files)
		if err != nil {
			return nil, err
		}
	}

	return &dto.ResubmitResponse{
		Param: param,
		Extension: &dto.Extension{
			AppType:  app.appType,
			UploadId: taskKey,
		},
	}, nil
}

func (s *jobServiceImpl) checkJobAndJobAttr(ctx context.Context, req *dto.ResubmitRequest) (*model.Job, *model.JobAttr, error) {
	jobId := snowflake.MustParseString(req.JobId)
	exist, job, err := s.jobDao.GetJobDetail(ctx, jobId)
	if err != nil {
		return nil, nil, err
	}
	if !exist {
		return nil, nil, status.Errorf(errcode.ErrJobNotExist, "job: [%v] not exist", req.JobId)
	}

	if job.RawState != consts.JobStateCompleted && job.RawState != consts.JobStateFailed && job.RawState != consts.JobStateTerminated && job.State != consts.JobStateBurstFailed {
		return nil, nil, status.Errorf(errcode.ErrJobStatusNotSupportResubmit, "job status: [%v] not support resubmit", job.RawState)
	}

	exist, jobAttr, err := s.jobAttrDao.GetJobAttrByKey(ctx, jobId, consts.JobAttrKeySubmitParams)
	if err != nil {
		return nil, nil, err
	}
	if !exist {
		return nil, nil, status.Errorf(errcode.ErrJobLastSubmitParamNotExist, "job: [%v] last submit param not exist", req.JobId)
	}
	return job, jobAttr, nil
}

func (s *jobServiceImpl) uploadFiles(ctx context.Context, username, tempUploadDir string, taskKey *string, jobDirsParam, jobFilesParam []string, isTemPath bool) error {
	if !isTemPath {
		jobDirs, jobFiles := make([]string, 0), make([]string, 0)
		for _, v := range jobDirsParam {
			jobDirs = append(jobDirs, filepath.Join(tempUploadDir, v))
		}
		for _, v := range jobFilesParam {
			jobFiles = append(jobFiles, filepath.Join(tempUploadDir, v))
		}
		jobDirsParam, jobFilesParam = jobDirs, jobFiles
	}

	resp, err := s.rpc.Storage.SubmitUploadHpcFileTask(ctx, &storage.SubmitUploadHpcFileTaskReq{
		SrcDirPaths:  jobDirsParam,
		SrcFilePaths: jobFilesParam,
		DestDirPath:  tempUploadDir,
		CurrentPath:  username,
		Overwrite:    true,
		UserName:     username,
		Cross:        isTemPath,
	})
	if err != nil {
		return err
	}
	*taskKey = resp.TaskKey

	return nil
}

func (s *jobServiceImpl) linkFiles(ctx context.Context, username, workdir, tempUploadDir string, jobDirsParam, jobFilesParam []string) error {
	jobWorkDir := strings.TrimRight(util.ConvertWorkDir(workdir, true), "/")
	pathSplits := strings.Split(jobWorkDir, "/")
	if len(pathSplits) >= 5 {
		jobWorkDir = strings.Join(pathSplits[:5], "/")
	}

	jobDirs, jobFiles, filterFiles := make([]string, 0), make([]string, 0), make([]string, 0)
	for _, v := range jobDirsParam {
		jobDirs = append(jobDirs, filepath.Join(jobWorkDir, v))
	}
	for _, v := range jobFilesParam {
		jobFiles = append(jobFiles, filepath.Join(jobWorkDir, v))
	}

	_, err := s.rpc.Storage.HardLink(ctx, &storage.HardLinkReq{
		SrcFilePaths: jobFiles,
		SrcDirPaths:  jobDirs,
		FilterPaths:  filterFiles,
		CurrentPath:  jobWorkDir,
		DstPath:      tempUploadDir,
		Overwrite:    true,
		Cross:        true,
		IsCloud:      false,
		UserName:     username,
	})
	if err != nil {
		return err
	}

	return nil
}
