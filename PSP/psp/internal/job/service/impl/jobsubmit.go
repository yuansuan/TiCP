package impl

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	apijobcreate "github.com/yuansuan/ticp/common/project-root-api/job/v1/jobcreate"
	"google.golang.org/grpc/status"

	"github.com/yuansuan/ticp/PSP/psp/internal/common"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi"
	openapijob "github.com/yuansuan/ticp/PSP/psp/internal/common/openapi/job"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/oplog"
	pbapp "github.com/yuansuan/ticp/PSP/psp/internal/common/proto/app"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/approve"
	pbnotice "github.com/yuansuan/ticp/PSP/psp/internal/common/proto/notice"
	pbproject "github.com/yuansuan/ticp/PSP/psp/internal/common/proto/project"
	pbstorage "github.com/yuansuan/ticp/PSP/psp/internal/common/proto/storage"
	pbconfig "github.com/yuansuan/ticp/PSP/psp/internal/common/proto/sysconfig"
	"github.com/yuansuan/ticp/PSP/psp/internal/job/config"
	"github.com/yuansuan/ticp/PSP/psp/internal/job/consts"
	"github.com/yuansuan/ticp/PSP/psp/internal/job/dao/model"
	"github.com/yuansuan/ticp/PSP/psp/internal/job/dto"
	"github.com/yuansuan/ticp/PSP/psp/internal/job/service/dataloader"
	"github.com/yuansuan/ticp/PSP/psp/internal/job/util"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
	"github.com/yuansuan/ticp/PSP/psp/pkg/tracelog"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/ginutil"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/strutil"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/timeutil"
)

type jobInfo struct {
	userID         snowflake.ID
	jobSetID       snowflake.ID
	projectID      snowflake.ID
	appID          string
	outJobID       string
	jobName        string
	appName        string
	userName       string
	queueName      string
	jobSetName     string
	projectName    string
	computeType    string
	fieldMap       map[string]string
	enableResidual bool
	enableSnapshot bool
}

type appInfo struct {
	name            string
	command         string
	outAppID        string
	appType         string
	computeType     string
	isCloud         bool
	enableResidual  bool
	enableSnapshot  bool
	schedulerParams map[string]string
	defaultFields   map[string]string
}

type projectInfo struct {
	id       snowflake.ID
	name     string
	filePath string
}

type submitResult struct {
	isBatch  bool
	isCloud  bool
	jobID    string
	outJobID string
	jobName  string
}

// JobSubmit 作业提交
func (s *jobServiceImpl) JobSubmit(ctx *gin.Context, param *dto.SubmitParam) ([]string, error) {
	var batchID snowflake.ID
	isBatch := len(param.MainFiles) > 1
	if isBatch {
		if !param.WorkDir.IsTemp {
			return nil, errors.New("batch submit job is not support")
		}

		batchID = s.sid.Generate()
	}

	result := make([]*submitResult, 0, len(param.MainFiles))
	jobIDs := make([]string, 0, len(param.MainFiles))

	isOpenapi := ""
	if param.IsOpenApi {
		isOpenapi = "【OPENAPI】"
	}
	for _, mainFile := range param.MainFiles {
		data, err := s.jobSubmit(ctx, param.AppID, param.ProjectID, param.UserName, param.QueueName, mainFile, param.UserID, batchID, isBatch, param.WorkDir, param.Fields, param.MainFiles)
		if err != nil {
			return nil, err
		}
		result = append(result, data)
		jobIDs = append(jobIDs, data.jobID)
		oplog.GetInstance().SaveAuditLogInfo(ctx, approve.OperateTypeEnum_JOB_MANAGER, fmt.Sprintf("%s用户%v提交求解作业[%v]", isOpenapi, ginutil.GetUserName(ctx), data.jobName))
	}

	if param.WorkDir.IsTemp && isBatch {
		tracelog.Info(ctx, fmt.Sprintf("%v delete job submit template directory [%v] after submit job", param.UserName, param.WorkDir.Path))
		go s.deletePath(ctx, param.UserName, param.WorkDir.Path, true, result[0].isCloud)
	}

	return jobIDs, nil
}

