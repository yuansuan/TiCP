package impl

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pkg/errors"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	pb "github.com/yuansuan/ticp/common/project-root-api/proto/license"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/license/collector"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/license/dao"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/license/dao/models"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/license/rpc"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/consts"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"xorm.io/xorm"
)

// LicenseServer server
type LicenseServer struct {
	dao dao.LicenseManagerDao
	pb.UnimplementedLicenseManagerServiceServer
}

// NewLicenseServer new
func NewLicenseServer(engine *xorm.Engine) *LicenseServer {
	server := &LicenseServer{
		dao: dao.NewLicenseImpl(engine),
	}
	return server
}

// AcquireLicenses 使用license
func (l *LicenseServer) AcquireLicenses(ctx context.Context, in *pb.ConsumeRequest) (*pb.ConsumeResponse, error) {
	var minusResults []*pb.ConsumeResult
	for _, minusInfo := range in.Info {
		logging.Default().Infof("start acquire license, consume info: %s", minusInfo.String())
		licM, err := l.dao.GetLicenseManager(ctx, snowflake.ID(minusInfo.LicManagerId))
		if err != nil {
			errorMsg := fmt.Sprintf("get license manager err, error: %v", err)
			logging.Default().Errorf(errorMsg)
			return nil, status.Error(consts.InvalidArgument, errorMsg)
		}
		if licM == nil {
			result := &pb.ConsumeResult{
				Status: pb.LicenseStatus_UNCONFIGURED,
			}
			minusResults = append(minusResults, result)
			logging.Default().Warnf("license manager not found, consume license info: %s", minusInfo.String())
			continue
		}
		var licList []*models.LicenseInfoExt
		for _, licId := range minusInfo.LicIds {
			if en, ok := isLicConfigured(licId, licM.Licenses); ok {
				licList = append(licList, en)
			}
		}
		if len(minusInfo.LicIds) == 0 {
			licList = licM.Licenses
		}
		// 未配置license
		if len(licList) == 0 {
			result := &pb.ConsumeResult{
				Status: pb.LicenseStatus_UNCONFIGURED,
			}
			minusResults = append(minusResults, result)
			if len(minusInfo.LicIds) > 0 {
				logging.Default().Errorf("error license config, %d needed in solution, but not in license manager, custom license info: %s",
					minusInfo.AppId, minusInfo.String())
			} else {
				logging.Default().Warnf("license unconfiured, custom info: %s", minusInfo.String())
			}
			continue
		}

		licListContent, err := json.Marshal(licList)
		if err != nil {
			logging.Default().Warnf("marshal LicenseInfoExt array fail, error: %v", err)
		}
		logging.Default().Infof("backup license info: %s", string(licListContent))

		// app未发布
		if !licM.Status.Published() {
			result := &pb.ConsumeResult{
				Status: pb.LicenseStatus_UNPUBLISH,
			}
			minusResults = append(minusResults, result)
			logging.Default().Warnf("license config not published, consume license info: %s", minusInfo.String())
			continue
		}
		// 已经获取过license的job不需要再获取
		existed, records, err := l.dao.IsJobUsed(ctx, snowflake.ID(minusInfo.JobId))
		if err == nil && existed {
			licId := records[0].LicenseId
			existed, selectedLic, err := l.dao.GetLicenseInfoByID(ctx, licId)
			if err != nil {
				logging.Default().Warnf("get licnese info fail, id: %s, error: %s", licId.String(), err.Error())
				return nil, err
			} else if !existed {
				logging.Default().Warnf("licnese info not existed, licesne id: %s, consume license: %s",
					licId.String(), minusInfo.String())
				return nil, errors.New("license not existed")
			} else {
				consumeResult, err := generateEnoughConsumeResult(minusInfo.JobId, minusInfo.HpcEndpoint, selectedLic)
				if err != nil {
					logging.Default().Warnf("generate consume result fail, jobID: %d, error: %v", minusInfo.JobId, err)
					return nil, status.Error(codes.Internal, err.Error())
				}
				minusResults = append(minusResults, consumeResult)
				logging.Default().Warnf("find job has acquired license, consume info: %s", minusInfo.String())
			}
			continue
		}
		if len(licList[0].Modules) == 0 {
			logging.Default().Errorf("no license module config, license id: %s, consume info: %s", licList[0].Id, minusInfo)
			return nil, errors.New(fmt.Sprintf("no license module config, license id: %s", licList[0].Id))
		}
		defaultModule := licList[0].Modules[0].ModuleName
		requiredLicNum, err := getRequiredLicenseNum(licM.ComputeRule, defaultModule, minusInfo.Cpus)
		if err != nil {
			logging.Default().Warnf("compute required license num fail, custom license info: %s, error: %s",
				minusInfo.String(), err.Error())
			return nil, err
		}
		logging.Default().Infof("required license num: %v, custom license info: %s", requiredLicNum, minusInfo.String())
		// 已发布的，过滤出适合
		suitableLicenses, err := filterSuitableLicense(licList, requiredLicNum, minusInfo.HpcEndpoint)
		if err != nil {
			logging.Default().Warnf("filter suitable license fail, custom license info: %s, error: %s",
				minusInfo.String(), err.Error())
			return nil, err
		}
		// 没有合适的
		if len(suitableLicenses) == 0 {
			result := &pb.ConsumeResult{
				Status: pb.LicenseStatus_NOTENOUTH,
			}
			minusResults = append(minusResults, result)
			logging.Default().Infof("license not enough, custom info: %s", minusInfo.String())
			continue
		}

		if len(minusInfo.LicIds) == 0 {
			// 未限定license供应商，则排序
			sort.Slice(suitableLicenses, func(i, j int) bool {
				return suitableLicenses[i].Weight <= suitableLicenses[j].Weight
			})
		}

		selectedLicExt := suitableLicenses[0]
		if !in.OnlyQuery { // 仅查询不消耗license
			// 计算并记录消耗数
			err = l.consumeLicense(ctx, minusInfo, selectedLicExt, requiredLicNum)
			if err != nil {
				return nil, err
			}
		}
		consumeResult, err := generateEnoughConsumeResult(minusInfo.JobId,
			minusInfo.HpcEndpoint, &selectedLicExt.LicenseInfo)
		if err != nil {
			logging.Default().Warnf("generate consume result fail, jobID: %d, error: %v", minusInfo.JobId, err)
			return nil, status.Error(codes.Internal, err.Error())
		}
		minusResults = append(minusResults, consumeResult)
		logging.Default().Infof("acquired license successfully, job id: %d, license id: %s",
			minusInfo.JobId, selectedLicExt.LicenseInfo.Id)
	}
	resp := &pb.ConsumeResponse{
		Result: minusResults,
	}
	return resp, nil
}

