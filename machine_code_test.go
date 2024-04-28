package goid

import (
	"errors"
	"net"
	"os"
	"testing"
)

func mockInterfaceAddrsSuccess() ([]net.Addr, error) {
	return []net.Addr{
		&net.IPNet{IP: net.ParseIP("127.0.0.1"), Mask: net.CIDRMask(8, 32)},
		&net.IPNet{IP: net.ParseIP("192.168.1.2"), Mask: net.CIDRMask(24, 32)},
	}, nil
}

func mockInterfaceAddrsLoopbackOnly() ([]net.Addr, error) {
	return []net.Addr{
		&net.IPNet{IP: net.ParseIP("127.0.0.1"), Mask: net.CIDRMask(8, 32)},
	}, nil
}

func mockInterfaceAddrsError() ([]net.Addr, error) {
	return nil, errors.New("test error")
}

func TestGetLocalIP_Success(t *testing.T) {
	netInterfaceAddrs = mockInterfaceAddrsSuccess
	defer func() {
		netInterfaceAddrs = net.InterfaceAddrs
	}()

	ip, err := GetLocalIP()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if ip != "192.168.1.2" {
		t.Errorf("unexpected ip: got %v, want 192.168.1.2", ip)
	}
}

func TestGetLocalIP_LoopbackOnly(t *testing.T) {
	netInterfaceAddrs = mockInterfaceAddrsLoopbackOnly
	defer func() {
		netInterfaceAddrs = net.InterfaceAddrs
	}()

	ip, err := GetLocalIP()
	if err == nil || err.Error() != "no non-loopback IP address found" {
		t.Errorf("expected error: no non-loopback IP address found, got %v", err)
	}
	if ip != "" {
		t.Errorf("expected no ip, got %v", ip)
	}
}

func TestGetLocalIP_Error(t *testing.T) {
	netInterfaceAddrs = mockInterfaceAddrsError
	defer func() {
		netInterfaceAddrs = net.InterfaceAddrs
	}()

	ip, err := GetLocalIP()
	if err == nil || err.Error() != "test error" {
		t.Errorf("expected error: test error, got %v", err)
	}
	if ip != "" {
		t.Errorf("expected no ip, got %v", ip)
	}
}

func TestGenerateMachineCode(t *testing.T) {
	defer func() { osHostname = os.Hostname }()
	osHostname = func() (string, error) {
		return "1234", nil
	}

	netInterfaceAddrs = mockInterfaceAddrsSuccess
	defer func() {
		netInterfaceAddrs = net.InterfaceAddrs
	}()

	type args struct {
		bits int8
	}
	tests := []struct {
		name     string
		args     args
		wantCode int
		wantErr  bool
	}{
		{
			name: "GenerateMachineCode 8",
			args: args{
				bits: 8,
			},
			wantCode: 18,
			wantErr:  false,
		},
		{
			name: "GenerateMachineCode 32",
			args: args{
				bits: 32,
			},
			wantCode: 1352976146,
			wantErr:  false,
		},
		{
			name: "GenerateMachineCode 6",
			args: args{
				bits: 6,
			},
			wantCode: 18,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCode, err := GenerateMachineCode(tt.args.bits)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateMachineCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotCode != tt.wantCode {
				t.Errorf("GenerateMachineCode() gotCode = %v, want %v", gotCode, tt.wantCode)
			}
		})
	}
}
