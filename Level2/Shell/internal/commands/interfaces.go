package commands

import "flag"

type Command interface {
	Execute(params []string, handler Handler, fs *flag.FlagSet) (string, error)
}

type Handler interface {
	setNext(handler Handler)
	Handle([]string, *flag.FlagSet) (string, error)
	Next() Handler
}
