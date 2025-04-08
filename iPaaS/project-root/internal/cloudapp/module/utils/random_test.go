package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandomPassword(t *testing.T) {
	pass := RandomPassword(2)
	assert.Len(t, pass, 2)

	pass = RandomPassword(4)
	assert.Contains(t, upperLetters, string(pass[0]))
	assert.Contains(t, lowerLetters, string(pass[1]))
	assert.Contains(t, numbers, string(pass[2]))
	assert.Contains(t, symbolLetters, string(pass[3]))
}
