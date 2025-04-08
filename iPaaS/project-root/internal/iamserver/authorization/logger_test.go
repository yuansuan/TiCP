package authorization

import (
	"fmt"
	reflect "reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ory/ladon"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/store"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/store/dao"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/store/fake"

	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
	mockFactory *store.MockFactory

	mockPolicyAudit *store.MockPolicyAuditStore
	audits          []*dao.PolicyAudit

	policies []*dao.Policy
}

func (s *Suite) SetupSuite() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()
	s.audits = fake.FakePolicyAudits(10)
	s.policies = fake.FakePolicies(10)

	s.mockFactory = store.NewMockFactory(ctrl)

	s.mockPolicyAudit = store.NewMockPolicyAuditStore(ctrl)
	s.mockFactory.EXPECT().PolicyAudits().Return(s.mockPolicyAudit).AnyTimes()
}

func TestNewAuditLogger(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFactory := store.NewMockFactory(ctrl)

	type args struct {
		factory store.Factory
	}
	tests := []struct {
		name string
		args args
		want *AuditLogger
	}{
		{
			name: "default",
			args: args{
				factory: mockFactory,
			},
			want: &AuditLogger{mockFactory},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAuditLogger(tt.args.factory); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAuditLogger() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPolicyAudit(t *testing.T) {
	suite.Run(t, new(Suite))
}

func subjectEq(expected *dao.PolicyAudit) gomock.Matcher {
	return subjectMatcher{expected}
}

// implement Matcher interface, work for gomock, which can't compare slice of pointer
type subjectMatcher struct {
	expected *dao.PolicyAudit
}

func (m subjectMatcher) Matches(x interface{}) bool {
	actual, ok := x.(*dao.PolicyAudit)
	if !ok {
		return false
	}
	return actual.Subject == m.expected.Subject
}

func (m subjectMatcher) String() string {
	return fmt.Sprintf("is equal to %v", m.expected)
}

func (s *Suite) Test_AuditLogger_LogRejectedAccessRequest() {
	s.mockPolicyAudit.EXPECT().Create(gomock.Any(), subjectEq(s.audits[0])).Return(nil)

	type fields struct {
		store store.Factory
	}

	type args struct {
		r *ladon.Request
		p ladon.Policies
		d ladon.Policies
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "default",
			fields: fields{
				store: s.mockFactory,
			},
			args: args{
				r: &ladon.Request{Subject: "YS_CloudCompute"},
				p: ladon.Policies{},
				d: ladon.Policies{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			a := &AuditLogger{
				store: tt.fields.store,
			}
			a.LogRejectedAccessRequest(tt.args.r, tt.args.p, tt.args.d)
		})
	}
}

func (s *Suite) Test_AuditLogger_LogGrantedAccessRequest() {
	s.mockPolicyAudit.EXPECT().Create(gomock.Any(), subjectEq(s.audits[0])).Return(nil)

	type fields struct {
		store store.Factory
	}

	type args struct {
		r *ladon.Request
		p ladon.Policies
		d ladon.Policies
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "default",
			fields: fields{
				store: s.mockFactory,
			},
			args: args{
				r: &ladon.Request{Subject: "YS_CloudCompute"},
				p: ladon.Policies{},
				d: ladon.Policies{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			a := &AuditLogger{
				store: tt.fields.store,
			}
			a.LogGrantedAccessRequest(tt.args.r, tt.args.p, tt.args.d)
		})
	}
}
