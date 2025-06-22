package commands

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
)

type Echo struct {
	handler Handler
}

type EchoEscape struct {
	NextHandler Handler
}

func (e *EchoEscape) setNext(handler Handler) {
	e.NextHandler = handler
}

func (e *EchoEscape) Handle(params []string, flags *flag.FlagSet) ([]string, error) {
	curr := strings.Join(params, " ")
	if flags.Lookup("e") == nil || flags.Lookup("e").Value.String() != "true" {
		return params, nil
	}
	if !(strings.HasPrefix(curr, `"`) && strings.HasSuffix(curr, `"`)) {
		curr = `"` + curr + `"`
	}
	processed2, err := strconv.Unquote(curr)
	if err != nil {
		return nil, fmt.Errorf("invalid escape sequence: %v\ninput: %q", err, curr)
	}
	return []string{processed2}, nil
}

func (e *EchoEscape) Next() Handler {
	return e.NextHandler
}

type EchoOmit struct {
	NextHandler Handler
}

func (e *EchoOmit) setNext(handler Handler) {
	e.NextHandler = handler
}

func (e *EchoOmit) Handle(params []string, flags *flag.FlagSet) ([]string, error) {
	if flags.Lookup("n") != nil && flags.Lookup("n").Value.String() == "true" {
		return []string{strings.TrimSuffix(strings.Join(params, " "), "\n")}, nil
	}
	return params, nil
}

func (e *EchoOmit) Next() Handler {
	return e.NextHandler
}

type EchoNormalMode struct {
	NextHandler Handler
}

func (e *EchoNormalMode) setNext(handler Handler) {
	e.NextHandler = handler
}

func (e *EchoNormalMode) Handle(params []string, flags *flag.FlagSet) ([]string, error) {
	builder := strings.Builder{}
	builder.WriteString(strings.Join(params, " "))
	ans := builder.String()
	if countOfStars(ans)%2 != 0 {
		return nil, fmt.Errorf("number of stars is incorrect")
	} else {
		ans = strings.ReplaceAll(ans, `"`, "")
	}
	return []string{fmt.Sprintf("%s\n", ans)}, nil
}

func (e *EchoNormalMode) Next() Handler {
	return e.NextHandler
}

func countOfStars(s string) int {
	runes, i := []rune(s), 0
	for _, r := range runes {
		if r == '"' {
			i++
		}
	}
	return i
}

func (e *Echo) Execute(params []string, handler Handler, fs *flag.FlagSet) ([]string, error) {
	tmpData := params
	for handler != nil {
		if val, err := handler.Handle(tmpData, fs); err == nil {
			tmpData = val
			handler = handler.Next()
		} else {
			return nil, err
		}
	}
	return tmpData, nil
}
