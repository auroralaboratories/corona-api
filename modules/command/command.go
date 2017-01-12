package command

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/mattn/go-shellwords"
)

const (
	DEFAULT_SHELL_WRAPPER = "bash -c '%s'"
)

type Command struct {
	Key         string `json:"key"`
	ShellWrap   string `json:"shellwrap"`
	CommandLine string `json:"command"`
	Detach      bool   `json:"detach,omitempty"`

	cmd  string
	args []string
}

type CommandResult struct {
	Status   int      `json:"status"`
	Error    bool     `json:"error"`
	Output   []string `json:"output"`
	Detached bool     `json:"detached"`
}

func (self *Command) Init() error {
	if self.CommandLine == `` {
		return fmt.Errorf("Cannot initialize command with an empty command line string")
	}

	if self.ShellWrap == `` {
		self.ShellWrap = DEFAULT_SHELL_WRAPPER
	}

	parser := shellwords.NewParser()
	parser.ParseEnv = true

	fullCmdLine := fmt.Sprintf(self.ShellWrap, self.CommandLine)

	if args, err := parser.Parse(fullCmdLine); err == nil {
		self.cmd = args[0]
		self.args = args[1:]
	} else {
		return err
	}

	return nil
}

func (self *Command) Execute(arguments ...string) (CommandResult, error) {
	command := exec.Command(self.cmd, self.args...)
	output := make([]byte, 0)
	outputLines := make([]string, 0)

	var err error

	if self.Detach {
		// self.Program.SysProcAttr.Setpgid = true
		err = command.Start()
	} else {
		output, err = command.CombinedOutput()
	}

	if len(output) > 0 {
		sep := "\n"

		for _, line := range bytes.Split(output, []byte(sep[:])) {
			if str := strings.TrimSpace(string(line[:])); str != `` {
				outputLines = append(outputLines, str)
			}
		}
	}

	return CommandResult{
		Error:    (err != nil),
		Output:   outputLines,
		Detached: self.Detach,
	}, err
}
