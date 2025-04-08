package application

// import (
// 	"context"
// 	"github.com/golang/mock/gomock"
// 	"github.com/stretchr/testify/suite"
// 	"testing"
// 	dbModels "github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao/models"
// 	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao/store"
// 	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao/store/fake"
// 	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
// )

// type Suite struct {
// 	suite.Suite
// 	mockFactory *store.MockFactory

// 	apps         []*dbModels.Application
// 	mockAppStore *store.MockApplicationStore
// }

// func (s *Suite) SetupSuite() {
// 	ctrl := gomock.NewController(s.T())
// 	defer ctrl.Finish()

// 	s.apps = fake.FakeApplications(10)
// 	s.mockFactory = store.NewMockFactory(ctrl)
// 	s.mockAppStore = store.NewMockApplicationStore(ctrl)
// 	s.mockFactory.EXPECT().Applications().AnyTimes().Return(s.mockAppStore)
// }

// func TestPolicy(t *testing.T) {
// 	suite.Run(t, new(Suite))
// }

// func (s *Suite) Test_applicationService_AddApp() {
// 	s.mockAppStore.EXPECT().AddApp(gomock.Any(), gomock.Eq(s.apps[0])).Return(nil)
// 	type fields struct {
// 		store store.Factory
// 	}
// 	type args struct {
// 		ctx context.Context
// 		app *dbModels.Application
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		wantErr bool
// 	}{
// 		{
// 			name: "default",
// 			fields: fields{
// 				store: s.mockFactory,
// 			},
// 			args: args{
// 				ctx: context.TODO(),
// 				app: s.apps[0],
// 			},
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		s.T().Run(tt.name, func(t *testing.T) {
// 			ss := &applicationService{
// 				store: tt.fields.store,
// 			}
// 			if _, err := ss.AddApp(tt.args.ctx, tt.args.app); (err != nil) != tt.wantErr {
// 				t.Errorf("applicationService.AddApp() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }

// func (s *Suite) Test_applicationService_GetApp() {
// 	s.mockAppStore.EXPECT().GetApp(gomock.Any(), gomock.Eq(s.apps[0].Ysid)).Return(s.apps[0], nil)
// 	type fields struct {
// 		store store.Factory
// 	}
// 	type args struct {
// 		ctx context.Context
// 		id  snowflake.ID
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		want    *dbModels.Application
// 		wantErr bool
// 	}{
// 		{
// 			name: "default",
// 			fields: fields{
// 				store: s.mockFactory,
// 			},
// 			args: args{
// 				ctx: context.TODO(),
// 				id:  s.apps[0].Ysid,
// 			},
// 			want:    s.apps[0],
// 			wantErr: false,
// 		},
// 		{
// 			name: "not found",
// 			fields: fields{
// 				store: s.mockFactory,
// 			},
// 			args: args{
// 				ctx: context.TODO(),
// 				id:  100,
// 			},
// 			want:    nil,
// 			wantErr: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		s.T().Run(tt.name, func(t *testing.T) {
// 			ss := &applicationService{
// 				store: tt.fields.store,
// 			}
// 			got, err := ss.GetApp(tt.args.ctx, tt.args.id)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("applicationService.GetApp() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if got != nil && got.Ysid != tt.want.Ysid {
// 				t.Errorf("applicationService.GetApp() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
