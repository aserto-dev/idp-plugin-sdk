package version

import (
	"runtime"
	"time"

	"github.com/aserto-dev/go-grpc/aserto/common/info/v1"
)

type GetVersionFunc func() (string, string, string)

// GetInfo gets version stamp information.
func GetBuildInfo(versionFunc GetVersionFunc) *info.BuildInfo {
	ver, date, commit := versionFunc()

	if ver == "" {
		ver = "0.0.0"
	}

	if date == "" {
		date = time.Now().UTC().Format(time.RFC3339)
	}

	if commit == "" {
		commit = "undefined"
	}

	return &info.BuildInfo{
		Version: ver,
		Date:    date,
		Commit:  commit,
		Os:      runtime.GOOS,
		Arch:    runtime.GOARCH,
	}
}
