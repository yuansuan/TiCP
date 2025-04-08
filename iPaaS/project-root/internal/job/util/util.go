// Copyright (C) 2018 LambdaCal Inc.

package util

import (
	"fmt"
	"math"
	"regexp"
	"strings"
	"time"

	"github.com/yuansuan/ticp/common/project-root-api/job/v1/admin/app/update"
	v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/consts"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao/models"
)

var (
	DefaultLocation, _ = time.LoadLocation("Asia/Shanghai")
)

const (
	UserIdKeyInHeader = "x-ys-user-id"
)

// CorrectionFunc to correct the cores and coresPerNode
// return the corrected cores and coresPerNode
type CorrectionFunc func(cores, coresPerNode int64) (int64, int64)

// AllocNodes use coresPerNode and cores to calculate the number of nodes
// cores: total cores
// coresPerNode: cores per node
// math.Round函数: 将一个浮点数 四舍五入 到最接近的整数值
// round: whether to round, if true, use math.Round, else use math.Floor
// correctionFuncs: correction functions, to correct the cores and coresPerNode, will be applied before calculation
// return: coresPerNodeCorrected, (nodes * coresPerNodeCorrected) as coresAllocated, error
func AllocNodes(cores, coresPerNode int64, round bool, correctionFuncs ...CorrectionFunc) (int64, int64, error) {
	// Apply correction functions
	for _, correctionFunc := range correctionFuncs {
		cores, coresPerNode = correctionFunc(cores, coresPerNode)
	}

	if coresPerNode == 0 {
		return 0, 0, fmt.Errorf("coresPerNode is zero!")
	}

	if cores < 0 {
		return 0, 0, fmt.Errorf("cores not positive!")
	}

	if cores == 0 {
		return 0, 0, nil
	}

	nodes := cores / coresPerNode
	if nodes == 0 {
		nodes = 1
	} else {
		nodesF := float64(cores) / float64(coresPerNode)
		if round {
			nodes = int64(math.Round(nodesF))
		} else {
			nodes = int64(math.Floor(nodesF))
		}
	}

	return coresPerNode, nodes * coresPerNode, nil
}

// WithSharedNode 共享节点时，单节点核数不进行取整
// e.g. cores=20, coresPerNode=24.
// WithSharedNode will return 20, 20.
// This value will be passed to the standard-compute job submission parameters.
// Standard-compute will use these two parameters to generate two variables,`nTasksPerNode` and `nodes`,
// 翻译：sc会用这2个参数生成2个变量：`nTasksPerNode` and `nodes`
// through the EnsureNTaskPerNode and OccupiedNodesNum functions,respectively.

// These variables are then passed as environment variables to SLURM(ntasks_per_node=20 and nodes=1) or PBSPro(number_of_cpu=20 and nodes=1) for job submission.
// The scheduler determines how many nodes and cores a job occupies based on these two values.
func WithSharedNode(cores, coresPerNode int64) (int64, int64) {
	if cores < coresPerNode {
		// 单节点时，不进行取整，即coresPerNode=cores
		return cores, cores
	}
	return cores, coresPerNode
}

// CoresRange 核数范围, 用户期望的最小和最大核数
type CoresRange struct {
	MinExpectedCores int64
	MaxExpectedCores int64
}

// CalculateResourceUsage 计算资源用量
func CalculateResourceUsage(coresRange CoresRange, coresPerNode, availableCores int64, shared bool) (int64, int64, error) {
	correctionFuncs := []CorrectionFunc{}
	if shared { // 共享节点, 单节点核数不进行取整
		correctionFuncs = append(correctionFuncs, WithSharedNode)
	}

	// 最大期望核数, AllocNodes with Round, 可用核数四舍五入
	perOfMax, maxExpectedCores, err := AllocNodes(coresRange.MaxExpectedCores, coresPerNode, true, correctionFuncs...)
	if err != nil {
		return 0, 0, err
	}

	// 最大可用核数, AllocNodes with Floor, 可用核数不能四舍五入, 向下取整
	perOfAva, maxAvailableCores, err := AllocNodes(availableCores, coresPerNode, false, correctionFuncs...)
	if err != nil {
		return 0, 0, err
	}

	// 最小期望核数, AllocNodes with Round, 可用核数四舍五入
	perOfMin, minExpectedCores, err := AllocNodes(coresRange.MinExpectedCores, coresPerNode, true, correctionFuncs...)
	if err != nil {
		return 0, 0, err
	}

	// 根据最大可用核数和期望核数选择合适的值
	// max < ava, => max
	if maxAvailableCores > maxExpectedCores {
		return perOfMax, maxExpectedCores, nil
	}
	// min < ava < max, => ava
	if maxAvailableCores > minExpectedCores {
		return perOfAva, maxAvailableCores, nil
	}
	// ava < min, => min, will be filtered
	return perOfMin, minExpectedCores, nil
}

