package server

import "testing"

func TestNormalizeIPv4(t *testing.T) {
	type args struct {
		address string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "Valid IPv4 address",
			args:    args{address: "192.168.1.1"},
			want:    "192.168.1.1",
			wantErr: false,
		},
		{
			name:    "Invalid IP address",
			args:    args{address: "invalid_ip"},
			want:    "",
			wantErr: true,
		},
		{
			name:    "IPv6 address mapped to IPv4",
			args:    args{address: "::ffff:192.168.1.1"},
			want:    "192.168.1.1",
			wantErr: false,
		},
		{
			name:    "Loopback address",
			args:    args{address: "::1"},
			want:    "127.0.0.1",
			wantErr: false,
		},
		{
			name:    "Not an IPv4-mapped IPv6 address",
			args:    args{address: "2001:db8::1"},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NormalizeIPv4(tt.args.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("NormalizeIPv4() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NormalizeIPv4() got = %v, want %v", got, tt.want)
			}
		})
	}
}
