package pulseaudio

import (
    "fmt"
    "os/exec"
    "github.com/shutterstock/go-stockutil/stringutil"
    "github.com/auroralaboratories/corona-api/modules/soundctl/backends"

    log "github.com/Sirupsen/logrus"
    "strings"
)

type Output struct {
    backends.BaseOutput

    sinkIndex int
}

func (self *Output) Mute() error {
    return self.runPactlSubcommand(`set-sink-mute`, self.sinkIndex, 1)
}


func (self *Output) Unmute() error {
    return self.runPactlSubcommand(`set-sink-mute`, self.sinkIndex, 0)
}


func (self *Output) ToggleMute() error {
    return self.runPactlSubcommand(`set-sink-mute`, self.sinkIndex, `toggle`)
}


func (self *Output) GetVolume() (float64, error) {
    return -1.0, fmt.Errorf("Not implemented")
}


func (self *Output) SetVolume(factor float64) error {
    return self.runPactlSubcommand(`set-sink-volume`, self.sinkIndex, fmt.Sprintf("%d%%", int(factor * 100.0)))
}


func (self *Output) IncreaseVolume(factor float64) error {
    if volume, err := self.GetVolume(); err == nil {
        return self.SetVolume(volume + factor)
    }else{
        return err
    }
}


func (self *Output) DecreaseVolume(factor float64) error {
    if volume, err := self.GetVolume(); err == nil {
        return self.SetVolume(volume - factor)
    }else{
        return err
    }
}


func (self *Output) runPactlSubcommand(subcommand string, args ...interface{}) error {
    cmdArgs := []string{ subcommand }

    for _, arg := range args {
        if v, err := stringutil.ToString(arg); err == nil {
            cmdArgs = append(cmdArgs, v)
        }else{
            return err
        }
    }

    command := exec.Command(`pactl`, cmdArgs...)

    log.Debugf("pulseaudio output cmd: %s", strings.Join(command.Args, ` `))

    return command.Run()
}