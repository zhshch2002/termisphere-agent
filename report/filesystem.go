//go:build darwin || freebsd || windows

package report

import "errors"

func Filesystem() ([]FilesystemReport, error) {
	return nil, errors.New("platform not supported")
}
