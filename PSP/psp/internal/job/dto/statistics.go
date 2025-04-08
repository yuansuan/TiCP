package dto

import "github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"

type JobComputeTypeListRequest struct{}

type ComputeTypeName struct {
	ComputeType string `json:"compute_type"` // 计算类型
	ShowName    string `json:"show_name"`    // 展示名称
}

type JobComputeTypeListResponse struct {
	ComputeTypes []*ComputeTypeName `json:"compute_types"` // 计算类型列表
}

type JobSetNameListRequest struct {
	ComputeType string `json:"compute_type" form:"compute_type" enums:"local,cloud"` // 计算类型
}

type JobSetNameListResponse struct {
	JobSetNames []string `json:"job_set_names"` // 作业集名称列表
}

type JobAppNameListRequest struct {
	ComputeType string `json:"compute_type" form:"compute_type" enums:"local,cloud"` // 计算类型
}

type JobAppNameListResponse struct{}

type JobUserNameListRequest struct {
	ComputeType string `json:"compute_type" form:"compute_type" enums:"local,cloud"` // 计算类型
}

type JobUserNameListResponse struct{}

type JobQueueNameListRequest struct {
	ComputeType string `json:"compute_type" form:"compute_type" enums:"local,cloud"` // 计算类型
}

type JobQueueNameListResponse struct{}

type GetJobStatisticsTotalCPUTimeRequest struct {
	QueryType   string   `json:"query_type" form:"query_type" enums:"app,user"`        // 统计维度
	ComputeType string   `json:"compute_type" form:"compute_type" enums:"local,cloud"` // 计算类型
	Names       []string `json:"names" form:"names"`                                   // 统计维度对应的值列表
	ProjectIDs  []string `json:"project_ids" form:"project_ids"`                       // 所属项目 ID 列表
	StartTime   int64    `json:"start_time" form:"start_time"`                         // 开始时间(按作业提交时间匹配: 提供秒级时间戳)
	EndTime     int64    `json:"end_time" form:"end_time"`                             // 结束时间(按作业提交时间匹配: 提供秒级时间戳)
}

type GetJobStatisticsTotalCPUTimeResponse struct {
	CPUTime string `json:"cpu_time"` // 核时(小时)
}

type GetJobStatisticsOverviewRequest struct {
	QueryType   string   `json:"query_type" form:"query_type" enums:"app,user"`        // 统计维度
	ComputeType string   `json:"compute_type" form:"compute_type" enums:"local,cloud"` // 计算类型
	Names       []string `json:"names" form:"names"`                                   // 统计维度对应的值列表
	ProjectIDs  []string `json:"project_ids" form:"project_ids"`                       // 所属项目 ID 列表
	StartTime   int64    `json:"start_time" form:"start_time"`                         // 开始时间(按作业提交时间匹配: 提供秒级时间戳)
	EndTime     int64    `json:"end_time" form:"end_time"`                             // 结束时间(按作业提交时间匹配: 提供秒级时间戳)
	PageIndex   int      `json:"page_index" form:"page_index"`                         // 分页索引
	PageSize    int      `json:"page_size" form:"page_size"`                           // 分页大小
}

type StatisticsOverview struct {
	Id          string `json:"id"`                              // 编号
	UId         string `json:"u_id"`                            // 唯一 Id
	Name        string `json:"name"`                            // 名称
	ComputeType string `json:"computeType" enums:"local,cloud"` // 计算类型
	ProjectName string `json:"project_name"`                    // 所属项目
	CPUTime     string `json:"cpu_time"`                        // 核时(小时)
}

type GetJobStatisticsOverviewResponse struct {
	Overviews []*StatisticsOverview `json:"overviews"`          // 总览列表
	Total     int64                 `json:"total" form:"total"` // 总数
}

type GetJobStatisticsDetailRequest struct {
	QueryType   string   `json:"query_type" form:"query_type" enums:"app,user"`        // 统计维度
	ComputeType string   `json:"compute_type" form:"compute_type" enums:"local,cloud"` // 计算类型
	Names       []string `json:"names" form:"names"`                                   // 统计维度对应的值列表
	ProjectIDs  []string `json:"project_ids" form:"project_ids"`                       // 所属项目 ID 列表
	StartTime   int64    `json:"start_time" form:"start_time"`                         // 开始时间(按作业提交时间匹配: 提供秒级时间戳)
	EndTime     int64    `json:"end_time" form:"end_time"`                             // 结束时间(按作业提交时间匹配: 提供秒级时间戳)
	PageIndex   int      `json:"page_index" form:"page_index"`                         // 分页索引
	PageSize    int      `json:"page_size" form:"page_size"`                           // 分页大小
}

type GetJobStatisticsDetailResponse struct {
	JobDetails []*JobDetailInfo `json:"job_details"`        // 作业详情
	Total      int64            `json:"total" form:"total"` // 总数
}

type GetTop5ProjectInfoRequest struct {
	Start int64 `json:"start" form:"start"` // 开始时间
	End   int64 `json:"end" form:"end"`     // 结束时间
}

type GetTop5ProjectInfoResponse struct {
	Projects []string  `json:"projects"` // 项目信息
	Users    []int64   `json:"users"`    // 用户信息
	Jobs     []int64   `json:"jobs"`     // 作业信息
	CpuTimes []float64 `json:"cputimes"` // 作业总核数信息
}

type GetJobStatisticsExportRequest struct {
	QueryType   string   `json:"query_type" form:"query_type" enums:"app,user"`        // 统计维度
	ComputeType string   `json:"compute_type" form:"compute_type" enums:"local,cloud"` // 计算类型
	Names       []string `json:"names" form:"names"`                                   // 统计维度对应的值列表
	ProjectIDs  []string `json:"project_ids" form:"project_ids"`                       // 所属项目 ID 列表
	StartTime   int64    `json:"start_time" form:"start_time"`                         // 开始时间(按作业提交时间匹配: 提供秒级时间戳)
	EndTime     int64    `json:"end_time" form:"end_time"`                             // 结束时间(按作业提交时间匹配: 提供秒级时间戳)
	ShowType    string   `json:"show_type" form:"show_type" enums:"overview,detail"`   // 对应页面展示类型
}

type GetJobStatisticsExportResponse struct{}

type ProjectCPUTime struct {
	ProjectId   snowflake.ID `json:"project_id" xorm:"project_id"`     // 项目 ID
	ProjectName string       `json:"project_name" xorm:"project_name"` // 项目名称
	CPUTime     float64      `json:"cpu_time" xorm:"cpu_time"`         // 核时(小时)
}

type ProjectJobCount struct {
	ProjectId   snowflake.ID `json:"project_id" xorm:"project_id"`     // 项目 ID
	ProjectName string       `json:"project_name" xorm:"project_name"` // 项目名称
	Count       int64        `json:"count" xorm:"count"`               // 作业数
}
