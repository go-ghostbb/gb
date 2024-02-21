package polaris

import "testing"

// Test_trimAndReplace Test trimAndReplace
func Test_trimAndReplace(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test_trimAndReplace-0",
			args: args{key: "/service/default/default/ghostbb-provider-0-tcp/latest/127.0.0.1:9000"},
			want: "service-default-default-ghostbb-provider-0-tcp-latest-127.0.0.1:9000",
		},
		{
			name: "Test_trimAndReplace-1",
			args: args{key: "/service/default/default/ghostbb-provider-1-tcp/latest/127.0.0.1:9001"},
			want: "service-default-default-ghostbb-provider-1-tcp-latest-127.0.0.1:9001",
		},
		{
			name: "Test_trimAndReplace-2",
			args: args{key: "/service/default/default/ghostbb-provider-2-tcp/latest/127.0.0.1:9002"},
			want: "service-default-default-ghostbb-provider-2-tcp-latest-127.0.0.1:9002",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := trimAndReplace(tt.args.key); got != tt.want {
				t.Errorf("trimAndReplace() = %v, want %v", got, tt.want)
			}
		})
	}
}
