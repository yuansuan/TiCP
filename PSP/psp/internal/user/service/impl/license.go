package impl

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/viper"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/config"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/osutil"
	"google.golang.org/grpc/status"

	"github.com/yuansuan/ticp/PSP/psp/internal/common"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/internal/user/consts"
	"github.com/yuansuan/ticp/PSP/psp/internal/user/dto"
	"github.com/yuansuan/ticp/PSP/psp/internal/user/service"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/timeutil"
)

type LicenseServiceImpl struct{}

func NewLicenseService() service.LicenseService {
	return &LicenseServiceImpl{}
}

func (s *LicenseServiceImpl) GetMachineID(ctx context.Context) (string, error) {
	systemSerialNumber, err := GetSystemSerialNumber(ctx)
	if err != nil {
		return "", err
	}

	if systemSerialNumber == "" {
		return "", status.Errorf(errcode.ErrAuthGetMachineIDEmpty, "get machine id empty")
	}

	id := EncryptData(systemSerialNumber, consts.EncryptHashSha256)
	return id, nil
}

func (s *LicenseServiceImpl) GetLicense(ctx context.Context) (*dto.License, error) {
	return OperatorLicenseSetting(nil, false)
}

func (s *LicenseServiceImpl) UpdateLicense(ctx context.Context, license *dto.License) error {
	expired, err := CheckExpiry(license)
	if err != nil {
		return err
	}

	if expired {
		return status.Errorf(errcode.ErrAuthLicenseHasExpiredOrNotExist, "license has expired or not exist, err: %v", err)
	}

	_, err = OperatorLicenseSetting(license, true)
	if err != nil {
		return err
	}

	return nil
}

func (s LicenseServiceImpl) CheckLicenseExpired(ctx context.Context) error {
	license, err := OperatorLicenseSetting(nil, false)
	if err != nil {
		return err
	}

	expired, err := CheckExpiry(license)
	if err != nil {
		return err
	}

	if expired {
		return fmt.Errorf("license has expired or not exist")
	}

	return nil
}

func CheckExpiry(license *dto.License) (bool, error) {
	if license == nil {
		return true, nil
	}

	if license.Name != consts.LicenseProductName || license.Version != consts.LicenseProductVersion {
		return true, nil
	}

	serialNumber, err := GetSystemSerialNumber(context.Background())
	if err != nil {
		return true, err
	}

	if license.MachineID != EncryptData(serialNumber, consts.EncryptHashSha256) {
		return true, nil
	}

	data := fmt.Sprintf("%v<%v>%v<%v", license.MachineID, license.Name, license.Version, license.Expiry)
	key := EncryptData(data, consts.EncryptHashSha512)
	if key != license.Key {
		return true, nil
	}

	if license.Expiry == "" {
		return true, nil
	}

	expiry := fmt.Sprintf("%v%v%v", license.Expiry, common.Blank, common.EndTimeFormat)
	expiryDatetime, err := timeutil.ParseTimeWithFormat(expiry, common.DatetimeFormat)
	if err != nil {
		return true, err
	}

	if time.Now().Sub(expiryDatetime) >= 0 {
		return true, nil
	}

	return false, nil
}

func OperatorLicenseSetting(license *dto.License, write bool) (*dto.License, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(pwd, config.ConfigDir)
	licenseFileName := fmt.Sprintf("%v%v%v", consts.LicenseConfigName, common.Dot, common.Yaml)
	filePath := filepath.Join(configPath, licenseFileName)
	_, err = os.Stat(filePath)
	if err != nil {
		return nil, err
	}

	newViper := viper.New()
	newViper.SetConfigType(common.Yaml)
	newViper.SetConfigName(consts.LicenseConfigName)
	newViper.AddConfigPath(configPath)
	err = newViper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	if write {
		newViper.Set(consts.LicenseName, license.Name)
		newViper.Set(consts.LicenseVersion, license.Version)
		newViper.Set(consts.LicenseExpiry, license.Expiry)
		newViper.Set(consts.LicenseMachineID, license.MachineID)
		newViper.Set(consts.LicenseKey, license.Key)

		err := newViper.WriteConfig()
		if err != nil {
			return nil, fmt.Errorf("license setting write err: %v", err)
		}

		return nil, nil
	}

	licenseData := &dto.License{
		Name:      newViper.GetString(consts.LicenseName),
		Version:   newViper.GetString(consts.LicenseVersion),
		Expiry:    newViper.GetString(consts.LicenseExpiry),
		MachineID: newViper.GetString(consts.LicenseMachineID),
		Key:       newViper.GetString(consts.LicenseKey),
	}

	if licenseData.Expiry != "" {
		// 计算当前距过期还有多少天
		expiry := fmt.Sprintf("%v%v%v", licenseData.Expiry, common.Blank, common.EndTimeFormat)
		expiryDatetime, err := timeutil.ParseTimeWithFormat(expiry, common.DatetimeFormat)
		if err != nil {
			return nil, err
		}

		intervalDay := timeutil.GetTimeIntervalDay(expiryDatetime, time.Now())
		licenseData.AvailableDays = intervalDay
	}

	return licenseData, nil
}

func EncryptData(data, sha string) string {
	hash := hmac.New(sha512.New, []byte(consts.HashSecurityKey))
	if sha == consts.EncryptHashSha256 {
		hash = hmac.New(sha256.New, []byte(consts.HashSecurityKey))
	}

	hash.Write([]byte(data))
	key := fmt.Sprintf("%x", hash.Sum(nil))

	return key
}

func GetSystemSerialNumber(ctx context.Context) (string, error) {
	uuid, err := executeCommand(ctx, consts.SystemUUIDCommand)
	if err != nil {
		return "", err
	}

	serialNumber, err := executeCommand(ctx, consts.SystemSerialNumberCommand)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s+%s", serialNumber, uuid), nil
}

func executeCommand(ctx context.Context, command string) (string, error) {
	stdout, stderr, err := osutil.CommandHelper.BashWithCurrent(ctx, command)
	if err != nil {
		return "", fmt.Errorf("execute command [%v] err: %v, stderr: %v", command, err, stderr)
	}

	return strings.ReplaceAll(stdout, "\n", ""), nil
}
