package v1

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	gomock "github.com/golang/mock/gomock"
	iam_api "github.com/yuansuan/ticp/common/project-root-iam/iam-api"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/pkg/common/snowflake"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/store"

	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/store/dao"
)

func (s *Suite) Test_listPolicy() {
	s.mockPolicyStore.EXPECT().List(gomock.Any(), gomock.Eq("4x7nt47MXA9"), gomock.Eq(0), gomock.Eq(10)).Return(s.policies, nil)

	type fields struct {
		store store.Factory
	}

	type args struct {
		ctx    context.Context
		useID  string
		offset int
		limit  int
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*dao.Policy
		wantErr bool
	}{
		{
			name: "test_list_policy",
			fields: fields{
				store: s.mockFactory,
			},
			args: args{
				ctx:    context.Background(),
				useID:  "4x7nt47MXA9",
				offset: 0,
				limit:  10,
			},
			want:    s.policies,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			r := &policyService{
				store: tt.fields.store,
			}
			got, err := r.ListPolicy(tt.args.ctx, tt.args.useID, tt.args.offset, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("policyService.ListPolicy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("policyService.ListPolicy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func (s *Suite) Test_getByPolicyName() {
	s.mockPolicyStore.EXPECT().Get(gomock.Any(), gomock.Eq("4x7nt47MXA9"), gomock.Eq("YS_CloudStorageAllAccess")).Return(s.policies[0], nil)

	type fields struct {
		store store.Factory
	}

	type args struct {
		ctx        context.Context
		userID     string
		policyName string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *dao.Policy
		wantErr bool
	}{
		{
			name: "test_get_policy",
			fields: fields{
				store: s.mockFactory,
			},
			args: args{
				ctx:        context.Background(),
				userID:     "4x7nt47MXA9",
				policyName: "YS_CloudStorageAllAccess",
			},
			want:    s.policies[0],
			wantErr: false,
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			r := &policyService{
				store: tt.fields.store,
			}
			got, err := r.GetByPolicyName(tt.args.ctx, tt.args.userID, tt.args.policyName)
			if (err != nil) != tt.wantErr {
				t.Errorf("policyService.getByPolicyName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("policyService.getByPolicyName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func (s *Suite) Test_addPolicy() {
	s.mockPolicyStore.EXPECT().Create(gomock.Any(), s.policies[0]).Return(false, nil)

	type fields struct {
		store store.Factory
	}

	type args struct {
		ctx    context.Context
		userID string
		req    *iam_api.AddPolicyRequest
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "test_add_policy",
			fields: fields{
				store: s.mockFactory,
			},
			args: args{
				ctx:    context.Background(),
				userID: "4x7nt47MXA9",
				req: &iam_api.AddPolicyRequest{
					PolicyName: "YS_CloudStorageAllAccess",
					Version:    "v1",
					Resources:  []string{"yrn:ys:cs::4x7nt47MXA9:path/<.*>"},
					Actions:    []string{"<.*>"},
					Effect:     "allow",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			r := &policyService{
				store: tt.fields.store,
			}
			if err := r.AddPolicy(tt.args.ctx, tt.args.userID, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("policyService.addPolicy() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func idEq(exptected snowflake.ID) gomock.Matcher {
	return idMatcher{exptected}
}

type idMatcher struct {
	expected snowflake.ID
}

func (m idMatcher) Matches(x interface{}) bool {
	actual, ok := x.(snowflake.ID)
	if !ok {
		return false
	}
	return actual == m.expected
}

func (m idMatcher) String() string {
	return fmt.Sprintf("is %v", m.expected)
}

func (s *Suite) Test_updatePolicy() {
	s.mockPolicyStore.EXPECT().Get(gomock.Any(), gomock.Eq("4x7nt47MXA9"), gomock.Eq("YS_HpcStorageAllAccess")).Return(s.policies[1], nil)
	s.mockPolicyStore.EXPECT().Update(gomock.Any(), idEq(1673584046262718465), gomock.Any()).Return(nil)

	type fields struct {
		store store.Factory
	}

	type args struct {
		ctx    context.Context
		userID string
		req    *iam_api.AddPolicyRequest
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "test_update_policy",
			fields: fields{store: s.mockFactory},
			args: args{
				ctx:    context.Background(),
				userID: "4x7nt47MXA9",
				req: &iam_api.AddPolicyRequest{
					PolicyName: "YS_HpcStorageAllAccess",
					Resources:  []string{"yrn:ys:cc::4x7nt47MXA9:path/<.*>"},
					Actions:    []string{"<.*>"},
					Effect:     "allow",
					Version:    "v1",
				},
			},
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			r := &policyService{
				store: tt.fields.store,
			}
			if err := r.UpdatePolicy(tt.args.ctx, tt.args.userID, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("roleService.UpdateRole() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
