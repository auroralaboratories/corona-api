package main

import (
    "github.com/cloudfoundry/gosigar"
)

type SystemStatsMemory struct {
    Total        uint64
    Used         uint64
    Free         uint64
    CachedUsed   uint64
    CachedFree   uint64
    PercentUsed  float32
    PercentCache float32
}

func (self *SystemPlugin) GetMemoryStats() (stats SystemStatsMemory) {
    stats     = SystemStatsMemory{}
    s_memory := sigar.Mem{}
    s_memory.Get()

    stats.Total        = s_memory.Total
    stats.Used         = s_memory.ActualUsed
    stats.Free         = s_memory.ActualFree
    stats.CachedUsed   = s_memory.Used
    stats.CachedFree   = s_memory.Free
    stats.PercentUsed  = (float32(stats.Used) / float32(stats.Total)) * 100.0
    stats.PercentCache = (float32(stats.CachedUsed) / float32(stats.Total)) * 100.0


    return
}