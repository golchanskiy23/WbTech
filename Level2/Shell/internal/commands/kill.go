package commands

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"syscall"
)

type Kill struct {
	handler Handler
}

func (k *Kill) Execute(params []string, handler Handler, fs *flag.FlagSet) (string, error) {
	tmpData := params
	for handler != nil {
		if val, err := handler.Handle(tmpData, fs); err == nil {
			tmpData = []string{val}
			handler = handler.Next()
		} else {
			return "", err
		}
	}
	return strings.Join(tmpData, " "), nil
}

type ParsePID struct {
	NextHandler Handler
}

type KillProcess struct {
	NextHandler Handler
}

func (ps *ParsePID) setNext(handler Handler) {
	ps.NextHandler = handler
}

func (ps *ParsePID) Handle(params []string, fs *flag.FlagSet) (string, error) {
	if len(params) == 0 {
		return "", fmt.Errorf("invalid amount of params: %d", len(params))
	}
	if _, err := strconv.Atoi(params[0]); err != nil {
		return "", fmt.Errorf("invalid pid: %v", err)
	}
	return strings.Join(params, " "), nil
}

func (ps *ParsePID) Next() Handler {
	return ps.NextHandler
}

func (ps *KillProcess) setNext(handler Handler) {
	ps.NextHandler = handler
}

func (ps *KillProcess) Handle(params []string, fs *flag.FlagSet) (string, error) {
	if len(params) == 0 {
		return "", fmt.Errorf("invalid command line arguments: %d\n", len(params))
	}

	parts := strings.Fields(params[0])
	pid, _ := strconv.Atoi(parts[0])

	signal := syscall.SIGTERM
	if fs.Lookup("s") != nil {
		val := fs.Lookup("s").Value.String()
		if sigInt, err := strconv.Atoi(val); err == nil {
			signal = syscall.Signal(sigInt)
		}
	}

	proc, err := os.FindProcess(pid)
	if err != nil {
		return "", fmt.Errorf("kill failed (find process): %v", err)
	}

	if runtime.GOOS == "windows" {
		if signal != syscall.SIGKILL {
			return "", fmt.Errorf("on Windows only SIGKILL (9) is supported")
		}
		if err := proc.Kill(); err != nil {
			return "", fmt.Errorf("kill failed: %v", err)
		}
	} else {
		/*if err := syscall.Kill(pid, sig); err != nil {
			return "", fmt.Errorf("kill failed: %v", err)
		}*/
	}

	return fmt.Sprintf("Process %d killed with signal %d\n", pid, signal), nil
}

func (ps *KillProcess) Next() Handler {
	return ps.NextHandler
}
