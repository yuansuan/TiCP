package mysql

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/store/dao"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Test_cleanExpireSecretInDB(t *testing.T) {

	dsn := "lambdacal:1234yskj@tcp(0.0.0.0:3306)/ys_base_iam?charset=utf8&parseTime=true&loc=Local"

	var err error
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	ds := &datastore{db: db}
	go cleanExpireSecretInDB(ds)

	// wait for the goroutine to start
	time.Sleep(100 * time.Millisecond)

	// add a secret that should expire
	secret := &dao.Secret{
		AccessKeyId:     "testOnly1",
		AccessKeySecret: "testOnly1",
		Expiration:      time.Now().Add(-1 * time.Hour),
	}
	err = ds.Secrets().Create(context.Background(), secret)
	assert.NoError(t, err)

	// wait for the goroutine to clean up the expired secret
	time.Sleep(2 * time.Minute)

	// check that the secret was deleted
	_, err = ds.Secrets().Get(context.Background(), secret.AccessKeyId)
	assert.NoError(t, err)
}
