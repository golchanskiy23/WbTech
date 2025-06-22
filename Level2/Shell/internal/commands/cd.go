package commands

import (
	"flag"
	"fmt"
	"os"
)

type Cd struct {
	handler Handler
}

func (cd *Cd) Execute(params []string, handler Handler, fs *flag.FlagSet) ([]string, error) {
	return handler.Handle(params, fs)
}

type CdWithoutParams struct {
	NextHandler Handler
}

func (cd *CdWithoutParams) setNext(handler Handler) {
	cd.NextHandler = handler
}

func getCurrentDirectory() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not get current directory: %v", err)
	}
	return homeDir, nil
}

func (cd *CdWithoutParams) Handle(params []string, flags *flag.FlagSet) ([]string, error) {
	var target string
	if len(params) == 1 {
		if params[0] == "." {
			currDir, err := os.Getwd()
			if err != nil {
				return nil, fmt.Errorf("could not get current directory: %v\n", err)
			}
			target = currDir
		} else if params[0] == "~" {
			if homeDir, err := getCurrentDirectory(); err != nil {
				return nil, fmt.Errorf("could not get current directory: %v\n", err)
			} else {
				target = homeDir
			}
		}
	} else if len(params) == 0 {
		if homeDir, err := getCurrentDirectory(); err != nil {
			return nil, fmt.Errorf("could not get current directory: %v\n", err)
		} else {
			target = homeDir
		}
	} else {
		return nil, fmt.Errorf("invalid command line arguments: %d\n", len(params))
	}

	if err := os.Chdir(target); err != nil {
		return nil, fmt.Errorf("could not find such directory: %s, execution ended with error : %v\n", target, err)
	}
	return []string{fmt.Sprintf("Current directory is: %s\n", target)}, nil
}

func (cd *CdWithoutParams) Next() Handler {
	return cd.NextHandler
}
