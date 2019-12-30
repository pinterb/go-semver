package version

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var (
	// version indicates which semantic version of the binary is running. It's
	// obviously more human-friendly than a git hash. Since this version could be
	// manually set, you could have different build artifacts with the same
	// semantic version. For this reason, use extreme CAUTION when basing any build
	// or deployment automation on this value. And in fact, you should consider
	// treating this as informational only.
	// This value is set by the build scripts.
	version = "dev"

	// commit indicates which git hash the binary was built off of.
	// It's probably the version you should track as part of your CI/CD pipeline.
	// This value is set by the build scripts.
	commit = "none"

	// date indicates when the application was built (i.e. compiled)
	date = "unknown"
)

func Subcommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf(`version: %s
git hash: %s
built on: %s
go version: %s
go compiler: %s
platform: %s/%s
`, version, commit, date, runtime.Version(), runtime.Compiler, runtime.GOOS, runtime.GOARCH)
		},
	}
}
