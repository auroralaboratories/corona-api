package types

type IOutput interface {
	Initialize(IBackend) error
	GetName() string
	SetName(string)
	GetProperty(string, interface{}) interface{}
	SetProperty(string, interface{})
	Mute() error
	Unmute() error
	ToggleMute() error
	GetVolume() (float64, error)
	SetVolume(float64) error
	IncreaseVolume(float64) error
	DecreaseVolume(float64) error
}
