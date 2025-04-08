package util

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/yuansuan/ticp/common/go-kit/logging"
	schema "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/yuansuan/ticp/PSP/psp/internal/common"
	pb "github.com/yuansuan/ticp/PSP/psp/internal/common/proto/job"
	"github.com/yuansuan/ticp/PSP/psp/internal/job/consts"
	"github.com/yuansuan/ticp/PSP/psp/internal/job/dao/model"
	"github.com/yuansuan/ticp/PSP/psp/internal/job/dto"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/floatutil"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/strutil"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/timeutil"
)

func ConvertCmd(command string) string {
	lines := strings.Split(command, "\n")

	cmdLines := make([]string, 0, len(lines))
	for _, line := range lines {
		trimLine := strings.TrimSpace(line)
		if len(trimLine) > 0 {
			cmdLines = append(cmdLines, trimLine)
		}
	}

	return strings.Join(cmdLines, ";")
}

func ConvertJob2ListInfo(job *model.Job) *dto.JobListInfo {
	if job != nil {
		return &dto.JobListInfo{
			Id:             job.Id.String(),
			ProjectId:      job.ProjectId.String(),
			AppId:          job.AppId.String(),
			UserId:         job.UserId.String(),
			JobSetId:       job.JobSetId.String(),
			OutJobId:       job.OutJobId,
			RealJobId:      job.RealJobId,
			Type:           job.Type,
			Name:           job.Name,
			Queue:          job.Queue,
			ProjectName:    job.ProjectName,
			State:          job.State,
			RawState:       job.RawState,
			DataState:      job.DataState,
			AppName:        job.AppName,
			UserName:       job.UserName,
			JobSetName:     job.JobSetName,
			Priority:       strconv.Itoa(job.Priority),
			CpusAlloc:      strconv.Itoa(job.CpusAlloc),
			MemAlloc:       strconv.Itoa(job.MemAlloc),
			ExecDuration:   strconv.Itoa(job.ExecDuration),
			EnableResidual: job.VisAnalysis[consts.JobVisAnalysisResidual],
			EnableSnapshot: job.VisAnalysis[consts.JobVisAnalysisSnapshot],
			SubmitTime:     timeutil.DefaultFormatTime(job.SubmitTime),
			PendTime:       timeutil.DefaultFormatTime(job.PendTime),
			StartTime:      timeutil.DefaultFormatTime(job.StartTime),
			EndTime:        timeutil.DefaultFormatTime(job.EndTime),
			SuspendTime:    timeutil.DefaultFormatTime(job.SuspendTime),
		}
	}

	return nil
}

func ConvertJob2Detail(job *model.Job) *dto.JobDetailInfo {
	if job != nil {
		workDir := ConvertWorkDir(job.WorkDir, false)
		exitCode := ConvertJobExitCode(job.ExitCode)

		return &dto.JobDetailInfo{
			Id:             job.Id.String(),
			AppId:          job.AppId.String(),
			UserId:         job.UserId.String(),
			JobSetId:       job.JobSetId.String(),
			OutJobId:       job.OutJobId,
			RealJobId:      job.RealJobId,
			ProjectId:      job.ProjectId.String(),
			Type:           job.Type,
			Name:           job.Name,
			Queue:          job.Queue,
			State:          job.State,
			RawState:       job.RawState,
			DataState:      job.DataState,
			ExitCode:       exitCode,
			ProjectName:    job.ProjectName,
			AppName:        job.AppName,
			UserName:       job.UserName,
			JobSetName:     job.JobSetName,
			ClusterName:    job.ClusterName,
			WorkDir:        workDir,
			ExecHosts:      job.ExecHosts,
			StateReason:    job.Reason,
			Priority:       strconv.Itoa(job.Priority),
			CpusAlloc:      strconv.Itoa(job.CpusAlloc),
			MemAlloc:       strconv.Itoa(job.MemAlloc),
			ExecDuration:   strconv.Itoa(job.ExecDuration),
			ExecHostNum:    strconv.Itoa(job.ExecHostNum),
			FileFilterRegs: make([]string, 0),
			Timelines:      make([]*dto.JobTimeLine, 0),
			SubmitTime:     timeutil.DefaultFormatTime(job.SubmitTime),
			PendTime:       timeutil.DefaultFormatTime(job.PendTime),
			StartTime:      timeutil.DefaultFormatTime(job.StartTime),
			EndTime:        timeutil.DefaultFormatTime(job.EndTime),
			SuspendTime:    timeutil.DefaultFormatTime(job.SuspendTime),
		}
	}

	return nil
}

