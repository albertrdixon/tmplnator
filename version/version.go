package version

import (
	"fmt"
)

const (
	CodeVersion = "v0.0.1"
)

var Build string

func RuntimeVersion(version string, build string) string {
	var vers string
	if build != "" && len(version) > 4 && version[len(version)-4:] == "-dev" {
		vers = fmt.Sprintf("%s-%s", version, build)
	} else {
		vers = version
	}
	return vers
}
