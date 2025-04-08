package dao

import (
	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"
	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/pkg/with"
	"golang.org/x/net/context"
	"google.golang.org/grpc/status"
	"xorm.io/xorm"

	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/consts"

	"strings"

	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/dao/models"
)

// UserDao UserDao
type UserDao struct {
}

// NewUserDao NewUserDao
func NewUserDao() *UserDao {
	return &UserDao{}
}

// Add Add
func (d *UserDao) Add(ctx context.Context, user *models.SsoUser) error {
	if user.Ysid == 0 {
		return status.Error(consts.ErrHydraLcpLackYsID, "user must have a ysid")
	}
	err := with.DefaultSession(ctx, func(db *xorm.Session) error {
		_, err := db.Nullable("email", "phone", "wechat_union_id").Insert(user)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return status.Errorf(consts.ErrHydraLcpDBDuplicatedEntry, "duplicated entry, err: %v", err.Error())
		}
		return status.Error(consts.ErrHydraLcpDBOpFail, err.Error())
	}
	return nil
}

// Get Get
func (d *UserDao) Get(ctx context.Context, user *models.SsoUser) (ok bool, err error) {
	err = with.DefaultSession(ctx, func(db *xorm.Session) error {
		ok, err = db.Get(user)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return false, status.Errorf(consts.ErrHydraLcpDBOpFail, "failed to get user, err: %v", err.Error())
	}
	return ok, nil
}

// GetBatch GetBatch
func (d *UserDao) GetBatch(ctx context.Context, ysid []int64) ([]*models.SsoUser, error) {
	if len(ysid) == 0 {
		return nil, status.Errorf(consts.ErrHydraLcpYsidEmpty, "no user id input")
	}

	session := boot.MW.DefaultSession(ctx)
	defer session.Close()

	var userList []*models.SsoUser
	err := session.In("ysid", ysid).Find(&userList)
	if err != nil {
		return nil, status.Errorf(consts.ErrHydraLcpDBOpFail, "failed to get all user information by user id, err: %v", err.Error())
	}

	if len(userList) == 0 {
		return nil, status.Errorf(consts.ErrHydraLcpDBUserNotExist, "user does not exist")
	}

	return userList, nil
}

func (d *UserDao) List(ctx context.Context, offset, limit int, name string) ([]*models.SsoUser, int64, error) {
	session := boot.MW.DefaultSession(ctx)
	defer session.Close()

	var userList []*models.SsoUser
	session.Where("name like ?", "%"+name+"%")
	session.Desc("create_time")
	err := session.Limit(limit, offset).Find(&userList)
	if err != nil {
		return nil, 0, status.Errorf(consts.ErrHydraLcpDBOpFail, "failed to get all user information, err: %v", err.Error())
	}

	// total like count
	total, err := session.Where("name like ?", "%"+name+"%").Desc("create_time").Count(&models.SsoUser{})
	if err != nil {
		return nil, 0, status.Errorf(consts.ErrHydraLcpDBOpFail, "failed to count all user information, err: %v", err.Error())
	}
	return userList, total, nil
}

func (d *UserDao) ListByName(ctx context.Context, name string) ([]*models.SsoUser, error) {
	session := boot.MW.DefaultSession(ctx)
	defer session.Close()

	var userList []*models.SsoUser
	err := session.Where("name = ?", name).Find(&userList)
	if err != nil {
		return nil, status.Errorf(consts.ErrHydraLcpDBOpFail, "failed to get all user information by name, err: %v", err.Error())
	}
	return userList, nil
}

// GetFromID GetFromID
func (d *UserDao) GetFromID(ctx context.Context, id int64) (user models.SsoUser, err error) {
	session := boot.MW.DefaultSession(ctx)
	defer session.Close()

	user = models.SsoUser{Ysid: id}

	ok, err := session.Get(&user)
	if err != nil {
		return user, status.Errorf(consts.ErrHydraLcpDBOpFail, "failed to get user(ysid: %v), err: %v", id, err.Error())
	}

	if !ok {
		return user, status.Errorf(consts.ErrHydraLcpDBUserNotExist, "user(ysid: %v) does not exist", id)
	}

	return user, nil
}

// Exists Exists
func (d *UserDao) Exists(ctx context.Context, ssoUser models.SsoUser) (err error) {
	session := boot.MW.DefaultSession(ctx)
	defer session.Close()

	ok, err := session.Exist(&ssoUser)
	if err != nil {
		return status.Errorf(consts.ErrHydraLcpDBOpFail, "failed to query user, err : %v", err.Error())
	}

	if !ok {
		return status.Errorf(consts.ErrHydraLcpDBUserNotExist, "user does not exist")
	}

	return nil
}

// GetFromEmail GetFromEmail
func (d *UserDao) GetFromEmail(ctx context.Context, email string) (user models.SsoUser, err error) {
	session := boot.MW.DefaultSession(ctx)
	defer session.Close()

	user = models.SsoUser{Email: email}

	ok, err := session.Get(&user)
	if err != nil {
		return user, status.Errorf(consts.ErrHydraLcpDBOpFail, "failed to get user(email: %v), err: %v", email, err.Error())
	}

	if !ok {
		return user, status.Errorf(consts.ErrHydraLcpDBUserNotExist, "user(email: %v) does not exist", email)
	}

	return user, nil
}

// GetFromPhone GetFromPhone
func (d *UserDao) GetFromPhone(ctx context.Context, phone string) (user models.SsoUser, err error) {
	session := boot.MW.DefaultSession(ctx)
	defer session.Close()

	user = models.SsoUser{Phone: phone}

	ok, err := session.Get(&user)
	if err != nil {
		return user, status.Errorf(consts.ErrHydraLcpDBOpFail, "failed to get user(phone: %v), err: %v", phone, err.Error())
	}

	if !ok {
		return user, status.Errorf(consts.ErrHydraLcpDBUserNotExist, "user(phone: %v) does not exist", phone)
	}

	return user, nil
}

// Update Update
func (d *UserDao) Update(ctx context.Context, user models.SsoUser) error {
	session := boot.MW.DefaultSession(ctx)
	defer session.Close()

	if user.Ysid == 0 {
		return status.Error(consts.ErrHydraLcpLackYsID, "user must have a ysid")
	}

	_, err := session.ID(user.Ysid).Update(user)
	if err != nil {
		return status.Error(consts.ErrHydraLcpDBOpFail, err.Error())
	}

	return err
}

// UpdateName UpdateName
func (d *UserDao) UpdateName(ctx context.Context, user models.SsoUser) error {
	session := boot.MW.DefaultSession(ctx)
	defer session.Close()

	if user.Ysid == 0 {
		return status.Error(consts.ErrHydraLcpLackYsID, "user must have a ysid")
	}

	_, err := session.ID(user.Ysid).Cols("name").Update(user)
	if err != nil {
		return status.Error(consts.ErrHydraLcpDBOpFail, err.Error())
	}

	return err
}

// UpdatePhone UpdatePhone
func (d *UserDao) UpdatePhone(ctx context.Context, user models.SsoUser) error {
	session := boot.MW.DefaultSession(ctx)
	defer session.Close()

	if user.Ysid == 0 {
		return status.Error(consts.ErrHydraLcpLackYsID, "user must have a ysid")
	}

	_, err := session.ID(user.Ysid).Cols("phone").Update(user)
	if err != nil {
		return status.Error(consts.ErrHydraLcpDBOpFail, err.Error())
	}

	return err
}

// UpdateEmail UpdateEmail
func (d *UserDao) UpdateEmail(ctx context.Context, user models.SsoUser) error {
	session := boot.MW.DefaultSession(ctx)
	defer session.Close()

	if user.Ysid == 0 {
		return status.Error(consts.ErrHydraLcpLackYsID, "user must have a ysid")
	}

	_, err := session.ID(user.Ysid).Cols("email").Update(user)
	if err != nil {
		return status.Error(consts.ErrHydraLcpDBOpFail, err.Error())
	}

	return err
}

// UpdateWechatInfo UpdateWechatInfo
func (d *UserDao) UpdateWechatInfo(ctx context.Context, user models.SsoUser) error {
	session := boot.MW.DefaultSession(ctx)
	defer session.Close()

	if user.Ysid == 0 {
		return status.Error(consts.ErrHydraLcpLackYsID, "user must have a ysid")
	}

	// check wechat info
	if len(user.WechatUnionId) == 0 || len(user.WechatOpenId) == 0 {
		return status.Error(consts.ErrHydraLcpLackWechatInfo, "wechat_union_id and wechat_open_id must not be empty")
	}

	// check wechat info exists
	existUnion, _ := session.Exist(&models.SsoUser{
		WechatUnionId: user.WechatUnionId,
	})
	if existUnion {
		return status.Errorf(consts.ErrHydraLcpWechatInfoExist, "wechat info already exist, union id is %v", user.WechatUnionId)
	}
	existOpen, _ := session.Exist(&models.SsoUser{
		WechatOpenId: user.WechatOpenId,
	})
	if existOpen {
		return status.Errorf(consts.ErrHydraLcpWechatInfoExist, "wechat info already exist, open id is %v", user.WechatOpenId)
	}

	_, err := session.ID(user.Ysid).Cols("wechat_union_id", "wechat_open_id", "wechat_nick_name").Update(user)
	if err != nil {
		return status.Error(consts.ErrHydraLcpDBOpFail, err.Error())
	}

	return err
}

// UpdateRealName UpdateRealName
func (d *UserDao) UpdateRealName(ctx context.Context, user models.SsoUser) error {
	session := boot.MW.DefaultSession(ctx)
	defer session.Close()

	if user.Ysid == 0 {
		return status.Error(consts.ErrHydraLcpLackYsID, "user must have a ysid")
	}

	_, err := session.ID(user.Ysid).Cols("real_name").Update(user)
	if err != nil {
		return status.Error(consts.ErrHydraLcpDBOpFail, err.Error())
	}

	return err
}

// UpdateHeadimg UpdateHeadimg
func (d *UserDao) UpdateHeadimg(ctx context.Context, user models.SsoUser) error {
	session := boot.MW.DefaultSession(ctx)
	defer session.Close()

	if user.Ysid == 0 {
		return status.Error(consts.ErrHydraLcpLackYsID, "user must have a ysid")
	}

	_, err := session.ID(user.Ysid).Cols("headimg_url").Update(user)
	if err != nil {
		return status.Error(consts.ErrHydraLcpDBOpFail, err.Error())
	}

	return err
}

// GetReferedUsers GetReferedUsers
func (d *UserDao) GetReferedUsers(ctx context.Context, refererID int64) (referedList []*models.SsoUser, total int, err error) {
	session := boot.MW.DefaultSession(ctx)
	defer session.Close()

	if refererID == 0 {
		return nil, int(0), status.Error(consts.ErrHydraLcpLackYsID, "user must have a referer ysid")
	}
	session.Table("sso_user").Alias("u")
	session.Where("u.user_referer = ?", refererID)

	if err != nil {
		logging.GetLogger(ctx).Error("err count refered users", "error", err, "referer_id", refererID)
		return referedList, int(0), status.Error(consts.ErrHydraLcpDBOpFail, "err count refered users")
	}

	err = session.Find(&referedList)
	if err != nil {
		logging.GetLogger(ctx).Error("err get refered users", "error", err, "referer_id", refererID)
		return referedList, int(0), status.Error(consts.ErrHydraLcpDBOpFail, "err get refered users")
	}
	total = len(referedList)
	logging.Default().Debugf("referedList: %v, total %v ", referedList, total)
	return referedList, total, nil
}