// AddAppImagePrefix 增加app image前缀
func AddAppImagePrefix(appImage string, isLocal bool) string {
	if isLocal {
		return consts.LocalImagePrefix + appImage
	}

	return consts.AppImagePrefix + appImage
}

// ParseYsID 解析路径中的ys_id
func ParseYsID(path string) string {
	re := regexp.MustCompile(`(?m)http[s]{0,1}:\/\/[^\/]+\/([^\/]+)`)
	match := re.FindStringSubmatch(path)
	if len(match) > 1 {
		return match[1]
	}
	return ""
}

// ParsePath 解析路径中的ysid后的路径部分
func ParsePath(path string) string {
	re := regexp.MustCompile(`(?m)http[s]{0,1}:\/\/[^\/]+\/[^\/]+(.*)`)
	match := re.FindStringSubmatch(path)
	if len(match) > 1 {
		return match[1]
	}
	return ""
}

// AddPrefixSlash 为路径增加'/'前缀
func AddPrefixSlash(path string) string {
	if strings.HasPrefix(path, "/") {
		return path
	}

	return "/" + path
}

// AddSuffixSlash 为路径增加'/'后缀
func AddSuffixSlash(path string) string {
	if strings.HasSuffix(path, "/") {
		return path
	}

	return path + "/"
}

// ParseYsIDWithOutDomain 解析不带域名路径中的ys_id
func ParseYsIDWithOutDomain(path string) string {
	path = AddPrefixSlash(path) // 这样比较好解析
	re := regexp.MustCompile(`(?m)^\/([^\/]+)`)
	match := re.FindStringSubmatch(path)
	if len(match) > 1 {
		return match[1]
	}
	return ""
}

// ParsePathWithOutDomain 解析不带域名路径中的ysid后的路径部分
func ParsePathWithOutDomain(path string) string {
	path = AddPrefixSlash(path) // 这样比较好解析
	re := regexp.MustCompile(`(?m)^\/[^\/]+(.*)`)
	match := re.FindStringSubmatch(path)
	if len(match) > 1 {
		return match[1]
	}
	return ""
}

// ZoneResource 分区资源
type ZoneResource struct {
	Queue         string
	Zone          string
	IsDefault     bool // 是否是该分区的默认队列
	CPU           int64
	Mem           int64
	CoresPerNode  int64
	ReservedCores int64 // 预留核数
	FilterReason  string
	TotalNodeNum  int64 //队列的总节点数
}

// AddFilterReason 增加过滤原因
func (z *ZoneResource) AddFilterReason(reason string) {
	if z.FilterReason != "" {
		z.FilterReason += "; "
	}
	z.FilterReason += reason
}

// PrintFilterReason 打印过滤原因
func (z *ZoneResource) PrintFilterReason() string {
	if z.FilterReason == "" {
		return ""
	}

	return fmt.Sprintf("[%s:%s] %s", z.Zone, z.Queue, z.FilterReason)
}

func EmptyTime(t time.Time) bool {
	return t == time.Time{}
}

func IsAccountInArrears(accountDetail *v20230530.AccountDetail) bool {
	return accountDetail.IsOverdrawn || accountDetail.AccountBalance+accountDetail.CreditQuotaAmount < 0
}

func IsAccountFrozen(accountDetail *v20230530.AccountDetail) bool {
	return accountDetail.FrozenStatus
}

func IsAppPublished(app *models.Application) bool {
	return app.PublishStatus == string(update.Published)
}

// IsNonZeroExitCode 是否退出码非零
func IsNonZeroExitCode(exitCode string) (bool, error) {
	parts := strings.Split(exitCode, ":")
	if len(parts) != 2 {
		return false, fmt.Errorf("invalid exit code: %s", exitCode)
	}

	exitStatus := parts[0]
	return exitStatus != "0", nil
}

func ParseTime(t, format string) (time.Time, error) {
	return time.ParseInLocation(format, t, DefaultLocation)
}
