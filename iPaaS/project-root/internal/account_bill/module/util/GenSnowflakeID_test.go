package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenID(t *testing.T) {

	t.Run("Normal AccountID", func(t *testing.T) {

		id := GenID()
		t.Logf("id : %v", id)
		assert.NotNil(t, id)
	})
}
