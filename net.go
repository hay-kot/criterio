package criterio

import (
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
)

// NetIP validates that a string is a valid IP address (IPv4 or IPv6).
var NetIP Validator[string] = func(val string) error {
	if net.ParseIP(val) == nil {
		return fmt.Errorf("must be a valid IP address")
	}
	return nil
}

// NetIPv4 validates that a string is a valid IPv4 address.
var NetIPv4 Validator[string] = func(val string) error {
	ip := net.ParseIP(val)
	if ip == nil || ip.To4() == nil {
		return fmt.Errorf("must be a valid IPv4 address")
	}
	return nil
}

// NetIPv6 validates that a string is a valid IPv6 address.
var NetIPv6 Validator[string] = func(val string) error {
	ip := net.ParseIP(val)
	if ip == nil || ip.To4() != nil {
		return fmt.Errorf("must be a valid IPv6 address")
	}
	return nil
}

// NetCIDR validates that a string is valid CIDR notation.
var NetCIDR Validator[string] = func(val string) error {
	_, _, err := net.ParseCIDR(val)
	if err != nil {
		return fmt.Errorf("must be valid CIDR notation")
	}
	return nil
}

var hostnameRegex = regexp.MustCompile(`^([a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?\.)*[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?$`)

// NetHost validates that a string is a valid hostname.
var NetHost Validator[string] = func(val string) error {
	if len(val) > 253 {
		return fmt.Errorf("must be a valid hostname")
	}
	if !hostnameRegex.MatchString(val) {
		return fmt.Errorf("must be a valid hostname")
	}
	// Check each label length
	for _, label := range strings.Split(val, ".") {
		if len(label) > 63 {
			return fmt.Errorf("must be a valid hostname")
		}
	}
	return nil
}

// NetPort validates that an integer is a valid port number (1-65535).
var NetPort Validator[int] = func(val int) error {
	if val < 1 || val > 65535 {
		return fmt.Errorf("must be a valid port number (1-65535)")
	}
	return nil
}

// NetPortStr validates that a string represents a valid port number (1-65535).
var NetPortStr Validator[string] = func(val string) error {
	port, err := strconv.Atoi(val)
	if err != nil || port < 1 || port > 65535 {
		return fmt.Errorf("must be a valid port number (1-65535)")
	}
	return nil
}
