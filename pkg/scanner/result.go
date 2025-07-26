package scanner

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/SayerLinux/sub/pkg/utils"
)

// Result represents a subdomain scan result
type Result struct {
	Subdomain string
	IP        string
	Found     bool
	Timestamp time.Time
}

// ServiceResult represents a service scan result
type ServiceResult struct {
	Subdomain string
	Port      int
	Service   string
	Info      string
	Timestamp time.Time
}

// FileResult represents a file extraction result
type FileResult struct {
	Subdomain string
	FilePath  string
	Success   bool
	Size      int64
	Timestamp time.Time
}

// ResultManager manages scan results
type ResultManager struct {
	results        []Result
	serviceResults []ServiceResult
	fileResults    []FileResult
	outputPath     string
	outputDir      string
	logger         *utils.Logger
	mutex          sync.Mutex
}

// NewResultManager creates a new result manager
func NewResultManager(outputPath string, outputDir string, logger *utils.Logger) *ResultManager {
	return &ResultManager{
		results:        []Result{},
		serviceResults: []ServiceResult{},
		fileResults:    []FileResult{},
		outputPath:     outputPath,
		outputDir:      outputDir,
		logger:         logger,
	}
}

// AddResult adds a subdomain result
func (rm *ResultManager) AddResult(subdomain string, ip string, found bool) {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	result := Result{
		Subdomain: subdomain,
		IP:        ip,
		Found:     found,
		Timestamp: time.Now(),
	}

	rm.results = append(rm.results, result)
	rm.logger.Result(subdomain, found, ip)
}

// AddServiceResult adds a service result
func (rm *ResultManager) AddServiceResult(subdomain string, port int, service string, info string) {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	result := ServiceResult{
		Subdomain: subdomain,
		Port:      port,
		Service:   service,
		Info:      info,
		Timestamp: time.Now(),
	}

	rm.serviceResults = append(rm.serviceResults, result)
	rm.logger.ServiceResult(subdomain, port, service, info)
}

// AddFileResult adds a file extraction result
func (rm *ResultManager) AddFileResult(subdomain string, filePath string, success bool, size int64) {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	result := FileResult{
		Subdomain: subdomain,
		FilePath:  filePath,
		Success:   success,
		Size:      size,
		Timestamp: time.Now(),
	}

	rm.fileResults = append(rm.fileResults, result)
	rm.logger.FileResult(subdomain, filePath, success, size)
}

// GetResults returns all subdomain results
func (rm *ResultManager) GetResults() []Result {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	return rm.results
}

// GetFoundResults returns only found subdomain results
func (rm *ResultManager) GetFoundResults() []Result {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	var foundResults []Result
	for _, result := range rm.results {
		if result.Found {
			foundResults = append(foundResults, result)
		}
	}

	return foundResults
}

// SaveResults saves the results to a file
func (rm *ResultManager) SaveResults() error {
	if rm.outputPath == "" {
		return nil
	}

	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	// Create the output file
	file, err := utils.CreateOutputFile(rm.outputPath, "# Sub Tool Results - Generated on "+time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		return err
	}
	defer file.Close()

	// Write found subdomains to the file
	for _, result := range rm.results {
		if result.Found {
			_, err := fmt.Fprintf(file, "%s,%s\n", result.Subdomain, result.IP)
			if err != nil {
				return fmt.Errorf("failed to write to output file: %v", err)
			}
		}
	}

	rm.logger.Success("Results saved to %s", rm.outputPath)
	return nil
}

// SaveServiceResults saves the service results to a file
func (rm *ResultManager) SaveServiceResults() error {
	if rm.outputDir == "" {
		return nil
	}

	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	// Ensure the output directory exists
	if err := utils.EnsureDirectory(rm.outputDir); err != nil {
		return err
	}

	// Create the services output file
	servicesPath := filepath.Join(rm.outputDir, "services.txt")
	file, err := utils.CreateOutputFile(servicesPath, "# Sub Tool Service Results - Generated on "+time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		return err
	}
	defer file.Close()

	// Write service results to the file
	for _, result := range rm.serviceResults {
		_, err := fmt.Fprintf(file, "%s:%d,%s,%s\n", result.Subdomain, result.Port, result.Service, result.Info)
		if err != nil {
			return fmt.Errorf("failed to write to services file: %v", err)
		}
	}

	rm.logger.Success("Service results saved to %s", servicesPath)
	return nil
}

// SaveFileResults saves the file extraction results to a file
func (rm *ResultManager) SaveFileResults() error {
	if rm.outputDir == "" {
		return nil
	}

	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	// Ensure the output directory exists
	if err := utils.EnsureDirectory(rm.outputDir); err != nil {
		return err
	}

	// Create the files output file
	filesPath := filepath.Join(rm.outputDir, "files.txt")
	file, err := utils.CreateOutputFile(filesPath, "# Sub Tool File Extraction Results - Generated on "+time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		return err
	}
	defer file.Close()

	// Write file results to the file
	for _, result := range rm.fileResults {
		if result.Success {
			_, err := fmt.Fprintf(file, "%s,%s,%d\n", result.Subdomain, result.FilePath, result.Size)
			if err != nil {
				return fmt.Errorf("failed to write to files file: %v", err)
			}
		}
	}

	rm.logger.Success("File results saved to %s", filesPath)
	return nil
}

// SaveAllResults saves all results
func (rm *ResultManager) SaveAllResults() error {
	// Save subdomain results
	if err := rm.SaveResults(); err != nil {
		return err
	}

	// Save service results
	if err := rm.SaveServiceResults(); err != nil {
		return err
	}

	// Save file results
	if err := rm.SaveFileResults(); err != nil {
		return err
	}

	return nil
}

// GenerateSummary generates a summary of the scan results
func (rm *ResultManager) GenerateSummary() string {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	// Count found subdomains
	foundCount := 0
	for _, result := range rm.results {
		if result.Found {
			foundCount++
		}
	}

	// Build summary
	var sb strings.Builder

	sb.WriteString("\n=== Sub Tool Scan Summary ===\n")
	sb.WriteString(fmt.Sprintf("Total subdomains scanned: %d\n", len(rm.results)))
	sb.WriteString(fmt.Sprintf("Subdomains found: %d\n", foundCount))
	sb.WriteString(fmt.Sprintf("Services discovered: %d\n", len(rm.serviceResults)))
	sb.WriteString(fmt.Sprintf("Files extracted: %d\n", len(rm.fileResults)))

	// Add output file information
	if rm.outputPath != "" {
		sb.WriteString(fmt.Sprintf("Results saved to: %s\n", rm.outputPath))
	}

	if rm.outputDir != "" {
		sb.WriteString(fmt.Sprintf("Service results saved to: %s\n", filepath.Join(rm.outputDir, "services.txt")))
		sb.WriteString(fmt.Sprintf("File results saved to: %s\n", filepath.Join(rm.outputDir, "files.txt")))
	}

	return sb.String()
}