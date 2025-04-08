package env

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"

	"github.com/joho/godotenv"
	"go.uber.org/zap/zapcore"
)

const (
	// ModeLocal local develop enviorment
	ModeLocal = iota
	// ModeDev ModeDev
	ModeDev
	// ModeTest ModeTest
	ModeTest
	// ModeStage ModeStage
	ModeStage
	// ModeProd ModeProd
	ModeProd
)

// Levels is log level for uber/zap
const (
	LevelDebug = iota - 1
	LevelInfo
	LevelWarn
	LevelError
)

const EnvNameMode = "YS_MODE"
const EnvNameModeAlias = "YS_ENV"
const EnvNameLogLevel = "YS_LOG_LEVEL"

var mut sync.Mutex

// env类型映射,白名单列表, 默认为：MODE_DEV, LEVEL_OFF
var (
	ModeMap = map[string]int{
		"dev":   ModeDev,
		"local": ModeLocal,
		"test":  ModeTest,
		"stage": ModeStage,
		"prod":  ModeProd,
	}
	LogLevelMap = map[string]int{
		"debug": int(zapcore.DebugLevel),
		"info":  int(zapcore.InfoLevel),
		"warn":  int(zapcore.WarnLevel),
		"error": int(zapcore.ErrorLevel),
	}
)

// EnvType 定义Env文件类型
type EnvType struct {
	Mode     int
	LogLevel int
	// Type means Mode name
	Type string
}

// env env
var (
	Env = &EnvType{}
)

// GinMode GinMode
func GinMode(mode int) string {
	m := "debug"
	switch mode {
	case ModeDev:
		m = "debug"
	case ModeLocal:
		m = "debug"
	case ModeTest:
		m = "test"
	case ModeStage:
		m = "release"
	case ModeProd:
		m = "release"
	default:
		panic("unknown mode: " + strconv.Itoa(mode))
	}
	return m
}

// ModeName ModeName
func ModeName(mode int) string {
	for k, v := range ModeMap {
		if v == mode {
			return k
		}
	}
	return ""
}

// LogLevelName LogLevelName
func LogLevelName(logLevel int) string {
	for k, v := range LogLevelMap {
		if v == logLevel {
			return k
		}
	}
	return ""
}

// ImportDockerSecretAsEnv :
// docker's secret is under /run/secrets and can't bind as env
// we import secrets into env which name is in "*_env_*" pattern
// content of secret should be XXX=YYY, multi-line is supported
// it will overload exist env
func ImportDockerSecretAsEnv() {
	files, err := ioutil.ReadDir("/run/secrets")
	if err != nil {
		return
	}

	envFiles := []string{}

	for _, file := range files {
		if strings.Index(file.Name(), "_env_") >= 0 {
			envFiles = append(envFiles, path.Join("/run/secrets", file.Name()))
		}
	}

	if len(envFiles) == 0 {
		return
	}

	if err := godotenv.Overload(envFiles...); err != nil {
		panic(err)
	}
}

// InitEnv init env.Env, path to .env file is optional
func InitEnv(path string) {
	if path != "" {
		err := loadEnvFile(path)
		if err != nil && !os.IsNotExist(err) {
			// only panic when error is not NotExist
			panic("InitEnv: " + err.Error())
		}
	}

	if mode := os.Getenv(EnvNameMode); mode != "" {
		v, ok := ModeMap[mode]
		if !ok {
			panic(fmt.Sprintf("invalid mode from env %s: %q", EnvNameMode, mode))
		}
		Env.Mode = v
	} else if mode := os.Getenv(EnvNameModeAlias); mode != "" {
		v, ok := ModeMap[mode]
		if !ok {
			panic(fmt.Sprintf("invalid mode from env %s: %q", EnvNameMode, mode))
		}
		Env.Mode = v
	}
	Env.Type = ModeName(Env.Mode)
	if logLevel := os.Getenv(EnvNameLogLevel); logLevel != "" {
		// not a must for normal running
		Env.LogLevel = LogLevelMap[logLevel]
	}
}

func loadEnvFile(path string) error {
	mut.Lock()
	defer mut.Unlock()

	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	r := bufio.NewReader(f)
	envsString := ""
	for {
		b, _, err := r.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			err := "reading Env file: " + err.Error()
			panic(err)
		}

		s := strings.TrimSpace(string(b))
		if strings.Index(s, "#") == 0 || strings.Index(s, "=") == -1 {
			continue
		}

		index := strings.Index(s, "=")
		key := strings.TrimSpace(s[:index])
		if len(key) == 0 {
			continue
		}
		value := strings.TrimSpace(s[index+1:])

		// 内部匿名方法, 截取文本内容
		findValue := func(value, separator string) string {
			pos := strings.Index(value, separator)
			if pos > -1 {
				value = value[0:pos]
			}
			return strings.TrimSpace(value)
		}

		value = findValue(value, `#`)
		value = findValue(value, `//`)
		if len(value) == 0 {
			continue
		}

		// set os env
		// the priority of `.env` is less than environment variables.
		if currentValue := os.Getenv(key); currentValue == "" {
			os.Setenv(key, value)
			envsString = fmt.Sprintf("%s  %s=%v\n", envsString, key, value)
		} else {
			envsString = fmt.Sprintf("%s  %s=%v\n", envsString, key, currentValue)
		}
	}
	log.Printf("Load env from: %s\n%s", path, envsString)
	return nil
}
