package v1

import (
	"context"
	"database/sql/driver"
	"fmt"
	"regexp"

	"testing"

	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/authorization"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/store"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/store/dao"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/store/fake"
	mysqlstore "github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/store/mysql"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type AnyID struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyID) Match(v driver.Value) bool {
	_, ok := v.(int64)
	return ok
}

type AnyString struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyString) Match(v driver.Value) bool {
	_, ok := v.(string)
	return ok
}

type AnyPolicyShadow struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyPolicyShadow) Match(v driver.Value) bool {
	_, ok := v.(string)
	return ok
}

func policySliceEq(expected []*dao.Policy) gomock.Matcher {
	return policySliceMatcher{expected}
}

// implement Matcher interface, work for gomock, which can't compare slice of pointer
type policySliceMatcher struct {
	expected []*dao.Policy
}

func (m policySliceMatcher) Matches(x interface{}) bool {
	actual, ok := x.([]*dao.Policy)
	if !ok || len(actual) != len(m.expected) {
		return false
	}
	for i, p := range actual {
		if p.Policy.Policy.Resources[0] != m.expected[i].Policy.Policy.Resources[0] {
			return false
		}
	}
	return true
}

func (m policySliceMatcher) String() string {
	return fmt.Sprintf("is equal to %v", m.expected)
}

func policyRoleSliceEq(expected []*dao.RolePolicyRelation) gomock.Matcher {
	return policyRoleSliceMatcher{expected}
}

// implement Matcher interface, work for gomock, which can't compare slice of pointer
type policyRoleSliceMatcher struct {
	expected []*dao.RolePolicyRelation
}

func (m policyRoleSliceMatcher) Matches(x interface{}) bool {
	actual, ok := x.([]*dao.RolePolicyRelation)
	if !ok || len(actual) != len(m.expected) {
		return false
	}
	for i, p := range actual {
		if p.PolicyId != m.expected[i].PolicyId || p.RoleId != m.expected[i].RoleId {
			return false
		}
	}
	return true
}

func (m policyRoleSliceMatcher) String() string {
	return fmt.Sprintf("is equal to %v", m.expected)
}

func roleEq(expected *dao.Role) gomock.Matcher {
	return roleMatcher{expected}
}

// implement Matcher interface, work for gomock, which can't compare slice of pointer
type roleMatcher struct {
	expected *dao.Role
}

func (m roleMatcher) Matches(x interface{}) bool {
	actual, ok := x.(*dao.Role)
	if !ok {
		return false
	}
	return actual.RoleName == m.expected.RoleName && actual.UserId == m.expected.UserId
}

func (m roleMatcher) String() string {
	return fmt.Sprintf("is equal to %v", m.expected)
}

type Suite struct {
	suite.Suite
	mockFactory *store.MockFactory

	mockPolicyStore *store.MockPolicyStore
	policies        []*dao.Policy

	mockRoleStore *store.MockRoleStore
	roles         []*dao.Role

	mockSecretStore *store.MockSecretStore
	secrets         []*dao.Secret

	mockRolePolicyRelationStore *store.MockRolePolicyRelationStore
	rolePolicyRelations         []*dao.RolePolicyRelation

	mockPolicyCheck *authorization.MockPolicyCheck

	db   *gorm.DB
	mock sqlmock.Sqlmock
}

func (s *Suite) SetupSuite() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	s.policies = fake.FakePolicies(10)
	s.roles = fake.FakeRoles(10)
	s.secrets = fake.FakeSecrets(10)
	s.rolePolicyRelations = fake.FakeRolePolicyRelations(10)
	s.mockFactory = store.NewMockFactory(ctrl)
	s.mockPolicyStore = store.NewMockPolicyStore(ctrl)
	s.mockFactory.EXPECT().Policies().AnyTimes().Return(s.mockPolicyStore)

	s.mockRoleStore = store.NewMockRoleStore(ctrl)
	s.mockFactory.EXPECT().Roles().AnyTimes().Return(s.mockRoleStore)

	s.mockSecretStore = store.NewMockSecretStore(ctrl)
	s.mockFactory.EXPECT().Secrets().AnyTimes().Return(s.mockSecretStore)

	s.mockRolePolicyRelationStore = store.NewMockRolePolicyRelationStore(ctrl)
	s.mockFactory.EXPECT().RolePolicyRelations().AnyTimes().Return(s.mockRolePolicyRelationStore)

	s.mockPolicyCheck = authorization.NewMockPolicyCheck(ctrl)

	LoadDefaultRole()
}

