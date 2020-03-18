package subsystems

type MemorySubsystem struct{}

func (s *MemorySubsystem) Name() string {
	return "memory"
}

func (s *MemorySubsystem) Set(cgroupPath string, res *ResourceConfig) error {
	return cgroupSet(s.Name(), cgroupPath, "memory.limit_in_bytes", res.MemoryLimit)
}

func (s *MemorySubsystem) Apply(cgroupPath string, pid int) error {
	return cgroupApply(s.Name(), cgroupPath, pid)
}

func (s *MemorySubsystem) Remove(cgroupPath string) error {
	return cgroupRemove(s.Name(), cgroupPath)
}
