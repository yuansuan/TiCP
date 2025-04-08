package impl

import (
	"context"
	"encoding/csv"
	"sort"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/PSP/psp/internal/common"
	"github.com/yuansuan/ticp/PSP/psp/internal/visual/consts"
	"github.com/yuansuan/ticp/PSP/psp/internal/visual/dao"
	"github.com/yuansuan/ticp/PSP/psp/internal/visual/dao/model"
	"github.com/yuansuan/ticp/PSP/psp/internal/visual/dto"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/csvutil"
)

func (s *VisualService) DurationStatistic(ctx context.Context, appIDs []string, startTime, endTime string) ([]*dto.DurationStatistic, error) {

	appIds := make([]snowflake.ID, 0)
	for _, v := range appIDs {
		appIds = append(appIds, snowflake.MustParseString(v))
	}
	statisticsData, err := s.sessionDao.DurationStatistics(ctx, appIds, startTime, endTime)
	if err != nil {
		return nil, err
	}

	statistics := make([]*dto.DurationStatistic, 0)
	for _, v := range statisticsData {
		statistic := &dto.DurationStatistic{
			AppID:    v.SoftwareID.String(),
			AppName:  v.SoftwareName,
			Duration: strconv.FormatFloat(float64(v.Duration)/float64(consts.UnixMilliToHour), 'f', 6, 64),
		}
		statistics = append(statistics, statistic)
	}
	return statistics, nil
}

func (s *VisualService) ListHistoryDuration(ctx context.Context, appIDs []string, startTime, endTime string, pageIndex, pageSize int) ([]*dto.HistoryDuration, int64, error) {

	total, historiesData, err := s.sessionDao.ListHistory(ctx, appIDs, startTime, endTime, pageIndex, pageSize)
	if err != nil {
		return nil, 0, err
	}

	list := make([]*dto.HistoryDuration, 0)
	for _, v := range historiesData {
		history := &dto.HistoryDuration{
			ID:       v.ID.String(),
			AppName:  v.SoftwareName,
			Duration: strconv.FormatFloat(float64(v.Duration)/float64(consts.UnixMilliToHour), 'f', 6, 64),
		}
		if v.StartTime != nil {
			history.StartTime = *v.StartTime
		}
		if v.EndTime != nil {
			history.EndTime = *v.EndTime
		}
		list = append(list, history)
	}
	return list, total, nil
}

func (s *VisualService) SessionUsageDurationStatistic(ctx context.Context, startTime, endTime int64) (*dto.SessionUsageDurationStatisticResponse, error) {
	logger := logging.GetLogger(ctx)

	originDataForSoftware, err := s.sessionDao.SessionStatistics(ctx, startTime, endTime, consts.ReportTypeSessionUsageDuration, consts.DimensionTypeSoftware)
	if err != nil {
		logger.Errorf("get session usage duration statistic for software failed, err: %v", err)
		return nil, err
	}

	originDataForUser, err := s.sessionDao.SessionStatistics(ctx, startTime, endTime, consts.ReportTypeSessionUsageDuration, consts.DimensionTypeUser)
	if err != nil {
		logger.Errorf("get session usage duration statistic for user failed, err: %v", err)
		return nil, err
	}

	return &dto.SessionUsageDurationStatisticResponse{
		UsageDurationBySoftwre: &dto.OriginStatisticData{
			Name:         "使用时长",
			OriginalData: convertStatisticItem(originDataForSoftware),
		},
		UsageDurationByUser: &dto.OriginStatisticData{
			Name:         "使用时长",
			OriginalData: convertStatisticItem(originDataForUser),
		},
	}, nil
}

func (s *VisualService) ExportUsageDurationStatistic(ctx *gin.Context, startTime, endTime int64) error {
	logger := logging.GetLogger(ctx)

	start := time.Now()
	err := csvutil.ExportCSVFilesToZip(ctx, "3D可视化-会话使用时长统计", []*csvutil.ExportCSVFileInfo{
		{
			CSVFileName: "会话使用时长统计-软件维度",
			CSVHeaders:  []string{"软件名称", "使用时长(小时)"},
			FillCSVData: func(w *csv.Writer) error {
				originDataForSoftware, err := s.sessionDao.SessionStatistics(ctx, startTime, endTime, consts.ReportTypeSessionUsageDuration, consts.DimensionTypeSoftware)
				if err != nil {
					logger.Errorf("get session usage duration statistic for software failed, err: %v", err)
					return err
				}

				sort.Slice(originDataForSoftware, func(i, j int) bool {
					return originDataForSoftware[i].Value > originDataForSoftware[j].Value
				})

				for _, v := range originDataForSoftware {
					rowData := make([]string, 0)
					rowData = append(rowData, csvutil.CSVContentWithTab(v.Key))
					rowData = append(rowData, csvutil.CSVContentWithTab(strconv.FormatFloat(v.Value, 'f', 2, 64)))
					_ = w.Write(rowData)
				}
				w.Flush()

				return nil
			},
		},
		{
			CSVFileName: "会话使用时长统计-用户维度",
			CSVHeaders:  []string{"用户名称", "使用时长(小时)"},
			FillCSVData: func(w *csv.Writer) error {
				originDataForUser, err := s.sessionDao.SessionStatistics(ctx, startTime, endTime, consts.ReportTypeSessionUsageDuration, consts.DimensionTypeUser)
				if err != nil {
					logger.Errorf("get session usage duration statistic for user failed, err: %v", err)
					return err
				}

				sort.Slice(originDataForUser, func(i, j int) bool {
					return originDataForUser[i].Value > originDataForUser[j].Value
				})

				for _, v := range originDataForUser {
					rowData := make([]string, 0)
					rowData = append(rowData, csvutil.CSVContentWithTab(v.Key))
					rowData = append(rowData, csvutil.CSVContentWithTab(strconv.FormatFloat(v.Value, 'f', 2, 64)))
					_ = w.Write(rowData)
				}
				w.Flush()

				return nil
			},
		},
	})
	if err != nil {
		return err
	}

	logging.GetLogger(ctx).Infof("export session usage duration statistic cost time: %v", time.Since(start))

	return nil
}

