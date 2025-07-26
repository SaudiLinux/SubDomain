package cmd

import (
	"fmt"
	"os"

	"github.com/SayerLinux/sub/pkg/scanner"
	"github.com/spf13/cobra"
)

// NewRootCmd creates the root command for the application
func NewRootCmd() *cobra.Command {
	// Add scan command
	scanCmd := NewScanCmd()
	var (
		target      string
		wordlist    string
		threads     int
		outputFile  string
		verbose     bool
		showVersion bool
	)

	rootCmd := &cobra.Command{
		Use:   "sub",
		Short: "Sub is a tool for discovering hidden subdomains",
		Long: `Sub is a powerful tool written in Go that helps discover hidden subdomains 
	and extract private and real files from the target.
	Developed by SayerLinux (SaudiSayer@gmail.com)`,
		Run: func(cmd *cobra.Command, args []string) {
			if showVersion {
				fmt.Println("Sub v1.0.0")
				return
			}

			if target == "" {
				fmt.Println("\033[1;31m[!] Error: Target domain is required\033[0m")
				cmd.Help()
				os.Exit(1)
			}

			if wordlist == "" {
				fmt.Println("\033[1;33m[!] Warning: No wordlist specified, using default wordlist\033[0m")
				wordlist = "./wordlists/default.txt"
			}

			// Create scanner configuration
			config := scanner.Config{
				Target:     target,
				Wordlist:   wordlist,
				Threads:    threads,
				OutputFile: outputFile,
				Verbose:    verbose,
			}

			// Start scanning
			scanner := scanner.NewScanner(config)
			scanner.Start()
		},
	}

	// Add flags
	rootCmd.Flags().StringVarP(&target, "target", "t", "", "Target domain to scan (required)")
	rootCmd.Flags().StringVarP(&wordlist, "wordlist", "w", "", "Path to wordlist file")
	rootCmd.Flags().IntVarP(&threads, "threads", "c", 50, "Number of concurrent threads")
	rootCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file to save results")
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	rootCmd.Flags().BoolVarP(&showVersion, "version", "", false, "Show version information")

	// Add subcommands
	rootCmd.AddCommand(scanCmd)

	return rootCmd
}