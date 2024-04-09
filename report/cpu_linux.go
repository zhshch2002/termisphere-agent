package report

import (
	"fmt"
	"os"
	"strconv"
	"time"
)
import "strings"

func CPU(d time.Duration) ([]CpuReport, error) {
	stat, err := os.ReadFile("/proc/stat")
	if err != nil {
		return nil, err
	}

	var res []CpuReport
	for _, line := range strings.Split(string(stat), "\n") {
		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}
		if !strings.HasPrefix(parts[0], "cpu") {
			continue
		}
		if parts[0] == "cpu" {
			continue
		}

		cpuID, err := strconv.ParseInt(parts[0][3:], 10, 64)
		if err != nil {
			return nil, err
		}
		for len(res) <= int(cpuID) {
			res = append(res, CpuReport{})
		}

		res[cpuID], err = parseCPU(parts[1:])
		if err != nil {
			return nil, err
		}
	}

	time.Sleep(d)

	stat, err = os.ReadFile("/proc/stat")
	if err != nil {
		return nil, err
	}

	for _, line := range strings.Split(string(stat), "\n") {
		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}
		if !strings.HasPrefix(parts[0], "cpu") {
			continue
		}
		if parts[0] == "cpu" {
			continue
		}

		cpuID, err := strconv.ParseInt(parts[0][3:], 10, 64)
		if err != nil {
			return nil, err
		}

		var data CpuReport
		data, err = parseCPU(parts[1:])
		if err != nil {
			return nil, err
		}

		res[cpuID] = CpuReport{
			US: data.US - res[cpuID].US,
			SY: data.SY - res[cpuID].SY,
			NI: data.NI - res[cpuID].NI,
			ID: data.ID - res[cpuID].ID,
			WA: data.WA - res[cpuID].WA,
			HI: data.HI - res[cpuID].HI,
			SI: data.SI - res[cpuID].SI,
			ST: data.ST - res[cpuID].ST,
		}
	}

	return res, nil
}

func parseCPU(fields []string) (CpuReport, error) {
	if len(fields) < 8 {
		return CpuReport{}, fmt.Errorf("not enough fields in cpu stat")
	}

	us, err := strconv.ParseUint(fields[0], 10, 64)
	if err != nil {
		return CpuReport{}, err
	}

	sy, err := strconv.ParseUint(fields[1], 10, 64)
	if err != nil {
		return CpuReport{}, err
	}

	ni, err := strconv.ParseUint(fields[2], 10, 64)
	if err != nil {
		return CpuReport{}, err
	}

	id, err := strconv.ParseUint(fields[3], 10, 64)
	if err != nil {
		return CpuReport{}, err
	}

	wa, err := strconv.ParseUint(fields[4], 10, 64)
	if err != nil {
		return CpuReport{}, err
	}

	hi, err := strconv.ParseUint(fields[5], 10, 64)
	if err != nil {
		return CpuReport{}, err
	}

	si, err := strconv.ParseUint(fields[6], 10, 64)
	if err != nil {
		return CpuReport{}, err
	}

	st, err := strconv.ParseUint(fields[7], 10, 64)
	if err != nil {
		return CpuReport{}, err
	}

	return CpuReport{US: us, SY: sy, NI: ni, ID: id, WA: wa, HI: hi, SI: si, ST: st}, nil
}
