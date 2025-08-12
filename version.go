package gobase

// Version information - this file can be updated by build scripts
var (
	Version     = "v0.0.2"
	Name        = "GoBase"
	Description = "A Django-inspired ORM and database toolkit for Go"
)

// GetVersion returns the current version
func GetVersion() string {
	return Version
}

// GetFullVersion returns version with name and description
func GetFullVersion() string {
	return Name + " " + Version + "\n" + Description
}
