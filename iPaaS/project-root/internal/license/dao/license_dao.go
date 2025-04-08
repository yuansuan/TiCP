package dao

import (
	"context"
	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql" // register mysql
	"github.com/pkg/errors"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/common/project-root-api/license/v1/license_manager/publish"
	"github.com/yuansuan/ticp/common/project-root-api/proto/idgen"
	pb "github.com/yuansuan/ticp/common/project-root-api/proto/license"
	dbModels "github.com/yuansuan/ticp/iPaaS/project-root/internal/license/dao/models"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/common/with"
	"strings"
	"time"
	"xorm.io/xorm"
)

var (
	joinEntity     = &dbModels.JoinEntity{}
	licenseManager = &dbModels.LicenseManager{}
	licenseInfo    = &dbModels.LicenseInfo{}
)

// LicenseImpl impl
type LicenseImpl struct {
	engine *xorm.Engine
}

// NewLicenseImpl impl
func NewLicenseImpl(engine *xorm.Engine) *LicenseImpl {
	return &LicenseImpl{
		engine: engine,
	}
}

func SyncTable() {
	err := with.DefaultSession(context.Background(), func(db *xorm.Session) error {
		if err := db.Sync2(new(dbModels.LicenseInfo)); err != nil {
			logging.Default().Fatalf("SyncLicenseInfoTableFail, Error: %s", err.Error())
		}
		if err := db.Sync2(new(dbModels.LicenseJob)); err != nil {
			logging.Default().Fatalf("SyncLicenseJobTableFail, Error: %s", err.Error())
		}
		if err := db.Sync2(new(dbModels.LicenseManager)); err != nil {
			logging.Default().Fatalf("SyncLicenseManagerTableFail, Error: %s", err.Error())
		}
		if err := db.Sync2(new(dbModels.ModuleConfig)); err != nil {
			logging.Default().Fatalf("SyncLicenseManagerTableFail, Error: %s", err.Error())
		}
		return nil
	})

	if err != nil {
		logging.Default().Fatalf("sync table fail, err: %v", err)
	}
	logging.Default().Info("SyncTableOk")
}

// ListAllLicenses query all license
type ListAllLicenses struct {
	Provider  string
	BeginTime time.Time
	EndTime   time.Time
	Page      *pb.Page
}

// LicenseInfo info
type LicenseInfo struct {
	ManagerID snowflake.ID
}

func (l *LicenseImpl) ListLicenseManagers(ctx context.Context, in *ListAllLicenses) (list []*dbModels.JoinEntity, total int64, err error) {
	// TODO 未分页
	db := l.engine.Context(ctx)
	defer func(db *xorm.Session) {
		err = db.Close()
		if err != nil {
			logging.Default().Warnf("close session error, err: %v", err)
		}
	}(db)

	db.Table(joinEntity.TableName()).
		Join("LEFT", "license_info", "license_manager.id = license_info.manager_id").
		Join("LEFT", "module_config", "license_info.id = module_config.license_id")

	if in != nil && len(in.Provider) > 0 {
		db.And("provider like ?", "%"+in.Provider+"%")
	}

	if in != nil && !in.BeginTime.IsZero() && !in.EndTime.IsZero() {
		db.And("begin_time >=?", in.BeginTime).Where("end_time <=?", in.EndTime)
	}
	total, queryErr := db.FindAndCount(&list)
	if queryErr != nil {
		return nil, -1, errors.Wrap(queryErr, "list dao")
	}

	return list, total, nil
}

func (l *LicenseImpl) GetLicenseManager(ctx context.Context, lmId snowflake.ID) (*dbModels.LicenseManagerExt, error) {
	var entities []*dbModels.JoinEntity
	err := with.DefaultSession(ctx, func(db *xorm.Session) error {
		queryErr := db.Table(joinEntity.TableName()).
			Join("LEFT", "license_info", "license_manager.id = license_info.manager_id").
			Join("LEFT", "module_config", "license_info.id = module_config.license_id").
			Where("license_manager.id = ?", lmId).
			Find(&entities)
		return errors.Wrap(queryErr, "single dao")
	})
	return dbModels.ToOneLicenseManagerExt(entities), err
}

