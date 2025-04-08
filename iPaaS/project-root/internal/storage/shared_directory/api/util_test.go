package api

import (
	"testing"
)

func Test_genUserName(t *testing.T) {
	un, _ := generateUserName()
	t.Log(un)
}

func Test_generateRandomPassword(t *testing.T) {
	p, err := generateRandomPassword(DefaultPWLength)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(p)
}
