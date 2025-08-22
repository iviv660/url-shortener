package validate

import (
	"context"
	"errors"
	"net"
	"net/netip"
	"net/url"
	"strings"
)

var (
	errEmpty                = errors.New("url is empty")
	errTooLong              = errors.New("url is too long")
	errInvalidURL           = errors.New("invalid url")
	errUnsupportedScheme    = errors.New("unsupported scheme (only http/https)")
	errMissingHost          = errors.New("missing host")
	errPrivateOrLocalTarget = errors.New("url points to private/local address")
)

const MaxURLLength = 2000

func Validate(ctx context.Context, raw string) (string, error) {
	if strings.TrimSpace(raw) == "" {
		return "", errEmpty
	}

	if len(raw) > MaxURLLength {
		return "", errTooLong
	}

	raw = strings.TrimSpace(raw)
	raw = strings.ReplaceAll(raw, "\n", "")
	raw = strings.ReplaceAll(raw, "\r", "")

	u, err := url.ParseRequestURI(raw)
	if err != nil {
		return "", errInvalidURL
	}

	switch strings.ToLower(u.Scheme) {
	case "http", "https":
	default:
		return "", errUnsupportedScheme
	}

	if u.Host == "" {
		return "", errMissingHost
	}

	host := strings.ToLower(u.Host)
	if h, p, err := net.SplitHostPort(host); err == nil {
		host = h
		u.Host = net.JoinHostPort(strings.ToLower(h), p)
	} else {
		u.Host = host
	}

	u.Fragment = ""

	if host == "localhost" || strings.HasSuffix(host, ".localhost") {
		return "", errPrivateOrLocalTarget
	}
	if ip, ok := parseIP(host); ok && isPrivateOrSpecialIP(ip) {
		return "", errPrivateOrLocalTarget
	}

	return u.String(), nil
}

func parseIP(host string) (netip.Addr, bool) {
	ip, err := netip.ParseAddr(host)
	if err != nil {
		return netip.Addr{}, false
	}
	return ip, true
}

func isPrivateOrSpecialIP(ip netip.Addr) bool {
	if ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() ||
		ip.IsMulticast() || ip.IsUnspecified() || ip.IsPrivate() {
		return true
	}
	// fc00::/7 (ULA IPv6)
	if ip.Is6() {
		b := ip.AsSlice()
		if len(b) > 0 && (b[0]&0xFE) == 0xFC {
			return true
		}
	}
	return false
}
