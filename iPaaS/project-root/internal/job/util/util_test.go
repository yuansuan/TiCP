// Copyright (C) 2018 LambdaCal Inc.

package util

import (
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common/utils"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
)

type AllocNodesSuite struct {
	suite.Suite
}

func TestAllocNodesSuite(t *testing.T) {
	suite.Run(t, new(AllocNodesSuite))
}

func (s *AllocNodesSuite) TestAllocNodes() {
	s.Run("Round", func() {
		s.Run("case 1", func() {
			coresPerNode, cores, err := AllocNodes(1, 1, true)
			s.Equal(int64(1), coresPerNode)
			s.Equal(int64(1), cores)
			s.Nil(err)
		})

		s.Run("case 2", func() {
			coresPerNode, cores, err := AllocNodes(2, 1, true)
			s.Equal(int64(1), coresPerNode)
			s.Equal(int64(2), cores)
			s.Nil(err)
		})

		s.Run("case 3", func() {
			coresPerNode, cores, err := AllocNodes(1, 2, true)
			s.Equal(int64(2), coresPerNode)
			s.Equal(int64(2), cores)
			s.Nil(err)
		})

		s.Run("case 4", func() {
			coresPerNode, cores, err := AllocNodes(3, 2, true)
			s.Equal(int64(2), coresPerNode)
			s.Equal(int64(4), cores)
			s.Nil(err)
		})

		s.Run("correction 1", func() {
			coresPerNode, cores, err := AllocNodes(2, 2, true, func(cores, coresPerNode int64) (int64, int64) {
				return cores, coresPerNode / 2
			})
			s.Equal(int64(1), coresPerNode)
			s.Equal(int64(2), cores)
			s.Nil(err)
		})

		s.Run("correction 2", func() {
			coresPerNode, cores, err := AllocNodes(1, 1, true, func(cores, coresPerNode int64) (int64, int64) {
				return cores, coresPerNode * 3
			})
			s.Equal(int64(3), coresPerNode)
			s.Equal(int64(3), cores)
			s.Nil(err)
		})
	})

	s.Run("Floor", func() {
		s.Run("case 1", func() {
			coresPerNode, cores, err := AllocNodes(1, 1, false)
			s.Equal(int64(1), coresPerNode)
			s.Equal(int64(1), cores)
			s.Nil(err)
		})

		s.Run("case 2", func() {
			coresPerNode, cores, err := AllocNodes(2, 1, false)
			s.Equal(int64(1), coresPerNode)
			s.Equal(int64(2), cores)
			s.Nil(err)
		})

		s.Run("case 3", func() {
			coresPerNode, cores, err := AllocNodes(1, 2, false)
			s.Equal(int64(2), coresPerNode)
			s.Equal(int64(2), cores)
			s.Nil(err)
		})

		s.Run("case 4", func() {
			coresPerNode, cores, err := AllocNodes(3, 2, false)
			s.Equal(int64(2), coresPerNode)
			s.Equal(int64(2), cores)
			s.Nil(err)
		})

		s.Run("correction 1", func() {
			coresPerNode, cores, err := AllocNodes(2, 2, false, func(cores, coresPerNode int64) (int64, int64) {
				return cores, coresPerNode / 2
			})
			s.Equal(int64(1), coresPerNode)
			s.Equal(int64(2), cores)
			s.Nil(err)
		})

		s.Run("correction 2", func() {
			coresPerNode, cores, err := AllocNodes(1, 1, false, func(cores, coresPerNode int64) (int64, int64) {
				return cores, coresPerNode * 3
			})
			s.Equal(int64(3), coresPerNode)
			s.Equal(int64(3), cores)
			s.Nil(err)
		})
	})

	s.Run("special", func() {
		// 允许不取整, cores小于coresPerNode
		s.Run("special 1", func() {
			coresPerNode, cores, err := AllocNodes(10, 56, true, WithSharedNode)
			s.Equal(int64(10), coresPerNode)
			s.Equal(int64(10), cores)
			s.Nil(err)
		})

		// 允许不取整, cores大于coresPerNode
		s.Run("special 2", func() {
			coresPerNode, cores, err := AllocNodes(100, 56, true, WithSharedNode)
			s.Equal(int64(56), coresPerNode)
			s.Equal(int64(112), cores)
			s.Nil(err)
		})
	})

	s.Run("error", func() {
		s.Run("coresPerNode zero error", func() {
			coresPerNode, cores, err := AllocNodes(1, 0, true)
			s.Equal(int64(0), coresPerNode)
			s.Equal(int64(0), cores)
			s.NotNil(err)
		})

		s.Run("cores not positive error", func() {
			coresPerNode, cores, err := AllocNodes(-1, 1, true)
			s.Equal(int64(0), coresPerNode)
			s.Equal(int64(0), cores)
			s.NotNil(err)
		})
	})
}

