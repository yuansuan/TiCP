package dao

import (
	"fmt"
	"time"

	"google.golang.org/grpc/status"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/internal/user/dao/model"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
)

type CertificateDaoImpl struct {
}

func NewCertificateDaoImpl() *CertificateDaoImpl {
	return &CertificateDaoImpl{}
}

func (dao *CertificateDaoImpl) Add(cert *model.OpenapiUserCertificate) (snowflake.ID, error) {
	session := GetSession()
	defer session.Close()

	cert.CreatedAt = time.Now()
	cert.UpdatedAt = time.Now()

	_, exist, err := dao.GetByUserID(cert.UserId)
	if err != nil {
		return 0, err
	}

	if exist {
		msg := fmt.Sprintf("user:[%v] exist cert", cert.UserId)
		return 0, status.Error(errcode.ErrUserOpenapiCertAlreadyExist, msg)
	}

	node, err := snowflake.GetInstance()
	if err != nil {
		msg := fmt.Sprintf("create user openapi cert failed, err:[%v]", err)
		return 0, status.Error(errcode.ErrUserOpenapiCertCreatedFailed, msg)
	}

	cert.Id = node.Generate()
	_, err = session.Insert(cert)
	if err != nil {
		msg := fmt.Sprintf("create user openapi cert failed, err:[%v]", err)
		return 0, status.Error(errcode.ErrUserOpenapiCertCreatedFailed, msg)
	}

	return cert.Id, nil
}

func (dao *CertificateDaoImpl) DelByUserID(userID snowflake.ID) error {
	session := GetSession()
	defer session.Close()

	_, err := session.Where("user_id = ?", userID).Delete(&model.OpenapiUserCertificate{})

	if err != nil {
		msg := fmt.Sprintf("delete user openapi cert failed, userID:[%v], err:[%v]", userID, err)
		return status.Error(errcode.ErrUserOpenapiCertDeleteFailed, msg)
	}
	return err
}

func (dao *CertificateDaoImpl) GetByUserID(userID snowflake.ID) (cert model.OpenapiUserCertificate, exist bool, err error) {
	session := GetSession()
	defer session.Close()

	exist, err = session.Where("user_id = ?", userID).Get(&cert)

	if err != nil {
		err = status.Error(errcode.ErrUserOpenapiCertGetFailed, err.Error())
	}

	return
}

func (dao *CertificateDaoImpl) CheckCert(certificate string) (*model.User, bool, error) {
	session := GetSession()
	defer session.Close()

	user := &model.User{}
	session.NoAutoCondition(true)
	exist, err := session.Select("user.id, user.name").Table("openapi_user_certificate").Alias("cert").
		Join("INNER", "user", "user.id = cert.user_id").
		Where("cert.certificate = ? and user.enable_openapi = 1 and user.is_deleted = '0001-01-01 00:00:00'", certificate).Get(user)

	return user, exist, err
}
