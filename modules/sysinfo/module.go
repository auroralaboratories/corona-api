package main

type SystemPlugin struct {
	BasePlugin
}

type SystemStats struct {
	CPU    SystemStatsCpu
	Memory SystemStatsMemory
}

func (self *SystemPlugin) Init() (err error) {
	return
}

func (self *SystemPlugin) GetAllStats() (stats SystemStats, err error) {
	stats = SystemStats{
		CPU:    self.GetCpuStats(),
		Memory: self.GetMemoryStats(),
	}

	return
}
