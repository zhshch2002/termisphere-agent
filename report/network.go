//go:build darwin || freebsd || windows

package report

import (
	"errors"
	"time"
)

func Network(d time.Duration) ([]NetworkReport, error) {
	return nil, errors.New("platform not supported")
}
