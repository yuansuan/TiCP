package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJob_JobID(t *testing.T) {
	t.Run("correct", func(t *testing.T) {
		var id int64 = 123456

		assert.Equal(t, "123456", (&Job{Id: id}).JobID())
	})

	t.Run("bad job id", func(t *testing.T) {
		var id int64 = 654321

		assert.NotEqual(t, "123456", (&Job{Id: id}).JobID())
	})
}
