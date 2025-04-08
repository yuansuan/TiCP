package mount

import "testing"

func Test_trim(t *testing.T) {
	type args struct {
		o *Options
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test",
			args: args{o: &Options{Server: "10.0.1.123:8081"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			trim(tt.args.o)
		})
	}
}
