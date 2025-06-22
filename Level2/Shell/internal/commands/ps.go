package commands

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Ps struct {
	handler Handler
}

func (ps *Ps) Execute(params []string, handler Handler, fs *flag.FlagSet) ([]string, error) {
	return handler.Handle(params, fs)
}

type PsWithoutFlags struct {
	NextHandler Handler
}

func (ps *PsWithoutFlags) setNext(handler Handler) {
	ps.NextHandler = handler
}

func (ps *PsWithoutFlags) Handle([]string, *flag.FlagSet) ([]string, error) {
	files, err := os.ReadDir("/proc")
	if err != nil {
		return nil, err
	}

	var processes []string
	for _, file := range files {
		if pid, err := strconv.Atoi(file.Name()); err == nil {
			statusPath := fmt.Sprintf("/proc/%d/status", pid)
			content, err := os.ReadFile(statusPath)
			if err != nil {
				continue
			}
			lines := strings.Split(string(content), "\n")
			var name string
			for _, line := range lines {
				if strings.HasPrefix(line, "Name:") {
					name = strings.TrimSpace(strings.TrimPrefix(line, "Name:"))
					break
				}
			}
			processes = append(processes, fmt.Sprintf("%5d %s", pid, name))
		}
	}
	return processes, nil
}

func (ps *PsWithoutFlags) Next() Handler {
	return ps.NextHandler
}
