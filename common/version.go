package common

import (
	"fmt"
	"runtime"

	"github.com/urfave/cli/v2"
)

var AppVersion AppVersionInfo

var NAME = "fringe-runner"
var VERSION = "development version"
var REVISION = "HEAD"
var BRANCH = "HEAD"
var BUILT = "unknown"

type AppVersionInfo struct {
	Name         string `json:"name"`
	Version      string `json:"version"`
	Revision     string `json:"revision"`
	Branch       string `json:"branch"`
	GOVersion    string `json:"go_version"`
	BuiltAt      string `json:"built_at"`
	OS           string `json:"os"`
	Architecture string `json:"architecture"`
}

func (v *AppVersionInfo) Printer(c *cli.Context) {
	fmt.Print(v.Extended())
}

func (v *AppVersionInfo) ShortLine() string {
	return fmt.Sprintf("%s (%s)", v.Version, v.Revision)
}

func (v *AppVersionInfo) Extended() string {
	version := fmt.Sprintf("Version:      %s\n", v.Version)
	version += fmt.Sprintf("Git revision: %s\n", v.Revision)
	version += fmt.Sprintf("Git branch:   %s\n", v.Branch)
	version += fmt.Sprintf("GO version:   %s\n", v.GOVersion)
	version += fmt.Sprintf("Built:        %s\n", v.BuiltAt)
	version += fmt.Sprintf("OS/Arch:      %s/%s\n", v.OS, v.Architecture)

	return version
}

func DefaultUserAgent() string {
	return fmt.Sprintf("%s/%s", NAME, VERSION)
}

func init() {
	AppVersion = AppVersionInfo{
		Name:         NAME,
		Version:      VERSION,
		Revision:     REVISION,
		Branch:       BRANCH,
		GOVersion:    runtime.Version(),
		BuiltAt:      BUILT,
		OS:           runtime.GOOS,
		Architecture: runtime.GOARCH,
	}
}
