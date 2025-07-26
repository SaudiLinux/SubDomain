package main

import (
	"fmt"
	"os"

	"github.com/SayerLinux/sub/cmd"
)

func main() {
	fmt.Println("\n\033[1;32m[+] Sub Domain Tool - By SayerLinux\033[0m")
	fmt.Println("\033[1;32m[+] Email: SaudiSayer@gmail.com\033[0m\n")

	rootCmd := cmd.NewRootCmd()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}