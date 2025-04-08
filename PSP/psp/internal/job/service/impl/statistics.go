package impl

import (
	"compress/gzip"
	"context"
	"encoding/csv"
	"fmt"
	"net/url"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"google.golang.org/protobuf/types/known/timestamppb"

	mainconfig "github.com/yuansuan/ticp/PSP/psp/cmd/config"
	"github.com/yuansuan/ticp/PSP/psp/internal/common"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/project"
	"github.com/yuansuan/ticp/PSP/psp/internal/job/consts"
	"github.com/yuansuan/ticp/PSP/psp/internal/job/dto"
	"github.com/yuansuan/ticp/PSP/psp/internal/job/util"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/csvutil"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/floatutil"
)

func (s *jobServiceImpl) GetTop5ProjectInfo(ctx context.Context, start, end int64) (*dto.GetTop5ProjectInfoResponse, error) {
	logger := logging.GetLogger(ctx)
	runningProjectIds, err := s.rpc.Project.GetRunningProjectIdsByTime(ctx, &project.GetRunningProjectIdsByTimeRequest{
		TimePoint: timestamppb.New(time.UnixMilli(start)),
	})
	if err != nil {
		logger.Errorf("get running project ids by time: [%v] err: %v", start, err)
		return nil, err
	}

	runningProjectSnowflakeIds := make([]snowflake.ID, 0, len(runningProjectIds.ProjectIds))
	for _, v := range runningProjectIds.ProjectIds {
		runningProjectSnowflakeIds = append(runningProjectSnowflakeIds, snowflake.ID(v))
	}

	top5ProjectInfos, err := s.jobDao.GetTop5ProjectByCPUTime(ctx, runningProjectSnowflakeIds, start, end)
	if err != nil {
		logger.Errorf("get top 5 project by cpu time: [%v, %v] err: %v", start, end, err)
		return nil, err
	}
	top5ProjectIdsInID := make([]snowflake.ID, 0, len(top5ProjectInfos))
	top5ProjectIdsInStr := make([]string, 0, len(top5ProjectInfos))
	for _, v := range top5ProjectInfos {
		top5ProjectIdsInID = append(top5ProjectIdsInID, v.ProjectId)
		top5ProjectIdsInStr = append(top5ProjectIdsInStr, v.ProjectId.String())
	}

	projectJobCountMap := make(map[snowflake.ID]int64)
	if len(top5ProjectIdsInID) > 0 {
		projectJobCount, err := s.jobDao.GetJobCountByProjectIds(ctx, top5ProjectIdsInID, start, end)
		if err != nil {
			logger.Errorf("get job count by project ids: [%+v] err: %v", top5ProjectIdsInID, err)
			return nil, err
		}
		for _, v := range projectJobCount {
			projectJobCountMap[v.ProjectId] = v.Count
		}
	}

	projectInfosMap := make(map[snowflake.ID]*project.GetProjectByIdResponse)
	if len(top5ProjectIdsInStr) > 0 {
		projectInfos, err := s.rpc.Project.GetProjectsDetailByIds(ctx, &project.GetProjectsDetailByIdsRequest{
			ProjectIds:         top5ProjectIdsInStr,
			IncludeMemberCount: true,
		})
		if err != nil {
			logger.Errorf("get projects detail by ids: [%+v] err: %v", top5ProjectIdsInStr, err)
			return nil, err
		}
		for _, v := range projectInfos.Projects {
			projectInfosMap[snowflake.MustParseString(v.ProjectId)] = v
		}
	}

	response := &dto.GetTop5ProjectInfoResponse{
		Projects: make([]string, 0),
		Users:    make([]int64, 0),
		Jobs:     make([]int64, 0),
		CpuTimes: make([]float64, 0),
	}

	if len(top5ProjectIdsInID) > 0 {
		for _, v := range top5ProjectInfos {
			response.Projects = append(response.Projects, v.ProjectName)
			if v, ok := projectInfosMap[v.ProjectId]; ok {
				response.Users = append(response.Users, v.MemberCount)
			} else {
				response.Users = append(response.Users, 0)
			}
			response.Jobs = append(response.Jobs, projectJobCountMap[v.ProjectId])
			response.CpuTimes = append(response.CpuTimes, v.CPUTime)
		}
	}
	if len(response.Projects) < 5 {
		runningProjectIdsInStr := make([]string, 0, len(runningProjectSnowflakeIds))
		for _, v := range runningProjectSnowflakeIds {
			runningProjectIdsInStr = append(runningProjectIdsInStr, v.String())
		}
		if len(runningProjectIdsInStr) > 0 {
			projectInfosRes, err := s.rpc.Project.GetProjectsDetailByIds(ctx, &project.GetProjectsDetailByIdsRequest{
				ProjectIds:         runningProjectIdsInStr,
				IncludeMemberCount: true,
			})
			if err != nil {
				logger.Errorf("get projects detail by ids: [%+v] err: %v", top5ProjectIdsInStr, err)
				return nil, err
			}

			projectInfos := projectInfosRes.Projects
			sort.Slice(projectInfos, func(i, j int) bool {
				return projectInfos[i].MemberCount > projectInfos[j].MemberCount
			})

			projectNameMap := make(map[string]bool, len(projectInfos))
			for _, v := range response.Projects {
				projectNameMap[v] = true
			}
			for _, v := range projectInfos {
				if _, ok := projectNameMap[v.ProjectName]; ok {
					continue
				}
				if len(response.Projects) == 5 {
					break
				}
				response.Projects = append(response.Projects, v.ProjectName)
				response.Users = append(response.Users, v.MemberCount)
				response.Jobs = append(response.Jobs, 0)
				response.CpuTimes = append(response.CpuTimes, 0)
			}
		}
	}

	return response, nil
}