// ReleaseLicense 增加license
func (l *LicenseServer) ReleaseLicense(ctx context.Context, in *pb.ReleaseRequest) (*empty.Empty, error) {
	jobID := snowflake.ID(in.JobId)
	res := &empty.Empty{}
	logging.Default().Infof("job starts to release license, jobID: %d", jobID)
	existed, records, err := l.dao.IsJobUsed(ctx, jobID)
	if err != nil {
		logging.Default().Warnf("release license fail when get job used, jobId: %s, error: %s", jobID.String(), err.Error())
		return res, err
	}
	if !existed {
		logging.Default().Infof("release license but no records about this job, jobId: %s", jobID.String())
		return res, nil
	}
	err = l.dao.ReleaseLicense(ctx, jobID, records)
	if err != nil {
		logging.Default().Warnf("release license fail, jobId: %s, error: %s", jobID.String(), err.Error())
		return res, err
	}
	logging.Default().Infof("release license successfully, job id: %s", jobID.String())
	return res, nil
}

// minusLicense dao操作
func (l *LicenseServer) consumeLicense(ctx context.Context, info *pb.ConsumeInfo,
	lic *models.LicenseInfoExt, required map[string]int) error {
	err := l.dao.AcquireLicense(ctx, snowflake.ID(info.JobId), rpc.GetInstance().IDGen, lic, required)
	if err != nil {
		logging.Default().Warnf("consume license fail, consume license info: %s, error: %s", info.String(), err.Error())
		return status.Error(consts.ErrServerInternal, err.Error())
	}
	return nil
}

// filterSuitableLicense 过滤适合的license
func filterSuitableLicense(licList []*models.LicenseInfoExt, required map[string]int,
	hpcEndpoint string) ([]*models.LicenseInfoExt, error) {
	var suitableLicenses []*models.LicenseInfoExt
	for _, v := range licList {
		if !isTimeOK(&v.LicenseInfo) || !isAuthed(&v.LicenseInfo) {
			logging.Default().Warnf("license expired or not authed, licenseid: %s", v.Id)
			continue
		}
		if !inAllowableHpcEndpoints(hpcEndpoint, v.AllowableHpcEndpoints) {
			logging.Default().Infof("requesting hpc endpoint is not the same as config hcp endpoint, "+
				"requesting hpc endpoint: %s, config hpc endpoint: %v, license manager id: %s, license provider id: %s",
				hpcEndpoint, v.AllowableHpcEndpoints, v.ManagerId, v.Id)
			continue
		}
		if ok, err := isLicenseEnough(v, required); err != nil {
			logging.Default().Warnf("get license enough info fail, licenseid: %s, required: %v, error: %s",
				v.Id, required, err.Error())
			return nil, errors.Wrap(err, "get license enough info fail")
		} else if ok {
			suitableLicenses = append(suitableLicenses, v)
		}
	}
	return suitableLicenses, nil
}

// inAllowableHpcEndpoint 判断hpcEndpoint是否在允许的列表中
func inAllowableHpcEndpoints(hpcEndpoint string, allowableHpcEndpoints []string) bool {
	for _, v := range allowableHpcEndpoints {
		if v == hpcEndpoint {
			return true
		}
	}
	return false
}

