//go:build darwin || freebsd || windows

package report

import (
	"os"
	"runtime"
	"strings"
)

func Arch() (string, error) {
	return strings.ToLower(runtime.GOARCH), nil
}

func Platform() (string, error) {
	return "darwin", nil
}

func Hostname() (string, error) {
	return os.Hostname()
}
