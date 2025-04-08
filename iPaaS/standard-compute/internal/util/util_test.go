package util

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/backend/job"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/dao/models"
)

func TestOccupiedNodesNum(t *testing.T) {
	assert.Equal(t, 1, OccupiedNodesNum(1, 10))
	assert.Equal(t, 2, OccupiedNodesNum(11, 10))
	assert.Equal(t, 3, OccupiedNodesNum(21, 10))
}

func TestReplaceEndpoint(t *testing.T) {
	type args struct {
		rawUrl   string
		endpoint string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			args: args{
				rawUrl:   "https://hpc-storage.domain/userId/file",
				endpoint: "http://127.0.0.1:8899",
			},
			want: "http://127.0.0.1:8899/userId/file",
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err == nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReplaceEndpoint(tt.args.rawUrl, tt.args.endpoint)
			if !tt.wantErr(t, err, fmt.Sprintf("ReplaceEndpoint(%v, %v)", tt.args.rawUrl, tt.args.endpoint)) {
				return
			}
			assert.Equalf(t, tt.want, got, "ReplaceEndpoint(%v, %v)", tt.args.rawUrl, tt.args.endpoint)
		})
	}
}

func TestParseRawStorageUrl(t *testing.T) {
	type args struct {
		rawUrl string
	}
	tests := []struct {
		name         string
		args         args
		wantEndpoint string
		wantPath     string
		wantErr      bool
	}{
		{
			name: "",
			args: args{
				rawUrl: "https://ysfortest.com:8888/4ZdQuVyZiDS/input/565tZRzTevS/08976f0f-d8af-47a3-b631-f15472464776/13- Forwardbraking3#.inp?123=312",
			},
			wantEndpoint: "https://ysfortest.com:8888",
			wantPath:     "/4ZdQuVyZiDS/input/565tZRzTevS/08976f0f-d8af-47a3-b631-f15472464776/13- Forwardbraking3#.inp?123=312",
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotEndpoint, gotPath, err := ParseRawStorageUrl(tt.args.rawUrl)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseRawStorageUrl() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotEndpoint != tt.wantEndpoint {
				t.Errorf("ParseRawStorageUrl() gotEndpoint = %v, want %v", gotEndpoint, tt.wantEndpoint)
			}
			if gotPath != tt.wantPath {
				t.Errorf("ParseRawStorageUrl() gotPath = %v, want %v", gotPath, tt.wantPath)
			}
		})
	}
}

func TestReplaceCommand(t *testing.T) {
	j := &job.Job{
		Job: &models.Job{
			Id: 123,
		},
	}

	preparedFilePath := "/tmp"
	tests := []struct {
		name     string
		cmd      string
		flag     string
		value    string
		expected string
	}{
		{
			name:     "replace flag with value",
			cmd:      "run #idflag",
			flag:     "#idflag",
			value:    "--id=xxx",
			expected: "run --id=xxx",
		},
		{
			name:     "replace flag with no value",
			cmd:      "echo #idflag",
			flag:     "#idflag",
			value:    "",
			expected: "echo ",
		},
		{
			name:     "replace flag with normal",
			cmd:      "hello;#YS_COMMAND_PREPARED",
			flag:     PreparedFlag,
			value:    PreparedCmd(j, preparedFilePath),
			expected: "hello;echo 'YS command prepared' > /tmp/123_prepared",
		},
		{
			name:     "replace flag with empty preparedFilePath",
			cmd:      "hello;#YS_COMMAND_PREPARED",
			flag:     PreparedFlag,
			value:    PreparedCmd(j, ""),
			expected: "hello;echo 'YS command prepared' > /tmp/123_prepared",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ReplaceCommand(tt.cmd, tt.flag, tt.value)
			assert.Equal(t, tt.expected, got)
		})
	}
}
