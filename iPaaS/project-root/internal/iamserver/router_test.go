package iamserver

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestIsAdminAllow(t *testing.T) {
	type args struct {
		c   *gin.Context
		tag string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "product account and admin url",
			args: args{
				c: &gin.Context{
					Request: &http.Request{
						URL: &url.URL{
							Path: "/admin",
						},
					},
				},
				tag: "YS_CloudStorage",
			},
			want: true,
		},
		{
			name: "product account and internal url",
			args: args{
				c: &gin.Context{
					Request: &http.Request{
						URL: &url.URL{
							Path: "/internal",
						},
					},
				},
				tag: "YS_CSP",
			},
			want: true,
		},
		{
			name: "non-product account and admin url",
			args: args{
				c: &gin.Context{
					Request: &http.Request{
						URL: &url.URL{
							Path: "/admin",
						},
					},
				},
				tag: "non-product",
			},
			want: false,
		},
		{
			name: "non-product account and internal url",
			args: args{
				c: &gin.Context{
					Request: &http.Request{
						URL: &url.URL{
							Path: "/internal",
						},
					},
				},
				tag: "non-product",
			},
			want: false,
		},
		{
			name: "product account and non-admin url",
			args: args{
				c: &gin.Context{
					Request: &http.Request{
						URL: &url.URL{
							Path: "/api",
						},
					},
				},
				tag: "product",
			},
			want: true,
		},
		{
			name: "non-product account and non-admin url",
			args: args{
				c: &gin.Context{
					Request: &http.Request{
						URL: &url.URL{
							Path: "/api/admin",
						},
					},
				},
				tag: "non-product",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isAdminAllow(tt.args.c, tt.args.tag)
			if got != tt.want {
				t.Errorf("isAdminAllow() = %v, want %v", got, tt.want)
			}
		})
	}
}
