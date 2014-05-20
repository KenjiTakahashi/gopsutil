package test

import (
	"github.com/shirou/gopsutil"
	"github.com/shirou/gopsutil/structs"
	"fmt"
	"testing"
)

func TestCpu_times(t *testing.T) {
	v, err := gopsutil.CPUTimes(false)
	if err != nil {
		t.Errorf("error %v", err)
	}
	if len(v) == 0 {
		t.Errorf("could not get CPUs ", err)
	}
	empty := structs.CPUTimesStat{}
	for _, vv := range v {
		if vv == empty {
			t.Errorf("could not get CPU User: %v", vv)
		}
	}
}

func TestCpu_counts(t *testing.T) {
	v, err := gopsutil.CPUCounts(true)
	if err != nil {
		t.Errorf("error %v", err)
	}
	if v == 0 {
		t.Errorf("could not get CPU counts: %v", v)
	}
}

func TestCPUTimeStat_String(t *testing.T) {
	v := structs.CPUTimesStat{
		CPU:    "cpu0",
		User:   100.1,
		System: 200.1,
		Idle:   300.1,
	}
	e := `{"cpu":"cpu0","user":100.1,"system":200.1,"idle":300.1,"nice":0,"iowait":0,"irq":0,"softirq":0,"steal":0,"guest":0,"guest_nice":0,"stolen":0}`
	if e != fmt.Sprintf("%v", v) {
		t.Errorf("CPUTimesStat string is invalid: %v", v)
	}
}

func TestCpuInfo(t *testing.T) {
	v, err := gopsutil.CPUInfo()
	if err != nil {
		t.Errorf("error %v", err)
	}
	for _, vv := range v {
		if vv.ModelName == "" {
			t.Errorf("could not get CPU Info: %v", vv)
		}
	}
}