func ConvertSubmitAdminJob(job *schema.AdminJobInfo) *model.Job {
	logger := logging.Default()

	if job != nil {
		var cpuAlloc, memAlloc int
		resource := job.AllocResource
		if resource != nil {
			cpuAlloc = resource.Cores
			memAlloc = resource.Memory
		}

		modelJob := &model.Job{
			OutJobId:     job.ID,
			RealJobId:    job.OriginJobID,
			Name:         job.Name,
			RawState:     job.JobState,
			Priority:     job.Priority,
			ExitCode:     job.ExitCode,
			CpusAlloc:    cpuAlloc,
			MemAlloc:     memAlloc,
			ExecDuration: job.ExecutionDuration,
			ExecHostNum:  job.ExecHostNum,
			Reason:       job.StateReason,
			WorkDir:      job.Workdir,
			ExecHosts:    job.ExecHosts,
			ClusterName:  job.Zone,
		}

		submitTime, err := timeutil.ParseJsonTime(job.CreateTime)
		if err == nil {
			modelJob.SubmitTime = submitTime
		} else {
			logger.Errorf("parse submit time [%v] err: %v", job.CreateTime, err)
		}

		return modelJob
	}

	return nil
}

func ConvertSubmitJob(job *schema.JobInfo) *model.Job {
	logger := logging.Default()

	if job != nil {
		var cpuAlloc, memAlloc int
		resource := job.AllocResource
		if resource != nil {
			cpuAlloc = resource.Cores
			memAlloc = resource.Memory
		}

		modelJob := &model.Job{
			OutJobId:     job.ID,
			Name:         job.Name,
			RawState:     job.JobState,
			ExitCode:     job.ExitCode,
			CpusAlloc:    cpuAlloc,
			MemAlloc:     memAlloc,
			ExecDuration: job.ExecutionDuration,
			Reason:       job.StateReason,
			WorkDir:      job.Workdir,
			ClusterName:  job.Zone,
		}

		submitTime, err := timeutil.ParseJsonTime(job.CreateTime)
		if err == nil {
			modelJob.SubmitTime = submitTime
		} else {
			logger.Errorf("parse submit time [%v] err: %v", job.CreateTime, err)
		}

		return modelJob
	}

	return nil
}

func ConvertJob2GRPCDetail(job *model.Job) *pb.JobDetailField {
	if job != nil {
		workDir := ConvertWorkDir(job.WorkDir, false)
		exitCode := ConvertJobExitCode(job.ExitCode)

		return &pb.JobDetailField{
			Id:            job.Id.String(),
			AppId:         job.AppId.String(),
			UserId:        job.UserId.String(),
			OutJobId:      job.OutJobId,
			RealJobId:     job.RealJobId,
			ProjectId:     job.ProjectId.String(),
			Type:          job.Type,
			Name:          job.Name,
			Queue:         job.Queue,
			State:         job.State,
			RawState:      job.RawState,
			ExitCode:      exitCode,
			AppName:       job.AppName,
			UserName:      job.UserName,
			ClusterName:   job.ClusterName,
			ProjectName:   job.ProjectName,
			WorkDir:       workDir,
			ExecHosts:     job.ExecHosts,
			Priority:      int64(job.Priority),
			CpusAlloc:     int64(job.CpusAlloc),
			MemAlloc:      int64(job.MemAlloc),
			ExecDuration:  int64(job.ExecDuration),
			ExecHostNum:   int64(job.ExecHostNum),
			SubmitTime:    timestamppb.New(job.SubmitTime),
			PendTime:      timestamppb.New(job.PendTime),
			StartTime:     timestamppb.New(job.StartTime),
			EndTime:       timestamppb.New(job.EndTime),
			SuspendTime:   timestamppb.New(job.SuspendTime),
			TerminateTime: timestamppb.New(job.TerminateTime),
		}
	}

	return nil
}

func ConvertStatisticsJobToDTOOverview(v *model.StatisticsJob, id, queryType string) *dto.StatisticsOverview {
	overview := &dto.StatisticsOverview{
		UId:         id,
		ComputeType: v.Type,
		ProjectName: v.ProjectName,
		CPUTime:     floatutil.NumberToFloatStr(v.CPUTime, common.DecimalPlaces),
	}

	switch queryType {
	case consts.JobStatisticsQueryTypeApp:
		overview.Id = v.AppId.String()
		overview.Name = v.AppName
	case consts.JobStatisticsQueryTypeUser:
		overview.Id = v.UserId.String()
		overview.Name = v.UserName
	}

	return overview
}

func ConvertStatisticsJobToDTODetail(v *model.StatisticsJob) *dto.JobDetailInfo {
	if v == nil || v.Job == nil {
		return nil
	}

	jobDetailInfo := ConvertJob2Detail(v.Job)
	jobDetailInfo.CPUTime = floatutil.NumberToFloatStr(v.CPUTime, common.DecimalPlaces)

	return jobDetailInfo
}

