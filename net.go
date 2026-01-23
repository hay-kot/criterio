package criterio

import (
	"errors"
	"net"
	"regexp"
	"strconv"
	"strings"
)

// Static error messages for network validators.
var (
	errNetIP       = errors.New("must be a valid IP address")
	errNetIPv4     = errors.New("must be a valid IPv4 address")
	errNetIPv6     = errors.New("must be a valid IPv6 address")
	errNetCIDR     = errors.New("must be valid CIDR notation")
	errNetHostname = errors.New("must be a valid hostname")
	errNetPort     = errors.New("must be a valid port number (1-65535)")
)

// Default regex patterns used by network validators.
// These can be modified to change validation behavior globally.
//
// WARNING: Modifying these affects all validations application-wide.
// Concurrent modification is not thread-safe. If you need to customize
// patterns, do so at program initialization before any validation occurs.
var (
	HostnameRegex = regexp.MustCompile(`^([a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?\.)*[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?$`)
)

// NetIP validates that a string is a valid IP address (IPv4 or IPv6).
func NetIP(val string) error {
	if net.ParseIP(val) == nil {
		return errNetIP
	}
	return nil
}

// NetIPv4 validates that a string is a valid IPv4 address.
func NetIPv4(val string) error {
	ip := net.ParseIP(val)
	if ip == nil || ip.To4() == nil {
		return errNetIPv4
	}
	return nil
}

// NetIPv6 validates that a string is a valid IPv6 address.
func NetIPv6(val string) error {
	ip := net.ParseIP(val)
	if ip == nil || ip.To4() != nil {
		return errNetIPv6
	}
	return nil
}

// NetCIDR validates that a string is valid CIDR notation.
func NetCIDR(val string) error {
	_, _, err := net.ParseCIDR(val)
	if err != nil {
		return errNetCIDR
	}
	return nil
}

// NetHost validates that a string is a valid hostname.
// Uses HostnameRegex which can be modified to change validation globally.
func NetHost(val string) error {
	if len(val) > 253 {
		return errNetHostname
	}
	if !HostnameRegex.MatchString(val) {
		return errNetHostname
	}
	// Check each label length
	for _, label := range strings.Split(val, ".") {
		if len(label) > 63 {
			return errNetHostname
		}
	}
	return nil
}

// NetPort validates that an integer is a valid port number (1-65535).
func NetPort(val int) error {
	if val < 1 || val > 65535 {
		return errNetPort
	}
	return nil
}

// NetPortStr validates that a string represents a valid port number (1-65535).
func NetPortStr(val string) error {
	port, err := strconv.Atoi(val)
	if err != nil || port < 1 || port > 65535 {
		return errNetPort
	}
	return nil
}
