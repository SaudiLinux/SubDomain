package utils

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"
)

// IsValidDomain checks if a domain is valid
func IsValidDomain(domain string) bool {
	// Simple validation to check if the domain has at least one dot
	// and doesn't contain invalid characters
	return strings.Contains(domain, ".") && 
		!strings.Contains(domain, " ") && 
		!strings.Contains(domain, "http://") && 
		!strings.Contains(domain, "https://")
}

// ResolveDomain resolves a domain to its IP address
func ResolveDomain(domain string) (string, error) {
	ips, err := net.LookupHost(domain)
	if err != nil {
		return "", err
	}

	if len(ips) == 0 {
		return "", fmt.Errorf("no IP addresses found")
	}

	return ips[0], nil
}

// CheckPort checks if a port is open on a host
func CheckPort(host string, port int, timeout time.Duration) bool {
	address := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// GetHTTPInfo gets information about an HTTP server
func GetHTTPInfo(url string, timeout time.Duration) (int, string, map[string][]string, error) {
	client := &http.Client{
		Timeout: timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Don't follow redirects
			return http.ErrUseLastResponse
		},
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, "", nil, err
	}

	// Set a user agent to avoid being blocked
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	req = req.WithContext(ctx)

	resp, err := client.Do(req)
	if err != nil {
		return 0, "", nil, err
	}
	defer resp.Body.Close()

	return resp.StatusCode, resp.Status, resp.Header, nil
}

// GetHTTPTitle gets the title of an HTTP page
func GetHTTPTitle(url string, timeout time.Duration) (string, error) {
	client := &http.Client{
		Timeout: timeout,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	// Set a user agent to avoid being blocked
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	req = req.WithContext(ctx)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read up to 8KB of the response body to look for the title
	buf := make([]byte, 8192)
	n, err := resp.Body.Read(buf)
	if err != nil && err.Error() != "EOF" {
		return "", err
	}

	body := string(buf[:n])

	// Extract title using a simple string search
	titleStart := strings.Index(body, "<title>")
	if titleStart == -1 {
		return "", nil
	}
	titleStart += 7 // Length of "<title>"

	titleEnd := strings.Index(body[titleStart:], "</title>")
	if titleEnd == -1 {
		return "", nil
	}

	title := body[titleStart : titleStart+titleEnd]
	return strings.TrimSpace(title), nil
}