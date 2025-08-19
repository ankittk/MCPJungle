package cmd

import (
	"fmt"
	"runtime/debug"

	"github.com/spf13/cobra"
)

var Version = "dev"

// GetVersion returns the CLI version. It prefers a build-time injected value,
// then falls back to the Go module build info when installed via
// `go install module@version`, and defaults to "dev" for local builds.
func GetVersion() string {
	if Version != "" && Version != "dev" {
		return Version
	}
	if info, ok := debug.ReadBuildInfo(); ok {
		if info.Main.Version != "" && info.Main.Version != "(devel)" {
			return info.Main.Version
		}
	}
	return "dev"
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Get version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("MCPJungle Version %s\n", GetVersion())
	},
	Annotations: map[string]string{
		"group": string(subCommandGroupBasic),
		"order": "7",
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
