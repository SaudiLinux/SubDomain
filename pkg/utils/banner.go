package utils

import (
	"fmt"
	"github.com/fatih/color"
)

// PrintBanner prints the tool banner
func PrintBanner() {
	blue := color.New(color.FgBlue).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	fmt.Println(blue(`
   _____       __    
  / ___/__  __/ /_   
  \__ \/ / / / __ \
 ___/ / /_/ / /_/ /
/____/\__,_/_.___/ 
`))
	fmt.Printf("%s %s %s\n", green("[*]"), "Sub Tool - Subdomain Discovery and File Extraction", green("[*]"))
	fmt.Printf("%s %s %s\n", yellow("[*]"), "Created by SayerLinux (SaudiSayer@gmail.com)", yellow("[*]"))
	fmt.Printf("%s %s %s\n\n", red("[*]"), "Use with caution and responsibility", red("[*]"))
}

// PrintVersion prints the version information
func PrintVersion(version string) {
	fmt.Printf("Sub Tool v%s\n", version)
	fmt.Println("Created by SayerLinux (SaudiSayer@gmail.com)")
}