// JobSubmit 作业提交
func (s *jobServiceImpl) jobSubmit(ctx context.Context, appID, projectID, userName, queueName, mainFile string, userID, jobSetId snowflake.ID,
	isBatch bool, workDir *dto.WorkDir, fields []*dto.Field, mainFiles []string) (*submitResult, error) {
	logger := logging.GetLogger(ctx)

	project, err := s.getAndCheckProject(ctx, userID.String(), userName, projectID)
	if err != nil {
		logger.Errorf("get and check project [%v] err: %v", projectID, err)
		return nil, err
	}

	projectFilePath := project.filePath

	app, err := s.GetAndCheckApp(ctx, appID)
	if err != nil {
		logger.Errorf("get and check app [%v] err: %v", appID, err)
		return nil, err
	}

	fieldMap, err := s.parseFields(fields, mainFile)
	if err != nil {
		logger.Errorf("job submit fields parse err: %v, fields: %v", err, fields)
		return nil, err
	}

	for k, v := range app.defaultFields {
		fieldMap[k] = v
	}

	var mainFileName, jobSetName string
	jobName := fieldMap[consts.JobName]
	if isBatch {
		jobSetName = jobName
		mainFileName = s.getMainFileName(mainFile)

		if strutil.IsEmpty(jobName) {
			jobName = mainFileName
		} else {
			jobName = fmt.Sprintf("%v_%v", jobName, mainFileName)
		}
	} else {
		if strutil.IsEmpty(jobName) {
			jobName = s.getMainFileName(mainFile)

			if strutil.IsEmpty(jobName) {
				jobName = app.appType
			}
		}
	}

	fieldMap[consts.JobName] = jobName

	jobConfig := config.GetConfig()
	var workDirName, subDirName string
	if !workDir.IsTemp {
		workDirName = filepath.Base(workDir.Path)
	} else {
		files := strings.Split(mainFile, "/")
		if len(files) > 1 {
			subDirName = files[0]
			fieldMap[consts.MainFile] = strings.Replace(mainFile, fmt.Sprintf("%v/", subDirName), "", 1)
		}

		workDirPrefix := strings.ReplaceAll(jobName, " ", "_")
		if consts.WorkDirSuffixDefault == jobConfig.WorkDir.Type {
			formatWorkDir, _ := s.formatDateWorkDir(jobConfig.WorkDir.Format)
			workDirName = fmt.Sprintf("%v_%v", formatWorkDir, workDirPrefix)
		} else {
			workDirName = fmt.Sprintf("%v_%v", workDirPrefix, strutil.RandString(7))
		}
	}

	computeType := app.computeType
	workspace := jobConfig.WorkDir.Workspace

	// record job upload file list info
	tmpWorkDirFilterPaths := removeElement(mainFiles, mainFile)
	tmpWorkDirFileInfos, err := s.getTmpWorkDirFileInfos(ctx, userName, workDir.Path, workDir.IsTemp)
	if err != nil {
		logger.Errorf("get tmp work dir file infos err: %v", err)
		return nil, err
	}

	if workDir.IsTemp {
		srcPath := workDir.Path
		dstPath := filepath.Join(userName, workspace, projectFilePath, workDirName)

		if err = s.initJobWorkDir(ctx, userName, srcPath, dstPath, tmpWorkDirFilterPaths, false, isBatch); err != nil {
			logger.Errorf("init %v job work directory err: %v", computeType, err)
			return nil, err
		}
	}

	if strutil.IsEmpty(queueName) {
		configResp, err := s.rpc.SysConfig.GetJobConfig(ctx, &pbconfig.GetJobConfigRequest{})
		if err != nil {
			logger.Errorf("get system job config err: %v", err)
			return nil, err
		}

		queueName = configResp.Queue
	}

	outJobID, err := s.openapiSubmit(ctx, fieldMap, app.command, app.outAppID, workDirName, jobName, userName, computeType,
		workspace, queueName, subDirName, projectFilePath, workDir.IsTemp, app.schedulerParams)
	if err != nil {
		if workDir.IsTemp {
			tempWorkDir := filepath.Join(workspace, projectFilePath, workDirName)

			if isBatch {
				if delErr := s.deletePath(ctx, userName, tempWorkDir, false, false); delErr != nil {
					logger.Errorf("rollback: delete linked workspace directory [%v] err: %v", tempWorkDir, delErr)
				}

				tracelog.Info(ctx, fmt.Sprintf("failed to batch submit [%v] job, rollback: %v delete work directory [%v]", computeType, userName, tempWorkDir))
			} else {
				srcPath := filepath.Join(userName, workspace, projectFilePath, workDirName)
				dstPath := workDir.Path
				if mvErr := s.fileMove(ctx, srcPath, dstPath, false); mvErr != nil {
					logger.Errorf("rollback: the [%v] workspace directory [%v] move to [%v] err: %v", computeType, srcPath, dstPath, mvErr)
				}

				tracelog.Info(ctx, fmt.Sprintf("failed to submit [%v] job, rollback: %v move work directory from [%v] to [%v]",
					computeType, userName, srcPath, dstPath))
			}
		}

		return nil, err
	}

	if strutil.IsEmpty(outJobID) {
		return nil, errors.New("openapi job id is empty")
	}

	dirInfos, fileInfos := s.getJobDirAndFileInfos(tmpWorkDirFilterPaths, workDir.Path, tmpWorkDirFileInfos)

	logger.Infof("submit %v job [%v] success", app.name, outJobID)

	jobID, err := s.saveJobInfo(ctx,
		&jobInfo{
			userID:         userID,
			jobSetID:       jobSetId,
			projectID:      project.id,
			appID:          appID,
			outJobID:       outJobID,
			jobName:        jobName,
			appName:        app.name,
			userName:       userName,
			queueName:      queueName,
			jobSetName:     jobSetName,
			projectName:    project.name,
			computeType:    computeType,
			fieldMap:       fieldMap,
			enableResidual: app.enableResidual,
			enableSnapshot: app.enableSnapshot,
		},
		&dto.JobSubmitParamInfo{
			Param: &dto.SubmitParam{
				AppID:     appID,
				ProjectID: projectID,
				UserID:    userID,
				UserName:  userName,
				QueueName: queueName,
				MainFiles: []string{mainFile},
				WorkDir:   workDir,
				Fields:    fields,
			},
			Dirs:  dirInfos,
			Files: fileInfos,
		})
	if err != nil {
		logger.Errorf("save submited job info err: %v", err)
		return nil, err
	}

	go s.sendSubmitMessage(userID.String(), jobID, jobName)

	return &submitResult{
		isBatch:  isBatch,
		isCloud:  false,
		jobID:    jobID,
		jobName:  jobName,
		outJobID: outJobID,
	}, nil
}