func (s *jobServiceImpl) GetJobCPUTimeTotal(ctx context.Context, queryType, computeType string, names, projectIds []string, startTime, endTime int64) (float64, error) {
	cpuTimeTotal, err := s.jobDao.GetJobCPUTimeTotal(ctx, queryType, computeType, names, projectIds, startTime, endTime)
	if err != nil {
		return 0.0, err
	}

	return cpuTimeTotal, nil
}

func (s *jobServiceImpl) GetJobStatisticsOverview(ctx context.Context, queryType, computeType string, names, projectIds []string, startTime, endTime int64, pageIndex, pageSize int) ([]*dto.StatisticsOverview, int64, error) {
	overviews, total, err := s.jobDao.GetJobStatisticsOverview(ctx, queryType, computeType, names, projectIds, startTime, endTime, pageIndex, pageSize)
	if err != nil {
		return nil, 0, err
	}

	ids := s.sid.Generates(int64(len(overviews)))
	statisticsOverviews := make([]*dto.StatisticsOverview, 0, len(overviews))
	for i, v := range overviews {
		overview := util.ConvertStatisticsJobToDTOOverview(v, ids[i].String(), queryType)
		statisticsOverviews = append(statisticsOverviews, overview)
	}

	return statisticsOverviews, total, nil
}

func (s *jobServiceImpl) GetJobStatisticsDetail(ctx context.Context, queryType, computeType string, names, projectIds []string, startTime, endTime int64, pageIndex, pageSize int) ([]*dto.JobDetailInfo, int64, error) {
	statisticsJobList, total, err := s.jobDao.GetJobStatisticsDetail(ctx, queryType, computeType, names, projectIds, startTime, endTime, pageIndex, pageSize)
	if err != nil {
		return nil, 0, err
	}

	jobDetails := make([]*dto.JobDetailInfo, 0, len(statisticsJobList))
	for _, v := range statisticsJobList {
		jobDetail := util.ConvertStatisticsJobToDTODetail(v)
		if jobDetail != nil {
			jobDetails = append(jobDetails, jobDetail)
		}
	}

	return jobDetails, total, nil
}

