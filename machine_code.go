package goid

import (
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"net"
	"os"
)

var netInterfaceAddrs = net.InterfaceAddrs
var osHostname = os.Hostname

func GetLocalIP() (ip string, err error) {
	addrs, err := netInterfaceAddrs()
	if err != nil {
		return
	}
	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ip = ipNet.IP.String()
				break
			}
		}
	}
	if ip == "" {
		err = errors.New("no non-loopback IP address found")
	}
	return
}

func GenerateMachineCode(bits int8) (code int, err error) {
	machineName, err := osHostname()
	if err != nil {
		return
	}
	ip, err := GetLocalIP()
	if err != nil {
		return
	}
	combinedString := machineName + "_" + ip
	hash := sha256.Sum256([]byte(combinedString))
	code = int(binary.BigEndian.Uint64(hash[:])) & int(1<<bits-1)
	return
}