func (s *VisualService) SessionCreateNumberStatistic(ctx context.Context, startTime, endTime int64) (*dto.SessionCreateNumberStatisticResponse, error) {
	logger := logging.GetLogger(ctx)

	originDataForSoftware, err := s.sessionDao.SessionStatistics(ctx, startTime, endTime, consts.ReportTypeSessionCreateNumber, consts.DimensionTypeSoftware)
	if err != nil {
		logger.Errorf("get session create number statistic for software failed, err: %v", err)
		return nil, err
	}

	originDataForUser, err := s.sessionDao.SessionStatistics(ctx, startTime, endTime, consts.ReportTypeSessionCreateNumber, consts.DimensionTypeUser)
	if err != nil {
		logger.Errorf("get session create number statistic for user failed, err: %v", err)
		return nil, err
	}

	return &dto.SessionCreateNumberStatisticResponse{
		CreateNumberBySoftwre: &dto.OriginStatisticData{
			Name:         "创建会话数",
			OriginalData: convertStatisticItem(originDataForSoftware),
		},
		CreateNumberByUser: &dto.OriginStatisticData{
			Name:         "创建会话数",
			OriginalData: convertStatisticItem(originDataForUser),
		},
	}, nil
}

func convertStatisticItem(data []*dao.StatisticItem) []*dto.StatisticItem {
	list := make([]*dto.StatisticItem, 0)

	for _, v := range data {
		item := &dto.StatisticItem{
			Key:   v.Key,
			Value: v.Value,
		}
		list = append(list, item)
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].Value > list[j].Value
	})

	return list
}

func (s *VisualService) SessionNumberStatusStatistic(ctx context.Context, startTime, endTime int64) (*dto.SessionNumberStatusStatisticResponse, error) {
	logger := logging.GetLogger(ctx)

	currentTime := time.Now()
	start, end := time.UnixMilli(startTime), time.UnixMilli(endTime)

	createdRecordMap := make(map[string]int64)
	runningRecordMap := make(map[string]int64)

	pageIndex, pageSize := common.DefaultPageIndex, common.DefaultMaxPageSize
	for {
		exportData, _, err := s.sessionDao.ListSession(ctx, "", nil, nil, nil, nil, startTime, endTime, pageIndex, pageSize)
		if err != nil {
			return nil, err
		}

		if len(exportData) == 0 {
			break
		}

		for _, v := range exportData {
			markSessionOperator(v, createdRecordMap, runningRecordMap, currentTime, start, end)
		}

		if len(exportData) < common.DefaultMaxPageSize {
			break
		}

		pageIndex++
	}

	statistic := &dto.SessionNumberStatusStatisticResponse{
		NumberStatus: []*dto.OriginStatisticDatas{
			{
				Name:         "创建会话数",
				OriginalData: make([]*dto.StatisticItems, 0),
			},
			{
				Name:         "运行会话数",
				OriginalData: make([]*dto.StatisticItems, 0),
			},
		},
	}
	if len(createdRecordMap) == 0 && len(runningRecordMap) == 0 {
		return statistic, nil
	}

	createdData, runningData := resolveMarkedData(createdRecordMap, runningRecordMap, startTime, endTime)
	statistic.NumberStatus[0].OriginalData = createdData
	statistic.NumberStatus[1].OriginalData = runningData

	logger.Infof("get session number status statistic cost time: %v", time.Since(currentTime))

	return statistic, nil
}

func markSessionOperator(v *model.Session, createdRecordMap, runningRecordMap map[string]int64, currentTime time.Time, startTime, endTime time.Time) {
	if v.CreateTime.Before(startTime) || v.CreateTime.After(endTime) {
		return
	}

	createdRecordMap[getCursorTime(v.CreateTime).Format(common.DatetimeFormatToHour)]++

	start, end := v.StartTime, v.EndTime
	if end.IsZero() {
		end = currentTime
	}

	for start.Before(end) {
		runningRecordMap[getCursorTime(start).Format(common.DatetimeFormatToHour)]++
		start = start.Add(time.Hour)
	}
}

func resolveMarkedData(createdRecordMap, runningRecordMap map[string]int64, startTime, endTime int64) ([]*dto.StatisticItems, []*dto.StatisticItems) {
	createdData, runningData := make([]*dto.StatisticItems, 0), make([]*dto.StatisticItems, 0)

	start, end := time.UnixMilli(startTime), time.UnixMilli(endTime)
	for start.Before(end) {
		cursor := getCursorTime(start)
		createdData = append(createdData, &dto.StatisticItems{
			Key:   cursor.UnixMilli(),
			Value: float64(createdRecordMap[cursor.Format(common.DatetimeFormatToHour)]),
		})
		runningData = append(runningData, &dto.StatisticItems{
			Key:   cursor.UnixMilli(),
			Value: float64(runningRecordMap[cursor.Format(common.DatetimeFormatToHour)]),
		})

		start = start.Add(time.Hour)
	}

	return createdData, runningData
}

func getCursorTime(ctime time.Time) time.Time {
	cursor := ctime.Truncate(time.Hour)
	if ctime.Add(-time.Hour).Before(ctime) {
		cursor = cursor.Add(time.Hour)
	}
	return cursor
}