func (s *jobServiceImpl) getTmpWorkDirFileInfos(ctx context.Context, userName, workDir string, isTempPath bool) ([]*pbstorage.File, error) {
	dirQueues, fileInfos := make([]string, 0), make([]*pbstorage.File, 0)
	err := s.getDirFileInfos(ctx, &dirQueues, &fileInfos, userName, workDir, isTempPath)
	if err != nil {
		return nil, err
	}

	for len(dirQueues) > 0 {
		dirQueuesLen := len(dirQueues)
		for i := 0; i < dirQueuesLen; i++ {
			err = s.getDirFileInfos(ctx, &dirQueues, &fileInfos, userName, dirQueues[i], isTempPath)
			if err != nil {
				return nil, err
			}
		}
		dirQueues = dirQueues[dirQueuesLen:]
	}

	return fileInfos, nil
}

func (s *jobServiceImpl) getDirFileInfos(ctx context.Context, dirQueues *[]string, fileInfos *[]*pbstorage.File, userName, workDir string, isTempPath bool) error {
	uploadFileListRes, err := s.rpc.Storage.List(ctx, &pbstorage.ListReq{
		Path:         workDir,
		UserName:     userName,
		ShowHideFile: true,
		Cross:        isTempPath,
	})
	if err != nil {
		return err
	}

	for _, file := range uploadFileListRes.Files {
		if file.IsDir {
			*dirQueues = append(*dirQueues, file.Path)
		} else {
			*fileInfos = append(*fileInfos, file)
		}
	}

	return nil
}

func (*jobServiceImpl) getJobDirAndFileInfos(tmpWorkDirFilterPaths []string, workDirPath string, tmpWorkDirFileInfos []*pbstorage.File) ([]string, []string) {
	dirInfos, fileInfos := make([]string, 0), make([]string, 0)
	tmpWorkDirFilterPathsMap := make(map[string]bool)
	for _, path := range tmpWorkDirFilterPaths {
		tmpUploadDirPath := filepath.Join(workDirPath, path)
		tmpWorkDirFilterPathsMap[tmpUploadDirPath] = true
	}
	for _, v := range tmpWorkDirFileInfos {
		if tmpWorkDirFilterPathsMap[v.Path] {
			continue
		}

		path := strings.TrimPrefix(v.Path, workDirPath)[1:]
		if v.IsDir {
			dirInfos = append(dirInfos, path)
		} else {
			fileInfos = append(fileInfos, path)
		}
	}
	return dirInfos, fileInfos
}

