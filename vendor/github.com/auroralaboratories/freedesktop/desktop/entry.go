package desktop

import (
	"os/exec"

	"github.com/mattn/go-shellwords"
)

type Entry struct {
	Key             string   `json:"key"`
	Type            string   `json:"type"`
	Version         string   `json:"version,omitempty"`
	Name            string   `json:"name"`
	GenericName     string   `json:"generic_name,omitempty"`
	Comment         string   `json:"comment,omitempty"`
	Icon            string   `json:"icon,omitempty"`
	TryExec         string   `json:"tryexec,omitempty"`
	Exec            string   `json:"exec,omitempty"`
	Path            string   `json:"path,omitempty"`
	StartupWMClass  string   `json:"startup_wm_class,omitempty"`
	URL             string   `json:"url,omitempty"`
	DesktopfilePath string   `json:"desktop_file_path,omitempty"`
	NoDisplay       bool     `json:"no_display,omitempty"`
	StartupNotify   bool     `json:"startup_notify,omitempty"`
	Hidden          bool     `json:"hidden,omitempty"`
	DBusActivatable bool     `json:"dbus_activateable,omitempty"`
	Terminal        bool     `json:"terminal,omitempty"`
	OnlyShowIn      []string `json:"only_show_in,omitempty"`
	NotShowIn       []string `json:"not_shown_in,omitempty"`
	Actions         []string `json:"actions,omitempty"`
	MimeType        []string `json:"mimetypes,omitempty"`
	Categories      []string `json:"categories,omitempty"`
	Implements      []string `json:"implements,omitempty"`
	Keywords        []string `json:"keywords,omitempty"`
}

func (self *Entry) Launch() error {
	parser := shellwords.NewParser()
	parser.ParseEnv = true

	if args, err := parser.Parse(self.Exec); err == nil {
		exe := args[0]
		args := args[1:len(args)]

		cmd := exec.Command(exe, args...)

		//  TODO: launch this with SysProcAttr.Setpgid
		return cmd.Start()
	} else {
		return err
	}
}
