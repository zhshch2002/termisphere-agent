package report

import (
	"os"
	"strconv"
	"strings"
	"time"
)

func Disk(d time.Duration) ([]DiskReport, error) {
	ls, err := os.ReadDir("/sys/block")
	if err != nil {
		return nil, err
	}

	lst := SlicesMap(ls, func(i os.DirEntry) string {
		return i.Name()
	})

	var res []DiskReport
	for _, dev := range lst {
		if strings.HasPrefix(dev, "loop") || strings.HasPrefix(dev, "ram") {
			continue
		}

		var disk DiskReport
		disk.Name = dev

		if model, err := os.ReadFile("/sys/block/" + dev + "/device/model"); err == nil {
			disk.Model = strings.TrimSpace(string(model))
		}

		if vendor, err := os.ReadFile("/sys/block/" + dev + "/device/vendor"); err == nil {
			disk.Vendor = strings.TrimSpace(string(vendor))
		}

		if serial, err := os.ReadFile("/sys/block/" + dev + "/device/serial"); err == nil {
			disk.Serial = strings.TrimSpace(string(serial))
		}

		if size, err := os.ReadFile("/sys/block/" + dev + "/size"); err == nil {
			if value, err := strconv.ParseUint(strings.TrimSpace(string(size)), 10, 64); err == nil {
				disk.Size = value * 512
			}
		}

		if stat, err := os.ReadFile("/sys/block/" + dev + "/stat"); err == nil {
			if fields := strings.Fields(string(stat)); len(fields) >= 7 {
				if read, err := strconv.ParseUint(fields[2], 10, 64); err == nil {
					disk.Read = read * 512
				}
				if write, err := strconv.ParseUint(fields[6], 10, 64); err == nil {
					disk.Write = write * 512
				}
			}
		}

		res = append(res, disk)
	}

	time.Sleep(d)

	for idx := range res {
		if stat, err := os.ReadFile("/sys/block/" + res[idx].Name + "/stat"); err == nil {
			if fields := strings.Fields(string(stat)); len(fields) >= 7 {
				if read, err := strconv.ParseUint(fields[2], 10, 64); err == nil {
					res[idx].ReadSpeed = (read*512 - res[idx].Read) / uint64(d.Seconds())
					res[idx].Read = read * 512
				}
				if write, err := strconv.ParseUint(fields[6], 10, 64); err == nil {
					res[idx].WriteSpeed = (write*512 - res[idx].Write) / uint64(d.Seconds())
					res[idx].Write = write * 512
				}
			}
		}

	}

	return res, nil
}

func SlicesMap[T any, U any](s []T, f func(T) U) []U {
	res := make([]U, 0, len(s))
	for _, v := range s {
		res = append(res, f(v))
	}
	return res
}