func (s *jobServiceImpl) getMainFileName(mainFile string) string {
	if !strutil.IsEmpty(mainFile) {
		fileName := filepath.Base(mainFile)
		fileExt := filepath.Ext(fileName)
		return fileName[:len(fileName)-len(fileExt)]
	}

	return ""
}

// parseFields 解析应用模版数据
func (s *jobServiceImpl) parseFields(fields []*dto.Field, mainFile string) (map[string]string, error) {
	envMap := make(map[string]string, consts.DefaultSize)

	if !strutil.IsEmpty(mainFile) {
		envMap[consts.MainFile] = mainFile
	}

	for _, field := range fields {
		if strutil.IsEmpty(field.ID) {
			continue
		}

		if consts.NodeSelectorType == field.Type && !strutil.IsEmpty(field.Value) {
			var nodeMap map[string]int64
			if err := json.Unmarshal([]byte(field.Value), &nodeMap); err != nil {
				return nil, fmt.Errorf("unmarshal custom_json_value_string [%v] err: %v", field.Value, err)
			}

			var result []string
			for k, v := range nodeMap {
				result = append(result, fmt.Sprintf("%v:%v", k, v))
			}

			if len(result) > 0 {
				nodeSelector := strings.Replace(strings.Trim(fmt.Sprint(result), "[]"), " ", ",", -1)
				envMap[consts.NodeSelector] = nodeSelector
			}

		} else if consts.DateTime == field.ID && !strutil.IsEmpty(field.Value) {
			t, err := timeutil.ParseJsonTime(field.Value)
			if err != nil {
				return nil, fmt.Errorf("field [%v] date format err: %v", field.Value, err)
			}

			unixTime := fmt.Sprintf("%v", timeutil.ParseTime(t))
			envMap[consts.DateTime] = unixTime

		} else if consts.MultipleType == field.Type && len(field.Values) > 0 {
			multiValues := ""
			for _, value := range field.Values {
				if strutil.IsEmpty(multiValues) {
					multiValues = value
				} else {
					multiValues += consts.Semicolon + value
				}
			}

			envMap[field.ID] = multiValues

		} else if consts.MultipleType != field.Type && !strutil.IsEmpty(field.Value) {
			envMap[field.ID] = field.Value
		}
	}

	return envMap, nil
}

// parseFields 解析应用模版默认数据
func (s *jobServiceImpl) parseDefaultFields(subForm *pbapp.SubForm) (map[string]string, error) {
	fieldMap := make(map[string]string)
	if subForm == nil {
		return fieldMap, nil
	}

	for _, section := range subForm.Section {
		for _, field := range section.Field {
			if field.Hidden {
				if consts.TextType == field.Type || consts.ListType == field.Type {
					if !strutil.IsEmpty(field.DefaultValue) {
						fieldMap[field.Id] = field.DefaultValue
					}
				} else if consts.MultipleType == field.Type && len(field.DefaultValues) > 0 {
					multiValues := ""
					for _, value := range field.DefaultValues {
						if strutil.IsEmpty(multiValues) {
							multiValues = value
						} else {
							multiValues += consts.Semicolon + value
						}
					}

					fieldMap[field.Id] = multiValues
				}
			}
		}
	}

	return fieldMap, nil
}

