package exorm

import "xorm.io/xorm"

type ExSession struct {
	db *xorm.Session
}

func (s *ExSession) Cond(cond bool, query interface{}, args ...interface{}) *ExSession {
	if cond {
		s.db.Where(query, args...)
	}
	return s
}

func (s *ExSession) CondIn(cond bool, column string, args ...interface{}) *ExSession {
	if cond {
		s.db.In(column, args...)
	}
	return s
}

func (s *ExSession) Raw() *xorm.Session {
	return s.db
}

func New(db *xorm.Session) *ExSession {
	return &ExSession{db: db}
}
