package version

import (
	"flag"
	"runtime"
)

var (
	Version   = "develop"
	GitCommit = ""
	BuildDate = ""
)

type BuildInfo struct {
	// Version is the current semver.
	Version string `json:"version,omitempty"`
	// GitCommit is the git sha1.
	GitCommit string `json:"gitCommit,omitempty"`
	// BuildDate is the build time.
	BuildDate string `json:"buildDate,omitempty"`
	// GoVersion is the version of the Go compiler used.
	GoVersion string `json:"goVersion,omitempty"`
}

func Get() BuildInfo {
	v := BuildInfo{
		Version:   Version,
		GitCommit: GitCommit,
		BuildDate: BuildDate,
		GoVersion: runtime.Version(),
	}

	// HACK(bacongobbler): strip out GoVersion during a test run for consistent test output
	if flag.Lookup("test.v") != nil {
		v.GoVersion = ""
	}
	return v
}
