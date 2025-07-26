package scanner

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
)

// ServiceInfo represents information about a service running on a subdomain
type ServiceInfo struct {
	Subdomain  string
	IP         string
	Port       int
	Service    string
	StatusCode int
	Title      string
	Server     string
}

// CheckCommonPorts checks common ports on a subdomain
func CheckCommonPorts(subdomain string, ip string) []ServiceInfo {
	commonPorts := []struct {
		Port    int
		Service string
	}{
		{80, "HTTP"},
		{443, "HTTPS"},
		{21, "FTP"},
		{22, "SSH"},
		{25, "SMTP"},
		{53, "DNS"},
		{8080, "HTTP-ALT"},
		{8443, "HTTPS-ALT"},
	}

	var results []ServiceInfo

	for _, port := range commonPorts {
		addr := fmt.Sprintf("%s:%d", ip, port.Port)
		conn, err := net.DialTimeout("tcp", addr, 2*time.Second)
		if err == nil {
			conn.Close()
			info := ServiceInfo{
				Subdomain: subdomain,
				IP:        ip,
				Port:      port.Port,
				Service:   port.Service,
			}

			// If HTTP or HTTPS, get more information
			if port.Port == 80 || port.Port == 8080 {
				getHTTPInfo(&info, false)
			} else if port.Port == 443 || port.Port == 8443 {
				getHTTPInfo(&info, true)
			}

			results = append(results, info)
		}
	}

	return results
}

// getHTTPInfo gets HTTP information from a service
func getHTTPInfo(info *ServiceInfo, isHTTPS bool) {
	var url string
	if isHTTPS {
		url = fmt.Sprintf("https://%s:%d", info.Subdomain, info.Port)
	} else {
		url = fmt.Sprintf("http://%s:%d", info.Subdomain, info.Port)
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	resp, err := client.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	info.StatusCode = resp.StatusCode
	info.Server = resp.Header.Get("Server")

	// Extract title if possible
	body, err := io.ReadAll(io.LimitReader(resp.Body, 10000))
	if err == nil {
		bodyStr := string(body)
		titleStart := strings.Index(bodyStr, "<title>")
		if titleStart != -1 {
			titleEnd := strings.Index(bodyStr[titleStart:], "</title>")
			if titleEnd != -1 {
				title := bodyStr[titleStart+7 : titleStart+titleEnd]
				info.Title = strings.TrimSpace(title)
			}
		}
	}
}

// ExtractFiles attempts to extract files from a subdomain
func ExtractFiles(subdomain string, outputDir string) error {
	// Common file paths to check
	commonPaths := []string{
		"/robots.txt",
		"/sitemap.xml",
		"/.git/HEAD",
		"/.env",
		"/wp-config.php",
		"/config.php",
		"/admin/",
		"/backup/",
		"/database/",
		"/api/",
	}

	// Create output directory for this subdomain
	subdomainDir := filepath.Join(outputDir, subdomain)
	err := os.MkdirAll(subdomainDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	// Check each path
	for _, path := range commonPaths {
		url := fmt.Sprintf("http://%s%s", subdomain, path)
		secureURL := fmt.Sprintf("https://%s%s", subdomain, path)

		// Try HTTP first
		if downloadFile(url, subdomainDir, path) {
			fmt.Printf("%s Found: %s\n", color.GreenString("[+]"), url)
			continue
		}

		// Try HTTPS if HTTP failed
		if downloadFile(secureURL, subdomainDir, path) {
			fmt.Printf("%s Found: %s\n", color.GreenString("[+]"), secureURL)
		}
	}

	return nil
}

// downloadFile downloads a file from a URL
func downloadFile(url string, outputDir string, path string) bool {
	client := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	resp, err := client.Get(url)
	if err != nil || resp.StatusCode >= 400 {
		if resp != nil {
			resp.Body.Close()
		}
		return false
	}
	defer resp.Body.Close()

	// Create file path
	filePath := filepath.Join(outputDir, strings.TrimPrefix(path, "/"))

	// Create directory if needed
	dir := filepath.Dir(filePath)
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return false
	}

	// Create file
	file, err := os.Create(filePath)
	if err != nil {
		return false
	}
	defer file.Close()

	// Write to file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return false
	}

	return true
}