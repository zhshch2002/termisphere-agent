package report

import (
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

func Network(d time.Duration) ([]NetworkReport, error) {
	ls, err := os.ReadDir("/sys/class/net")
	if err != nil {
		return nil, err
	}

	lst := SlicesMap(ls, func(i os.DirEntry) string {
		return i.Name()
	})

	var res []NetworkReport
	for _, dev := range lst {
		if dev == "lo" {
			continue
		}

		var network NetworkReport
		network.Name = dev
		network.IP = make([]string, 0)

		if speed, err := os.ReadFile("/sys/class/net/" + dev + "/speed"); err == nil {
			if value, err := strconv.ParseUint(strings.TrimSpace(string(speed)), 10, 64); err == nil {
				network.Speed = value * 1024 * 1024
			}
		}

		if mac, err := os.ReadFile("/sys/class/net/" + dev + "/address"); err == nil {
			network.MAC = strings.TrimSpace(string(mac))
		}

		if mtu, err := os.ReadFile("/sys/class/net/" + dev + "/mtu"); err == nil {
			if value, err := strconv.ParseUint(strings.TrimSpace(string(mtu)), 10, 64); err == nil {
				network.MTU = value
			}
		}

		if iface, err := net.InterfaceByName(dev); err == nil {
			if addrs, err := iface.Addrs(); err == nil {
				for _, addr := range addrs {
					network.IP = append(network.IP, addr.String())
				}
			}
		}

		if rx, err := os.ReadFile("/sys/class/net/" + dev + "/statistics/rx_bytes"); err == nil {
			if value, err := strconv.ParseUint(strings.TrimSpace(string(rx)), 10, 64); err == nil {
				network.RX = value
			}
		}

		if stat, err := os.ReadFile("/sys/class/net/" + dev + "/statistics/tx_bytes"); err == nil {
			if value, err := strconv.ParseUint(strings.TrimSpace(string(stat)), 10, 64); err == nil {
				network.TX = value
			}
		}

		res = append(res, network)
	}

	time.Sleep(d)

	for idx := range res {
		if rx, err := os.ReadFile("/sys/class/net/" + res[idx].Name + "/statistics/rx_bytes"); err == nil {
			if value, err := strconv.ParseUint(strings.TrimSpace(string(rx)), 10, 64); err == nil {
				res[idx].RXSpeed = (value - res[idx].RX) / uint64(d.Seconds())
				res[idx].RX = value
			}
		}

		if stat, err := os.ReadFile("/sys/class/net/" + res[idx].Name + "/statistics/tx_bytes"); err == nil {
			if value, err := strconv.ParseUint(strings.TrimSpace(string(stat)), 10, 64); err == nil {
				res[idx].TXSpeed = (value - res[idx].TX) / uint64(d.Seconds())
				res[idx].TX = value
			}
		}
	}

	return res, nil
}
