package rpc

import (
	"context"

	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/PSP/psp/internal/common"
	pb "github.com/yuansuan/ticp/PSP/psp/internal/common/proto/job"
	"github.com/yuansuan/ticp/PSP/psp/internal/job/dto"
	"github.com/yuansuan/ticp/PSP/psp/internal/job/util"
)

// GetCloudJobDetail 云作业详情
func (s *GRPCService) GetCloudJobDetail(ctx context.Context, req *pb.GetCloudJobDetailRequest) (*pb.GetCloudJobDetailResponse, error) {
	logger := logging.GetLogger(ctx)

	detail, err := s.JobService.GetJobDetailByOutID(ctx, req.JobId, common.Cloud)
	if err != nil {
		logger.Errorf("rpc get job detail err: %v", err)
		return nil, err
	}

	return &pb.GetCloudJobDetailResponse{
		Job: util.ConvertJob2GRPCDetail(detail),
	}, nil
}

func (s *GRPCService) GetJobCPUTimeMetric(ctx context.Context, req *pb.GetJobMetricRequest) (*pb.GetJobMetricResponse, error) {
	logger := logging.GetLogger(ctx)

	filer := &dto.JobMetricFiler{
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		TopSize:   int(req.TopSize),
	}
	runningTimeMetric, err := s.JobService.GetJobCPUTimeMetric(ctx, filer)
	if err != nil {
		logger.Errorf("rpc get job core running time metric err: %v", err)
		return nil, err
	}

	metric := util.ConvertJobCPUTimeMetric(runningTimeMetric)
	return metric, nil
}

func (s *GRPCService) GetJobCountMetric(ctx context.Context, req *pb.GetJobMetricRequest) (*pb.GetJobMetricResponse, error) {
	logger := logging.GetLogger(ctx)

	filer := &dto.JobMetricFiler{
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	}

	jobCountMetric, err := s.JobService.GetJobCountMetric(ctx, filer)
	if err != nil {
		logger.Errorf("rpc get job count running time metric err: %v", err)
		return nil, err
	}

	metric := util.ConvertJobCountMetric(jobCountMetric)
	return metric, nil
}

func (s *GRPCService) GetJobDeliverCount(ctx context.Context, req *pb.GetJobMetricRequest) (*pb.GetJobMetricResponse, error) {
	logger := logging.GetLogger(ctx)

	filer := &dto.JobMetricFiler{
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	}

	jobCountMetric, err := s.JobService.GetJobDeliverCount(ctx, filer)
	if err != nil {
		logger.Errorf("rpc get job count running time metric err: %v", err)
		return nil, err
	}

	metric := util.ConvertJobCountMetric(jobCountMetric)
	return metric, nil
}

func (s *GRPCService) GetJobWaitStatistic(ctx context.Context, req *pb.GetJobMetricRequest) (*pb.GetJobWaitTimeStatisticResponse, error) {
	logger := logging.GetLogger(ctx)

	filer := &dto.JobMetricFiler{
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	}

	jobWaitStatistic, err := s.JobService.GetJobWaitStatistic(ctx, filer)
	if err != nil {
		logger.Errorf("rpc get job wait statistic err: %v", err)
		return nil, err
	}

	metric := util.ConvertJobWaitStatisticMetric(jobWaitStatistic)
	return metric, nil
}

func (s *GRPCService) GetJobStatus(ctx context.Context, req *pb.GetJobStatusRequest) (*pb.GetJobStatusResponse, error) {
	logger := logging.GetLogger(ctx)

	jobStatusMap, err := s.JobService.GetJobStatus(ctx)
	if err != nil {
		logger.Errorf("rpc get job core running time metric err: %v", err)
		return nil, err
	}
	augmentedData(jobStatusMap)
	return &pb.GetJobStatusResponse{
		JobStatusMap: jobStatusMap,
	}, nil
}

// AugmentedData 补充作业状态信息
func augmentedData(jobStatusMap map[string]int64) {
	keys := util.MonitorJobStates
	if len(jobStatusMap) < 4 {
		for _, key := range keys {
			if _, ok := jobStatusMap[key]; !ok {
				jobStatusMap[key] = 0
			}
		}
	}
}
