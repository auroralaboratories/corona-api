package main

import (
    // "strings"
    // "strconv"
    "github.com/cloudfoundry/gosigar"
)

type SystemStatsCpuCore struct {
    sigar.Cpu
    SocketNumber uint
    CoreNumber   uint
    Speed        uint
    Percent      float32
}

type SystemStatsCpu struct {
    Cores        []SystemStatsCpuCore
    Make         string
    Model        string
    Flags        []string
    Speed        uint
}

func (self *SystemPlugin) GetCpuStats() (stats SystemStatsCpu) {
    stats     =  SystemStatsCpu{}
    s_cpuList := sigar.CpuList{}
    s_cpuList.Get()

    for i, cpu := range s_cpuList.List {
        core            := SystemStatsCpuCore{}
        core.User        = cpu.User
        core.Nice        = cpu.Nice
        core.Sys         = cpu.Sys
        core.Idle        = cpu.Idle
        core.Wait        = cpu.Wait
        core.Irq         = cpu.Irq
        core.SoftIrq     = cpu.SoftIrq
        core.Stolen      = cpu.Stolen
        core.CoreNumber  = uint(i)

        stats.Cores  = append(stats.Cores, core)
    }

    // if _, err := os.Stat("/proc/cpuinfo"); err == nil {
    //     file, err := os.Open("/proc/cpuinfo")

    //     if err == nil {
    //         scanner := bufio.NewScanner(file)

    //         current_core := 0

    //         for scanner.Scan() {
    //             line  := scanner.Text()
    //             parts := strings.SplitN(line, ":", -1)
    //             key   := strings.ToLower(strings.TrimSpace(parts[0]))
    //             value := strings.TrimSpace(parts[1])

    //             switch key {
    //             case 'processor':
    //                 if i, err := strconv.Atoi(value); err == nil{
    //                     current_core = i
    //                 }

    //             case 'cpu mhz':

    //             }
    //         }
    //     }
    // }

    return
}