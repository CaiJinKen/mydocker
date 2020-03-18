package subsystems

type CpusetSubsystem struct{}

func (s *CpusetSubsystem) Name() string {
	return "cpuset"
}

func (s *CpusetSubsystem) Set(cgroupPath string, res *ResourceConfig) error {
	return cgroupSet(s.Name(), cgroupPath, "cpuset.cpus", res.CpuSet)
}

func (s *CpusetSubsystem) Apply(cgroupPath string, pid int) error {
	return cgroupApply(s.Name(), cgroupPath, pid)
}

func (s *CpusetSubsystem) Remove(cgroupPath string) error {
	return cgroupRemove(s.Name(), cgroupPath)
}
