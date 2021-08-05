package zname

import "fmt"

// Version metadata set by ldflags during the build.
var (
	version string
	commit  string
	date    string
)

// Version returns a string with version metadata: version number, git sha and build date.
// It returns "development" if version variables are not set during the build.
func Version() string {
	if version == "" {
		return "development"
	}

	if len(commit) > 6 {
		commit = commit[:6]
	}

	return fmt.Sprintf("%s - revision %s built at %s", version, commit, date)
}
