package cloud

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestA(t *testing.T) {
	sp := &ScriptParams{
		ShareUsername:   "user1",
		SharePassword:   "pass1",
		ShareMountPaths: "mountSrc1=mountPoint1",
	}

	sp = sp.UpdateByMount("user2", "pass2", "mountSrc2", "mountPoint2")
	assert.Equal(t, "user1,user2", sp.ShareUsername.String())
	assert.Equal(t, "pass1,pass2", sp.SharePassword.String())
	assert.Equal(t, "mountSrc1=mountPoint1,mountSrc2=mountPoint2", sp.ShareMountPaths.String())

	sp = sp.UpdateByUMount("mountPoint1")
	assert.Equal(t, "user2", sp.ShareUsername.String())
	assert.Equal(t, "pass2", sp.SharePassword.String())
	assert.Equal(t, "mountSrc2=mountPoint2", sp.ShareMountPaths.String())
}
