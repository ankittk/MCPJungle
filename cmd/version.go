package cmd

import (
	"fmt"
	"runtime/debug"

	"github.com/spf13/cobra"
)

var (
	// Version can be overridden at build time using:
	// go build -ldflags="-X 'github.com/mcpjungle/mcpjungle/cmd.Version=v1.2.3'"
	Version = "dev"
)

// GetVersion returns the CLI version string.
//
// Priority order:
// 1. Build-time injected `Version` (via ldflags)
// 2. Go module build info (when installed with `go install module@version`)
// 3. Default: "dev" (local build)
func GetVersion() string {
	if Version != "" && Version != "dev" {
		return normalizeVersion(Version)
	}

	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" && info.Main.Version != "(devel)" {
		return normalizeVersion(info.Main.Version)
	}

	return "dev"
}

// normalizeVersion ensures a consistent version format:
// - If version starts with a digit (e.g., "1.2.3"), prefix with 'v' → "v1.2.3"
// - Leave values starting with 'v' or non-semver strings untouched
func normalizeVersion(v string) string {
	if v == "" {
		return v
	}
	if v[0] >= '0' && v[0] <= '9' {
		return "v" + v
	}
	return v
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("MCPJungle %s\n", GetVersion())
	},
	Annotations: map[string]string{
		"group": string(subCommandGroupBasic),
		"order": "7",
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
