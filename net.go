package criterio

import (
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
)

// NetIP returns a validator that checks if a string is a valid IP address (IPv4 or IPv6).
func NetIP() Validator[string] {
	return func(val string) error {
		if net.ParseIP(val) == nil {
			return fmt.Errorf("must be a valid IP address")
		}
		return nil
	}
}

// NetIPv4 returns a validator that checks if a string is a valid IPv4 address.
func NetIPv4() Validator[string] {
	return func(val string) error {
		ip := net.ParseIP(val)
		if ip == nil || ip.To4() == nil {
			return fmt.Errorf("must be a valid IPv4 address")
		}
		return nil
	}
}

// NetIPv6 returns a validator that checks if a string is a valid IPv6 address.
func NetIPv6() Validator[string] {
	return func(val string) error {
		ip := net.ParseIP(val)
		if ip == nil || ip.To4() != nil {
			return fmt.Errorf("must be a valid IPv6 address")
		}
		return nil
	}
}

// NetCIDR returns a validator that checks if a string is valid CIDR notation.
func NetCIDR() Validator[string] {
	return func(val string) error {
		_, _, err := net.ParseCIDR(val)
		if err != nil {
			return fmt.Errorf("must be valid CIDR notation")
		}
		return nil
	}
}

var hostnameRegex = regexp.MustCompile(`^([a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?\.)*[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?$`)

// NetHost returns a validator that checks if a string is a valid hostname.
func NetHost() Validator[string] {
	return func(val string) error {
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
}

// NetPort returns a validator that checks if an integer is a valid port number (1-65535).
func NetPort() Validator[int] {
	return func(val int) error {
		if val < 1 || val > 65535 {
			return fmt.Errorf("must be a valid port number (1-65535)")
		}
		return nil
	}
}

// NetPortStr returns a validator that checks if a string represents a valid port number (1-65535).
func NetPortStr() Validator[string] {
	return func(val string) error {
		port, err := strconv.Atoi(val)
		if err != nil || port < 1 || port > 65535 {
			return fmt.Errorf("must be a valid port number (1-65535)")
		}
		return nil
	}
}
