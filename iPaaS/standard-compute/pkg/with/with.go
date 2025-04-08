package with

import (
	"context"
	"errors"
	"time"

	"go.uber.org/multierr"
	"xorm.io/xorm"
)

type sessionKey struct{}

var SessionKey = sessionKey{}

func KeepSession(ctx context.Context, db *xorm.Session) context.Context {
	return context.WithValue(ctx, SessionKey, db)
}

type ormKey struct{}

var OrmKey = ormKey{}

func DefaultTransaction(ctx context.Context, action func(context.Context) error) error {
	db, ok := ctx.Value(OrmKey).(*xorm.Engine)
	if !ok {
		return errors.New("get xorm engine failed")
	}
	_, err := db.Transaction(func(db *xorm.Session) (interface{}, error) {
		return nil, action(KeepSession(ctx, db))
	})
	return err
}

func DefaultSession(ctx context.Context, action func(db *xorm.Session) error) (err error) {
	if session, ok := ctx.Value(SessionKey).(*xorm.Session); ok {
		return action(session)
	}

	db, ok := ctx.Value(OrmKey).(*xorm.Engine)
	if !ok {
		return errors.New("get xorm engine failed")
	}
	session := db.NewSession()
	defer func() { err = multierr.Append(err, session.Close()) }()
	return action(session)
}

func TimedLoop(ctx context.Context, timeout time.Duration, action func(context.Context) (bool, error)) error {
	for {
		if exit, err := action(ctx); exit {
			return err
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(timeout):
		}
	}
}
