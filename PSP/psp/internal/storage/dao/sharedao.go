package dao

import (
	"context"
	"fmt"
	"time"

	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"
	"google.golang.org/grpc/status"
	"xorm.io/xorm"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/internal/storage/dao/model"
	"github.com/yuansuan/ticp/PSP/psp/internal/storage/dto"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/strutil"
	"github.com/yuansuan/ticp/PSP/psp/pkg/xtype"
)

const (
	StateUnHandle = 1
	StateHandle   = 2
)

// ShareFileRecordDaoImpl ShareFileRecordDaoImpl
type ShareFileRecordDaoImpl struct {
}

func NewShareFileRecordDaoImpl() *ShareFileRecordDaoImpl {
	return &ShareFileRecordDaoImpl{}
}

func (dao *ShareFileRecordDaoImpl) ReadAll(userId int64) error {
	session := GetSession()
	defer session.Close()

	shareFileUser := &model.ShareFileUser{State: StateHandle}

	_, err := session.Where("user_id=?", userId).Cols("state").Update(shareFileUser)
	if err != nil {
		return err
	}

	return nil
}

func (dao *ShareFileRecordDaoImpl) Add(share model.ShareFileRecord) (id snowflake.ID, err error) {
	session := GetSession()
	defer session.Close()

	node, err := snowflake.GetInstance()
	if err != nil {
		msg := fmt.Sprintf("add share failed %v", err)
		return 0, status.Error(errcode.ErrFileShareFailed, msg)
	}

	share.Id = node.Generate()
	share.CreateTime = time.Now()
	share.UpdateTime = time.Now()

	_, err = session.Insert(&share)

	if err != nil {
		msg := fmt.Sprintf("add share failed %v", err)
		return 0, status.Error(errcode.ErrFileShareFailed, msg)
	}
	return share.Id, nil
}

func (dao *ShareFileRecordDaoImpl) GetFileUserList(userID int64, page *xtype.Page, filterFilePath string) (recordList []*dto.ShareRecordInfo, total int64, err error) {
	session := GetSession()
	defer session.Close()

	session.Table(&model.ShareFileUser{}).Alias("fu").Select("r.id, r.owner, r.file_path, r.create_time as share_time, fu.state, r.type as share_type").
		Join("LEFT OUTER", "share_file_record as r", "fu.share_record_id = r.id").Where("fu.user_id = ?", userID)

	if strutil.IsNotEmpty(filterFilePath) {
		session.Where("r.file_path like ?", "%"+filterFilePath+"%")
	}

	if page.Index > 0 {
		session.Limit(int(page.Size), int((page.Index-1)*page.Size))
	}

	total, err = session.Desc("r.create_time").FindAndCount(&recordList)

	return
}

func (dao *ShareFileRecordDaoImpl) AddShareFileUser(shareFileUser model.ShareFileUser) error {
	session := GetSession()
	defer session.Close()
	shareFileUser.State = StateUnHandle

	_, err := session.Insert(&shareFileUser)
	if err != nil {
		msg := fmt.Sprintf("add shareFileUser failed %v", err)
		return status.Error(errcode.ErrFileShareFailed, msg)
	}
	return err
}

func (dao *ShareFileRecordDaoImpl) Get(id snowflake.ID) (model.ShareFileRecord, bool, error) {
	session := GetSession()
	defer session.Close()

	shareInfo := model.ShareFileRecord{Id: id}
	ok, err := session.Get(&shareInfo)
	if err != nil {
		msg := fmt.Sprintf("get shareInfo failed %v", err)
		return shareInfo, false, status.Error(errcode.ErrUserGetFailed, msg)
	}

	return shareInfo, ok, nil
}

func (dao *ShareFileRecordDaoImpl) Count(userId int64, state int8) (int64, error) {
	session := GetSession()
	defer session.Close()

	session.Table(&model.ShareFileUser{}).Alias("fu").
		Join("LEFT OUTER", "share_file_record as r", "fu.share_record_id = r.id").Where("fu.user_id = ?", userId)

	if state > 0 {
		session.Where("state=?", state)
	}

	return session.Count(&model.ShareFileRecord{})
}

func (dao *ShareFileRecordDaoImpl) UpdateRecordState(userId int64, recordIds []snowflake.ID) error {
	session := GetSession()
	defer session.Close()

	shareFileUser := &model.ShareFileUser{State: StateHandle}

	_, err := session.Where("user_id=?", userId).In("share_record_id", recordIds).Cols("state").Update(shareFileUser)
	if err != nil {
		return err
	}

	return nil
}

// GetSession GetSession
func GetSession() *xorm.Session {
	ctx := context.TODO()
	return boot.MW.DefaultSession(ctx)
}
