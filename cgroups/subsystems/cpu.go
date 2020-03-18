package subsystems

type CpuSubsystem struct{}

func (s *CpuSubsystem) Name() string {
	return "cpu"
}

func (s *CpuSubsystem) Set(cgroupPath string, res *ResourceConfig) error {
	return cgroupSet(s.Name(), cgroupPath, "cpu.shares", res.CpuShare)
}

func (s *CpuSubsystem) Apply(cgroupPath string, pid int) error {
	return cgroupApply(s.Name(), cgroupPath, pid)
}

func (s *CpuSubsystem) Remove(cgroupPath string) error {
	return cgroupRemove(s.Name(), cgroupPath)
}
