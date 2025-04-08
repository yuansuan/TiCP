package env

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Env_empty(t *testing.T) {
	InitEnv("")
	t.Log(Env)
	assert.Equal(t, ModeLocal, Env.Mode)
	assert.Equal(t, "local", Env.Type)
	assert.Equal(t, LevelInfo, Env.LogLevel)
}

func Test_Env_ok(t *testing.T) {
	InitEnv("ok.env")
	t.Log(Env)
	assert.Equal(t, ModeTest, Env.Mode)
	assert.Equal(t, "test", Env.Type)
	assert.Equal(t, LevelWarn, Env.LogLevel)
}

func Test_Env_err(t *testing.T) {
	assert.Panics(t, func() {
		InitEnv("err.env")
	})
}

func Test_Env_os(t *testing.T) {
	os.Setenv(EnvNameMode, "prod")
	InitEnv("")
	t.Log(Env)
	assert.Equal(t, ModeProd, Env.Mode)
}
