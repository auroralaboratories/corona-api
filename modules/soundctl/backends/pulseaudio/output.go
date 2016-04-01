package pulseaudio

import (
	"github.com/auroralaboratories/corona-api/modules/soundctl/backends"
	"github.com/auroralaboratories/pulse"
)

type Output struct {
	backends.BaseOutput

	sink *pulse.Sink
}

func (self *Output) Mute() error {
	return self.sink.Mute()
}

func (self *Output) Unmute() error {
	return self.sink.Unmute()
}

func (self *Output) ToggleMute() error {
	return self.sink.ToggleMute()
}

func (self *Output) GetVolume() (float64, error) {
	if err := self.sink.Refresh(); err == nil {
		return self.sink.VolumeFactor, nil
	} else {
		return -1.0, err
	}
}

func (self *Output) SetVolume(factor float64) error {
	return self.sink.SetVolume(factor)
}

func (self *Output) IncreaseVolume(factor float64) error {
	return self.sink.IncreaseVolume(factor)
}

func (self *Output) DecreaseVolume(factor float64) error {
	return self.sink.DecreaseVolume(factor)
}
