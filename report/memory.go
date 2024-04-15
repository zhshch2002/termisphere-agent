//go:build darwin || freebsd || windows

package report

import "errors"

func Memory() (*MemoryReport, error) {
	return nil, errors.New("platform not supported")
}
