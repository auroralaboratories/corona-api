package pulseaudio

import (
    "fmt"

    "github.com/auroralaboratories/corona-api/modules/soundctl/backends"
    "github.com/auroralaboratories/corona-api/modules/soundctl/types"
    "github.com/auroralaboratories/pulse"
)

type Backend struct {
    backends.BaseBackend

    client *pulse.Client
    info   pulse.ServerInfo
}

func New() *Backend {
    rv := &Backend{
        BaseBackend: backends.BaseBackend{},
    }

    return rv
}


func (self *Backend) Refresh() error {
    if self.client == nil {
        if client, err := pulse.NewClient(`corona-api`); err == nil {
            self.client = client
        }else{
            return err
        }
    }

    if info, err := self.client.GetServerInfo(); err == nil {
        self.info = info
    }else{
        return err
    }

    if err := self.loadSinks(); err != nil {
        return err
    }

    return nil
}


func (self *Backend) GetCurrentOutput() (types.IOutput, error) {
    if defaultSink := self.info.DefaultSinkName; defaultSink != `` {
        if output, ok := self.GetOutputByName(defaultSink); ok {
            return output, nil
        }
    }

    return &Output{}, fmt.Errorf("No default output is currently selected")
}


func (self *Backend) SetCurrentOutput(index int) error {
    return fmt.Errorf("Not implemented")
}


func (self *Backend) loadSinks() error {
    if sinks, err := self.client.GetSinks(); err == nil {
        for _, sink := range sinks {
            newOutput := &Output{
                sink: &sink,
            }

            if err := newOutput.Initialize(self); err != nil {
                return err
            }

            newOutput.SetName(sink.Name)

            newOutput.SetProperty(`index`,       fmt.Sprintf("%d", sink.Index))
            newOutput.SetProperty(`volume`,      fmt.Sprintf("%f", sink.VolumeFactor))
            newOutput.SetProperty(`channels`,    fmt.Sprintf("%d", sink.Channels))
            newOutput.SetProperty(`description`, sink.Description)
            newOutput.SetProperty(`muted`,       fmt.Sprintf("%s", sink.Muted))

            if err := self.AddOutput(newOutput); err != nil {
                return err
            }
        }
    }else{
        return err
    }

    return nil
}