package types

type IBackend interface {
    Initialize() error
    Refresh() error
    GetName() string
    SetName(string)
    GetProperty(string, string) string
    SetProperty(string, string)
    AddOutput(output IOutput) error
    GetOutputs() []IOutput
    GetOutputsByProperty(string, string) []IOutput
    GetCurrentOutput() (IOutput, error)
    SetCurrentOutput(int) error
}