// saveJob 保存作业信息
func (s *jobServiceImpl) saveJobInfo(ctx context.Context, jobInfo *jobInfo, jobSubmitParam *dto.JobSubmitParamInfo) (string, error) {
	logger := logging.GetLogger(ctx)

	var job *model.Job
	openapiJob, err := openapijob.AdminGetJob(s.localAPI, jobInfo.outJobID)
	if err != nil {
		logger.Errorf("get openapi local job [%v] err: %v", jobInfo.outJobID, err)
		return "", err
	}

	job = util.ConvertSubmitAdminJob(openapiJob)
	job.Queue = jobInfo.queueName

	appID, err := snowflake.ParseString(jobInfo.appID)
	if err != nil {
		return "", err
	}

	job.AppId = appID
	job.JobSetId = jobInfo.jobSetID
	job.ProjectId = jobInfo.projectID
	job.Name = jobInfo.jobName
	job.AppName = jobInfo.appName
	job.UserId = jobInfo.userID
	job.UserName = jobInfo.userName
	job.JobSetName = jobInfo.jobSetName
	job.ProjectName = jobInfo.projectName
	job.Type = jobInfo.computeType
	job.OutJobId = jobInfo.outJobID
	job.State = consts.JobStateSubmitted
	job.VisAnalysis = util.SetVisAnalysisValue(jobInfo.enableResidual, jobInfo.enableSnapshot)

	has, dbJob, err := s.jobDao.GetJobByOutID(ctx, jobInfo.outJobID, jobInfo.computeType)
	if err != nil {
		logger.Errorf("get job [%v] err: %v", jobInfo.outJobID, err)
		return "", err
	}

	var jobID snowflake.ID
	if has {
		logger.Errorf("the submit job is exist, new job info: %+v, old job info: %+v", job, dbJob)
		return "", status.Error(errcode.ErrJobExist, "job has exist")
	} else {
		jobID, err = s.jobDao.InsertJob(ctx, job)
		if err != nil {
			logger.Errorf("save job [%+v] err: %v", job, err)
			return "", err
		}
	}

	tracelog.Info(ctx, fmt.Sprintf("submit job success, save job info: %v", spew.Sdump(job)))

	if err = s.saveJobEnvs(ctx, jobID, jobInfo.fieldMap); err != nil {
		logger.Errorf("save job [%v] env err: %v", jobID, err)
	}
	if err = s.saveJobParam(ctx, jobID, jobSubmitParam); err != nil {
		logger.Errorf("save job [%v] param err: %v", jobID, err)
	}

	dataloader.SaveJobTimeLine(ctx, job, s.jobTimelineDao)

	return jobID.String(), nil
}

// openapiSubmit 通过openapi提交作业
func (s *jobServiceImpl) openapiSubmit(ctx context.Context, fieldMap map[string]string, command, outAppID, workDir,
	jobName, userName, computeType, workspace, queue, subDirName, projectPath string, isTemp bool, schedulerParams map[string]string) (string, error) {
	logger := logging.GetLogger(ctx)

	var err error
	var cpuNum, memNum int
	cpuField := fieldMap[consts.CpuNum]
	if !strutil.IsEmpty(cpuField) {
		cpuNum, err = strconv.Atoi(cpuField)
		if err != nil {
			return "", err
		}
	}

	memField := fieldMap[consts.MemNum]
	if !strutil.IsEmpty(memField) {
		memNum, err = strconv.Atoi(memField)
		if err != nil {
			return "", err
		}
	}

	resource := &apijobcreate.Resource{
		Cores:  &cpuNum,
		Memory: &memNum,
	}

	platform := fieldMap[consts.Platform]
	if !strutil.IsEmpty(platform) {
		nodeNum, cpuNumPerNode := getSubmitNodeNum(ctx, platform, cpuNum)
		if nodeNum <= 1 {
			nodeNum = 1
			cpuNumPerNode = cpuNum
		}

		fieldMap[consts.NodeNum] = strconv.Itoa(nodeNum)
		fieldMap[consts.CpuNum] = strconv.Itoa(cpuNumPerNode)
		fieldMap[consts.OriginCpuNum] = strconv.Itoa(cpuNum)

		totalCpuNum := cpuNumPerNode * nodeNum
		resource.Cores = &totalCpuNum
	}

	schedulerParamMap := make(map[string]string, len(schedulerParams))
	for k, v := range schedulerParams {
		schedulerParamMap[k] = replaceSchedulerParams(v, fieldMap)
	}

	var openAPI *openapi.OpenAPI
	var apiUserID, apiEndPoint, submitCmd, zoneName string
	openAPI = s.localAPI
	submitCmd = command

	localConfig := s.localCfg.Settings
	zoneName = localConfig.Zone
	apiUserID = localConfig.UserId
	apiEndPoint = localConfig.HPCEndpoint

	inputPath := fmt.Sprintf("%v/%v", apiEndPoint, filepath.Join(apiUserID, userName, workspace, projectPath, workDir))
	if !isTemp {
		inputPath = fmt.Sprintf("%v/%v", apiEndPoint, filepath.Join(apiUserID, userName, workDir))
	}
	if !strutil.IsEmpty(subDirName) {
		inputPath = fmt.Sprintf("%v/%v", inputPath, subDirName)
	}

	storageType := common.StorageTypeHPC
	submitParams := &openapijob.SubmitParams{
		Name: jobName,
		Zone: zoneName,
		Params: apijobcreate.Params{
			Application: apijobcreate.Application{
				Command: submitCmd,
				AppID:   outAppID,
			},
			Resource: resource,
			EnvVars:  fieldMap,
			Input: &apijobcreate.Input{
				Type:   storageType,
				Source: inputPath,
			},
			TmpWorkdir:        false,
			SubmitWithSuspend: false,
		},
		SchedulerParams: schedulerParamMap,
	}

	if !strutil.IsEmpty(queue) {
		submitParams.Queue = queue
	}

	jobID, err := openapijob.SubmitJob(openAPI, submitParams)
	if err != nil {
		logger.Errorf("openapi submit [%v] job err: %v", computeType, err)
		return "", err
	}

	tracelog.Info(ctx, fmt.Sprintf("using openapi submit [%v] job, params: %v", computeType, spew.Sdump(submitParams)))

	return jobID, nil
}