func (l *LicenseImpl) DeleteLicenseManager(ctx context.Context, lmId snowflake.ID) (bool, bool, error) {
	var count int64
	var canDelete = true
	err := with.DefaultTransaction(ctx, func(ctx context.Context) error {
		return with.DefaultSession(ctx, func(db *xorm.Session) error {
			var licInfos []*dbModels.LicenseInfo
			err := db.Where("manager_id = ? ", lmId).ForUpdate().Find(&licInfos)
			if err != nil {
				return err
			}
			if len(licInfos) > 0 {
				canDelete = false
				return nil
			}
			count, err = db.Delete(&dbModels.LicenseManager{Id: lmId})
			return err
		})
	})
	return count == 1, canDelete, err
}

func (l *LicenseImpl) AddLicenseInfo(ctx context.Context, licInfo *dbModels.LicenseInfo) error {
	return with.DefaultSession(ctx, func(db *xorm.Session) error {
		_, err := db.Insert(licInfo)
		return errors.Wrap(err, "add license dao ")
	})
}

func (l *LicenseImpl) UpdateLicenseInfo(ctx context.Context, licInfo *dbModels.LicenseInfo) (suc bool, err error) {
	var count int64
	err = with.DefaultSession(ctx, func(db *xorm.Session) error {
		count, err = db.Where("id = ?", licInfo.Id).Cols("provider", "license_server", "mac_addr", "license_url",
			"license_port", "license_proxies", "license_num", "weight", "begin_time", "end_time", "auth", "license_type",
			"tool_path", "collector_type", "hpc_endpoint", "allowable_hpc_endpoints").Update(licInfo)
		return err
	})
	if err != nil {
		return false, err
	}
	return count == 1, err
}

func (l *LicenseImpl) DeleteLicenseInfo(ctx context.Context, licId snowflake.ID) (suc bool, canDelete bool, err error) {
	var count int64
	canDelete = true
	err = with.DefaultTransaction(ctx, func(ctx context.Context) error {
		return with.DefaultSession(ctx, func(db *xorm.Session) error {
			var mCfgs []*dbModels.ModuleConfig
			// 锁住moduleconfig
			err := db.Where("license_id = ?", licId).ForUpdate().Find(&mCfgs)
			if err != nil {
				return err
			}
			licUsed := false
			for _, cfg := range mCfgs {
				if cfg.Used > 0 {
					licUsed = true
					break
				}
			}
			if licUsed {
				canDelete = false
				return nil
			}
			count, err = db.Delete(&dbModels.LicenseInfo{Id: licId})
			if err != nil {
				return err
			}
			_, err = db.Where("license_id = ?", licId).Delete(&dbModels.ModuleConfig{})
			if err != nil {
				return err
			}
			return err
		})
	})
	return count >= 1, canDelete, err
}

func (l *LicenseImpl) SetLicenseServerStatus(ctx context.Context, licId snowflake.ID, status string) (suc bool, err error) {
	var affectRows int64
	err = with.DefaultSession(ctx, func(db *xorm.Session) error {
		affectRows, err = db.ID(licId).Cols("license_server_status").Update(&dbModels.LicenseInfo{LicenseServerStatus: status})
		return err
	})

	return affectRows > 0, err
}

func (l *LicenseImpl) ListModuleConfig(ctx context.Context, licId snowflake.ID) (list []*dbModels.ModuleConfig, err error) {
	err = with.DefaultSession(ctx, func(db *xorm.Session) error {
		err = db.Where("license_id = ?", licId).Find(&list)
		return errors.Wrap(err, "list module configs dao")
	})
	return
}

func (l *LicenseImpl) AddModuleConfig(ctx context.Context, cfg *dbModels.ModuleConfig) (suc bool, err error) {
	err = with.DefaultSession(ctx, func(db *xorm.Session) error {
		count, err := db.Insert(cfg)
		// work in mysql
		if err != nil && strings.Contains(strings.ToLower(err.Error()), "duplicate") {
			suc = false
			err = nil
		}
		suc = count == 1
		return errors.Wrap(err, "add license dao ")
	})
	return
}

