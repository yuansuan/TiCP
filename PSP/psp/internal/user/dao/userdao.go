package dao

import (
	"context"
	"fmt"
	"time"

	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/collection"
	"google.golang.org/grpc/status"
	"xorm.io/xorm"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/internal/user/dao/model"
	"github.com/yuansuan/ticp/PSP/psp/internal/user/util"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/strutil"
)

type UserDaoImpl struct {
}

func NewUserDaoImpl() *UserDaoImpl {
	return &UserDaoImpl{}
}

// Add Add user
func (dao UserDaoImpl) Add(user model.User) (id int64, err error) {
	session := GetSession()
	defer session.Close()

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	id, err = dao.GetIdByName(user.Name)
	if err == nil && id > 0 {
		msg := fmt.Sprintf("user exist name:%v, id:%v", user.Name, id)
		return 0, status.Error(errcode.ErrUserCreateFailedAlreadyExist, msg)
	}

	node, err := snowflake.GetInstance()
	if err != nil {
		msg := fmt.Sprintf("create user failed %v", err)
		return 0, status.Error(errcode.ErrUserCreatedFailed, msg)
	}

	user.Password = util.PasswdCrypto(user.Password)
	user.Id = node.Generate().Int64()
	_, err = session.Insert(&user)
	if err != nil {
		msg := fmt.Sprintf("create user failed %v", err)
		return 0, status.Error(errcode.ErrUserCreatedFailed, msg)
	}

	return user.Id, nil
}

// UpdateOrInsertUser ...
func (dao UserDaoImpl) UpdateOrInsertUser(user model.User) error {
	session := GetSession()
	defer session.Close()

	userModel := &model.User{Id: user.Id}
	exist, err := session.Where("id = ?", user.Id).Get(userModel)
	if err != nil {
		msg := fmt.Sprintf("Query user failed %v", err)
		return status.Error(errcode.ErrUserCreatedFailed, msg)
	}

	if exist {
		_, err = session.ID("id").Update(&user)
		if err != nil {
			msg := fmt.Sprintf("update user failed %v", err)
			return status.Error(errcode.ErrUserUpdateFailed, msg)
		}
	} else {
		_, err = dao.InsertUser(user)
		return err
	}

	return nil
}

// InsertUser ...
func (dao UserDaoImpl) InsertUser(user model.User) (*model.User, error) {
	session := GetSession()
	defer session.Close()

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	_, err := session.Insert(&user)

	if err != nil {
		msg := fmt.Sprintf("create user failed %v", err)
		return nil, status.Error(errcode.ErrUserCreatedFailed, msg)
	}

	return &user, nil
}

// GetByID GetByID
func (dao UserDaoImpl) GetByID(id int64) (model.User, bool, error) {
	return dao.Get(id)
}

// Get Get
func (dao UserDaoImpl) Get(id int64) (model.User, bool, error) {
	session := GetSession()
	defer session.Close()

	user := model.User{Id: id}
	ok, err := session.Get(&user)
	if err != nil {
		msg := fmt.Sprintf("get user failed %v", err)
		return user, false, status.Error(errcode.ErrUserGetFailed, msg)
	}

	return user, ok, nil
}

// GetIncludeDeleted GetIncludeDeleted
func (dao UserDaoImpl) GetIncludeDeleted(id int64) (model.User, bool, error) {
	session := GetSession()
	defer session.Close()

	user := model.User{Id: id}
	ok, err := session.Unscoped().Get(&user)
	if err != nil {
		msg := fmt.Sprintf("get user failed %v", err)
		return user, false, status.Error(errcode.ErrUserGetFailed, msg)
	}

	return user, ok, nil
}

// GetByUserID GetByUserID
func (dao UserDaoImpl) GetByUserID(id int64) (model.User, error) {
	session := GetSession()
	defer session.Close()

	user := model.User{Id: id}

	ok, err := session.Get(&user)
	if err != nil {
		msg := fmt.Sprintf("get user failed %v", err)
		return user, status.Error(errcode.ErrUserGetFailed, msg)
	}

	if !ok {
		msg := fmt.Sprintf("user not found %v", id)
		return user, status.Error(errcode.ErrUserNotFound, msg)
	}

	return user, nil
}

// BatchGet batch Get
func (dao UserDaoImpl) BatchGet(ids []int64) (users []*model.User, err error) {
	session := GetSession()
	defer session.Close()

	err = session.In("id", ids).Find(&users)
	if err != nil {
		msg := fmt.Sprintf("get user failed %v", err)
		return users, status.Error(errcode.ErrUserGetFailed, msg)
	}
	return users, nil
}

// GetID GetID
func (dao UserDaoImpl) GetIdByName(name string) (id int64, err error) {
	session := GetSession()
	defer session.Close()
	user := model.User{Name: name}

	ok, err := session.Get(&user)

	if err != nil {
		msg := fmt.Sprintf("get user failed %v %v", id, err)
		return -1, status.Error(errcode.ErrUserGetFailed, msg)
	}

	if !ok {
		msg := fmt.Sprintf("user not found %v", id)
		return -1, status.Error(errcode.ErrUserNotFound, msg)
	}

	return user.Id, nil
}

func (dao UserDaoImpl) GetUserByName(name string) (*model.User, error) {
	session := GetSession()
	defer session.Close()
	user := &model.User{Name: name}

	ok, err := session.Get(user)

	if err != nil {
		msg := fmt.Sprintf("get user failed %v %v", name, err)
		return nil, status.Error(errcode.ErrUserGetFailed, msg)
	}

	if !ok {
		msg := fmt.Sprintf("user not found %v", name)
		return nil, status.Error(errcode.ErrUserNotFound, msg)
	}

	return user, nil
}