func getSubmitNodeNum(ctx context.Context, env string, totalCpuNum int) (int, int) {
	logger := logging.GetLogger(ctx)

	cpuNumPerNode, err := parsePlatformCpuNumPerNode(env)
	if err != nil {
		logger.Errorf("parse platform node cpu num err: %v", err)
		return 0, 0
	}

	nodeNum := float64(totalCpuNum) / float64(cpuNumPerNode)
	return int(math.Ceil(nodeNum)), cpuNumPerNode
}

func parsePlatformCpuNumPerNode(env string) (int, error) {
	platformRegexp := config.GetConfig().PlatformRegexp
	reg := regexp.MustCompile(platformRegexp)

	match := reg.FindStringSubmatch(env)
	if len(match) > 1 {
		num, err := strconv.Atoi(match[1])
		if err != nil {
			return 0, err
		}

		return num, nil
	} else {
		return 0, errors.Errorf("failed to match [%v] by regexp [%v]", env, platformRegexp)
	}
}

func replaceSchedulerParams(schedulerParams string, envs map[string]string) string {
	for key, value := range envs {
		variable := "$" + key
		schedulerParams = strings.ReplaceAll(schedulerParams, variable, value)
	}

	return schedulerParams
}

func (s *jobServiceImpl) GetAndCheckApp(ctx context.Context, appID string) (*appInfo, error) {
	appResp, err := s.rpc.App.GetAppInfoById(ctx, &pbapp.GetAppInfoByIdRequest{AppId: appID})
	if err != nil {
		return nil, err
	}

	app := appResp.App
	if app == nil {
		return nil, fmt.Errorf("the app template [%v] is not found", appID)
	}

	if app.State != consts.AppStatePublished {
		return nil, status.Error(errcode.ErrAppTemplateHasUnpublished, "app is not published")
	}

	schedulerMap := make(map[string]string, len(app.SchedulerParam))
	for _, param := range app.SchedulerParam {
		schedulerMap[param.Key] = param.Value
	}

	fieldMap, err := s.parseDefaultFields(app.SubForm)
	if err != nil {
		return nil, err
	}

	return &appInfo{
		name:            app.Name,
		command:         app.Script,
		outAppID:        app.OutAppId,
		appType:         app.Type,
		computeType:     app.ComputeType,
		isCloud:         false,
		enableResidual:  app.EnableResidual,
		enableSnapshot:  app.EnableSnapshot,
		schedulerParams: schedulerMap,
		defaultFields:   fieldMap,
	}, nil
}

func (s *jobServiceImpl) getAndCheckProject(ctx context.Context, userId, userName, projectID string) (*projectInfo, error) {
	if strutil.IsEmpty(projectID) || projectID == common.PersonalProjectID.String() {
		return &projectInfo{
			id:       common.PersonalProjectID,
			name:     common.PersonalProjectName,
			filePath: common.PersonalProjectName,
		}, nil
	}

	if resp, err := s.rpc.Project.ExistsProjectMember(ctx, &pbproject.ExistsProjectMemberRequest{UserId: userId, ProjectId: projectID}); err != nil {
		return nil, err
	} else {
		if !resp.IsExist {
			return nil, status.Errorf(errcode.ErrJobCurrentProjectNotAccess, fmt.Sprintf("[%v] use project [%v] permission denied", userName, projectID))
		}
	}

	project, err := s.rpc.Project.GetProjectDetailById(ctx, &pbproject.GetProjectDetailByIdRequest{ProjectId: projectID})
	if err != nil {
		return nil, err
	}

	if project.State != common.ProjectRunning {
		return nil, status.Errorf(errcode.ErrJobCurrentProjectNotRunning, "project is not running, cannot submit job")
	}

	projectMember, err := s.rpc.Project.GetProjectMemberByProjectIdAndUserId(ctx, &pbproject.GetProjectMemberByProjectIdAndUserIdRequest{UserId: userId, ProjectId: projectID})
	if err != nil {
		return nil, err
	}

	projectFilePath := strings.Replace(projectMember.FilePath, fmt.Sprintf("/%v/%v", userName, common.WorkspaceFolderPath), "", 1)
	return &projectInfo{
		id:       snowflake.MustParseString(project.ProjectId),
		name:     project.ProjectName,
		filePath: projectFilePath,
	}, nil
}

