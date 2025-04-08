package zone

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/zonelist"
	v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/config"
)

func TestList(t *testing.T) {
	type args struct {
		cfg config.CustomT
	}
	tests := []struct {
		name string
		args args
		want zonelist.Data
	}{
		{
			name: "",
			args: args{
				cfg: config.CustomT{
					ChangeLicense: false,
					Zones: map[string]*v20230530.Zone{"A": {
						HPCEndpoint:     "HH",
						StorageEndpoint: "VV",
						CloudAppEnable:  false,
					}},
					SelfYsID: "",
					AK:       "",
					AS:       "",
				},
			},
			want: zonelist.Data{
				Zones: map[string]*v20230530.Zone{"A": {
					HPCEndpoint:     "HH",
					StorageEndpoint: "VV",
					CloudAppEnable:  false,
				}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := List(tt.args.cfg)
			assert.Equal(t, resp, tt.want)
		})
	}
}
