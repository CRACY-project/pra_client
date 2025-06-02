package sysinfo

import (
	"github.com/elastic/go-sysinfo"
	"github.com/elastic/go-sysinfo/types"
)

type SystemInfo struct {
	HostInfo   types.HostInfo
	MemoryInfo types.HostMemoryInfo
	Process    types.ProcessInfo
	CPUTime    types.CPUTimes
}

func GatherSystemInfo() (*SystemInfo, error) {
	host, err := sysinfo.Host()
	if err != nil {
		return nil, err
	}
	hostInfo := host.Info()
	memoryInfo, err := host.Memory()
	if err != nil {
		return nil, err
	}

	proc, err := sysinfo.Self()
	if err != nil {
		return nil, err
	}

	procInfo, err := proc.Info()
	if err != nil {
		return nil, err
	}

	cpuTime, err := host.CPUTime()
	if err != nil {
		return nil, err
	}

	return &SystemInfo{
		HostInfo:   hostInfo,
		MemoryInfo: *memoryInfo,
		Process:    procInfo,
		CPUTime:    cpuTime,
	}, nil
}
