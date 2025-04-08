package with

import (
	"context"

	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"
	"go.uber.org/multierr"
	"xorm.io/xorm"
)

type sessionKey struct{}

func KeepSession(ctx context.Context, db *xorm.Session) context.Context {
	return context.WithValue(ctx, sessionKey{}, db)
}

func DefaultTransaction(ctx context.Context, action func(context.Context) error) error {
	_, err := boot.MW.DefaultTransaction(ctx, func(db *xorm.Session) (interface{}, error) {
		return nil, action(KeepSession(ctx, db))
	})
	return err
}

func DefaultSession(ctx context.Context, action func(db *xorm.Session) error) (err error) {
	if session, ok := ctx.Value(sessionKey{}).(*xorm.Session); ok {
		return action(session)
	}

	session := boot.MW.DefaultSession(ctx)
	defer func() { err = multierr.Append(err, session.Close()) }()
	return action(session)
}
