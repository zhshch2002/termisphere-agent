package report

import (
	"cmx-termisphere-agent/utils"
	"os"
	"strconv"
	"strings"
	"time"
)

func BlockDevice(d time.Duration) ([]BlockDeviceReport, error) {
	ls, err := os.ReadDir("/sys/block")
	if err != nil {
		return nil, err
	}

	lst := utils.SlicesMap(ls, func(i os.DirEntry) string {
		return i.Name()
	})

	var res []BlockDeviceReport
	for _, dev := range lst {
		if strings.HasPrefix(dev, "loop") || strings.HasPrefix(dev, "ram") {
			continue
		}

		var blk BlockDeviceReport
		blk.Name = dev

		if model, err := os.ReadFile("/sys/block/" + dev + "/device/model"); err == nil {
			blk.Model = strings.TrimSpace(string(model))
		}

		if vendor, err := os.ReadFile("/sys/block/" + dev + "/device/vendor"); err == nil {
			blk.Vendor = strings.TrimSpace(string(vendor))
		}

		if serial, err := os.ReadFile("/sys/block/" + dev + "/device/serial"); err == nil {
			blk.Serial = strings.TrimSpace(string(serial))
		}

		if size, err := os.ReadFile("/sys/block/" + dev + "/size"); err == nil {
			if value, err := strconv.ParseUint(strings.TrimSpace(string(size)), 10, 64); err == nil {
				blk.Size = value * 512
			}
		}

		if blk.Size == 0 {
			continue
		}

		if stat, err := os.ReadFile("/sys/block/" + dev + "/stat"); err == nil {
			if fields := strings.Fields(string(stat)); len(fields) >= 7 {
				if read, err := strconv.ParseUint(fields[2], 10, 64); err == nil {
					blk.Read = read * 512
				}
				if write, err := strconv.ParseUint(fields[6], 10, 64); err == nil {
					blk.Write = write * 512
				}
			}
		}

		res = append(res, blk)
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
