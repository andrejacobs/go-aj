// Package buildinfo is used to include build time information into a compiled app.
//
// Please see github.com/andrejacobs/go-aj/examples/buildinfo/Makefile on how to use GO_LDFLAGS
// to set the information at compile time.
package buildinfo

import "fmt"

var (
	AppName       string // The name of the app (binary executable)
	GitCommitHash string // The hash of the last commit in the git repository
	Version       string // The version of the app
)

// Return the version information that can be shown to a user.
func VersionString() string {
	version := "v0.0.0"
	if Version != "" {
		version = Version
	}
	return fmt.Sprintf("%s %s", version, GitCommitHash)
}

// Return the name of the app
func UsageName() string {
	appName := "unknown"
	if AppName != "" {
		appName = AppName
	}
	return appName
}

// Return the name of the app and version info as displayed in the usage information
func UsageNameAndVersion() string {
	return fmt.Sprintf("%s version: %s", UsageName(), VersionString())
}
