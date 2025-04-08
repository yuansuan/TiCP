package authorization

import (
	"encoding/json"
	reflect "reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ory/ladon"
	"github.com/samber/lo"
	"github.com/stretchr/testify/suite"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/store"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/store/dao"
)

func TestNewAuthorizer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFactory := store.NewMockFactory(ctrl)

	type args struct {
		factory store.Factory
	}
	tests := []struct {
		name string
		args args
		want *Authorizer
	}{
		{
			name: "default",
			args: args{
				factory: mockFactory,
			},
			want: &Authorizer{
				l: &ladon.Ladon{
					AuditLogger: NewAuditLogger(mockFactory),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAuthorizer(tt.args.factory); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAuthorizer() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestAuthorizer(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) Test_Authorizer_DoPoliciesAllow() {
	s.mockPolicyAudit.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	type fields struct {
		l *ladon.Ladon
	}

	type args struct {
		request *ladon.Request
		p       []ladon.DefaultPolicy
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "YS_CloudCompute allow",
			fields: fields{
				l: &ladon.Ladon{
					AuditLogger: NewAuditLogger(s.mockFactory),
				},
			},
			args: args{
				request: &ladon.Request{
					Subject:  "YS_CloudCompute",
					Resource: "yrn:ys:cc::4x7nt47MXA9:path/*",
					Action:   "*",
					Context:  ladon.Context{},
				},
				p: convert(s.policies),
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "YS_CloudCompute(subject) allow start.sim(resource)",
			fields: fields{
				l: &ladon.Ladon{
					AuditLogger: NewAuditLogger(s.mockFactory),
				},
			},
			args: args{
				request: &ladon.Request{
					Subject:  "YS_CloudCompute",
					Resource: "yrn:ys:cc::4x7nt47MXA9:path/start.sim",
					Action:   "*",
					Context:  ladon.Context{},
				},
				p: convert(s.policies),
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "YS_CloudCompute deny err path(patd/) start.sim(resource)",
			fields: fields{
				l: &ladon.Ladon{
					AuditLogger: NewAuditLogger(s.mockFactory),
				},
			},
			args: args{
				request: &ladon.Request{
					Subject:  "YS_CloudCompute",
					Resource: "yrn:ys:cc::4x7nt47MXA9:patd/start.sim",
					Action:   "*",
					Context:  ladon.Context{},
				},
				p: convert(s.policies),
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "YS_CloudCompute deny 4x7nt47MXA8:path/start.sim(resource)",
			fields: fields{
				l: &ladon.Ladon{
					AuditLogger: NewAuditLogger(s.mockFactory),
				},
			},
			args: args{
				request: &ladon.Request{
					Subject:  "YS_CloudCompute",
					Resource: "yrn:ys:cc::4x7nt47MXA8:path/start.sim",
					Action:   "*",
					Context:  ladon.Context{},
				},
				p: convert(s.policies),
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "YS_CloudCompute123 deny 4x7nt47MXA9:path/start.sim(resource)",
			fields: fields{
				l: &ladon.Ladon{
					AuditLogger: NewAuditLogger(s.mockFactory),
				},
			},
			args: args{
				request: &ladon.Request{
					Subject:  "YS_CloudCompute123",
					Resource: "yrn:ys:cc::4x7nt47MXA9:path/start.sim",
					Action:   "*",
					Context:  ladon.Context{},
				},
				p: convert(s.policies),
			},
			want:    false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			a := &Authorizer{
				l: tt.fields.l,
			}
			got, err := a.DoPoliciesAllow(tt.args.request, tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("Authorizer.DoPoliciesAllow() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Authorizer.DoPoliciesAllow() = %v, want %v", got, tt.want)
			}
		})
	}
}

func convert(ps []*dao.Policy) []ladon.DefaultPolicy {
	// deep copy ps
	var policies []*dao.Policy
	for _, p := range ps {
		ps1 := p.Policy.String()
		policy := &dao.Policy{}
		json.Unmarshal([]byte(ps1), &policy.Policy)
		policies = append(policies, policy)
	}
	dups := lo.Map(policies, func(p *dao.Policy, index int) ladon.DefaultPolicy {
		return p.Policy.Policy
	})
	return dups
}
