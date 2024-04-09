package report

import (
	"os"
	"strings"
	"syscall"
)

func Filesystem() ([]FilesystemReport, error) {
	mounts, err := os.ReadFile("/proc/mounts")
	if err != nil {
		return nil, err
	}

	var res []FilesystemReport
	for _, line := range strings.Split(string(mounts), "\n") {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		if !strings.HasPrefix(fields[0], "/") {
			continue
		}

		var fs FilesystemReport
		fs.Mount = fields[1]
		fs.Device = fields[0]
		if len(fields) >= 3 {
			fs.Type = fields[2]
		}

		var stat syscall.Statfs_t
		if err := syscall.Statfs(fs.Mount, &stat); err == nil {
			fs.Total = stat.Blocks * uint64(stat.Bsize)
			fs.Free = stat.Bfree * uint64(stat.Bsize)
		}

		res = append(res, fs)
	}

	return res, nil
}
