package cli

import (
	"cmp"
	"runtime/debug"

	"github.com/spf13/cobra"

	"github.com/cego/aikido-dojo/internal/catalog"
)

type versionInfo struct {
	Version string      `json:"version"`
	Commit  string      `json:"commit,omitempty"`
	Go      string      `json:"go,omitempty"`
	Spec    specVersion `json:"spec"`
}

type specVersion struct {
	UpdatedAt string `json:"updated_at"`
	SHA256    string `json:"sha256"`
}

func (a *app) versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print aikido-dojo's version and the API spec snapshot its commands come from",
		Long: "version prints aikido-dojo's version, the commit it was built from when known, the Go version, " +
			"and the date and SHA-256 of the vendored API spec its commands were generated from.",
		Args: exactArgs(nil),
		RunE: func(cmd *cobra.Command, _ []string) error {
			return a.printJSON(cmd.Context(), versionOf(debug.ReadBuildInfo()))
		},
	}
}

// versionOf reads what Go stamps into the binary: the module version for go
// install, and for a build in a git checkout its tag or a pseudo-version of
// its commit, so a release needs no version flag of its own.
func versionOf(bi *debug.BuildInfo, ok bool) versionInfo {
	v := versionInfo{Version: "(unknown)", Spec: specVersion{UpdatedAt: catalog.SpecUpdatedAt, SHA256: catalog.SpecSHA256}}
	if !ok {
		return v
	}
	v.Version, v.Go = cmp.Or(bi.Main.Version, v.Version), bi.GoVersion
	for _, s := range bi.Settings {
		if s.Key == "vcs.revision" {
			v.Commit = s.Value
		}
	}
	return v
}
