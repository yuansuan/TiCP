package service

import (
	"context"
	"strings"

	"google.golang.org/grpc/status"

	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/pkg/snowflake"

	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/consts"

	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/dao"
	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/dao/models"
	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/util"

	hydra_lcp "github.com/yuansuan/ticp/iPaaS/sso/protos"
)

// UserService UserService
type UserService struct {
	userDao *dao.UserDao
}

// NewUserSrv NewUserSrv
func NewUserSrv() *UserService {
	userDao := dao.NewUserDao()
	return &UserService{
		userDao: userDao,
	}
}

// Add Add
func (s *UserService) Add(ctx context.Context, user *models.SsoUser, pwd string) (err error) {
	user.PwdHash, err = util.PwdGenerate(pwd)
	if err != nil {
		return err
	}

	err = s.userDao.Add(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

// Get Get
func (s *UserService) Get(ctx context.Context, u *models.SsoUser) (bool, error) {
	return s.userDao.Get(ctx, u)
}

func (s *UserService) GetFromPhone(ctx context.Context, phone string) (info models.SsoUser, err error) {
	return s.userDao.GetFromPhone(ctx, phone)
}

// GetBatch GetBatch
func (s *UserService) GetBatch(ctx context.Context, ysid []int64) ([]*models.SsoUser, error) {
	return s.userDao.GetBatch(ctx, ysid)
}

func (s *UserService) List(ctx context.Context, offset, limit int64, name string) ([]*models.SsoUser, int64, error) {
	offset = offset - 1
	return s.userDao.List(ctx, int(offset), int(limit), name)
}

// GetID GetID
func (s *UserService) GetID(ctx context.Context, email string, phone string, wechatUnionID string) (userID int64, err error) {
	u := models.SsoUser{Email: email, Phone: phone, WechatUnionId: wechatUnionID}

	ok, err := s.userDao.Get(ctx, &u)
	if err != nil {
		return 0, err
	}
	if !ok {
		return 0, status.Errorf(consts.ErrHydraLcpDBUserNotExist, "user does not exist")
	}
	return u.Ysid, nil
}

// CheckUserExists CheckUserExists
func (s *UserService) CheckUserExists(ctx context.Context, user models.SsoUser) (err error) {
	return s.userDao.Exists(ctx, user)
}

// VerifyPasswordByPhone VerifyPasswordByPhone
func (s *UserService) VerifyPasswordByPhone(ctx context.Context, phone string, pwd string) (userID int64, err error) {
	return s.verifyPassword(ctx, "", phone, 0, "", pwd)
}

// VerifyPasswordByEmail VerifyPasswordByEmail
func (s *UserService) VerifyPasswordByEmail(ctx context.Context, email string, pwd string) (userID int64, err error) {
	return s.verifyPassword(ctx, email, "", 0, "", pwd)
}

// VerifyPasswordByUserID VerifyPasswordByUserID
func (s *UserService) VerifyPasswordByUserID(ctx context.Context, ysid int64, pwd string) (userID int64, err error) {
	return s.verifyPassword(ctx, "", "", ysid, "", pwd)
}

// VerifyPasswordByName VerifyPasswordByName
// FIXME name不是唯一键，提供name+pwd的登录方式应该是不安全的
func (s *UserService) VerifyPasswordByName(ctx context.Context, name string, pwd string) (userID int64, err error) {
	return s.verifyPassword(ctx, "", "", 0, name, pwd)
}

// Update Update
func (s *UserService) Update(ctx context.Context, user models.SsoUser) error {
	return s.userDao.Update(ctx, user)
}

// UpdatePwd UpdatePwd
func (s *UserService) UpdatePwd(ctx context.Context, user models.SsoUser, pwd string) (err error) {
	user.PwdHash, err = util.PwdGenerate(pwd)
	if err != nil {
		return err
	}
	return s.Update(ctx, user)
}

// UpdateName UpdateName
func (s *UserService) UpdateName(ctx context.Context, user models.SsoUser) (info models.SsoUser, err error) {
	err = s.userDao.UpdateName(ctx, user)
	if nil != err {
		return models.SsoUser{}, err
	}

	info, err = s.userDao.GetFromID(ctx, user.Ysid)
	return info, err
}

// UpdatePhone UpdatePhone
func (s *UserService) UpdatePhone(ctx context.Context, user models.SsoUser) (info models.SsoUser, err error) {
	err = s.userDao.UpdatePhone(ctx, user)
	if nil != err {
		return models.SsoUser{}, err
	}

	info, err = s.userDao.GetFromID(ctx, user.Ysid)
	return info, err
}

// UpdateEmail UpdateEmail
func (s *UserService) UpdateEmail(ctx context.Context, user models.SsoUser) (info models.SsoUser, err error) {
	err = s.userDao.UpdateEmail(ctx, user)
	if nil != err {
		return models.SsoUser{}, err
	}

	info, err = s.userDao.GetFromID(ctx, user.Ysid)
	return info, err
}

// UpdateWechatInfo UpdateWechatInfo
func (s *UserService) UpdateWechatInfo(ctx context.Context, user models.SsoUser) (info models.SsoUser, err error) {
	err = s.userDao.UpdateWechatInfo(ctx, user)
	if nil != err {
		return models.SsoUser{}, err
	}

	info, err = s.userDao.GetFromID(ctx, user.Ysid)
	return info, err
}

// UpdateRealName UpdateRealName
func (s *UserService) UpdateRealName(ctx context.Context, user models.SsoUser) (info models.SsoUser, err error) {
	err = s.userDao.UpdateRealName(ctx, user)
	if nil != err {
		return models.SsoUser{}, err
	}

	info, err = s.userDao.GetFromID(ctx, user.Ysid)
	return info, err
}

// UpdateHeadimg UpdateHeadimg
func (s *UserService) UpdateHeadimg(ctx context.Context, user models.SsoUser) (info models.SsoUser, err error) {
	err = s.userDao.UpdateHeadimg(ctx, user)
	if nil != err {
		return models.SsoUser{}, err
	}

	info, err = s.userDao.GetFromID(ctx, user.Ysid)
	return info, err
}

// GetReferedUsers GetReferedUsers
func (s *UserService) GetReferedUsers(ctx context.Context, userID int64) (referedList []*models.SsoUser, total int, err error) {
	return s.userDao.GetReferedUsers(ctx, userID)
}

func (s *UserService) verifyPassword(ctx context.Context, email string, phone string, userid int64, name string, pwd string) (userID int64, err error) {
	u := models.SsoUser{Email: email, Phone: phone, Ysid: userid, Name: name}
	ok, err := s.Get(ctx, &u)
	if err != nil {
		return 0, err
	}
	if !ok {
		return 0, status.Errorf(consts.ErrHydraLcpDBUserNotExist, "user does not exist")
	}
	err = util.PwdVerify(pwd, u.PwdHash)
	if err != nil {
		return 0, err
	}

	return u.Ysid, nil
}

// ModelToProtoUserInfo ...
func (s *UserService) ModelToProtoUserInfo(m *models.SsoUser) *hydra_lcp.UserInfo {
	displayUsername := getDisplayUserName(m.Name)
	return &hydra_lcp.UserInfo{
		Ysid:            snowflake.ID(m.Ysid).String(),
		Name:            m.Name,
		Phone:           m.Phone,
		Email:           m.Email,
		WechatUnionId:   m.WechatUnionId,
		WechatOpenId:    m.WechatOpenId,
		WechatNickName:  m.WechatNickName,
		HeadimgUrl:      m.HeadimgUrl,
		RealName:        m.RealName,
		UserName:        m.Name,
		UserChannel:     m.UserChannel,
		UserSource:      m.UserSource,
		UserReferer:     snowflake.ID(m.UserReferer).String(),
		DisplayUserName: displayUsername,
		CreateTime:      ModelToProtoTime(&m.CreateTime),
		Company:         m.Company,
	}
}

// getDisplayUserName ...
func getDisplayUserName(userName string) string {
	t := strings.Split(userName, ".")
	num := len(t)
	switch num {
	case 1:
		return userName
	default:
		// 处理用户名中包含”.“的情况
		// @example:  dev.test.a1 返回 test.a1
		return strings.Join(t[1:], ".")
	}
}
