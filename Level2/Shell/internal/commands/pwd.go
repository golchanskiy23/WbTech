package commands

import (
	"flag"
	"fmt"
)

type Pwd struct {
	handler Handler
}

func (pwd *Pwd) Execute(params []string, handler Handler, fs *flag.FlagSet) (string, error) {
	return handler.Handle(params, fs)
}

type PwdWithoutParams struct {
	NextHandler Handler
}

func (pwd *PwdWithoutParams) setNext(handler Handler) {
	pwd.NextHandler = handler
}

func (pwd *PwdWithoutParams) Handle(params []string, flags *flag.FlagSet) (string, error) {
	if len(params) != 0 {
		return "", fmt.Errorf("incorrect number of params: %d\n", len(params))
	}
	homeDir, err := getCurrentDirectory()
	if err != nil {
		return "", fmt.Errorf("could not get current directory: %v\n", err)
	}
	return fmt.Sprintf("Root directory is: %s\n", homeDir), nil
}

func (pwd *PwdWithoutParams) Next() Handler {
	return pwd.NextHandler
}
