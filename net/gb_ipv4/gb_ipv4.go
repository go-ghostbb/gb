// Package gbipv4 provides useful API for IPv4 address handling.
package gbipv4

import (
	"encoding/binary"
	"fmt"
	gbregex "github.com/Ghostbb-io/gb/text/gb_regex"
	"net"
	"strconv"
)

// Ip2long converts ip address to an uint32 integer.
func Ip2long(ip string) uint32 {
	netIp := net.ParseIP(ip)
	if netIp == nil {
		return 0
	}
	return binary.BigEndian.Uint32(netIp.To4())
}

// Long2ip converts an uint32 integer ip address to its string type address.
func Long2ip(long uint32) string {
	ipByte := make([]byte, 4)
	binary.BigEndian.PutUint32(ipByte, long)
	return net.IP(ipByte).String()
}

// Validate checks whether given `ip` a valid IPv4 address.
func Validate(ip string) bool {
	return gbregex.IsMatchString(`^((25[0-5]|2[0-4]\d|[01]?\d\d?)\.){3}(25[0-5]|2[0-4]\d|[01]?\d\d?)$`, ip)
}

// ParseAddress parses `address` to its ip and port.
// Eg: 192.168.1.1:80 -> 192.168.1.1, 80
func ParseAddress(address string) (string, int) {
	match, err := gbregex.MatchString(`^(.+):(\d+)$`, address)
	if err == nil {
		i, _ := strconv.Atoi(match[2])
		return match[1], i
	}
	return "", 0
}

// GetSegment returns the segment of given ip address.
// Eg: 192.168.2.102 -> 192.168.2
func GetSegment(ip string) string {
	match, err := gbregex.MatchString(`^(\d{1,3})\.(\d{1,3})\.(\d{1,3})\.(\d{1,3})$`, ip)
	if err != nil || len(match) < 4 {
		return ""
	}
	return fmt.Sprintf("%s.%s.%s", match[1], match[2], match[3])
}
