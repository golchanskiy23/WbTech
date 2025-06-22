package internal

import (
	"flag"
	"fmt"
	"shell/internal/commands"
	"strings"
	"time"
)

type CmdMeta struct {
	Cmd   commands.Command
	Flags func(fs *flag.FlagSet)
	Chain func() commands.Handler
}

type Response struct {
	commands   commands.Command
	parameters []string
}

type PipeExecutor interface {
	ExecutePipeline(response []Response) ([]string, error)
}

type Shell struct {
	CurrentCommand commands.Command
	Registry       map[string]CmdMeta
}

func (shell *Shell) ExecutePipeline(response []Response) ([]string, error) {
	var data []string

	for i, resp := range response {
		meta, ok := shell.Registry[resp.parameters[0]]
		if !ok {
			return nil, fmt.Errorf("no such command: %s", resp.parameters[0])
		}

		fs := flag.NewFlagSet(resp.parameters[0], flag.ContinueOnError)
		meta.Flags(fs)
		if err := fs.Parse(resp.parameters[1:]); err != nil {
			return nil, fmt.Errorf("error parsing flags: %v", err)
		}

		handlerChain := meta.Chain()
		shell.CurrentCommand = meta.Cmd

		var result []string
		var err error

		if i == 0 {
			result, err = shell.CurrentCommand.Execute(fs.Args(), handlerChain, fs)
		} else {
			result, err = shell.CurrentCommand.Execute(data, handlerChain, fs)
		}

		if err != nil {
			return nil, fmt.Errorf("error executing %s: %v", resp.parameters[0], err)
		}

		data = result
	}
	return data, nil
}

func (s *Shell) CheckPipeline(arr []string) ([]Response, bool) {
	var ans []Response
	for _, str := range arr {
		splitted := strings.Split(str, " ")
		if val, ok := s.Registry[splitted[0]]; ok {
			ans = append(ans, Response{commands: val.Cmd, parameters: splitted})
		} else {
			return []Response{}, false
		}
	}
	return ans, true
}

func buildEchoChain() commands.Handler {
	return &commands.EchoEscape{
		NextHandler: &commands.EchoNormalMode{
			NextHandler: &commands.EchoOmit{
				NextHandler: nil,
			},
		},
	}
}

func buildKillChain() commands.Handler {
	return &commands.ParsePID{
		NextHandler: &commands.KillProcess{
			NextHandler: nil,
		},
	}
}

func buildWgetChain() commands.Handler {
	return &commands.WgetWithoutFlags{
		NextHandler: &commands.WgetWithOutput{
			NextHandler: nil,
		},
	}
}

func (s *Shell) InitShell() {
	s.Registry = map[string]CmdMeta{
		"echo": {
			Cmd: &commands.Echo{},
			Flags: func(fs *flag.FlagSet) {
				fs.Bool("n", false, "omit trailing newline")
				fs.Bool("e", false, "interpret escape sequences")
			},
			Chain: buildEchoChain,
		},
		"cd": {
			Cmd:   &commands.Cd{},
			Flags: func(fs *flag.FlagSet) {},
			Chain: func() commands.Handler { return &commands.CdWithoutParams{} },
		},
		"pwd": {
			Cmd:   &commands.Pwd{},
			Flags: func(fs *flag.FlagSet) {},
			Chain: func() commands.Handler { return &commands.PwdWithoutParams{} },
		},
		"ps": {
			Cmd:   &commands.Ps{},
			Flags: func(fs *flag.FlagSet) {},
			Chain: func() commands.Handler { return &commands.PsWithoutFlags{} },
		},
		"kill": {
			Cmd: &commands.Kill{},
			Flags: func(fs *flag.FlagSet) {
				fs.Int("s", 15, "default signal number")
			},
			Chain: buildKillChain,
		},
		"wget": {
			Cmd: &commands.Wget{},
			Flags: func(fs *flag.FlagSet) {
				fs.String("O", "", "")
			},
			Chain: buildWgetChain,
		},
		"telnet": {
			Cmd: &commands.Telnet{},
			Flags: func(fs *flag.FlagSet) {
				fs.Duration("timeout", 10*time.Second, "timeout for connection closing")
			},
			Chain: func() commands.Handler { return &commands.TelnetConnect{} },
		},
	}
}

func (shell *Shell) SetCommand(s string) {
	shell.CurrentCommand = shell.Registry[s].Cmd
}
