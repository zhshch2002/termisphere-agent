package report

import (
	"sync"
	"time"
)

func Fetch(d time.Duration, req Request) (res Report, err error) {
	var wg sync.WaitGroup

	if req.Arch {
		res.Arch, err = Arch()
		if err != nil {
			return Report{}, err
		}
	}

	if req.Platform {
		res.Platform, err = Platform()
		if err != nil {
			return Report{}, err
		}
	}

	if req.Hostname {
		res.Hostname, err = Hostname()
		if err != nil {
			return Report{}, err
		}
	}

	if req.CPU {
		wg.Add(1)
		go func() {
			defer wg.Done()
			res.CPU, err = CPU(d)
		}()
	}

	if req.Memory {
		res.Memory, err = Memory()
		if err != nil {
			return Report{}, err
		}
	}

	if req.BlockDevice {
		wg.Add(1)
		go func() {
			defer wg.Done()
			res.BlockDevice, err = BlockDevice(d)
		}()
	}

	if req.Filesystem {
		res.Filesystem, err = Filesystem()
		if err != nil {
			return Report{}, err
		}
	}

	if req.Network {
		wg.Add(1)
		go func() {
			defer wg.Done()
			res.Network, err = Network(d)
		}()
	}

	if wg.Wait(); err != nil {
		return Report{}, err
	}

	return res, nil
}

type Request struct {
	Arch        bool
	Platform    bool
	Hostname    bool
	CPU         bool
	Memory      bool
	BlockDevice bool
	Filesystem  bool
	Network     bool
}

type Report struct {
	Arch        string              `json:"arch,omitempty"`
	Platform    string              `json:"platform,omitempty"`
	Hostname    string              `json:"hostname,omitempty"`
	CPU         []CpuReport         `json:"cpu,omitempty"`
	Memory      *MemoryReport       `json:"memory,omitempty"`
	BlockDevice []BlockDeviceReport `json:"block_device,omitempty"`
	Filesystem  []FilesystemReport  `json:"filesystem,omitempty"`
	Network     []NetworkReport     `json:"network,omitempty"`
}

type CpuReport struct {
	US uint64 `json:"us"`
	SY uint64 `json:"sy"`
	NI uint64 `json:"ni"`
	ID uint64 `json:"id"`
	WA uint64 `json:"wa"`
	HI uint64 `json:"hi"`
	SI uint64 `json:"si"`
	ST uint64 `json:"st"`
}

type MemoryReport struct {
	Total     uint64 `json:"total"`
	Free      uint64 `json:"free"`
	Available uint64 `json:"available"`
	Buffer    uint64 `json:"buffer"`
	Cache     uint64 `json:"cache"`
	SwapTotal uint64 `json:"swap_total"`
	SwapFree  uint64 `json:"swap_free"`
	SwapCache uint64 `json:"swap_cache"`
}

type BlockDeviceReport struct {
	Name       string `json:"name"`
	Model      string `json:"model"`
	Vendor     string `json:"vendor"`
	Serial     string `json:"serial"`
	Size       uint64 `json:"size"`
	Read       uint64 `json:"read"`
	ReadSpeed  uint64 `json:"read_speed"`
	Write      uint64 `json:"write"`
	WriteSpeed uint64 `json:"write_speed"`
}

type FilesystemReport struct {
	Mount  string `json:"mount"`
	Device string `json:"device"`
	Type   string `json:"type"`
	Total  uint64 `json:"total"`
	Free   uint64 `json:"free"`
}

type NetworkReport struct {
	Name    string   `json:"name"`
	Speed   uint64   `json:"speed"`
	MAC     string   `json:"mac"`
	MTU     uint64   `json:"mtu"`
	IP      []string `json:"ip"`
	RX      uint64   `json:"rx"`
	RXSpeed uint64   `json:"rx_speed"`
	TX      uint64   `json:"tx"`
	TXSpeed uint64   `json:"tx_speed"`
}