func (s *jobServiceImpl) sendSubmitMessage(userID, jobID, jobName string) {
	ctx := context.Background()
	logger := logging.GetLogger(ctx)

	msg := &pbnotice.WebsocketMessage{
		UserId:  userID,
		Type:    common.JobEventType,
		Content: fmt.Sprintf("作业[编号:%v 名称:%v]提交成功", jobID, jobName),
	}

	if _, err := s.rpc.Notice.SendWebsocketMessage(ctx, msg); err != nil {
		logger.Errorf("job submit send ws message err: %v", err)
	}
}

func (s *jobServiceImpl) fileMove(ctx context.Context, srcPath, dstPath string, isCloud bool) error {
	req := &pbstorage.MvReq{
		Srcpath:   srcPath,
		Dstpath:   dstPath,
		Overwrite: true,
		IsCloud:   isCloud,
	}

	if _, err := s.rpc.Storage.Mv(ctx, req); err != nil {
		return err
	}

	return nil
}

func (s *jobServiceImpl) fileLink(ctx context.Context, userName, srcPath, dstPath string, filterPaths []string, isCloud bool) error {
	req := &pbstorage.HardLinkReq{
		CurrentPath: srcPath,
		SrcDirPaths: []string{srcPath},
		FilterPaths: filterPaths,
		DstPath:     dstPath,
		Overwrite:   true,
		Cross:       true,
		IsCloud:     isCloud,
		UserName:    userName,
	}

	logging.GetLogger(ctx).Infof("link file from [%v] to [%v], req: %+v", srcPath, dstPath, req)
	if _, err := s.rpc.Storage.HardLink(ctx, req); err != nil {
		return err
	}

	return nil
}

func (s *jobServiceImpl) saveJobEnvs(ctx context.Context, jobID snowflake.ID, envs map[string]string) error {
	envJSON, err := json.Marshal(envs)
	if err != nil {
		return err
	}

	if err = s.saveJobAttr(ctx, jobID, consts.JobAttrKeySubmitEnvs, string(envJSON)); err != nil {
		return err
	}

	return nil
}

func (s *jobServiceImpl) saveJobParam(ctx context.Context, jobId snowflake.ID, submitParam *dto.JobSubmitParamInfo) error {
	jsonData, err := json.Marshal(submitParam)
	if err != nil {
		return err
	}

	err = s.saveJobAttr(ctx, jobId, consts.JobAttrKeySubmitParams, string(jsonData))
	if err != nil {
		return err
	}

	return nil
}

func (s *jobServiceImpl) saveJobAttr(ctx context.Context, jobID snowflake.ID, key, value string) error {
	has, attr, err := s.jobAttrDao.GetJobAttrByKey(ctx, jobID, key)
	if err != nil {
		return err
	}

	if has {
		attr.Key = value
		if err = s.jobAttrDao.UpdateJobAttr(ctx, attr); err != nil {
			return err
		}
	} else {
		attrs := &model.JobAttr{
			JobId: jobID,
			Key:   key,
			Value: value,
		}
		if err = s.jobAttrDao.InsertJobAttr(ctx, attrs); err != nil {
			return err
		}
	}

	return nil
}

