package internal

import (
	"flag"
	"shell/internal/commands"
)

type CmdMeta struct {
	Cmd   commands.Command
	Flags func(fs *flag.FlagSet)
	Chain func() commands.Handler
}

type Shell struct {
	CurrentCommand commands.Command
	Registry       map[string]CmdMeta
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
	}
}

func (shell *Shell) SetCommand(s string) {
	shell.CurrentCommand = shell.Registry[s].Cmd
}