func (l *LicenseImpl) BatchAddModuleConfigs(ctx context.Context, moduleConfigs []*dbModels.ModuleConfig) (err error) {
	return with.DefaultSession(ctx, func(db *xorm.Session) error {
		count, err := db.Insert(moduleConfigs)
		if err != nil {
			return errors.Wrapf(err, "failed to batch add %d module configs", len(moduleConfigs))
		}
		if count != int64(len(moduleConfigs)) {
			return errors.Errorf("not all module configurations were inserted successfully: expected %d, got %d", len(moduleConfigs), count)
		}
		return nil
	})
}

func (l *LicenseImpl) UpdateModuleConfigTotal(ctx context.Context, cfg *dbModels.ModuleConfig) (suc bool, err error) {
	var count int64
	err = with.DefaultSession(ctx, func(db *xorm.Session) error {
		count, err = db.Where("id = ?", cfg.Id).Cols("module_name", "total").Update(cfg)
		return err
	})
	if err != nil {
		return false, err
	}
	return count == 1, err
}

func (l *LicenseImpl) UpdateModuleConfigActual(ctx context.Context, cfg *dbModels.ModuleConfig) (suc bool, err error) {
	var count int64
	err = with.DefaultSession(ctx, func(db *xorm.Session) error {
		count, err = db.Where("license_id = ? and module_name = ?", cfg.LicenseId, cfg.ModuleName).Cols("actual_used", "actual_total").Update(cfg)
		return err
	})
	if err != nil {
		return false, err
	}
	return count == 1, err
}

func (l *LicenseImpl) DeleteModuleConfig(ctx context.Context, id snowflake.ID) (suc bool, canDelete bool, err error) {
	canDelete = true
	err = with.DefaultTransaction(ctx, func(ctx context.Context) error {
		return with.DefaultSession(ctx, func(db *xorm.Session) error {
			var mCfg = dbModels.ModuleConfig{}
			// 锁住moduleconfig
			suc, err = db.Where("id = ?", id).ForUpdate().Get(&mCfg)
			if err != nil {
				return err
			}
			if !suc {
				return nil
			}
			canDelete = mCfg.Used == 0
			if !canDelete {
				return nil
			}
			_, err = db.Delete(&dbModels.ModuleConfig{Id: id})
			return err
		})
	})
	return
}

func (l *LicenseImpl) GetLicenseInfoByID(ctx context.Context, licId snowflake.ID) (existed bool, lic *dbModels.LicenseInfo, err error) {
	lic = &dbModels.LicenseInfo{}
	err = with.DefaultSession(ctx, func(db *xorm.Session) error {
		existed, err = db.Where("id = ?", licId.Int64()).ForUpdate().Get(lic)
		return err
	})
	return existed, lic, err
}

// LicenseInfoByAppIDs 批量查询AppIDs
func (l *LicenseImpl) LicenseInfoByAppIDs(ctx context.Context, appIDs []snowflake.ID) (entities []*dbModels.JoinEntity, err error) {
	err = with.DefaultSession(ctx, func(db *xorm.Session) error {
		queryErr := db.Table(joinEntity.TableName()).
			Join("INNER", "license_info", "license_manager.id = license_info.manager_id").
			Join("INNER", "module_config", "license_info.id = module_config.license_id").
			//Where("license_manager.app_id = ?", appID).
			In("license_manager.app_id = ?", appIDs).
			Find(&entities)
		return errors.Wrap(queryErr, "appIDs dao")
	})
	return
}

func (l *LicenseImpl) AddLicenseManager(ctx context.Context, lm *dbModels.LicenseManager) (err error) {
	return with.DefaultSession(ctx, func(db *xorm.Session) error {
		_, err := db.Insert(lm)
		return errors.Wrap(err, "add licenseManage dao ")
	})
}

