package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/pinterb/go-semver/cmd/semver/version"
	"github.com/pinterb/go-semver/pkg/crlf"
	"github.com/pinterb/go-semver/pkg/git"
	"github.com/pinterb/go-semver/pkg/semver"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var (
	gdir       string
	incr       string
	preid      string
	defv       string
	latestOnly bool

	incrdesc = fmt.Sprintf("Increment a valid version by the specified level. Level can %sbe one of: major, minor, patch, premajor, preminor, prepatch, %sor prerelease. If more than one version is provided, then %sthe most current version is incremented.", crlf.Linebreak, crlf.Linebreak, crlf.Linebreak)
	predesc  = fmt.Sprintf("Identifier to be used to prefix premajor, preminor, %sprepatch or prerelease version increments.", crlf.Linebreak)
)

func commandRoot() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "semver",
		Short: "A tool to manage semantic versioning of software",
		Long: `
The primary goal of semver is to make working with semantic versions
easier. Currently, its two primary functions are to a) validate lists
of raw versions; and b) increment a version.

A secondary goal is work well with git repositories. So while you may
pass one or more versions (as arguments) to semver, you can just as
easily use the tags from a local git repository. So semver can validate
git repository tags and perhaps most importantly, it can help manage
git tags by providing a clean interface for incrementing a current tag
to a valid next version.
`,
		Example: "semver --help",
		Run: func(cmd *cobra.Command, args []string) {
			if err := handleVersions(cmd, args); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(2)
			}
		},
		Args: func(cmd *cobra.Command, args []string) error {
			if err := validArgs(cmd, args); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(2)
			}
			return nil
		},
	}

	path, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	// add flags
	rootCmd.Flags().SortFlags = false

	rootCmd.Flags().StringVarP(&incr, "increment", "i", "", incrdesc)
	rootCmd.Flag("increment").NoOptDefVal = "patch"

	rootCmd.Flags().StringVar(&preid, "preid", "", predesc)

	rootCmd.Flags().StringVarP(&gdir, "repo-dir", "r", "", "Use tags from a local git repo as source of versions.")
	rootCmd.Flag("repo-dir").NoOptDefVal = path

	rootCmd.Flags().StringVarP(&defv, "default", "d", "", "Default version to use when no valid versions are provided")
	rootCmd.Flag("default").NoOptDefVal = "0.0.0"

	rootCmd.Flags().BoolVarP(&latestOnly, "latest-only", "l", false, "Only return the latest version")

	rootCmd.Flags().BoolP("help", "h", false, "Help for semver")

	// add subcommands
	rootCmd.AddCommand(version.Subcommand())

	return rootCmd
}

func main() {
	root := commandRoot()
	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(2)
	}

	// docs(root)
}

func validArgs(cmd *cobra.Command, args []string) error {
	if len(args) < 1 && gdir == "" {
		return errors.New("at least one version needs to be provided")
	}

	if len(args) > 1 && gdir != "" {
		return errors.New("versions are not allowed when specifying a git repository")
	}

	return nil
}

func handleVersions(cmd *cobra.Command, args []string) error {
	var v2 []string
	var err error
	if defv != "" {
		v2 = []string{defv}
	} else {
		v2 = []string{}
	}

	// use either passed in versions (i.e. args) or tags from git repo
	v2 = append(v2, args...)
	if gdir != "" {
		v, err := git.Tags(gdir)
		if err != nil {
			return err
		}
		v2 = append(v2, v...)
	}

	// get sorted list of valid versions
	valid, err := semver.SortedList(v2)
	if err != nil {
		return err
	}

	if len(valid) > 0 {
		if incr == "" {
			var fv string = strings.Join(valid, " ")
			if latestOnly {
				fv = valid[len(valid)-1]
			}
			fmt.Println(fv)
		} else {
			rt, err := semver.ToReleaseType(incr)
			if err != nil {
				return err
			}

			// increment current version
			nv, err := semver.Increment(valid[len(valid)-1], rt, preid)
			if err != nil {
				return err
			}
			fmt.Println(nv)
		}
	}

	return nil
}

// ONLY FOR DEVELOPMENT!
func docs(cmd *cobra.Command) {
	out := new(bytes.Buffer)
	err := doc.GenMarkdown(cmd, out)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(2)
	}

	path, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	path = fmt.Sprintf("%s/semver.md", path)
	err = ioutil.WriteFile(path, out.Bytes(), 0666)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(2)
	}
}
