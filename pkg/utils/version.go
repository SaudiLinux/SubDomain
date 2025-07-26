package utils

// Version information
const (
	// Version is the current version of the tool
	Version = "1.0.0"

	// Author is the author of the tool
	Author = "SayerLinux"

	// Email is the contact email of the author
	Email = "SaudiSayer@gmail.com"

	// Website is the website of the tool
	Website = "https://github.com/SayerLinux/sub"

	// Description is a short description of the tool
	Description = "Sub is a tool for discovering hidden subdomains and extracting files from targets"
)

// GetVersionInfo returns a formatted version string
func GetVersionInfo() string {
	return "Sub v" + Version + " - Created by " + Author + " (" + Email + ")"
}