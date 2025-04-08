package otelxorm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSpanName(t *testing.T) {
	t.Run("select", func(t *testing.T) {
		name := spanName("select * from apple limit 100 offset 1")
		assert.Equal(t, "SELECT - apple", name)
	})

	t.Run("select and join", func(t *testing.T) {
		name := spanName("select * from `apple` as `a` join `banana` as `b` limit 100 offset 1")
		assert.Equal(t, "SELECT - apple", name)
	})

	t.Run("insert", func(t *testing.T) {
		name := spanName("insert into apple(a, b) values(?, ?);")
		assert.Equal(t, "INSERT - apple", name)
	})

	t.Run("update", func(t *testing.T) {
		name := spanName("update apple set a = ?, b = ?")
		assert.Equal(t, "UPDATE - apple", name)
	})

	t.Run("delete", func(t *testing.T) {
		name := spanName("delete from apple where a = ? and b = ?")
		assert.Equal(t, "DELETE - apple", name)
	})
}
