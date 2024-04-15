//go:build darwin || freebsd || windows

package report

import (
	"errors"
	"time"
)

func CPU(d time.Duration) ([]CpuReport, error) {
	return nil, errors.New("platform not supported")
}