// Exist Exist
func (dao UserDaoImpl) Exist(id int64) bool {
	session := GetSession()
	defer session.Close()
	ok, _ := session.ID(id).Exist(&model.User{})
	if !ok {
		return false
	}
	return true
}

// Exists Exists
func (dao UserDaoImpl) Exists(ids ...int64) error {
	session := GetSession()
	defer session.Close()
	for _, id := range ids {
		ok, err := session.ID(id).Exist(&model.User{})

		if err != nil {
			msg := fmt.Sprintf("get user failed %v %v", id, err)
			return status.Error(errcode.ErrUserGetFailed, msg)
		}

		if !ok {
			msg := fmt.Sprintf("user not found %v", id)
			return status.Error(errcode.ErrUserNotFound, msg)
		}
	}

	return nil
}

// List List
func (dao UserDaoImpl) List(enabled bool, query string, order string, desc bool, page int64, pageSize int64, isInternal bool) (users []*model.User, total int64, err error) {

	session := GetSession()
	defer session.Close()
	if !collection.Contain([]string{"name", "id", "created_at"}, order) {
		msg := fmt.Sprintf("invalid order condition %v", order)
		return nil, 0, status.Error(errcode.ErrUserQueryInvalidOrderCondition, msg)
	}

	if enabled {
		session = session.Where("enabled = ?", 1)
	}

	if !isInternal {
		session = session.Where("is_internal != 1")
	}

	if desc {
		session.Where("name Like ?", "%"+query+"%").Desc(order)
	} else {
		session.Where("name Like ?", "%"+query+"%").Asc(order)
	}

	total, err = session.Limit(int(pageSize), int((page-1)*pageSize)).FindAndCount(&users)

	return users, total, err
}

// ListAll ListAll
func (dao UserDaoImpl) ListAll(pageIndex, pageSize int64) ([]*model.User, int64, error) {
	session := GetSession()
	defer session.Close()

	users := []*model.User{}

	total, err := session.Count(&model.User{})
	if err != nil {
		return nil, 0, status.Errorf(errcode.ErrUserQueryFailed, "query user failed %v", err)
	}

	err = session.Asc("id").Limit(int(pageSize), int((pageIndex-1)*pageSize)).Find(&users)
	return users, total, err
}

// Update Update
func (dao UserDaoImpl) Update(user model.User) error {
	session := GetSession()
	defer session.Close()
	_, exists, err := dao.Get(user.Id)
	if err != nil {
		return status.Error(errcode.ErrInternalServer, err.Error())
	}
	if !exists {
		return status.Errorf(errcode.ErrUserNotFound, "user not found %v", user.Id)
	}

	_, err = session.MustCols("enabled", "email", "mobile", "enable_openapi").ID(user.Id).Update(&user)
	if err != nil {
		msg := fmt.Sprintf("user (id: %v) update failed", user.Id)
		return status.Error(errcode.ErrUserUpdateFailed, msg)
	}

	return nil
}

// Delete Delete
func (dao UserDaoImpl) Delete(id int64) error {
	session := GetSession()
	defer session.Close()

	//user := model.User{IsDeleted: true}
	//_, err := session.Unscoped().Where("id = ?", id).MustCols("is_deleted").Update(&user)
	_, err := session.ID(id).Delete(&model.User{})

	if err != nil {
		msg := fmt.Sprintf("delete user failed %v %v", id, err)
		return status.Error(errcode.ErrUserDeleted, msg)
	}
	return err
}

// GetSession GetSession
func GetSession() *xorm.Session {
	ctx := context.TODO()
	return boot.MW.DefaultSession(ctx)
}

// ListAllUser List all user
func (dao UserDaoImpl) ListAllUser() ([]*model.User, error) {
	session := GetSession()
	defer session.Close()

	users := []*model.User{}
	if err := session.Find(&users); err != nil {
		return nil, err
	}
	return users, nil
}

// ListAllAdminUser List all admin user
func (dao UserDaoImpl) ListAllAdminUser() ([]*model.User, error) {
	session := GetSession()
	defer session.Close()

	users := []*model.User{}
	if err := session.Where("is_internal=?", 1).Find(&users); err != nil {
		return nil, err
	}
	return users, nil
}

// ListUserByName List user by names
func (dao UserDaoImpl) ListUserByName(names []string) ([]model.User, error) {
	session := GetSession()
	defer session.Close()

	users := []model.User{}
	if err := session.In("name", names).Find(&users); err != nil {
		return nil, err
	}
	return users, nil
}

// ListUserLikeName List user like names
func (dao UserDaoImpl) ListUserLikeName(filterName string, ids []int64) ([]model.User, error) {
	session := GetSession()
	defer session.Close()

	users := []model.User{}

	if ids != nil && len(ids) > 0 {
		session.In("id", ids)
	}

	if strutil.IsNotEmpty(filterName) {
		session.Where("name Like ?", "%"+filterName+"%")
	}

	if err := session.Where("enabled = 1").OrderBy("name").Find(&users); err != nil {
		return nil, err
	}
	return users, nil
}

// GetUserByUserId get user by user id
func (dao UserDaoImpl) GetUserByUserId(id snowflake.ID) (*model.User, error) {
	session := GetSession()
	defer session.Close()

	user := &model.User{}
	has, err := session.Where("user_id = ?", id).Get(user)
	if err != nil {
		return nil, err
	}

	if has {
		return user, nil
	}

	return nil, nil
}

// DeleteUserByID Delete user by ids
func (dao UserDaoImpl) DeleteUserByID(ids []int64) error {
	session := GetSession()
	defer session.Close()

	if _, err := session.In("id", ids).Delete(&model.User{}); err != nil {
		return err
	}

	if _, err := session.In("user_id", ids).Delete(&model.UserOrg{}); err != nil {
		return err
	}

	return nil
}
