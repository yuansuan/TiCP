package fake

import (
	"sync"
	"time"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao"
	dbModels "github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao/models"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao/store"
)

const ResourceCount = 1000

type datastore struct {
	sync.RWMutex
	apps []*dbModels.Application
}

func (ds *datastore) Applications() dao.ApplicationDao {
	return newApplication(ds)
}

func FakeApplications(count int) []*dbModels.Application {
	apps := make([]*dbModels.Application, 0, count)
	for i := 0; i < count; i++ {
		apps = append(apps, &dbModels.Application{
			ID:   getSnowflakeID(),
			Name: "fake-app",
		})
	}
	return apps
}

func getSnowflakeID() snowflake.ID {
	// current unix time
	unix := time.Now().Unix()
	return snowflake.ID(unix)
}

var (
	fakeFactory store.Factory
	once        sync.Once
)

func GetFakeFactoryOr() (store.Factory, error) {
	once.Do(func() {
		fakeFactory = &datastore{
			apps: FakeApplications(ResourceCount),
		}
	})
	return fakeFactory, nil
}
