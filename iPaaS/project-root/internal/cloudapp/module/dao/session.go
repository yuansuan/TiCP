package dao

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"xorm.io/xorm"

	schema "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	zone "github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/config"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/module/dao/models"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/module/utils"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/module/utils/exorm"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/common/with"
)

var (
	// ErrIbvSessionNotFound 找不到指定的会话
	ErrIbvSessionNotFound = errors.New("dao: session not found")

	_Session      = (*models.Session)(nil)
	_Instance     = (*models.Instance)(nil)
	_Hardware     = (*models.Hardware)(nil)
	_HardwareUser = (*models.HardwareUser)(nil)
	_Software     = (*models.Software)(nil)
	_SoftwareUser = (*models.SoftwareUser)(nil)
)

// CreateSession 创建一个会话记录
func CreateSession(ctx context.Context, session *models.Session) error {
	return with.DefaultSession(ctx, func(db *xorm.Session) error {
		_, err := db.Insert(session)
		return errors.Wrap(err, "dao")
	})
}

type ListSessionDetailParams struct {
	SessionIDs  []snowflake.ID
	UserID      snowflake.ID
	HardwareID  snowflake.ID
	SoftwareID  snowflake.ID
	Statuses    []string
	PageOffset  int
	PageSize    int
	Zone        zone.Zone
	UserIDs     []snowflake.ID
	WithDeleted bool
}

func ListSessionDetail(ctx context.Context, params *ListSessionDetailParams) (list []*models.SessionWithDetail, total int64, err error) {
	err = with.DefaultSession(ctx, func(db *xorm.Session) error {
		db.Table(_Session.TableName()).Alias("s").
			Join("INNER", []string{_Instance.TableName(), "i"}, "s.instance_id = i.id").
			Join("LEFT", []string{_Hardware.TableName(), "h"}, "i.hardware_id = h.id").
			Join("LEFT", []string{_Software.TableName(), "f"}, "i.software_id = f.id")

		total, err = exorm.New(db).
			Cond(params.UserID != 0, "s.user_id = ?", params.UserID).
			Cond(params.HardwareID != 0, "i.hardware_id = ?", params.HardwareID).
			Cond(params.SoftwareID != 0, "i.software_id = ?", params.SoftwareID).
			Cond(params.Zone != "", "s.zone = ?", params.Zone).
			Cond(!params.WithDeleted, "s.deleted = ?", false).
			CondIn(len(params.Statuses) != 0, "s.status", params.Statuses).
			CondIn(len(params.SessionIDs) != 0, "s.id", params.SessionIDs).
			CondIn(len(params.UserIDs) != 0, "s.user_id", params.UserIDs).
			Raw().
			OrderBy("s.id desc").
			Limit(params.PageSize, params.PageOffset).
			FindAndCount(&list)
		return errors.Wrap(err, "dao")
	})
	return
}

func GetSessionByInstancePk(ctx context.Context, instanceId snowflake.ID) (*models.Session, error) {
	var session models.Session
	err := with.DefaultSession(ctx, func(db *xorm.Session) error {
		exists, err := db.Where("instance_id = ?", instanceId).Get(&session)
		if err != nil {
			return errors.Wrap(err, "dao")
		} else if !exists {
			return ErrIbvSessionNotFound
		}
		return nil
	})
	return &session, err
}

func GetSessionDetailByRawInstanceId(ctx context.Context, rawInstanceId string) (*models.SessionWithDetail, error) {
	var session models.SessionWithDetail
	err := with.DefaultSession(ctx, func(db *xorm.Session) error {
		exists, err := db.Table(_Session.TableName()).Alias("s").
			Join("INNER", []string{_Instance.TableName(), "i"}, "s.instance_id = i.id").
			Join("LEFT", []string{_Hardware.TableName(), "h"}, "i.hardware_id = h.id").
			Join("LEFT", []string{_Software.TableName(), "f"}, "i.software_id = f.id").
			Where("i.instance_id = ?", rawInstanceId).Get(&session)
		if err != nil {
			return errors.Wrap(err, "dao")
		} else if !exists {
			return ErrIbvSessionNotFound
		}
		return nil
	})
	return &session, err
}

// ListWaitingCloseSessionByZone 列出所有等待关闭的会话（这里不做删除筛选）
func ListWaitingCloseSessionByZone(ctx context.Context, zoneName string) (list []*models.SessionWithDetail, _ error) {
	err := with.DefaultSession(ctx, func(db *xorm.Session) error {
		err := db.Table(_Session.TableName()).Alias("s").
			Join("INNER", []string{_Instance.TableName(), "i"}, "s.instance_id = i.id").
			Join("LEFT", []string{_Hardware.TableName(), "h"}, "i.hardware_id = h.id").
			Join("LEFT", []string{_Software.TableName(), "f"}, "i.software_id = f.id").
			Where("s.zone = ?", zoneName).
			Where("s.close_signal = ?", true).
			Where("s.status != ?", schema.SessionClosed.String()).
			Find(&list)
		return errors.Wrap(err, "dao")
	})
	return list, err
}

