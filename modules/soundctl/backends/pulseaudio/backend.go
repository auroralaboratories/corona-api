package pulseaudio

import (
    "bufio"
    "bytes"
    "fmt"
    "os/exec"
    "regexp"
    "strings"

    "github.com/shutterstock/go-stockutil/stringutil"
    "github.com/auroralaboratories/corona-api/modules/soundctl/backends"
    "github.com/auroralaboratories/corona-api/modules/soundctl/types"
)

type Backend struct {
    backends.BaseBackend
}

func New() *Backend {
    rv := &Backend{
        BaseBackend: backends.BaseBackend{},
    }

    return rv
}


func (self *Backend) Refresh() error {
    if err := self.loadInfo(); err != nil {
        return err
    }

    if err := self.loadSinks(); err != nil {
        return err
    }

    return nil
}


func (self *Backend) GetCurrentOutput() (types.IOutput, error) {
    if defaultSink := self.GetProperty(`default_sink`, ``); defaultSink != `` {
        if outputs := self.GetOutputsByProperty(`name`, defaultSink); len(outputs) == 1 {
            return outputs[0], nil
        }
    }

    return &Output{}, fmt.Errorf("No default output is currently selected")
}


func (self *Backend) SetCurrentOutput(index int) error {
    return fmt.Errorf("Not implemented")
}


func (self *Backend) loadInfo() error {
    if output, err := exec.Command(`pactl`, `info`).Output(); err == nil {
        scanner := bufio.NewScanner(bytes.NewReader(output))

        for scanner.Scan() {
            line := scanner.Text()

            parts := strings.SplitN(line, `: `, 2)

            if len(parts) == 2 {
                key        := strings.ToLower(strings.Replace(parts[0], ` `, `_`, -1))
                value      := strings.TrimSpace(parts[1])

                self.SetProperty(key, value)
            }
        }
    }else{
        return err
    }

    return nil
}


func (self *Backend) loadSinks() error {
    if output, err := exec.Command(`pactl`, `list`, `sinks`).Output(); err == nil {
        scanner := bufio.NewScanner(bytes.NewReader(output))

        var newOutput *Output

        for scanner.Scan() {
            line := scanner.Text()

            if rx, err := regexp.Compile(`^Sink #(\d+)$`); err == nil && rx.MatchString(line) {
                if newOutput != nil {
                    if err := self.AddOutput(newOutput); err != nil {
                        return err
                    }
                }

                newOutput = &Output{
                    sinkIndex: -1,
                }

                if err := newOutput.Initialize(self); err != nil {
                    return err
                }

                if matches := rx.FindStringSubmatch(line); len(matches) == 2 {
                    if v, err := stringutil.ConvertToInteger(matches[1]); err == nil {
                        newOutput.sinkIndex = int(v)
                        newOutput.SetProperty(`index`, matches[1])
                        newOutput.SetName(fmt.Sprintf("sink-%s", matches[1]))
                    }
                }
            }else if rx, err := regexp.Compile("^\t([^\\:\t]+):\\s+(.*)$"); err == nil && rx.MatchString(line) && newOutput.sinkIndex >= 0 {
                if matches := rx.FindStringSubmatch(line); len(matches) == 3 {
                    key        := strings.ToLower(strings.Replace(matches[1], ` `, `_`, -1))
                    value      := strings.TrimSpace(matches[2])

                    newOutput.SetProperty(key, value)
                }
            }

        }

        if newOutput != nil {
            if err := self.AddOutput(newOutput); err != nil {
                return err
            }
        }
    }else{
        return err
    }

    return nil
}