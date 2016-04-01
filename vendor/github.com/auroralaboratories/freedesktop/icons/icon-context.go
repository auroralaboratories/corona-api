package icons

const (
	DEFAULT_ICON_CONTEXT_THRESHOLD     = 2
	DEFAULT_ICON_CONTEXT_SCALABLE_SIZE = 16
	DEFAULT_ICON_CONTEXT_SCALABLE_MIN  = 12
	DEFAULT_ICON_CONTEXT_SCALABLE_MAX  = 4096
)

type IconContextType int

const (
	IconContextThreshold IconContextType = 0
	IconContextFixed                     = 1
	IconContextScalable                  = 2
)

type IconContext struct {
	Subdirectory string
	Size         int
	Name         string
	Type         IconContextType
	MaxSize      int // Type=IconContextScalable
	MinSize      int // Type=IconContextScalable
	Threshold    int // Type=IconContextThreshold
}

func (self *IconContext) IsValid() bool {
	return (self.Size > 0)
}
