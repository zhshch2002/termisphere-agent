//go:build darwin || freebsd || windows

package report

import (
	"errors"
	"time"
)

func BlockDevice(d time.Duration) ([]BlockDeviceReport, error) {
	return nil, errors.New("platform not supported")
}
