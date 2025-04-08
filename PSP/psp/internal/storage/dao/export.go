package dao

import (
	"github.com/yuansuan/ticp/PSP/psp/internal/storage/dao/model"
	"github.com/yuansuan/ticp/PSP/psp/internal/storage/dto"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
	"github.com/yuansuan/ticp/PSP/psp/pkg/xtype"
)

type ShareRecordDao interface {
	Add(share model.ShareFileRecord) (id snowflake.ID, err error)
	Get(id snowflake.ID) (model.ShareFileRecord, bool, error)
	AddShareFileUser(shareUser model.ShareFileUser) error
	GetFileUserList(id int64, page *xtype.Page, filterFilePath string) (recordList []*dto.ShareRecordInfo, total int64, err error)
	Count(userId int64, share int8) (int64, error)
	UpdateRecordState(userId int64, sids []snowflake.ID) error
	ReadAll(id int64) error
}