func ConvertJobCPUTimeMetric(metrics *dto.JobCPUTimeMetric) *pb.GetJobMetricResponse {
	if metrics == nil {
		return nil
	}

	appMetrics := make([]*pb.MetricKV, 0, len(metrics.AppMetrics))
	userMetrics := make([]*pb.MetricKV, 0, len(metrics.UserMetrics))

	for _, metric := range metrics.AppMetrics {
		metricKV := &pb.MetricKV{
			Key:   metric.GroupCol,
			Value: metric.CPUTime,
		}
		appMetrics = append(appMetrics, metricKV)
	}

	for _, metric := range metrics.UserMetrics {
		metricKV := &pb.MetricKV{
			Key:   metric.GroupCol,
			Value: metric.CPUTime,
		}
		userMetrics = append(userMetrics, metricKV)
	}

	return &pb.GetJobMetricResponse{
		AppMetrics:  appMetrics,
		UserMetrics: userMetrics,
	}
}

func ConvertJobCountMetric(metrics *dto.JobCountMetric) *pb.GetJobMetricResponse {
	if metrics == nil {
		return nil
	}

	return &pb.GetJobMetricResponse{
		AppMetrics:  convertJobQueryMetric2MetricKV(metrics.AppCountMetrics),
		UserMetrics: convertJobQueryMetric2MetricKV(metrics.UserCountMetrics),
	}
}

func ConvertJobWaitStatisticMetric(statistic *dto.JobWaitStatistic) *pb.GetJobWaitTimeStatisticResponse {
	if statistic == nil {
		return nil
	}

	return &pb.GetJobWaitTimeStatisticResponse{
		WaitNumStatisticTotal:  convertJobQueryMetric2MetricKV(statistic.JobWaitNumStatistic),
		WaitTimeStatisticAvg:   convertJobQueryMetric2MetricKV(statistic.JobWaitTimeStatisticAvg),
		WaitTimeStatisticMax:   convertJobQueryMetric2MetricKV(statistic.JobWaitTimeStatisticMax),
		WaitTimeStatisticTotal: convertJobQueryMetric2MetricKV(statistic.JobWaitTimeStatisticTotal),
	}
}

func convertJobQueryMetric2MetricKV(metrics []*dto.JobQueryResultMetric) []*pb.MetricKV {
	if len(metrics) == 0 {
		return []*pb.MetricKV{}
	}

	results := make([]*pb.MetricKV, 0, len(metrics))

	for _, metric := range metrics {
		metricKV := &pb.MetricKV{
			Key:   metric.Item,
			Value: metric.Count,
		}
		results = append(results, metricKV)
	}
	return results
}

func ConvertJobExitCode(exitCode string) string {
	if exitCode != "" {
		codes := strings.Split(exitCode, ":")
		if len(codes) == 2 {
			return codes[1]
		}
	}

	return exitCode
}

func ConvertWorkDir(workDir string, includeUsername bool) string {
	if strutil.IsEmpty(workDir) {
		return ""
	}

	var newWorkDir string
	if strings.HasPrefix(workDir, "http") {
		parsedURL, err := url.Parse(workDir)
		if err != nil {
			return workDir
		}

		pathSegments := strings.Split(parsedURL.Path, "/")

		if len(pathSegments) >= 3 {
			newPathSegments := pathSegments[3:]
			newWorkDir = strings.Join(newPathSegments, "/")
		}

		if includeUsername && len(pathSegments) >= 2 {
			newPathSegments := pathSegments[2:]
			newWorkDir = strings.Join(newPathSegments, "/")
		}
	} else {
		pathSegments := strings.Split(workDir, "/")

		if len(pathSegments) >= 2 {
			newPathSegments := pathSegments[2:]
			newWorkDir = strings.Join(newPathSegments, "/")
		}

		if includeUsername && len(pathSegments) >= 1 {
			newPathSegments := pathSegments[1:]
			newWorkDir = strings.Join(newPathSegments, "/")
		}
	}

	if strutil.IsEmpty(newWorkDir) {
		return workDir
	} else {
		return newWorkDir
	}
}

func ConvertJobResidual(residual *schema.Residual) *dto.JobResidualResponse {
	if residual == nil {
		return &dto.JobResidualResponse{}
	}

	varItems := make([]*dto.VarItem, 0, len(residual.Vars))
	for _, v := range residual.Vars {
		varItems = append(varItems, &dto.VarItem{Name: v.Name, Values: v.Values})
	}

	availableXvar := residual.AvailableXvar
	if availableXvar == nil {
		availableXvar = make([]string, 0)
	}

	return &dto.JobResidualResponse{
		AvailableXvar: availableXvar,
		Vars:          varItems,
	}
}

func SetVisAnalysisValue(enableResidual, enableSnapshot bool) map[string]bool {
	visAnalysisMap := make(map[string]bool)

	visAnalysisMap[consts.JobVisAnalysisResidual] = enableResidual
	visAnalysisMap[consts.JobVisAnalysisSnapshot] = enableSnapshot

	return visAnalysisMap
}