// UpdateLicenseManager ...
func (l *LicenseImpl) UpdateLicenseManager(ctx context.Context, lm *dbModels.LicenseManager) (suc bool, err error) {
	var count int64
	err = with.DefaultSession(ctx, func(db *xorm.Session) error {
		count, err = db.ID(lm.Id).Cols("os", "description", "compute_rule", "app_type", "status").Update(lm)
		return errors.Wrap(err, "update licenseManage dao ")
	})
	if err != nil {
		return false, err
	}
	return count == 1, err
}

// SelectByManagerID 根据manager_id查询
func (l *LicenseImpl) SelectByManagerID(ctx context.Context, managerID snowflake.ID) (fromDB []*dbModels.LicenseInfo, err error) {
	err = with.DefaultSession(ctx, func(db *xorm.Session) error {
		err = db.Where("manager_id = ?", managerID).Find(&fromDB)
		return errors.Wrap(err, "SelectByManagerID dao")
	})
	return
}

// SelectByAddressAndScId 根据licenseServer地址和sc_id查询
func (l *LicenseImpl) SelectByAddressAndScId(ctx context.Context, scId, host, port string) (licenseInfos []*dbModels.LicenseInfo, err error) {
	err = with.DefaultSession(ctx, func(db *xorm.Session) error {
		err = db.Where("sc_id = ?", scId).And("license_url = ?", host).And("license_port = ?", port).Find(&licenseInfos)
		return errors.Wrap(err, "SelectByAppIdAndScId dao")
	})
	return
}

// LicenseInfoPublished 查询已发布的license
func (l *LicenseImpl) LicenseInfoPublished(ctx context.Context) (entities []*dbModels.LicenseInfo, err error) {
	err = with.DefaultSession(ctx, func(db *xorm.Session) error {
		err = db.Table(licenseInfo.TableName()).
			Join("INNER", "license_manager", "license_info.manager_id = license_manager.id").
			And("license_manager.status = 1").And("license_info.Auth = 1").
			Find(&entities)
		if err != nil {
			return errors.Wrap(err, "Published dao")
		}
		return nil
	})
	return
}

// SelectByInfoIDs 获取info
func (l *LicenseImpl) SelectByInfoIDs(ctx context.Context, ids []snowflake.ID) (fromDB []*dbModels.ModuleConfig, err error) {
	err = with.DefaultSession(ctx, func(db *xorm.Session) error {
		err = db.In("license_id", ids).Find(&fromDB)
		return errors.Wrap(err, "SelectByInfoIDs dao")
	})
	return
}

// PublishLicense publish
func (l *LicenseImpl) PublishLicense(ctx context.Context, managerID snowflake.ID, status publish.Status, publishTime time.Time) (err error) {
	return with.DefaultSession(ctx, func(db *xorm.Session) error {
		_, err := db.ID(managerID).Cols("status", "publish_time").Update(&dbModels.LicenseManager{Status: status, PublishTime: publishTime})
		return errors.Wrap(err, "publish licenseManage dao ")
	})
}

// IsAppUsed app
func (l *LicenseImpl) IsAppUsed(ctx context.Context, appID int64) (get bool, err error) {
	var manager dbModels.LicenseManager
	err = with.DefaultSession(ctx, func(db *xorm.Session) error {
		get, err = db.In("app_id", appID).Get(&manager)

		if err != nil {
			return errors.Wrap(err, "find licenseManager error")
		}
		return nil
	})
	if err != nil {
		return false, err
	}
	if !get {
		return false, nil
	}
	return true, nil
}

