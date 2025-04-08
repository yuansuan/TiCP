package mock

import (
	"github.com/DATA-DOG/go-sqlmock"
	"xorm.io/xorm"
	"xorm.io/xorm/core"
)

// Engine mock xorm engine
// 模拟xorm engine, 用于测试代码内部调用数据库时，模拟数据库的返回值
func Engine() (*xorm.Engine, sqlmock.Sqlmock) {
	db, sqlMock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}

	engine, err := xorm.NewEngineWithDB("mysql", "", core.FromDB(db))
	if err != nil {
		panic(err)
	}

	engine.ShowSQL(true) // 打印sql语句

	return engine, sqlMock
}