func SessionStarting(ctx context.Context, id snowflake.ID) error {
	return with.DefaultSession(ctx, func(db *xorm.Session) error {
		_, err := db.ID(id).
			Where("status = ?", schema.SessionPending).
			Cols("status", "start_time").
			Update(&models.Session{
				Status:    schema.SessionStarting,
				StartTime: utils.PNow(),
			})
		return errors.Wrap(err, "dao")
	})
}

func UpdateSessionStatus(ctx context.Context, id snowflake.ID, from []schema.SessionStatus, to schema.SessionStatus) error {
	return with.DefaultSession(ctx, func(db *xorm.Session) error {
		_, err := db.ID(id).
			In("status", from).
			Cols("status").
			Update(&models.Session{
				Status: to,
			})
		return errors.Wrap(err, "dao")
	})
}

// SessionUserClosing 标记指定会话将要关闭（置位close_signal）
func SessionUserClosing(ctx context.Context, userId snowflake.ID, sessionId snowflake.ID) (bool, error) {
	exist := false
	err := with.DefaultSession(ctx, func(db *xorm.Session) error {
		affected, err := db.ID(sessionId).
			Where("user_id = ?", userId).
			In("status", []schema.SessionStatus{schema.SessionStarting, schema.SessionStarted}).
			Cols("close_signal", "status", "user_close_time").
			Update(&models.Session{
				CloseSignal:   true,
				Status:        schema.SessionClosing,
				UserCloseTime: time.Now(),
			})
		if affected != 0 {
			exist = true
		}

		return errors.Wrap(err, "dao")
	})

	return exist, err
}

// SessionSysClosing 标记指定会话将要关闭（置位close_signal）
func SessionSysClosing(ctx context.Context, id snowflake.ID, reason string) error {
	return with.DefaultSession(ctx, func(db *xorm.Session) error {
		_, err := db.ID(id).
			Where("status != ?", schema.SessionClosed).
			Cols("close_signal", "status", "exit_reason").
			Update(&models.Session{
				CloseSignal: true,
				ExitReason:  reason,
				Status:      schema.SessionClosing,
			})
		return errors.Wrap(err, "dao")
	})
}

// SessionClosed 标记会话已关闭
func SessionClosed(ctx context.Context, id snowflake.ID, reason string) error {
	return with.DefaultSession(ctx, func(db *xorm.Session) error {
		var session models.Session
		exists, err := db.ID(id).Cols("exit_reason").Get(&session)
		if err != nil {
			return err
		} else if !exists {
			return ErrIbvSessionNotFound
		}

		_, err = db.ID(id).
			Cols("status", "exit_reason", "end_time").
			Update(&models.Session{
				Status:     schema.SessionClosed,
				ExitReason: joinExitReason(session.ExitReason, reason),
				EndTime:    utils.PNow(),
			})
		return errors.Wrap(err, "dao")
	})
}

// DeleteSession 删除会话 软
func DeleteSession(ctx context.Context, id snowflake.ID) (deleted bool, _ error) {
	return deleted, with.DefaultSession(ctx, func(db *xorm.Session) error {
		affectedRows, err := db.ID(id).
			Cols("deleted").
			Where("status = ?", schema.SessionClosed).
			Update(&models.Session{Deleted: true})

		deleted = affectedRows != 0
		return errors.Wrap(err, "dao")
	})
}

// GetSession 根据ID获取会话实例
func GetSession(ctx context.Context, userId, sessionId snowflake.ID, lock bool) (sess *models.Session, exist bool, err error) {
	err = with.DefaultSession(ctx, func(db *xorm.Session) error {
		if lock {
			db.ForUpdate()
		}

		if userId != snowflake.ID(0) {
			db.Where("user_id = ?", userId)
		}

		sess = new(models.Session)
		if exists, qErr := db.ID(sessionId).
			Get(sess); qErr != nil {
			return errors.Wrap(qErr, "dao")
		} else if !exists {
			exist = false
			return nil
		}
		exist = true

		return nil
	})

	return
}

// joinExitReason 拼接关闭原因
func joinExitReason(base, reason string) string {
	if len(base) == 0 || len(reason) == 0 {
		return base
	}
	return base + "|" + reason
}

func UpdateSession(ctx context.Context, sessionUpdate models.Session, cols []string) (err error) {
	err = with.DefaultSession(ctx, func(db *xorm.Session) error {
		_, err = db.ID(sessionUpdate.Id).
			Cols(cols...).
			Update(&sessionUpdate)
		return errors.Wrap(err, "dao")
	})
	return
}

