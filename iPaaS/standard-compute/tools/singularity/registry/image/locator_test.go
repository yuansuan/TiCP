package image

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFromString_Success(t *testing.T) {
	table := []struct {
		Target string
		Name   string
		Tag    string
		Hash   string
	}{
		{"aa:bb", "aa", "bb", ""},
		{"aa:bb@cc", "aa", "bb", "cc"},
		{"a:b@c", "a", "b", "c"},
		{"a:1@c", "a", "1", "c"},
		{"a:1@2", "a", "1", "2"},
		{"a:1.1@2", "a", "1.1", "2"},
	}

	for _, testcase := range table {
		if l, err := FromString(testcase.Target); assert.NoError(t, err) {
			assert.Equal(t, testcase.Name, l.Name())
			assert.Equal(t, testcase.Tag, l.Tag())
			assert.Equal(t, testcase.Hash, l.Hash())
		}
	}
}

func TestFromString_Error(t *testing.T) {
	tables := []string{
		"aa:", "aa:bb@", ":bb", "aa:bb:cc",
		"@cc", "aa:bb@cc@", "aa:bb@cc:",
		"aa:@cc", "aa.xx",
	}

	for _, testcase := range tables {
		_, err := FromString(testcase)
		assert.ErrorIs(t, err, ErrInvalidLocatorString)
	}
}