func TestAddAppImagePrefix(t *testing.T) {
	type args struct {
		appImage string
		isLocal  bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "case 1",
			args: args{
				appImage: "appImage",
				isLocal:  true,
			},
			want: "local:appImage",
		},
		{
			name: "case 2",
			args: args{
				appImage: "appImage",
				isLocal:  false,
			},
			want: "image:appImage",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AddAppImagePrefix(tt.args.appImage, tt.args.isLocal)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetUserId(t *testing.T) {
	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)

	ctxNormal, _ := gin.CreateTestContext(w)
	ctxNormal.Request = &http.Request{
		Header: make(http.Header),
		URL:    &url.URL{},
	}
	ctxNormal.Request.Header.Set("x-ys-user-id", "4RXv3DvUe1u")

	ctxError, _ := gin.CreateTestContext(w)
	ctxError.Request = &http.Request{
		Header: make(http.Header),
		URL:    &url.URL{},
	}
	ctxError.Request.Header.Set("x-ys-user-id", "1663475500737122304")

	ctxEmpty, _ := gin.CreateTestContext(w)
	ctxEmpty.Request = &http.Request{
		Header: make(http.Header),
		URL:    &url.URL{},
	}

	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name    string
		args    args
		want    snowflake.ID
		wantErr bool
	}{
		{
			name: "case normal",
			args: args{
				c: ctxNormal,
			},
			want:    1663475500737122304,
			wantErr: false,
		},
		{
			name: "case error",
			args: args{
				c: ctxError,
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "case empty",
			args: args{
				c: ctxEmpty,
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := utils.GetUserID(tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestParseYsID(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "case 1",
			args: args{
				path: "https://shanhe-storage.com/ys_id/csp_project1/",
			},
			want: "ys_id",
		},
		{
			name: "case 2",
			args: args{
				path: "http://172.20.1.247:8898/AABB/path1",
			},
			want: "AABB",
		},
		{
			name: "case 3",
			args: args{
				path: "http://10.0.4.55:8899/Ysser123/path2/",
			},
			want: "Ysser123",
		},
		{
			name: "case 4",
			args: args{
				path: "https://wuxi-storage.yuansuan.cn/OSO12223/path3/input",
			},
			want: "OSO12223",
		},
		{
			name: "case 5",
			args: args{
				path: "https://wuxi-storage.yuansuan.cn/OSO12223/",
			},
			want: "OSO12223",
		},
		{
			name: "case 6",
			args: args{
				path: "https://wuxi-storage.yuansuan.cn/OSO12223",
			},
			want: "OSO12223",
		},
		{
			name: "case error",
			args: args{
				path: "oms://wuxi-storage.yuansuan.cn",
			},
			want: "",
		},
		{
			name: "case error2",
			args: args{
				path: "https://wuxi-storage.yuansuan.cn",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseYsID(tt.args.path)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestParsePath(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "case 1",
			args: args{
				path: "https://shanhe-storage.com/ys_id/csp_project1/",
			},
			want: "/csp_project1/",
		},
		{
			name: "case 2",
			args: args{
				path: "http://172.20.1.247:8898/AABB/path1",
			},
			want: "/path1",
		},
		{
			name: "case 3",
			args: args{
				path: "http://10.0.4.55:8899/Ysser123/path2/",
			},
			want: "/path2/",
		},
		{
			name: "case 4",
			args: args{
				path: "https://wuxi-storage.yuansuan.cn/OSO12223/path3/input",
			},
			want: "/path3/input",
		},
		{
			name: "case 5",
			args: args{
				path: "https://wuxi-storage.yuansuan.cn/OSO12223/path3/input/",
			},
			want: "/path3/input/",
		},
		{
			name: "case error",
			args: args{
				path: "oms://wuxi-storage.yuansuan.cn",
			},
			want: "",
		},
		{
			name: "case error2",
			args: args{
				path: "https://wuxi-storage.yuansuan.cn",
			},
			want: "",
		},
		{
			name: "case error3",
			args: args{
				path: "https://wuxi-storage.yuansuan.cn/OSO12223",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParsePath(tt.args.path)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestParseYsIDWithOutDomain(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "case 1",
			args: args{
				path: "/ys_id/csp_project1/",
			},
			want: "ys_id",
		},
		{
			name: "case 2",
			args: args{
				path: "/AABB/path1",
			},
			want: "AABB",
		},
		{
			name: "case 3",
			args: args{
				path: "/Ysser123/path2/",
			},
			want: "Ysser123",
		},
		{
			name: "case 4",
			args: args{
				path: "/OSO12223/path3/input",
			},
			want: "OSO12223",
		},
		{
			name: "case 5",
			args: args{
				path: "/OSO12223/",
			},
			want: "OSO12223",
		},
		{
			name: "case 6",
			args: args{
				path: "/OSO12223",
			},
			want: "OSO12223",
		},
		{
			name: "case error",
			args: args{
				path: "///",
			},
			want: "",
		},
		{
			name: "case error2",
			args: args{
				path: "",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseYsIDWithOutDomain(tt.args.path)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestParsePathWithOutDomain(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "case 1",
			args: args{
				path: "/ys_id/csp_project1/",
			},
			want: "/csp_project1/",
		},
		{
			name: "case 2",
			args: args{
				path: "/AABB/path1",
			},
			want: "/path1",
		},
		{
			name: "case 3",
			args: args{
				path: "/Ysser123/path2/",
			},
			want: "/path2/",
		},
		{
			name: "case 4",
			args: args{
				path: "/OSO12223/path3/input",
			},
			want: "/path3/input",
		},
		{
			name: "case 5",
			args: args{
				path: "/OSO12223/path3/input/",
			},
			want: "/path3/input/",
		},
		{
			name: "case error",
			args: args{
				path: "",
			},
			want: "",
		},
		{
			name: "case error2",
			args: args{
				path: "///",
			},
			want: "",
		},
		{
			name: "case error3",
			args: args{
				path: "/OSO12223",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParsePathWithOutDomain(tt.args.path)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestIsNonZeroExitCode(t *testing.T) {
	testCases := []struct {
		exitCode string
		expected bool
	}{
		{"0:0", false},
		{"255:0", true},
		{"1:0", true},
		{"123:0", true},
		{"42:0", true},
		{"0:9", false},
		{"123:9", true},
		{"invalid", false}, // Testing for an invalid exit code
	}

	for _, tc := range testCases {
		t.Run(tc.exitCode, func(t *testing.T) {
			result, err := IsNonZeroExitCode(tc.exitCode)
			if err != nil {
				assert.Error(t, err, "For exitCode %s", tc.exitCode)
				return
			}
			assert.Equal(t, tc.expected, result, "For exitCode %s", tc.exitCode)
		})
	}
}
