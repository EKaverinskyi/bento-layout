// Package version holds the build version of bento-layout.
package version

// Version is overridden at build time:
//
//	go build -ldflags "-X github.com/EduardKav/bento-layout/cmd/internal/version.Version=v0.1.0"
var Version = "dev"
