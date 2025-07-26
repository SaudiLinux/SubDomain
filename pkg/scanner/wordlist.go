package scanner

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/SayerLinux/sub/pkg/utils"
)

// DefaultWordlistPath is the path to the default wordlist
const DefaultWordlistPath = "wordlists/default.txt"

// WordlistManager handles wordlist operations
type WordlistManager struct {
	wordlistPath string
	wordlist     []string
	logger       *utils.Logger
}

// NewWordlistManager creates a new wordlist manager
func NewWordlistManager(wordlistPath string, logger *utils.Logger) *WordlistManager {
	return &WordlistManager{
		wordlistPath: wordlistPath,
		logger:       logger,
	}
}

// Load loads the wordlist
func (wm *WordlistManager) Load() error {
	// If no wordlist path is provided, use the default
	if wm.wordlistPath == "" {
		// Check if the default wordlist exists
		if _, err := os.Stat(DefaultWordlistPath); os.IsNotExist(err) {
			// Try to find the default wordlist in the executable directory
			execPath, err := os.Executable()
			if err != nil {
				return fmt.Errorf("failed to get executable path: %v", err)
			}
			execDir := filepath.Dir(execPath)
			defaultPath := filepath.Join(execDir, DefaultWordlistPath)

			if _, err := os.Stat(defaultPath); os.IsNotExist(err) {
				return fmt.Errorf("default wordlist not found at %s or %s", DefaultWordlistPath, defaultPath)
			}

			wm.wordlistPath = defaultPath
		} else {
			wm.wordlistPath = DefaultWordlistPath
		}
	}

	wm.logger.Debug("Loading wordlist from %s", wm.wordlistPath)
	words, err := utils.LoadWordlist(wm.wordlistPath)
	if err != nil {
		return err
	}

	wm.wordlist = words
	wm.logger.Info("Loaded %d words from wordlist", len(wm.wordlist))
	return nil
}

// GetWordlist returns the loaded wordlist
func (wm *WordlistManager) GetWordlist() []string {
	return wm.wordlist
}

// GenerateSubdomains generates subdomains for a target domain
func (wm *WordlistManager) GenerateSubdomains(domain string) []string {
	var subdomains []string

	// Clean the domain (remove http://, https://, trailing slashes)
	domain = strings.TrimPrefix(domain, "http://")
	domain = strings.TrimPrefix(domain, "https://")
	domain = strings.Split(domain, "/")[0]

	// Generate subdomains
	for _, word := range wm.wordlist {
		subdomains = append(subdomains, fmt.Sprintf("%s.%s", word, domain))
	}

	return subdomains
}