func isLicenseEnough(v *models.LicenseInfoExt, required map[string]int) (bool, error) {
	if v.IsSelfLic() || v.IsConsignedLic() {
		return isSelfOwnedLicenseEnough(v, required)
	} else if v.IsOthersLic() {
		return isOtherOwnedLicenseEnough(v, required)
	} else {
		logging.Default().Errorf("unsupported license type: %d", v.LicenseType)
	}
	return false, errors.New("unsupported license type")
}

func isOtherOwnedLicenseEnough(lic *models.LicenseInfoExt, required map[string]int) (bool, error) {
	components, err := GetRemainLicense(&lic.LicenseInfo)
	if err != nil {
		return false, err
	}
	saved := make(map[string]int, len(components))
	for k, v := range components {
		saved[k] = int(v.Total - v.Used)
	}
	return isRemainEnough(saved, required), nil
}

func GetRemainLicense(licInfo *models.LicenseInfo) (map[string]collector.Component, error) {
	licenseUrl := licInfo.LicenseUrl
	if licInfo.LicensePort > 0 {
		if licInfo.CollectorType == "dsli" {
			licenseUrl = fmt.Sprintf("%s %d", licInfo.LicenseUrl, licInfo.LicensePort)
		} else {
			licenseUrl = fmt.Sprintf("%d@%s", licInfo.LicensePort, licInfo.LicenseUrl)
		}
	}
	licCollectInfo := collector.LicenseCollectInfo{
		LicensePath:   licenseUrl,
		LmstatPath:    licInfo.ToolPath,
		HpcEndpoint:   licInfo.HpcEndpoint,
		CollectorType: licInfo.CollectorType,
	}
	c, err := collector.NewCollector(&licCollectInfo)
	if err != nil {
		logging.Default().Warnf("NewCollectorFail, Error: %s, CollectorInfo: %+v", err.Error(), licCollectInfo)
		return nil, err
	}
	collectRuntimeErr := c.Collect()
	components := c.GetComponents()
	return components, collectRuntimeErr
}

func isSelfOwnedLicenseEnough(lic *models.LicenseInfoExt, required map[string]int) (bool, error) {
	saved := make(map[string]int)
	actualSaved := make(map[string]int)
	for _, m := range lic.Modules {
		//使用剩余数量
		saved[m.ModuleName] = m.Total - m.Used
		//实际剩余数量
		actualSaved[m.ModuleName] = m.ActualTotal - m.ActualUsed
	}
	savedEnough := isRemainEnough(saved, required)
	actualSavedEnough := isRemainEnough(actualSaved, required)
	return savedEnough && actualSavedEnough, nil
}

func isRemainEnough(saved, required map[string]int) bool {
	for feature, reqNum := range required {
		if remain, ok := saved[feature]; !ok || reqNum > remain {
			return false
		}
	}
	return true
}

func getRequiredLicenseNum(computeLicNumRule string, defaultModule string, cpus int64) (map[string]int, error) {
	// TODO license数量会不会是float？
	res := make(map[string]int)
	// 从description 获取license规则
	shell := computeLicNumRule
	// 执行sh 获取license数量
	// shell 文件格式 option="${1}" \ncase ${option} in ..
	// cpus数为shell 入参数
	cmd := exec.Command("/bin/sh", "-c", shell, "--", strconv.FormatInt(cpus, 10))
	cmd.Env = []string{fmt.Sprintf("CPU=%d", cpus)}
	output, err := cmd.CombinedOutput()
	if err != nil {
		logging.Default().Infof("Error executing sh error: %v", err)
		return res, err
	}

	outputString := strings.TrimSpace(string(output))
	err = json.Unmarshal([]byte(outputString), &res)
	if err != nil {
		if num, err := parseInt(outputString); err == nil {
			res[defaultModule] = num
			return res, nil
		}
		logging.Default().Warnf("get required license num fail, json unmarshal fail, outputstring: %s, error: %s",
			outputString, err.Error())
		return nil, err
	}
	return res, nil
}

func parseInt(str string) (int, error) {
	re := regexp.MustCompile(`\d+`)
	numberString := re.FindString(str)
	return strconv.Atoi(numberString)
}

// isAuthed 是否auth
func isAuthed(v *models.LicenseInfo) bool {
	return v.Auth == 1
}

// isTimeOK 时间是否在范围内
func isTimeOK(v *models.LicenseInfo) bool {
	begin := v.BeginTime
	end := v.EndTime
	now := time.Now()
	return now.After(begin) && now.Before(end)
}

func isLicConfigured(licId string, backup []*models.LicenseInfoExt) (*models.LicenseInfoExt, bool) {
	for _, en := range backup {
		if en.LicenseInfo.Id.String() == licId {
			return en, true
		}
	}
	return nil, false
}

func generateEnoughConsumeResult(jobId int64, hpcEndpoint string, lic *models.LicenseInfo) (*pb.ConsumeResult, error) {
	server, err := lic.GetLicenseProxies(hpcEndpoint)
	if err != nil {
		return nil, err
	}
	return &pb.ConsumeResult{
		JobId:       jobId,
		ServerUrl:   server,
		LicenseEnvs: []string{server},
		Status:      pb.LicenseStatus_ENOUGH,
	}, nil
}
