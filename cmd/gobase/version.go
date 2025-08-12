package main

// Version information for GoBase CLI
var (
	Version     = "v0.0.1-alpha"
	Name        = "GoBase CLI"
	Description = "A Django-inspired ORM and database toolkit for Go"
	BuildDate   = "unknown"
	GitCommit   = "unknown"
)

// getVersionInfo returns formatted version information
func getVersionInfo() string {
	return Name + " " + Version + "\n" + Description
}

// getFullVersionInfo returns detailed version information
func getFullVersionInfo() string {
	return Name + " " + Version + "\n" +
		Description + "\n" +
		"Build Date: " + BuildDate + "\n" +
		"Git Commit: " + GitCommit
}
