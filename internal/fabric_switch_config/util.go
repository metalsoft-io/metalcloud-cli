package fabric_switch_config

import (
	"fmt"
	"net"
	"strconv"
)

// ConfigError is a rule/validation violation in the configuration or device
// tags. It mirrors the Python ConfigError: it aborts the whole computation.
type ConfigError struct{ msg string }

func (e *ConfigError) Error() string { return e.msg }

func configErrorf(format string, args ...any) error {
	return &ConfigError{msg: fmt.Sprintf(format, args...)}
}

func itoa(n int64) string { return strconv.FormatInt(n, 10) }

// logicalToSwp returns the breakout sub-port name for a logical split index
// (1=swp1s0, 2=swp1s1, ..., 128=swp64s1).
func logicalToSwp(logical int) string {
	return fmt.Sprintf("swp%ds%d", (logical+1)/2, 1-(logical%2))
}

// ipv4ToUint parses a dotted IPv4 string into a uint32. Returns ok=false for
// non-IPv4 input.
func ipv4ToUint(s string) (uint32, bool) {
	ip := net.ParseIP(s)
	if ip == nil {
		return 0, false
	}
	ip4 := ip.To4()
	if ip4 == nil {
		return 0, false
	}
	return uint32(ip4[0])<<24 | uint32(ip4[1])<<16 | uint32(ip4[2])<<8 | uint32(ip4[3]), true
}

func uintToIpv4(v uint32) string {
	return fmt.Sprintf("%d.%d.%d.%d", byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
}

// ipv4Network is a parsed CIDR (base address + prefix length).
type ipv4Network struct {
	base      uint32
	prefixLen int
}

func parseIpv4Network(cidr string) (ipv4Network, error) {
	ip, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return ipv4Network{}, err
	}
	ip4 := ipNet.IP.To4()
	if ip4 == nil || ip.To4() == nil {
		return ipv4Network{}, fmt.Errorf("not an IPv4 network: %s", cidr)
	}
	ones, _ := ipNet.Mask.Size()
	base := uint32(ip4[0])<<24 | uint32(ip4[1])<<16 | uint32(ip4[2])<<8 | uint32(ip4[3])
	return ipv4Network{base: base, prefixLen: ones}, nil
}

// lastAddress returns the broadcast (last) address of the network.
func (n ipv4Network) lastAddress() uint32 {
	hostBits := 32 - n.prefixLen
	if hostBits <= 0 {
		return n.base
	}
	return n.base | (uint32(1)<<uint(hostBits) - 1)
}

// containsSubnet reports whether a /31 at base 'addr' fits entirely within n.
func (n ipv4Network) containsSubnet(addr uint32, prefixLen int) bool {
	if prefixLen < n.prefixLen {
		return false
	}
	hostBits := 32 - prefixLen
	last := addr | (uint32(1)<<uint(hostBits) - 1)
	return addr >= n.base && last <= n.lastAddress()
}
