// +build linux

package linux

import (
	"github.com/shirou/gopsutil/structs"

	"syscall"
)

func VirtualMemory() (*structs.VirtualMemoryStat, error) {
	sysinfo := &syscall.Sysinfo_t{}

	if err := syscall.Sysinfo(sysinfo); err != nil {
		return nil, err
	}

	ret := &structs.VirtualMemoryStat{
		Total:   uint64(sysinfo.Totalram),
		Free:    uint64(sysinfo.Freeram),
		Shared:  uint64(sysinfo.Sharedram),
		Buffers: uint64(sysinfo.Bufferram),
	}

	// TODO: platform independent
	ret.Available = ret.Free + ret.Buffers + ret.Cached

	ret.Used = ret.Total - ret.Free
	ret.UsedPercent = float64(ret.Total-ret.Available) / float64(ret.Total) * 100.0

	/*
		kern := buffers + cached
		ret.ActualFree = ret.Free + kern
		ret.ActualUsed = ret.Used - kern
	*/

	return ret, nil
}

func SwapMemory() (*structs.SwapMemoryStat, error) {
	sysinfo := &syscall.Sysinfo_t{}

	if err := syscall.Sysinfo(sysinfo); err != nil {
		return nil, err
	}
	ret := &structs.SwapMemoryStat{
		Total: uint64(sysinfo.Totalswap),
		Free:  uint64(sysinfo.Freeswap),
	}
	ret.Used = ret.Total - ret.Free
	ret.UsedPercent = float64(ret.Total-ret.Free) / float64(ret.Total) * 100.0

	return ret, nil
}