// GetSessionDetailsBySessionID 获取session detail
func GetSessionDetailsBySessionID(ctx context.Context, userID, sessionID snowflake.ID) (*models.SessionWithDetail, bool, error) {
	sessionDetail := &models.SessionWithDetail{}
	exists := true
	err := with.DefaultSession(ctx, func(db *xorm.Session) error {
		db = db.Table(_Session.TableName()).Alias("s").
			Join("INNER", []string{_Instance.TableName(), "i"}, "s.instance_id = i.id").
			Join("LEFT", []string{_Hardware.TableName(), "h"}, "i.hardware_id = h.id").
			Join("LEFT", []string{_Software.TableName(), "f"}, "i.software_id = f.id").
			Where("s.id = ?", sessionID).
			Where("s.deleted = ?", false)

		if userID != 0 {
			db = db.Where("s.user_id = ?", userID)
		}

		exist, err := db.Get(sessionDetail)

		if err != nil {
			return errors.Wrap(err, "dao")
		}
		if !exist {
			exists = false
			return nil
		}

		return nil
	})
	if err != nil {
		return nil, exists, err
	}

	return sessionDetail, exists, nil
}

func GetSessionDetailsBySessionIDWithLock(ctx context.Context, userID, sessionID snowflake.ID) (*models.SessionWithDetail, bool, error) {
	sessionDetail := &models.SessionWithDetail{}
	exists := true

	err := with.DefaultSession(ctx, func(db *xorm.Session) error {
		db.ForUpdate()

		db = db.Table(_Session.TableName()).Alias("s").
			Join("INNER", []string{_Instance.TableName(), "i"}, "s.instance_id = i.id").
			Join("LEFT", []string{_Hardware.TableName(), "h"}, "i.hardware_id = h.id").
			Join("LEFT", []string{_Software.TableName(), "f"}, "i.software_id = f.id")

		if userID != snowflake.ID(0) {
			db.Where("s.user_id = ?", userID)
		}

		exist, err := db.Where("s.id = ?", sessionID).
			Where("s.deleted = ?", false).
			Get(sessionDetail)
		if err != nil {
			return errors.Wrap(err, "dao")
		}
		if !exist {
			exists = false
			return nil
		}

		return nil
	})
	if err != nil {
		return nil, exists, err
	}

	return sessionDetail, exists, nil
}

var shouldPaidSessionStatuses = []schema.SessionStatus{schema.SessionStarting, schema.SessionStarted,
	schema.SessionClosing, schema.SessionClosed,
	schema.SessionPoweringOff, schema.SessionPowerOff, schema.SessionPoweringOn, schema.SessionRebooting}

// ListSessionShouldDoPostPaid 认为 STARTING | STARTED | CLOSING | CLOSED状态的session, 有账户的，且session未完全结束支付的，需要被计费。
func ListSessionShouldDoPostPaid(ctx context.Context) ([]*models.SessionWithDetail, error) {
	sessionDetails := make([]*models.SessionWithDetail, 0)
	err := with.DefaultSession(ctx, func(db *xorm.Session) error {
		return db.Table(_Session.TableName()).Alias("s").
			Join("INNER", []string{_Instance.TableName(), "i"}, "s.instance_id = i.id").
			Join("LEFT", []string{_Hardware.TableName(), "h"}, "i.hardware_id = h.id").
			Join("LEFT", []string{_Software.TableName(), "f"}, "i.software_id = f.id").
			In("s.status", shouldPaidSessionStatuses).
			Where("s.is_paid_finished = ?", false).
			Where("s.account_id != ?", 0).
			Where("s.charge_type = ?", schema.PostPaid).
			// 存在场景：调用 RunInstance 接口失败（网络抖动/ak替换等其余情况），数据库里记录已经创建，此时不会更新start_time
			// 故计费时应过滤掉start_time为空的数据
			Where("s.start_time IS NOT NULL").
			Find(&sessionDetails)
	})
	if err != nil {
		return nil, errors.Wrap(err, "dao")
	}

	return sessionDetails, nil
}

func MarkSessionPaidFinished(ctx context.Context, sessionId snowflake.ID) error {
	return with.DefaultSession(ctx, func(db *xorm.Session) error {
		_, err := db.ID(sessionId).
			Where("is_paid_finished = ?", false).
			Cols("is_paid_finished").
			Update(&models.Session{
				IsPaidFinished: true,
			})
		return err
	})
}

func GetSessionWithInstanceBySessionId(ctx context.Context, sessionId, userId snowflake.ID, lock bool) (*models.SessionWithInstance, bool, error) {
	si := new(models.SessionWithInstance)
	var exist bool
	err := with.DefaultSession(ctx, func(db *xorm.Session) error {
		if lock {
			db.ForUpdate()
		}

		db = db.Table(_Session.TableName()).Alias("s").
			Join("INNER", []string{_Instance.TableName(), "i"}, "s.instance_id = i.id").
			Where("s.id = ?", sessionId)
		if userId != 0 {
			db = db.Where("s.user_id = ?", userId)
		}

		var e error
		exist, e = db.Get(si)
		return e
	})
	if err != nil {
		return nil, false, errors.Wrap(err, "dao")
	}

	return si, exist, nil
}