// AcquireLicense minus
func (l *LicenseImpl) AcquireLicense(ctx context.Context, jobID snowflake.ID, idGenClient idgen.IdGenClient,
	lic *dbModels.LicenseInfoExt, required map[string]int) (err error) {
	err = with.DefaultTransaction(ctx, func(ctx context.Context) error {
		// 查询jobID是否被分配
		exit, _, queryError := isJobUsed(ctx, jobID)
		if queryError != nil {
			logging.Default().Warnf("judge license used fail, jobID: %s, error: %s", jobID, err.Error())
			return queryError
		}
		if exit {
			logging.Default().Warnf("job has acquire license, jobID: %s", jobID)
			return nil
		}
		for moduleName, reqNum := range required {
			reply, err := idGenClient.GenerateID(ctx, &idgen.GenRequest{})
			if err != nil {
				return err
			}
			recordID := snowflake.ID(reply.Id)
			moduleID := moduleNameToID(moduleName, lic.Modules)
			var remainingLicenseNum int
			remainingLicenseNum, err = lic.GetRemainingLicenseNum(moduleName)
			if err != nil {
				logging.Default().Warnf("get remaining license num fail, moduleName: %s, error: %s", moduleName, err.Error())
				return err
			}
			if !lic.IsOthersLic() {
				// 减license
				err = acquireLicenseNum(ctx, reqNum, moduleID, remainingLicenseNum)
				if err != nil {
					return err
				}
			}
			// 绑定joID和moduleID
			err = insertLicenseJob(ctx, jobID, lic.Id, moduleID, recordID, reqNum)
			if err != nil {
				return err
			}
		}

		return nil
	})
	return
}

// BatchUpdateModuleConfig 批量修改
func (l *LicenseImpl) BatchUpdateModuleConfig(ctx context.Context, modules []*dbModels.ModuleConfig) error {
	return with.DefaultSession(ctx, func(db *xorm.Session) error {
		dbSQL := goqu.Dialect("mysql")
		moduleSQL := dbSQL.Insert("module_config").Rows(modules)
		sql, _, err := moduleSQL.ToSQL()
		if err != nil {
			return errors.Wrap(err, "update module sql error")
		}
		// Upsert实现
		sql += " on DUPLICATE KEY UPDATE actual_total = values(actual_total), " +
			"actual_used = values(actual_used)"
		_, err = db.Exec(sql)
		return errors.Wrap(err, "update module dao")
	})
}

func moduleNameToID(name string, modules []*dbModels.ModuleConfig) snowflake.ID {
	for _, m := range modules {
		if m.ModuleName == name {
			return m.Id
		}
	}
	logging.Default().Errorf("module name not found, name: %s, modules: %v", name, modules)
	return snowflake.Zero()
}

// insertLicenseJob insert
func insertLicenseJob(ctx context.Context, jobID, licID, moduleID, id snowflake.ID, num int) (err error) {
	err = with.DefaultSession(ctx, func(db *xorm.Session) error {
		licenseJob := &dbModels.LicenseJob{
			Id:        id,
			ModuleId:  moduleID,
			JobId:     jobID,
			Used:      1,
			Licenses:  int64(num),
			LicenseId: licID,
		}
		_, err = db.Insert(licenseJob)
		if err != nil {
			err = errors.Wrap(err, "insert licenseJob dao error")
		}
		return nil
	})
	return err
}

// updateLicenseJob update
func updateLicenseJob(ctx context.Context, jobID snowflake.ID) (err error) {
	err = with.DefaultSession(ctx, func(db *xorm.Session) error {
		licenseJob := dbModels.LicenseJob{
			Used: 2,
		}
		_, err = db.Where("job_id = ?", jobID).Cols("used").Update(&licenseJob)
		if err != nil {
			errors.Wrap(err, "update licenseJob dao error")
		}
		return nil
	})
	return err
}

func acquireLicenseNum(ctx context.Context, requiredLicenseNum int, moduleID snowflake.ID, remainingLicenseNum int) (err error) {
	err = with.DefaultSession(ctx, func(db *xorm.Session) error {
		sql := "update `module_config` set used = used + ? where id = ? and used = ?"
		_, err = db.Exec(sql, requiredLicenseNum, moduleID, remainingLicenseNum)
		if err != nil {
			return errors.Wrap(err, "minus module dao error")
		}
		return nil
	})
	return err
}

