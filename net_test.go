package criterio

import (
	"testing"
)

func TestNetIP(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"valid IPv4", "192.168.1.1", false},
		{"valid IPv6", "2001:db8::1", false},
		{"valid IPv6 full", "2001:0db8:0000:0000:0000:0000:0000:0001", false},
		{"invalid IP", "999.999.999.999", true},
		{"not an IP", "hello", true},
		{"empty string", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NetIP()(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("NetIP()(%q) error = %v, wantErr %v", tt.value, err, tt.wantErr)
			}
		})
	}
}

func TestNetIPv4(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"valid IPv4", "192.168.1.1", false},
		{"valid IPv4 localhost", "127.0.0.1", false},
		{"valid IPv4 all zeros", "0.0.0.0", false},
		{"IPv6 address", "2001:db8::1", true},
		{"invalid IP", "999.999.999.999", true},
		{"partial IPv4", "192.168.1", true},
		{"not an IP", "hello", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NetIPv4()(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("NetIPv4()(%q) error = %v, wantErr %v", tt.value, err, tt.wantErr)
			}
		})
	}
}

func TestNetIPv6(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"valid IPv6 short", "2001:db8::1", false},
		{"valid IPv6 full", "2001:0db8:0000:0000:0000:0000:0000:0001", false},
		{"valid IPv6 localhost", "::1", false},
		{"IPv4 address", "192.168.1.1", true},
		{"invalid IP", "hello", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NetIPv6()(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("NetIPv6()(%q) error = %v, wantErr %v", tt.value, err, tt.wantErr)
			}
		})
	}
}

func TestNetCIDR(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"valid IPv4 CIDR", "192.168.1.0/24", false},
		{"valid IPv4 CIDR /32", "192.168.1.1/32", false},
		{"valid IPv6 CIDR", "2001:db8::/32", false},
		{"missing prefix", "192.168.1.0", true},
		{"invalid prefix", "192.168.1.0/33", true},
		{"not CIDR", "hello", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NetCIDR()(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("NetCIDR()(%q) error = %v, wantErr %v", tt.value, err, tt.wantErr)
			}
		})
	}
}

func TestNetHost(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"simple hostname", "localhost", false},
		{"domain", "example.com", false},
		{"subdomain", "www.example.com", false},
		{"with hyphen", "my-server.example.com", false},
		{"with numbers", "server1.example.com", false},
		{"starts with hyphen", "-example.com", true},
		{"ends with hyphen", "example-.com", true},
		{"empty string", "", true},
		{"too long", string(make([]byte, 254)), true},
		{"with underscore", "my_server.example.com", true},
		{"with space", "my server.com", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NetHost()(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("NetHost()(%q) error = %v, wantErr %v", tt.value, err, tt.wantErr)
			}
		})
	}
}

func TestNetPort(t *testing.T) {
	tests := []struct {
		name    string
		value   int
		wantErr bool
	}{
		{"valid port 80", 80, false},
		{"valid port 443", 443, false},
		{"valid port 1", 1, false},
		{"valid port 65535", 65535, false},
		{"port 0", 0, true},
		{"negative port", -1, true},
		{"port too high", 65536, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NetPort()(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("NetPort()(%d) error = %v, wantErr %v", tt.value, err, tt.wantErr)
			}
		})
	}
}

func TestNetPortStr(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"valid port 80", "80", false},
		{"valid port 443", "443", false},
		{"valid port 1", "1", false},
		{"valid port 65535", "65535", false},
		{"port 0", "0", true},
		{"negative port", "-1", true},
		{"port too high", "65536", true},
		{"not a number", "abc", true},
		{"empty string", "", true},
		{"float", "80.5", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NetPortStr()(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("NetPortStr()(%q) error = %v, wantErr %v", tt.value, err, tt.wantErr)
			}
		})
	}
}
