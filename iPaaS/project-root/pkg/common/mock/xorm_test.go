package mock

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"xorm.io/xorm"
)

// xormExample 示例
type xormExample struct {
	ID   int64  `xorm:"id pk autoincr"`
	Name string `xorm:"name"`
	Age  int    `xorm:"age"`
}

func (e *xormExample) TableName() string {
	return "xorm_example"
}

func getXormExample(engine *xorm.Engine, id int64) *xormExample {
	example := &xormExample{}
	_, err := engine.ID(id).Get(example)
	if err != nil {
		panic(err)
	}
	return example
}

func transactionExample(engine *xorm.Engine, id int64) *xormExample {
	session := engine.NewSession()
	defer session.Close()

	err := session.Begin()
	if err != nil {
		panic(err)
	}

	example := &xormExample{}
	has, err := session.ID(id).Get(example)
	if err != nil {
		err := session.Rollback()
		if err != nil {
			panic(err)
		}
		panic(err)
	}

	if !has {
		err := session.Rollback()
		if err != nil {
			panic(err)
		}
		panic("not found")
	}

	err = session.Commit()
	if err != nil {
		panic(err)
	}

	return example
}

func TestGetXormExample(t *testing.T) {
	mockEngine, mocksql := Engine()

	mocksql.ExpectQuery("SELECT (.+) FROM `xorm_example` WHERE `id`=(.+) LIMIT 1").
		WithArgs(10000).
		WillReturnRows(mocksql.NewRows([]string{"id", "name", "age"}).
			AddRow(10000, "xx", 18))

	example := getXormExample(mockEngine, 10000)

	assert.Equal(t, int64(10000), example.ID)
	assert.Equal(t, "xx", example.Name)
	assert.Equal(t, 18, example.Age)
}

func TestTransactionExample(t *testing.T) {
	mockEngine, mocksql := Engine()

	mocksql.ExpectBegin()
	mocksql.ExpectQuery("SELECT (.+) FROM `xorm_example` WHERE `id`=(.+) LIMIT 1").
		WithArgs(10000).
		WillReturnRows(mocksql.NewRows([]string{"id", "name", "age"}).
			AddRow(10000, "xx", 18))
	mocksql.ExpectCommit()

	example := transactionExample(mockEngine, 10000)

	assert.Equal(t, int64(10000), example.ID)
	assert.Equal(t, "xx", example.Name)
	assert.Equal(t, 18, example.Age)

	/* ------------------------------- not fount test ------------------------------- */

	mocksql.ExpectBegin()
	mocksql.ExpectQuery("SELECT (.+) FROM `xorm_example` WHERE `id`=(.+) LIMIT 1").
		WithArgs(99999).
		WillReturnRows(mocksql.NewRows([]string{"id", "name", "age"}))
	mocksql.ExpectRollback()

	// assert.PanicsWithError(t, "not found", func() {
	assert.PanicsWithValue(t, "not found", func() {
		transactionExample(mockEngine, 99999)
	}, "The code did not panic")
}
