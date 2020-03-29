package subsystems

import "encoding/json"

type ResourceConfig struct {
	MemoryLimit string `json:"memory_limit"`
	CpuShare    string `json:"cpu_share"`
	CpuSet      string `json:"cpu_set"`
}

func (r *ResourceConfig) String() string {
	bts, err := json.Marshal(r)
	if err != nil {
		return ""
	}
	return string(bts)
}

type Subsystem interface {
	Name() string
	Set(path string, res *ResourceConfig) error
	Apply(path string, pid int) error
	Remove(path string) error
}

var SubsystemsIns = []Subsystem{
	&CpusetSubsystem{},
	&CpuSubsystem{},
	&MemorySubsystem{},
}
