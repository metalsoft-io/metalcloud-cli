// Package fabric_template_config ports the configure_freeform.py and
// configure_bgp.py scripts: it computes the per-device freeform / BGP-underlay /
// EVPN-overlay / PFC profile variables from the same topology+p2p plan the
// fabric_switch_config engine produces, then registers the device-configuration
// templates (from .j2 bodies) and one variables-carrying profile per switch,
// idempotently. Template rendering happens server-side (Nunjucks); the optional
// verification uses the engine's stateless render endpoint (there is no local
// Jinja2 render, unlike the Python --verify-render parity gate).
package fabric_template_config

import (
	"net"

	fsc "github.com/metalsoft-io/metalcloud-cli/internal/fabric_switch_config"
)

// switchPositions is the fixed order switches are processed in.
var switchPositions = []string{"leaf", "spine", "super_spine"}

// deviceRecord is a fabric device plus the current-state fields the template
// engine reads (asn / loopback come from the configure_switches writes;
// customVariables is reconciled by the BGP run). The embedded fsc.Device is what
// the plan computation consumes.
type deviceRecord struct {
	fsc.Device
	Asn                 *int64
	LoopbackAddressIpv4 *string
	CustomVariables     map[string]interface{}
	Revision            string
}

// hostOf returns a device's planned hostname (fallback: its current identifier).
func hostOf(dev *fsc.Device, state *fsc.DesiredState, records map[int64]*deviceRecord) string {
	if d, ok := state.ByDevice[dev.Id]; ok && d.Hostname != nil && *d.Hostname != "" {
		return *d.Hostname
	}
	if rec, ok := records[dev.Id]; ok && rec.IdentifierString != "" {
		return rec.IdentifierString
	}
	return dev.Label()
}

// loopbackOf returns a device's loopback /32: the plan's assignment when the
// config manages loopbacks, else whatever the device record carries.
func loopbackOf(dev *fsc.Device, state *fsc.DesiredState, records map[int64]*deviceRecord) string {
	if d, ok := state.ByDevice[dev.Id]; ok && d.LoopbackIp != nil && *d.LoopbackIp != "" {
		return *d.LoopbackIp
	}
	if rec, ok := records[dev.Id]; ok && rec.LoopbackAddressIpv4 != nil {
		return *rec.LoopbackAddressIpv4
	}
	return ""
}

func ipToUint(s string) uint32 {
	ip := net.ParseIP(s).To4()
	if ip == nil {
		return 0
	}
	return uint32(ip[0])<<24 | uint32(ip[1])<<16 | uint32(ip[2])<<8 | uint32(ip[3])
}

func uintToIP(v uint32) string {
	return net.IPv4(byte(v>>24), byte(v>>16), byte(v>>8), byte(v)).String()
}

// subnetHostAddrs returns the two addresses of a /31 (interfaceA = the gateway =
// the smaller/even address = [0]; interfaceB = [1]).
func subnetHostAddrs(s *fsc.Subnet) (string, string) {
	base := ipToUint(s.NetworkAddress)
	return uintToIP(base), uintToIP(base + 1)
}