func TestPolicy(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) Test_service_addRoleForPlatform() {
	s.mockRoleStore.EXPECT().Create(gomock.Any(), roleEq(s.roles[2])).Return(false, nil)
	type fields struct {
		store store.Factory
	}
	type args struct {
		c        *gin.Context
		roleYrn  string
		roleName string
		userID   string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "test",
			fields: fields{
				store: s.mockFactory,
			},
			args: args{
				c:        nil,
				roleYrn:  "yrn:ys:iam::4x7nt47MXA9:role/YS_CloudComputeRole",
				roleName: "YS_CloudComputeRole",
				userID:   "4x7nt47MXA9",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			svc := &service{
				store: tt.fields.store,
				p:     s.mockPolicyCheck,
			}
			if err := svc.addRoleForPlatform(tt.args.c, tt.args.roleYrn, tt.args.roleName, tt.args.userID); (err != nil) != tt.wantErr {
				t.Errorf("service.addRoleForPlatform() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func (s *Suite) Test_service_loadManagedPolicies() {
	s.mockPolicyStore.EXPECT().BatchCreate(gomock.Any(), policySliceEq([]*dao.Policy{s.policies[0], s.policies[1]})).Return(nil)
	type fields struct {
		store store.Factory
	}
	type args struct {
		c        *gin.Context
		roleName string
		userID   string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "test",
			fields: fields{
				store: s.mockFactory,
			},
			args: args{
				c:        nil,
				roleName: "YS_CloudComputeRole",
				userID:   "4x7nt47MXA9",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			svc := &service{
				store: tt.fields.store,
				p:     s.mockPolicyCheck,
			}
			if err := svc.loadManagedPolicies(tt.args.c, tt.args.roleName, tt.args.userID); (err != nil) != tt.wantErr {
				t.Errorf("service.loadManagedPolicies() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func (s *Suite) Test_attachPolicyToRole() {
	gomock.InOrder(
		s.mockPolicyStore.EXPECT().ListByNameAndUserId(gomock.Any(), gomock.Any(), gomock.Any()).Return([]*dao.Policy{s.policies[0], s.policies[1]}, nil),
		s.mockRoleStore.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(s.roles[2], nil),
	)
	s.mockRolePolicyRelationStore.EXPECT().CreateBatch(gomock.Any(), policyRoleSliceEq([]*dao.RolePolicyRelation{s.rolePolicyRelations[0], s.rolePolicyRelations[1]})).Return(nil)

	type fields struct {
		store store.Factory
	}
	type args struct {
		c        *gin.Context
		roleName string
		userID   string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "test",
			fields: fields{
				store: s.mockFactory,
			},
			args: args{
				c:        nil,
				roleName: "YS_CloudComputeRole",
				userID:   "4x7nt47MXA9",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			s := &service{
				store: tt.fields.store,
				p:     s.mockPolicyCheck,
			}
			if err := s.attachPolicyToRole(tt.args.c, tt.args.roleName, tt.args.userID); (err != nil) != tt.wantErr {
				t.Errorf("service.attachPolicyToRole() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// case1 add role failed
func (s *Suite) Test_createPlatformRoleCase1() {
	db, mock, err := sqlmock.New()
	if err != nil {
		s.T().Fatal(err)
	}
	s.mock = mock
	defer db.Close()

	s.db, err = gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		s.T().Fatal(err)
	}

	s.mock.ExpectBegin()

	s.mock.ExpectExec(
		regexp.QuoteMeta("insert into role")).
		WithArgs(AnyID{}, "4x7nt47MXA9", "YS_CloudComputeRole", AnyString{}, AnyPolicyShadow{}).
		WillReturnError(fmt.Errorf("mock rollback"))

	s.mock.ExpectRollback()

	factory, err := mysqlstore.SetMySQLFactory(s.db)
	if err != nil {
		s.T().Fatal(err)
	}

	svc := &service{
		store: factory,
	}

	err = svc.createPlatformRole(context.Background(), "4x7nt47MXA9", "YS_CloudComputeRole", "4x7nt47MXA9")
	if err == nil {
		s.T().Fatal("expected error")
	}
	if err := s.mock.ExpectationsWereMet(); err != nil {
		s.T().Fatal(err)
	}
}

// case2 add role success , add policies failed
func (s *Suite) Test_createPlatformRoleCase2() {
	db, mock, err := sqlmock.New()
	if err != nil {
		s.T().Fatal(err)
	}
	s.mock = mock
	defer db.Close()

	s.db, err = gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		s.T().Fatal(err)
	}

	s.mock.ExpectBegin()

	s.mock.ExpectExec(
		regexp.QuoteMeta("insert into role")).
		WithArgs(AnyID{}, "4x7nt47MXA9", "YS_CloudComputeRole", AnyString{}, AnyPolicyShadow{}).
		WillReturnResult(sqlmock.NewResult(0, 1))

	s.mock.ExpectExec(
		// regexp.QuoteMeta("insert into policy (id, userId, policyName, statementShadow) values (?,?, ?, ?), (?,?, ?, ?)")).
		regexp.QuoteMeta("insert into policy")).
		WithArgs(
			AnyID{}, "4x7nt47MXA9", "YS_CloudStorageAllAccess", AnyPolicyShadow{},
			AnyID{}, "4x7nt47MXA9", "YS_HpcStorageAllAccess", AnyPolicyShadow{}).
		WillReturnError(fmt.Errorf("mock case2 rollback"))

	s.mock.ExpectRollback()

	factory, err := mysqlstore.SetMySQLFactory(s.db)
	if err != nil {
		s.T().Fatal(err)
	}

	svc := &service{
		store: factory,
	}

	err = svc.createPlatformRole(context.Background(), "4x7nt47MXA9", "YS_CloudComputeRole", "4x7nt47MXA9")
	if err == nil {
		s.T().Fatal("expected error")
	}
	if err := s.mock.ExpectationsWereMet(); err != nil {
		s.T().Fatal(err)
	}
}

// case3 add role success , add policies success, find policies failed
func (s *Suite) Test_createPlatformRoleCase3() {
	db, mock, err := sqlmock.New()
	if err != nil {
		s.T().Fatal(err)
	}
	s.mock = mock
	defer db.Close()

	s.db, err = gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		s.T().Fatal(err)
	}

	s.mock.ExpectBegin()

	s.mock.ExpectExec(
		regexp.QuoteMeta("insert into role")).
		WithArgs(AnyID{}, "4x7nt47MXA9", "YS_CloudComputeRole", AnyString{}, AnyPolicyShadow{}).
		WillReturnResult(sqlmock.NewResult(0, 1))

	s.mock.ExpectExec(
		regexp.QuoteMeta("insert into policy")).
		WithArgs(
			AnyID{}, "4x7nt47MXA9", "YS_CloudStorageAllAccess", AnyPolicyShadow{},
			AnyID{}, "4x7nt47MXA9", "YS_HpcStorageAllAccess", AnyPolicyShadow{}).
		WillReturnResult(sqlmock.NewResult(0, 2))

	s.mock.ExpectQuery("^SELECT (.+) FROM policy").WillReturnError(fmt.Errorf("mock case3 rollback"))

	s.mock.ExpectRollback()

	factory, err := mysqlstore.SetMySQLFactory(s.db)
	if err != nil {
		s.T().Fatal(err)
	}

	svc := &service{
		store: factory,
	}

	err = svc.createPlatformRole(context.Background(), "4x7nt47MXA9", "YS_CloudComputeRole", "4x7nt47MXA9")
	if err == nil {
		s.T().Fatal("expected error")
	}
	if err := s.mock.ExpectationsWereMet(); err != nil {
		s.T().Fatal(err)
	}
}

// case4 add role success , add policies success, find policies success, find role failed
func (s *Suite) Test_createPlatformRoleCase4() {
	db, mock, err := sqlmock.New()
	if err != nil {
		s.T().Fatal(err)
	}
	s.mock = mock
	defer db.Close()

	s.db, err = gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		s.T().Fatal(err)
	}

	s.mock.ExpectBegin()

	s.mock.ExpectExec(
		regexp.QuoteMeta("insert into role")).
		WithArgs(AnyID{}, "4x7nt47MXA9", "YS_CloudComputeRole", AnyString{}, AnyPolicyShadow{}).
		WillReturnResult(sqlmock.NewResult(0, 1))

	s.mock.ExpectExec(
		regexp.QuoteMeta("insert into policy")).
		WithArgs(
			AnyID{}, "4x7nt47MXA9", "YS_CloudStorageAllAccess", AnyPolicyShadow{},
			AnyID{}, "4x7nt47MXA9", "YS_HpcStorageAllAccess", AnyPolicyShadow{}).
		WillReturnResult(sqlmock.NewResult(0, 2))

	policyRows := sqlmock.NewRows([]string{"id", "policyName"}).
		AddRow(1122333, "YS_CloudStorageAllAccess").
		AddRow(4455666, "YS_HpcStorageAllAccess")

	s.mock.ExpectQuery("^SELECT (.+) FROM policy").WillReturnRows(policyRows)

	s.mock.ExpectQuery("SELECT").WithArgs("4x7nt47MXA9", "YS_CloudComputeRole").WillReturnError(fmt.Errorf("mock case4 rollback"))
	s.mock.ExpectRollback()

	factory, err := mysqlstore.SetMySQLFactory(s.db)
	if err != nil {
		s.T().Fatal(err)
	}

	svc := &service{
		store: factory,
	}

	err = svc.createPlatformRole(context.Background(), "4x7nt47MXA9", "YS_CloudComputeRole", "4x7nt47MXA9")
	if err == nil {
		s.T().Fatal("expected error")
	}
	if err := s.mock.ExpectationsWereMet(); err != nil {
		s.T().Fatal(err)
	}
}

// case5 add role success , add policies success, find policies success, find role success, attach policies to role failed
func (s *Suite) Test_createPlatformRoleCase5() {
	db, mock, err := sqlmock.New()
	if err != nil {
		s.T().Fatal(err)
	}
	s.mock = mock
	defer db.Close()

	s.db, err = gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		s.T().Fatal(err)
	}

	s.mock.ExpectBegin()

	s.mock.ExpectExec(
		regexp.QuoteMeta("insert into role")).
		WithArgs(AnyID{}, "4x7nt47MXA9", "YS_CloudComputeRole", AnyString{}, AnyPolicyShadow{}).
		WillReturnResult(sqlmock.NewResult(0, 1))

	s.mock.ExpectExec(
		regexp.QuoteMeta("insert into policy")).
		WithArgs(
			AnyID{}, "4x7nt47MXA9", "YS_CloudStorageAllAccess", AnyPolicyShadow{},
			AnyID{}, "4x7nt47MXA9", "YS_HpcStorageAllAccess", AnyPolicyShadow{}).
		WillReturnResult(sqlmock.NewResult(0, 2))

	policyRows := sqlmock.NewRows([]string{"id", "policyName"}).
		AddRow(1122333, "YS_CloudStorageAllAccess").
		AddRow(4455666, "YS_HpcStorageAllAccess")

	s.mock.ExpectQuery("^SELECT (.+) FROM policy").WillReturnRows(policyRows)

	roleRow := sqlmock.NewRows([]string{"id"}).
		AddRow(778899)
	s.mock.ExpectQuery("SELECT").WithArgs("4x7nt47MXA9", "YS_CloudComputeRole").WillReturnRows(roleRow)

	s.mock.ExpectExec(
		regexp.QuoteMeta("insert into role_policy_relation")).
		WithArgs(
			AnyID{}, 778899, 1122333,
			AnyID{}, 778899, 4455666).
		WillReturnError(fmt.Errorf("mock case5 rollback"))

	s.mock.ExpectRollback()

	factory, err := mysqlstore.SetMySQLFactory(s.db)
	if err != nil {
		s.T().Fatal(err)
	}

	svc := &service{
		store: factory,
	}

	err = svc.createPlatformRole(context.Background(), "4x7nt47MXA9", "YS_CloudComputeRole", "4x7nt47MXA9")
	if err == nil {
		s.T().Fatal("expected error")
	}
	if err := s.mock.ExpectationsWereMet(); err != nil {
		s.T().Fatal(err)
	}
}

// case6 add role success , add policies success, find policies success, find role success, attach policies to role success
func (s *Suite) Test_createPlatformRoleCase6() {
	db, mock, err := sqlmock.New()
	if err != nil {
		s.T().Fatal(err)
	}
	s.mock = mock
	defer db.Close()

	s.db, err = gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		s.T().Fatal(err)
	}

	s.mock.ExpectBegin()

	s.mock.ExpectExec(
		regexp.QuoteMeta("insert into role")).
		WithArgs(AnyID{}, "4x7nt47MXA9", "YS_CloudComputeRole", AnyString{}, AnyPolicyShadow{}).
		WillReturnResult(sqlmock.NewResult(0, 1))

	s.mock.ExpectExec(
		regexp.QuoteMeta("insert into policy")).
		WithArgs(
			AnyID{}, "4x7nt47MXA9", "YS_CloudStorageAllAccess", AnyPolicyShadow{},
			AnyID{}, "4x7nt47MXA9", "YS_HpcStorageAllAccess", AnyPolicyShadow{}).
		WillReturnResult(sqlmock.NewResult(0, 2))

	policyRows := sqlmock.NewRows([]string{"id", "policyName"}).
		AddRow(1122333, "YS_CloudStorageAllAccess").
		AddRow(4455666, "YS_HpcStorageAllAccess")

	s.mock.ExpectQuery("^SELECT (.+) FROM policy").WillReturnRows(policyRows)

	roleRow := sqlmock.NewRows([]string{"id"}).
		AddRow(778899)

	s.mock.ExpectQuery("SELECT").WithArgs("4x7nt47MXA9", "YS_CloudComputeRole").WillReturnRows(roleRow)

	s.mock.ExpectExec(
		regexp.QuoteMeta("insert into role_policy_relation")).
		WithArgs(
			AnyID{}, 778899, 1122333,
			AnyID{}, 778899, 4455666).
		WillReturnResult(sqlmock.NewResult(0, 2))

	s.mock.ExpectCommit()

	factory, err := mysqlstore.SetMySQLFactory(s.db)
	if err != nil {
		s.T().Fatal(err)
	}

	svc := &service{
		store: factory,
	}

	err = svc.createPlatformRole(context.Background(), "4x7nt47MXA9", "YS_CloudComputeRole", "4x7nt47MXA9")

	if err != nil {
		s.Assert().Equal(nil, err)
	}

	if err := s.mock.ExpectationsWereMet(); err != nil {
		s.T().Fatal(err)
	}
}

func (s *Suite) Test_isAllowedSTS() {
	s.mockPolicyCheck.EXPECT().DoPoliciesAllow(gomock.Any(), gomock.Any()).Return(true, nil)
	s.mockPolicyCheck.EXPECT().DoPoliciesAllow(gomock.Any(), gomock.Any()).Return(false, nil)
	s.mockPolicyCheck.EXPECT().DoPoliciesAllow(gomock.Any(), gomock.Any()).Return(false, nil)
	type fields struct {
		store store.Factory
	}
	type args struct {
		r      *dao.Role
		userID string
		yrn    string
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "total match",
			args: args{
				r:      s.roles[2],
				userID: "YS_CloudCompute",
				yrn:    "yrn:ys:iam::4x7nt47MXA9:role/YS_CloudComputeRole",
			},
			want: true,
		},
		{
			name: "user not match",
			args: args{
				r:      s.roles[2],
				userID: "4x7nt47MXA9",
				yrn:    "yrn:ys:iam::4x7nt47MXA9:role/YS_CloudComputeRole",
			},
			want: false,
		},
		{
			name: "resource not match",
			args: args{
				r:      s.roles[2],
				userID: "YS_CloudCompute",
				yrn:    "ys:iam:::user/4x7nt47MXA9",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			svc := &service{
				store: tt.fields.store,
				p:     s.mockPolicyCheck,
			}
			if got, _ := svc.isAllowSTS(tt.args.r, tt.args.userID, tt.args.yrn); got != tt.want {
				t.Errorf("isAllowedSTS() = %v, want %v", got, tt.want)
			}
		})
	}
}
