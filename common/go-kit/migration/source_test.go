package migration

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/yuansuan/ticp/common/go-kit/migration/example"
)

func Test_parseSourceFilename(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		want1   migrationType
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "1.init.up.sql success",
			args: args{
				filename: "1.init.up.sql",
			},
			want:  1,
			want1: up,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err == nil
			},
		},
		{
			name: "1.init1.up.sql success",
			args: args{
				filename: "1.init1.up.sql",
			},
			want:  1,
			want1: up,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err == nil
			},
		},
		{
			name: "1.init.down.sql success",
			args: args{
				filename: "1.init.down.sql",
			},
			want:  1,
			want1: down,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err == nil
			},
		},
		{
			name: "1.init.abcd.sql fail",
			args: args{
				filename: "1.init.abcd.sql",
			},
			want:  0,
			want1: "",
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err != nil
			},
		},
		{
			name: "1.test.unknown.up.sql fail",
			args: args{
				filename: "1.test.unknown.up.sql",
			},
			want:  0,
			want1: "",
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err != nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := parseSourceFilename(tt.args.filename)
			if !tt.wantErr(t, err, fmt.Sprintf("parseSourceFilename(%v)", tt.args.filename)) {
				return
			}
			assert.Equalf(t, tt.want, got, "parseSourceFilename(%v)", tt.args.filename)
			assert.Equalf(t, tt.want1, got1, "parseSourceFilename(%v)", tt.args.filename)
		})
	}
}

func TestInitSourceNodeLinkedList(t *testing.T) {
	s := &source{
		root: &sourceNode{
			version: 0,
		},
	}

	exist := s.findAndInsert(3, down, []byte("sql script 3 down"), "3.table3.down.sql")
	assert.False(t, exist)
	exist = s.findAndInsert(3, up, []byte("sql script 3 up"), "3.table3.up.sql")
	assert.False(t, exist)
	exist = s.findAndInsert(2, down, []byte("sql script 2 down"), "2.table2.down.sql")
	assert.False(t, exist)
	exist = s.findAndInsert(1, up, []byte("sql script 1 up"), "1.init.up.sql")
	assert.False(t, exist)
	exist = s.findAndInsert(1, down, []byte("sql script 1 down"), "1.init.down.sql")
	assert.False(t, exist)
	exist = s.findAndInsert(2, up, []byte("sql script 2 up"), "2.table1.up.sql")
	assert.False(t, exist)

	exist = s.findAndInsert(2, up, []byte("sql script 2 up"), "2.table1.up.sql")
	assert.True(t, exist)

	curr := s.root
	for i := 1; i <= 3; i++ {
		curr = curr.next
		assert.NotNil(t, curr)
		assert.NotNil(t, curr.upgrade)
		assert.NotNil(t, curr.downgrade)
		assert.Equal(t, i, curr.version)

		if curr.next != nil {
			assert.Equal(t, i+1, curr.next.version)
		}
	}
}

func TestInitMigrationSource(t *testing.T) {
	s, err := newSource(example.Mysql, Mysql)
	assert.NoError(t, err)
	assert.NotNil(t, s)
}

func TestSubLinkList(t *testing.T) {
	s := &source{
		root: &sourceNode{
			version: 0,
		},
	}

	nodeArgs := []struct {
		version  int
		mType    migrationType
		script   []byte
		filename string
	}{
		{1, up, []byte("sql script 1 up"), "1.init.up.sql"},
		{1, down, []byte("sql script 1 down"), "1.init.down.sql"},
		{2, up, []byte("sql script 2 up"), "2.table2.up.sql"},
		{2, down, []byte("sql script 2 down"), "2.table2.down.sql"},
		{3, up, []byte("sql script 3 up"), "3.table3.up.sql"},
		{3, down, []byte("sql script 3 down"), "3.table3.down.sql"},
	}

	for _, arg := range nodeArgs {
		s.findAndInsert(arg.version, arg.mType, arg.script, arg.filename)
	}

	head, tail, err := s.subLinkList(1, 2)
	assert.NoError(t, err)
	assert.Equal(t, 1, head.version)
	assert.Equal(t, 2, tail.version)
	assert.Equal(t, tail, head.next)
	assert.Nil(t, head.next.next)
	assert.Nil(t, tail.next)

	head, tail, err = s.subLinkList(0, 3)
	assert.NoError(t, err)
	assert.Equal(t, 0, head.version)
	assert.Equal(t, 3, tail.version)
}
