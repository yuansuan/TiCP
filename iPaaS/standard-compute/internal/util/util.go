package util

import (
	"bytes"
	"context"
	"fmt"
	"math"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
	"time"

	"github.com/pkg/errors"

	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/backend/job"
)

const (
	StateKeyInGinCtx = "state"
)

func RenderScript(_ context.Context, j *job.Job, scriptTpl string, providerType string) ([]byte, error) {
	var b bytes.Buffer
	var tmpl = template.Must(template.New("script").Parse(scriptTpl))
	err := tmpl.Execute(&b, j)
	return b.Bytes(), errors.Wrap(err, providerType+" render")
}

func EnvPair(key string, value interface{}) string {
	return fmt.Sprintf("%v=%v", key, value)
}

func RenderEnvVars(_ context.Context, j *job.Job, providerType string) ([]byte, error) {
	var b bytes.Buffer
	for _, v := range j.EnvVars {
		_, err := b.WriteString(v + "\n")
		if err != nil {
			return nil, errors.Wrap(err, providerType+" render: renderEnvVars")
		}
	}
	return b.Bytes(), nil
}

// ParseDomain 解析路径中域名的部分
func ParseDomain(path string) (string, error) {
	var re = regexp.MustCompile(`(?m)(http[s]{0,1}:\/\/[^\/]+[/]{0,1})`)
	match := re.FindStringSubmatch(path)
	if len(match) > 1 {
		return strings.TrimSuffix(match[1], "/"), nil
	}
	return "", fmt.Errorf("no domain found in path")
}

// ParseRawStorageUrl rawUrl should be like "schema://host/path"
func ParseRawStorageUrl(rawUrl string) (endpoint, path string, err error) {
	endpoint, err = ParseDomain(rawUrl)
	if err != nil {
		return
	}

	path = strings.TrimPrefix(rawUrl, endpoint)
	return
}

func ReplaceEndpoint(rawUrl, endpoint string) (string, error) {
	oldEndpoint, err := ParseDomain(rawUrl)
	if err != nil {
		return "", err
	}
	newEndpoint, err := ParseDomain(endpoint)
	if err != nil {
		return "", err
	}

	return strings.Replace(rawUrl, oldEndpoint, newEndpoint, 1), nil
}

func OccupiedNodesNum(requestCores int, nTasksPerNode int) int {
	return int(math.Ceil(float64(requestCores) / float64(nTasksPerNode)))
}

func PTime(t time.Time) *time.Time {
	return &t
}

const (
	PreparedFlag = "#YS_COMMAND_PREPARED"
	preparedCmd  = "echo 'YS command prepared' > %s"
	preparedFile = "%d_prepared"
)

func PreparedCmd(j *job.Job, preparedFilePath string) string {
	if preparedFilePath == "" {
		preparedFilePath = "/tmp"
	}
	file := PreparedFile(j, preparedFilePath)
	return fmt.Sprintf(preparedCmd, file)
}

func PreparedFile(j *job.Job, preparedFilePath string) string {
	return filepath.Join(preparedFilePath, fmt.Sprintf("%d", j.Id), fmt.Sprintf(preparedFile, j.Id))
}

func ReplaceCommand(cmd, flag, value string) string {
	return strings.ReplaceAll(cmd, flag, value)
}
