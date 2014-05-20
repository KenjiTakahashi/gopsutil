package structs

import (
	"encoding/json"
	"os"
)


type Process os.Process


type ProcessInterface interface {
}

type OpenFilesStat struct {
	Path string `json:"path"`
	Fd   uint64 `json:"fd"`
}

type MemoryInfoStat struct {
	RSS uint64 `json:"rss"` // bytes
	VMS uint64 `json:"vms"` // bytes
}

type RlimitStat struct {
	Resource int32 `json:"resource"`
	Soft     int32 `json:"soft"`
	Hard     int32 `json:"hard"`
}

type IOCountersStat struct {
	ReadCount  int32 `json:"read_count"`
	WriteCount int32 `json:"write_count"`
	ReadBytes  int32 `json:"read_bytes"`
	WriteBytes int32 `json:"write_bytes"`
}

type NumCtxSwitchesStat struct {
	Voluntary   int32 `json:"voluntary"`
	Involuntary int32 `json:"involuntary"`
}

/*
func (p Process) String() string {
	s, _ := json.Marshal(p)
	return string(s)
}
*/

func (o OpenFilesStat) String() string {
	s, _ := json.Marshal(o)
	return string(s)
}

func (m MemoryInfoStat) String() string {
	s, _ := json.Marshal(m)
	return string(s)
}

func (r RlimitStat) String() string {
	s, _ := json.Marshal(r)
	return string(s)
}

func (i IOCountersStat) String() string {
	s, _ := json.Marshal(i)
	return string(s)
}

func (p NumCtxSwitchesStat) String() string {
	s, _ := json.Marshal(p)
	return string(s)
}