func releaseRemainingLicenseNum(ctx context.Context, releasedNum int, moduleID snowflake.ID) (err error) {
	err = with.DefaultSession(ctx, func(db *xorm.Session) error {
		sql := "update `module_config` set used = used - ? where id = ? and used - ? >= 0"
		ret, err := db.Exec(sql, releasedNum, moduleID, releasedNum)
		if err != nil {
			return errors.Wrap(err, "minus module dao error")
		}
		affectRows, err := ret.RowsAffected()
		if err != nil {
			return errors.Wrap(err, "minus module dao error")
		}
		if affectRows <= 0 {
			logging.Default().Warnf("the number of affected rows is 0, moduleID: %s, releasedNum: %d", moduleID.String(), releasedNum)
		}
		return nil
	})
	return err
}

// IsJobUsed job
func (l *LicenseImpl) IsJobUsed(ctx context.Context, jobID snowflake.ID) (existed bool, licenseJob []*dbModels.LicenseJob, err error) {
	return isJobUsed(ctx, jobID)
}

// isJobUsed 内部函数使用
func isJobUsed(ctx context.Context, jobID snowflake.ID) (existed bool, records []*dbModels.LicenseJob, err error) {
	err = with.DefaultSession(ctx, func(db *xorm.Session) error {
		return db.Where("job_id = ? and used = 1", jobID).Find(&records)
	})
	existed = len(records) > 0
	return
}

// IsLicenseUsed 取消发布时判断是否有license在使用
func (l *LicenseImpl) IsLicenseUsed(ctx context.Context, moduleID snowflake.ID) (licenseJobs []*dbModels.LicenseJob, err error) {
	return isLicenseUsed(ctx, moduleID)
}

// isLicenseUsed 内部函数使用
func isLicenseUsed(ctx context.Context, moduleID snowflake.ID) (licenseJobs []*dbModels.LicenseJob, err error) {
	err = with.DefaultSession(ctx, func(db *xorm.Session) error {
		err = db.Where("module_id = ?", moduleID).Where("used = 1").Find(&licenseJobs)
		return errors.Wrapf(err, "isLicenseUsed error %d", moduleID)
	})
	return
}

// ReleaseLicense plus
func (l *LicenseImpl) ReleaseLicense(ctx context.Context, jobID snowflake.ID, licenseJob []*dbModels.LicenseJob) (err error) {
	err = with.DefaultTransaction(ctx, func(ctx context.Context) error {
		licId := licenseJob[0].LicenseId
		existed, licInfo, err := l.GetLicenseInfoByID(ctx, licId)
		if err != nil {
			logging.Default().Warnf("get license info fail, licId: %s, jobId: %s", licId.String(), jobID.String())
			return errors.Wrap(err, "get license info fail")
		} else if !existed {
			logging.Default().Warnf("license not existed, licId: %s, jobId: %s", licId.String(), jobID.String())
			return errors.Wrap(err, "license not existed")
		}
		if !licInfo.IsOthersLic() {
			// 非外部license需要更新计数
			for _, record := range licenseJob {
				// 归还license
				err = releaseRemainingLicenseNum(ctx, int(record.Licenses), record.ModuleId)
				if err != nil {
					return err
				}
			}
		}
		// 更新记录
		err = updateLicenseJob(ctx, jobID)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

// LicenseUsed 查询license使用情况（企业纬度）
func (l *LicenseImpl) LicenseUsed(ctx context.Context, managerId, size, index int64) (results []*dbModels.LicenseUsedResult, total int64, err error) {
	err = with.DefaultSession(ctx, func(sess *xorm.Session) error {
		sess.Table("license_job").
			Join("Left", "module_config", "license_job.module_id = module_config.id").
			Join("Left", "license_info", "license_info.id = module_config.license_id").
			Join("Left", "license_manager", "license_manager.id = license_info.manager_id")

		sess.Where("license_manager.id=?", managerId)
		sess.GroupBy("license_job.job_id")
		limitSize, limitOffset := int(size), int((index-1)*size)

		sess.Select(
			"license_job.id ," +
				"license_job.job_id ," +
				"sum(license_job.licenses) as licenses," +
				"license_job.create_time ," +
				"license_manager.app_id ")

		//决定排序方式
		sess = sess.OrderBy("create_time desc")
		total, err = sess.Limit(limitSize, limitOffset).FindAndCount(&results)
		if err != nil {
			return err
		}
		return nil
	})
	return results, total, err
}
