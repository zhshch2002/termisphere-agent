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
	return "linux", nil
}

func Hostname() (string, error) {
	return os.Hostname()
}
