package types

type IBackend interface {
	Initialize() error
	Reset()
	Refresh() error
	GetName() string
	SetName(string)
	GetProperty(string, interface{}) interface{}
	SetProperty(string, interface{})
	AddOutput(output IOutput) error
	GetOutputs() []IOutput
	GetOutputByName(string) (IOutput, bool)
	GetOutputsByProperty(string, interface{}) []IOutput
	GetCurrentOutput() (IOutput, error)
	SetCurrentOutput(int) error
}
