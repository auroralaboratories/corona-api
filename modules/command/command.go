package command

import (
    "bytes"
    "fmt"
    "os/exec"
    "strings"

    "github.com/mattn/go-shellwords"
)

type Command struct {
    Key         string     `json:"key"`
    Shell       string     `json:"shell"`
    CommandLine string     `json:"command"`
    Program     *exec.Cmd  `json:"-"`
}

type CommandResult struct {
    Status  int        `json:"status"`
    Error   bool       `json:"error"`
    Output  []string   `json:"output"`
}

func (self *Command) Init() error {
    if self.CommandLine == `` {
        return fmt.Errorf("Cannot initialize command with an empty command line string")
    }

    parser := shellwords.NewParser()
    parser.ParseEnv = true

    if args, err := parser.Parse(self.CommandLine); err == nil {
        self.Program = exec.Command(args[0], args[1:len(args)]...)
    }else{
        return err
    }

    return nil
}


func (self *Command) Execute(arguments ...string) (CommandResult, error) {
    output, err := self.Program.CombinedOutput()

    sep := "\n"
    outputLines := make([]string, 0)

    for _, line := range bytes.Split(output, []byte(sep[:])) {
        if str := strings.TrimSpace(string(line[:])); str != `` {
            outputLines = append(outputLines, str)
        }
    }

    return CommandResult{
        Error:  (err != nil),
        Output: outputLines,
    }, err
}