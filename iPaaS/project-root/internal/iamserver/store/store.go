package store

//go:generate mockgen -self_package=github.com/yuansuan/ticp/common/project-root-iam/internal/iamserver/store -destination mock_store.go -package store github.com/yuansuan/ticp/common/project-root-iam/internal/iamserver/store Factory,RoleStore,SecretStore,PolicyStore,RolePolicyRelationStore,PolicyAuditStore

// Factory defines the iam platform storage interface.
type Factory interface {
	Secrets() SecretStore
	Policies() PolicyStore
	Roles() RoleStore
	PolicyAudits() PolicyAuditStore
	RolePolicyRelations() RolePolicyRelationStore
	MigrateDatabase() error
	Close() error
}
