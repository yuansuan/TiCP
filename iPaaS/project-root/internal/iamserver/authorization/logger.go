package authorization

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/ory/ladon"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/store"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/store/dao"
)

// AuditLogger outputs and cache information about granting or rejecting policies.
type AuditLogger struct {
	store store.Factory
}

// NewAuditLogger creates a AuditLogger with default parameters.
func NewAuditLogger(factory store.Factory) *AuditLogger {
	return &AuditLogger{
		store: factory,
	}
}

// LogRejectedAccessRequest write rejected subject access to log.
func (a *AuditLogger) LogRejectedAccessRequest(r *ladon.Request, p ladon.Policies, d ladon.Policies) {
	logging.Default().Infof("subject access review rejected, request: %+v, deciders: %+v", r, d)
	record := rejected(r, p, d)
	audit := &dao.PolicyAudit{
		Subject:      r.Subject,
		PolicyShadow: record.String(),
	}
	if err := a.store.PolicyAudits().Create(context.Background(), audit); err != nil {
		logging.Default().Errorf("create reject policy audit failed, error: %s", err.Error())
	}
}

// LogGrantedAccessRequest write granted subject access to log.
func (a *AuditLogger) LogGrantedAccessRequest(r *ladon.Request, p ladon.Policies, d ladon.Policies) {
	logging.Default().Infof("subject access review granted, request: %+v, deciders: %+v", r, d)
	conclusion := fmt.Sprintf("policies %s allow access", joinPoliciesNames(d))
	rstring, pstring, dstring := convertToString(r, p, d)
	record := AnalyticsRecord{
		TimeStamp:  time.Now().Unix(),
		Username:   r.Subject,
		Effect:     ladon.AllowAccess,
		Conclusion: conclusion,
		Request:    rstring,
		Policies:   pstring,
		Deciders:   dstring,
	}
	audit := &dao.PolicyAudit{
		Subject:      r.Subject,
		PolicyShadow: record.String(),
	}
	if err := a.store.PolicyAudits().Create(context.Background(), audit); err != nil {
		logging.Default().Errorf("create access policy audit failed, error: %s", err.Error())
	}
}

func rejected(r *ladon.Request, p ladon.Policies, d ladon.Policies) *AnalyticsRecord {
	var conclusion string
	if len(d) > 1 {
		allowed := joinPoliciesNames(d[0 : len(d)-1])
		denied := d[len(d)-1].GetID()
		conclusion = fmt.Sprintf("policies %s allow access, but %s forcefully deny access", allowed, denied)
	} else if len(d) == 1 {
		denied := d[len(d)-1].GetID()
		conclusion = fmt.Sprintf("policy %s forcefully deny access", denied)
	} else {
		conclusion = "no policy allow access"
	}
	request, policies, decisions := convertToString(r, p, d)
	record := &AnalyticsRecord{
		TimeStamp:  time.Now().Unix(),
		Username:   r.Subject,
		Effect:     ladon.DenyAccess,
		Conclusion: conclusion,
		Request:    request,
		Policies:   policies,
		Deciders:   decisions,
	}
	return record
}

func joinPoliciesNames(policies ladon.Policies) string {
	var names []string
	for _, policy := range policies {
		names = append(names, policy.GetID())
	}
	return strings.Join(names, ", ")
}

func convertToString(r *ladon.Request, p ladon.Policies, d ladon.Policies) (string, string, string) {
	rbytes, _ := json.Marshal(r)
	pbytes, _ := json.Marshal(p)
	dbytes, _ := json.Marshal(d)

	return string(rbytes), string(pbytes), string(dbytes)
}

// AnalyticsRecord encodes the details of a authorization request.
type AnalyticsRecord struct {
	TimeStamp  int64     `json:"timestamp"`
	Username   string    `json:"username"`
	Effect     string    `json:"effect"`
	Conclusion string    `json:"conclusion"`
	Request    string    `json:"request"`
	Policies   string    `json:"policies"`
	Deciders   string    `json:"deciders"`
	ExpireAt   time.Time `json:"expireAt"`
}

// String()
func (a *AnalyticsRecord) String() string {
	return fmt.Sprintf("username: %s, effect: %s, conclusion: %s, request: %s, policies: %s, deciders: %s",
		a.Username, a.Effect, a.Conclusion, a.Request, a.Policies, a.Deciders)
}
