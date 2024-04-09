package report

import (
	"os"
	"strconv"
	"strings"
)

func Memory() (*MemoryReport, error) {
	info, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return nil, err
	}

	var res MemoryReport
	for _, line := range strings.Split(string(info), "\n") {
		if strings.HasPrefix(line, "MemTotal") {
			res.Total = parseMem(line)
		} else if strings.HasPrefix(line, "MemFree") {
			res.Free = parseMem(line)
		} else if strings.HasPrefix(line, "MemAvailable") {
			res.Available = parseMem(line)
		} else if strings.HasPrefix(line, "Buffers") {
			res.Buffer = parseMem(line)
		} else if strings.HasPrefix(line, "Cached") {
			res.Cache = parseMem(line)
		} else if strings.HasPrefix(line, "SwapTotal") {
			res.SwapTotal = parseMem(line)
		} else if strings.HasPrefix(line, "SwapFree") {
			res.SwapFree = parseMem(line)
		} else if strings.HasPrefix(line, "SwapCached") {
			res.SwapCache = parseMem(line)
		}
	}

	return &res, nil
}

func parseMem(line string) uint64 {
	fields := strings.Fields(line)
	if len(fields) != 3 {
		return 0
	}

	v, err := strconv.ParseUint(fields[1], 10, 64)
	if err != nil {
		return 0
	}

	return v * 1024
}