func (s *jobServiceImpl) createWorkspace(ctx context.Context, userName string, isCloud bool) error {
	logger := logging.GetLogger(ctx)

	workspace := config.GetConfig().WorkDir.Workspace
	resp, err := s.rpc.Storage.Exist(ctx, &pbstorage.ExistReq{
		Paths:    []string{workspace},
		UserName: userName,
		Cross:    false,
		IsCloud:  isCloud,
	})
	if err != nil {
		logger.Errorf("check path [%v] exist err: %v", workspace, err)
		return err
	}

	if !resp.IsExist[0] {
		createDirReq := &pbstorage.CreateDirReq{
			Path:     workspace,
			UserName: userName,
			Cross:    false,
			IsCloud:  isCloud,
		}

		if _, err = s.rpc.Storage.CreateDir(ctx, createDirReq); err != nil {
			logger.Errorf("create workspace directory [%v] err: %v", workspace, err)
			return err
		}
	}

	return nil
}

func (s *jobServiceImpl) initJobWorkDir(ctx context.Context, userName, srcPath, dstPath string, filterPaths []string, isCloud, isBatch bool) error {
	logger := logging.GetLogger(ctx)

	if err := s.createWorkspace(ctx, userName, isCloud); err != nil {
		return err
	}

	if isBatch {
		if err := s.fileLink(ctx, userName, srcPath, dstPath, filterPaths, isCloud); err != nil {
			logger.Errorf("the temp path [%v] link to [%v] err: %v", srcPath, dstPath, err)
			return err
		}
		tracelog.Info(ctx, fmt.Sprintf("before batch submit job, %v link files from %v to %v, params: filterPaths=%+v, isCloud=%v",
			userName, srcPath, dstPath, filterPaths, isCloud))
	} else {
		if err := s.fileMove(ctx, srcPath, dstPath, isCloud); err != nil {
			logger.Errorf("the temp path [%v] move to [%v] err: %v", srcPath, dstPath, err)
			return err
		}
		tracelog.Info(ctx, fmt.Sprintf("before batch submit job, %v move files from %v to %v, params: isCloud=%v",
			userName, srcPath, dstPath, isCloud))
	}

	return nil
}

func (s *jobServiceImpl) deleteLocalWorkDir(ctx context.Context, userName, path string) error {
	if err := s.deletePath(ctx, userName, path, false, false); err != nil {
		return err
	}

	return nil
}

func (s *jobServiceImpl) deletePath(ctx context.Context, userName string, path string, isCross, isCloud bool) error {
	logger := logging.GetLogger(ctx)

	logger.Infof("%v delete path [%v], isCross=%v, isCloud=%v", userName, path, isCross, isCloud)

	resp, err := s.rpc.Storage.Exist(ctx, &pbstorage.ExistReq{
		Paths:    []string{path},
		UserName: userName,
		Cross:    isCross,
		IsCloud:  isCloud,
	})
	if err != nil {
		logger.Errorf("check path [%v] exist err: %v", path, err)
		return err
	}

	if resp.IsExist[0] {
		createDirReq := &pbstorage.RmReq{
			Paths:    []string{path},
			UserName: userName,
			Cross:    isCross,
			IsCloud:  isCloud,
		}
		if _, err = s.rpc.Storage.Rm(ctx, createDirReq); err != nil {
			logger.Errorf("delete path [%v] err: %v", path, err)
			return err
		}
	}

	return nil
}

func (s *jobServiceImpl) createLocalJobWorkDir(ctx context.Context, userName, path string) error {
	if err := s.createWorkspace(ctx, userName, false); err != nil {
		return err
	}

	createDirReq := &pbstorage.CreateDirReq{
		Path:     path,
		UserName: userName,
		Cross:    false,
		IsCloud:  false,
	}
	if _, err := s.rpc.Storage.CreateDir(ctx, createDirReq); err != nil {
		return err
	}

	return nil
}

func (s *jobServiceImpl) formatDateWorkDir(format string) (string, string) {
	var formatWorkDir, lastDir string
	formats := strings.Split(format, "/")

	size := len(formats)
	if size > 1 {
		formatTimes := make([]string, 0, size)
		for _, s := range formats {
			currentTime := time.Now().Format(s)
			formatTimes = append(formatTimes, currentTime)
		}

		lastDir = formatTimes[size-1]
		formatWorkDir = strings.Join(formatTimes, "/")
	} else {
		currentTime := time.Now().Format(format)
		lastDir = currentTime
		formatWorkDir = currentTime
	}

	return formatWorkDir, lastDir
}

func removeElement(slice []string, element string) []string {
	newSlice := make([]string, len(slice))
	copy(newSlice, slice)

	for i, v := range newSlice {
		if v == element {
			return append(newSlice[:i], newSlice[i+1:]...)
		}
	}

	return newSlice
}