func (s *jobServiceImpl) GetJobStatisticsExport(ctx *gin.Context, queryType, computeType, showType string, names, projectIds []string, startTime, endTime int64) error {

	start := time.Now()

	csvName, csvHeaders, err := getCSVNameAndHeaders(queryType, showType)
	if err != nil {
		return err
	}

	exportFileName := fmt.Sprintf("%v-%v", url.QueryEscape(csvName), time.Now().Format(common.DateOnly))
	disposition := fmt.Sprintf("attachment; filename=%s.csv.gz", exportFileName)
	ctx.Header("Content-Type", "application/x-gzip")
	ctx.Header("Content-Disposition", disposition)

	gzWriter := gzip.NewWriter(ctx.Writer)
	defer gzWriter.Close()

	// Write BOM directly to the gzip writer
	bom := []byte{0xEF, 0xBB, 0xBF}
	_, err = gzWriter.Write(bom)
	if err != nil {
		return err
	}

	w := csv.NewWriter(gzWriter)
	defer w.Flush()

	_ = w.Write(csvHeaders)

	pageIndex, pageSize := common.DefaultPageIndex, common.CSVExportNumber
	competeTypeNameMap := mainconfig.Custom.Main.ComputeTypeNames

	switch showType {
	case consts.JobStatisticsShowTypeOverview:
		for {
			exportData, _, err := s.jobDao.GetJobStatisticsOverview(ctx, queryType, computeType, names, projectIds, startTime, endTime, pageIndex, pageSize)
			if err != nil {
				return err
			}

			if len(exportData) == 0 {
				break
			}

			for _, jobInfo := range exportData {
				rowData := make([]string, 0, len(csvHeaders))

				id, name := "", ""
				switch queryType {
				case consts.JobStatisticsQueryTypeApp:
					id = jobInfo.AppId.String()
					name = jobInfo.AppName
				case consts.JobStatisticsQueryTypeUser:
					id = jobInfo.UserId.String()
					name = jobInfo.UserName
				}

				rowData = append(rowData, id)
				rowData = append(rowData, name)
				rowData = append(rowData, getComputeTypeName(competeTypeNameMap, jobInfo.Type))
				rowData = append(rowData, jobInfo.ProjectName)
				rowData = append(rowData, csvutil.CSVContentWithTab(floatutil.NumberToFloatStr(jobInfo.CPUTime, common.DecimalPlaces)))

				_ = w.Write(rowData)
			}
			w.Flush()

			if len(exportData) < common.CSVExportNumber {
				break
			}

			pageIndex++
		}
	case consts.JobStatisticsShowTypeDetail:
		for {
			exportData, _, err := s.jobDao.GetJobStatisticsDetail(ctx, queryType, computeType, names, projectIds, startTime, endTime, pageIndex, pageSize)
			if err != nil {
				return err
			}

			if len(exportData) == 0 {
				break
			}

			for _, jobInfo := range exportData {
				rowData := make([]string, 0, len(csvHeaders))

				rowData = append(rowData, jobInfo.Id.String())
				rowData = append(rowData, jobInfo.Name)
				rowData = append(rowData, getComputeTypeName(competeTypeNameMap, jobInfo.Type))
				rowData = append(rowData, jobInfo.ProjectName)
				rowData = append(rowData, jobInfo.AppName)
				rowData = append(rowData, jobInfo.UserName)
				rowData = append(rowData, csvutil.CSVFormatTime(jobInfo.SubmitTime, common.DatetimeFormat, "--"))
				rowData = append(rowData, csvutil.CSVFormatTime(jobInfo.StartTime, common.DatetimeFormat, "--"))
				rowData = append(rowData, csvutil.CSVFormatTime(jobInfo.EndTime, common.DatetimeFormat, "--"))
				rowData = append(rowData, csvutil.CSVContentWithTab(floatutil.NumberToFloatStr(jobInfo.CPUTime, common.DecimalPlaces)))

				_ = w.Write(rowData)
			}
			w.Flush()

			if len(exportData) < common.CSVExportNumber {
				break
			}

			pageIndex++
		}
	default:
		return fmt.Errorf("when export data, show type not match")
	}

	end := time.Now()
	logging.GetLogger(ctx).Infof("export job statistics data using time: %fs", end.Sub(start).Seconds())

	return nil
}

func getCSVNameAndHeaders(queryType string, showType string) (string, []string, error) {
	csvName := ""
	csvHeaders := make([]string, 0)

	switch {
	case queryType == consts.JobStatisticsQueryTypeApp && showType == consts.JobStatisticsShowTypeOverview:
		csvName = "作业统计-应用统计总览"
		csvHeaders = []string{"应用编号", "应用名称", "计算类型", "项目名称", "核时(小时)"}
	case queryType == consts.JobStatisticsQueryTypeApp && showType == consts.JobStatisticsShowTypeDetail:
		csvName = "作业统计-应用统计详情"
		csvHeaders = []string{"作业编号", "作业名称", "计算类型", "项目名称", "应用名称", "用户名称", "提交时间", "开始时间", "结束时间", "核时(小时)"}
	case queryType == consts.JobStatisticsQueryTypeUser && showType == consts.JobStatisticsShowTypeOverview:
		csvName = "作业统计-用户统计总览"
		csvHeaders = []string{"用户编号", "用户名称", "计算类型", "项目名称", "核时(小时)"}
	case queryType == consts.JobStatisticsQueryTypeUser && showType == consts.JobStatisticsShowTypeDetail:
		csvName = "作业统计-用户统计详情"
		csvHeaders = []string{"作业编号", "作业名称", "计算类型", "项目名称", "应用名称", "用户名称", "提交时间", "开始时间", "结束时间", "核时(小时)"}
	default:
		return "", nil, fmt.Errorf("query type and show type combine not match")
	}

	return csvName, csvHeaders, nil
}

func getComputeTypeName(computeTypeNameMap map[string]string, computeType string) string {
	vName := fmt.Sprintf("[未配置]%v", computeType)
	if name, ok := computeTypeNameMap[computeType]; ok {
		vName = name
	}
	return vName
}
