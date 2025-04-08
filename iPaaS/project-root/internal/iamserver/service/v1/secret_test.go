package v1

import (
	context "context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/store"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/store/dao"
)

func Test_newSecrets(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFactory := store.NewMockFactory(ctrl)

	type args struct {
		s *svc
	}
	tests := []struct {
		name string
		args args
		want *secretService
	}{
		{
			name: "test",
			args: args{
				s: &svc{
					store: mockFactory,
				},
			},
			want: &secretService{
				store: mockFactory,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newSecrets(tt.args.s); got == nil {
				t.Errorf("newSecrets() = %v, want %v", got, tt.want)
			}
		})
	}
}

func (s *Suite) Test_secretService_Get() {
	s.mockSecretStore.EXPECT().Get(gomock.Any(), gomock.Eq("FL1E9NMAL7CPJUY7NJ5O")).Return(s.secrets[0], nil)

	type fields struct {
		store store.Factory
	}

	type args struct {
		ctx  context.Context
		akID string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *dao.Secret
		wantErr bool
	}{
		{
			name: "test_secret_get",
			fields: fields{
				store: s.mockFactory,
			},
			args: args{
				ctx:  context.Background(),
				akID: "FL1E9NMAL7CPJUY7NJ5O",
			},
			want:    s.secrets[0],
			wantErr: false,
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			s := &secretService{
				store: tt.fields.store,
			}
			got, err := s.Get(tt.args.ctx, tt.args.akID)
			if (err != nil) != tt.wantErr {
				t.Errorf("secretService.GetQuota() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Errorf("secretService.GetQuota() = %v, want %v", got, tt.want)
			}
		})
	}
}

func (s *Suite) Test_secretService_ListByParentUserID() {
	s.mockSecretStore.EXPECT().List(gomock.Any(), gomock.Eq("4x7nt47MXA9"), gomock.Eq(0), gomock.Eq(10)).Return(s.secrets, nil)

	type fields struct {
		store store.Factory
	}

	type args struct {
		ctx      context.Context
		userID   string
		offset   int
		pageSize int
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*dao.Secret
		wantErr bool
	}{
		{
			name: "test_secret_list_by_parent_user_id",
			fields: fields{
				store: s.mockFactory,
			},
			args: args{
				ctx:      context.Background(),
				userID:   "4x7nt47MXA9",
				offset:   0,
				pageSize: 10,
			},
			want:    s.secrets,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			s := &secretService{
				store: tt.fields.store,
			}
			got, err := s.ListByParentUserID(tt.args.ctx, tt.args.userID, tt.args.offset, tt.args.pageSize)
			if (err != nil) != tt.wantErr {
				t.Errorf("secretService.ListByParentUserID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Errorf("secretService.ListByParentUserID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func (s *Suite) Test_secretService_List() {
	s.mockSecretStore.EXPECT().ListAll(gomock.Any(), gomock.Eq(0), gomock.Eq(10)).Return(s.secrets, nil)

	type fields struct {
		store store.Factory
	}

	type args struct {
		ctx      context.Context
		offset   int
		pageSize int
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*dao.Secret
		wantErr bool
	}{
		{
			name: "test_secret_list",
			fields: fields{
				store: s.mockFactory,
			},
			args: args{
				ctx:      context.Background(),
				offset:   0,
				pageSize: 10,
			},
			want:    s.secrets,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			s := &secretService{
				store: tt.fields.store,
			}
			got, err := s.List(tt.args.ctx, tt.args.offset, tt.args.pageSize)
			if (err != nil) != tt.wantErr {
				t.Errorf("secretService.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Errorf("secretService.List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func (s *Suite) Test_secretService_GetByUserID() {
	s.mockSecretStore.EXPECT().GetByUserID(gomock.Any(), gomock.Eq("4x7nt47MXA9")).Return(s.secrets[0], nil)

	type fields struct {
		store store.Factory
	}

	type args struct {
		ctx    context.Context
		userID string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *dao.Secret
		wantErr bool
	}{
		{
			name: "test_secret_get_by_user_id",
			fields: fields{
				store: s.mockFactory,
			},
			args: args{
				ctx:    context.Background(),
				userID: "4x7nt47MXA9",
			},
			want:    s.secrets[0],
			wantErr: false,
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			s := &secretService{
				store: tt.fields.store,
			}
			got, err := s.GetByUserID(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("secretService.GetByUserID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Errorf("secretService.GetByUserID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func (s *Suite) Test_secretService_DeleteByParentUser() {
	s.mockSecretStore.EXPECT().DeleteByParentUser(gomock.Any(), gomock.Eq("FL1E9NMAL7CPJUY7NJ5O"), gomock.Eq("4x7nt47MXA9")).Return(nil)

	type fields struct {
		store store.Factory
	}

	type args struct {
		ctx    context.Context
		akID   string
		userID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "test_secret_delete_by_parent_user",
			fields: fields{
				store: s.mockFactory,
			},
			args: args{
				ctx:    context.Background(),
				akID:   "FL1E9NMAL7CPJUY7NJ5O",
				userID: "4x7nt47MXA9",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			s := &secretService{
				store: tt.fields.store,
			}
			if err := s.DeleteByParentUser(tt.args.ctx, tt.args.akID, tt.args.userID); (err != nil) != tt.wantErr {
				t.Errorf("secretService.DeleteByParentUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func (s *Suite) Test_secretService_AdminDelete() {
	s.mockSecretStore.EXPECT().Delete(gomock.Any(), gomock.Eq("FL1E9NMAL7CPJUY7NJ5O")).Return(nil)

	type fields struct {
		store store.Factory
	}

	type args struct {
		ctx  context.Context
		akID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "test_secret_admin_delete",
			fields: fields{
				store: s.mockFactory,
			},
			args: args{
				ctx:  context.Background(),
				akID: "FL1E9NMAL7CPJUY7NJ5O",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			s := &secretService{
				store: tt.fields.store,
			}
			if err := s.AdminDelete(tt.args.ctx, tt.args.akID); (err != nil) != tt.wantErr {
				t.Errorf("secretService.AdminDelete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
