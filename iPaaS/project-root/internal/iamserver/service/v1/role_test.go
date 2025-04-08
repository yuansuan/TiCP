package v1

import (
	"context"
	"reflect"
	"testing"

	gomock "github.com/golang/mock/gomock"
	iam_api "github.com/yuansuan/ticp/common/project-root-iam/iam-api"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/store"

	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/store/dao"
)

func (s *Suite) Test_listRole() {
	s.mockRoleStore.EXPECT().List(gomock.Any(), gomock.Eq("4x7nt47MXA9")).Return(s.roles, nil)

	type fields struct {
		store store.Factory
	}

	type args struct {
		ctx   context.Context
		useID string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*dao.Role
		wantErr bool
	}{
		{
			name: "test_list_role",
			fields: fields{
				store: s.mockFactory,
			},
			args: args{
				ctx:   context.Background(),
				useID: "4x7nt47MXA9",
			},
			want:    s.roles,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			r := &roleService{
				store: tt.fields.store,
			}
			got, err := r.ListRole(tt.args.ctx, tt.args.useID)
			if (err != nil) != tt.wantErr {
				t.Errorf("roleService.ListRole() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("roleService.ListRole() = %v, want %v", got, tt.want)
			}
		})
	}
}

func (s *Suite) Test_getRole() {
	s.mockRoleStore.EXPECT().Get(gomock.Any(), gomock.Eq("4x7nt47MXA9"), gomock.Eq("YS_CloudComputeRole")).Return(s.roles[0], nil)

	type fields struct {
		store store.Factory
	}

	type args struct {
		ctx      context.Context
		userID   string
		roleName string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *dao.Role
		wantErr bool
	}{
		{
			name: "test_get_role",
			fields: fields{
				store: s.mockFactory,
			},
			args: args{
				ctx:      context.Background(),
				userID:   "4x7nt47MXA9",
				roleName: "YS_CloudComputeRole",
			},
			want:    s.roles[0],
			wantErr: false,
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			r := &roleService{
				store: tt.fields.store,
			}
			got, err := r.GetRole(tt.args.ctx, tt.args.userID, tt.args.roleName)
			if (err != nil) != tt.wantErr {
				t.Errorf("roleService.GetRole() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.RoleName, tt.want.RoleName) {
				t.Errorf("roleService.GetRole() = %v, want %v", got, tt.want)
			}
		})
	}
}

func (s *Suite) Test_addRole() {
	s.mockRoleStore.EXPECT().Create(gomock.Any(), gomock.Eq(s.roles[0])).Return(false, nil)

	type fields struct {
		store store.Factory
	}

	type args struct {
		ctx    context.Context
		userID string
		req    *iam_api.AddRoleRequest
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "test_add_role",
			fields: fields{store: s.mockFactory},
			args: args{
				ctx:    context.Background(),
				userID: "4TiSxuPtJEm",
				req: &iam_api.AddRoleRequest{
					RoleName: "VIPBoxRole_4x7nt47MXA7",
					TrustPolicy: &iam_api.RolePolicy{
						Principals: []string{"4x7nt47MXA7"},
						Resources:  []string{"yrn:ys:iam::4TiSxuPtJEm:role/VIPBoxRole_4x7nt47MXA7"},
						Effect:     "allow",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			r := &roleService{
				store: tt.fields.store,
			}
			if err := r.AddRole(tt.args.ctx, tt.args.userID, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("roleService.AddRole() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func (s *Suite) Test_updateRole() {
	s.mockRoleStore.EXPECT().Get(gomock.Any(), gomock.Eq("4TiSxuPtJEm"), gomock.Eq("VIPBoxRole_4x7nt47MXA7")).Return(s.roles[0], nil)
	s.mockRoleStore.EXPECT().Update(gomock.Any(), gomock.Eq(s.roles[1]), gomock.Any()).Return(nil)

	type fields struct {
		store store.Factory
	}

	type args struct {
		ctx    context.Context
		userID string
		req    *iam_api.Role
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "test_update_role",
			fields: fields{store: s.mockFactory},
			args: args{
				ctx:    context.Background(),
				userID: "4TiSxuPtJEm",
				req: &iam_api.Role{
					RoleName: "VIPBoxRole_4x7nt47MXA7",
					TrustPolicy: &iam_api.RolePolicy{
						Principals: []string{"4x7nt47MXA7"},
						Resources:  []string{"yrn:ys:iam::4TiSxuPtJEm:role/VIPBoxRole_4x7nt47MXA7"},
						Actions:    []string{"STS:AssumeRole"},
						Effect:     "allow",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			r := &roleService{
				store: tt.fields.store,
			}
			if err := r.UpdateRole(tt.args.ctx, tt.args.userID, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("roleService.UpdateRole() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func (s *Suite) Test_patchPolicy() {
	s.mockRolePolicyRelationStore.EXPECT().Create(gomock.Any(), gomock.Eq(&dao.RolePolicyRelation{
		RoleId:   1673559263525474304,
		PolicyId: 1673584046262718464,
	})).Return(false, nil)

	type fields struct {
		store store.Factory
	}

	type args struct {
		ctx    context.Context
		role   *dao.Role
		policy *dao.Policy
	}
	test := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "test_patch_policy",
			fields: fields{store: s.mockFactory},
			args: args{
				ctx:    context.Background(),
				role:   &dao.Role{ID: 1673559263525474304},
				policy: &dao.Policy{ID: 1673584046262718464},
			},
			wantErr: false,
		},
	}

	for _, tt := range test {
		s.T().Run(tt.name, func(t *testing.T) {
			r := &roleService{
				store: tt.fields.store,
			}
			if err := r.PatchPolicy(tt.args.ctx, tt.args.role, tt.args.policy); (err != nil) != tt.wantErr {
				t.Errorf("roleService.patchPolicy() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func (s *Suite) Test_detachPolicy() {
	s.mockRolePolicyRelationStore.EXPECT().Delete(gomock.Any(), gomock.Eq(&dao.RolePolicyRelation{
		RoleId:   1673559263525474304,
		PolicyId: 1673584046262718464,
	})).Return(nil)

	type fields struct {
		store store.Factory
	}

	type args struct {
		ctx    context.Context
		role   *dao.Role
		policy *dao.Policy
	}
	test := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "test_detach_policy",
			fields: fields{store: s.mockFactory},
			args: args{
				ctx:    context.Background(),
				role:   &dao.Role{ID: 1673559263525474304},
				policy: &dao.Policy{ID: 1673584046262718464},
			},
			wantErr: false,
		},
	}

	for _, tt := range test {
		s.T().Run(tt.name, func(t *testing.T) {
			r := &roleService{
				store: tt.fields.store,
			}
			if err := r.DetachPolicy(tt.args.ctx, tt.args.role, tt.args.policy); (err != nil) != tt.wantErr {
				t.Errorf("roleService.patchPolicy() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
