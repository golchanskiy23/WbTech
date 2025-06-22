package commands

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type Wget struct {
	handler Handler
	Output  *string
}

type WgetWithoutFlags struct {
	NextHandler Handler
}

type WgetWithOutput struct {
	NextHandler Handler
}

func (wget *WgetWithoutFlags) setNext(handler Handler) {
	wget.NextHandler = handler
}

func getInfoFromWebPage(args []string) ([]string, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("usage: wget <url>")
	}
	req, err := http.NewRequest("GET", args[0], nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (GoWgetClient)")

	client := &http.Client{}
	resp, err := client.Do(req)
	ans := make([]string, 0)
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		ans = append(ans, scanner.Text())
	}
	if err != nil {
		return nil, fmt.Errorf("wget error: %v", err)
	}
	return ans, nil
}

func (wget *WgetWithoutFlags) Handle(arr []string, flag *flag.FlagSet) ([]string, error) {
	if flag.Lookup("O").Value.String() == "" {
		if ans, err := getInfoFromWebPage(arr); err != nil {
			return nil, fmt.Errorf("something went wrong when getting info from WebPage: %s", err)
		} else {
			return ans, nil
		}
	}
	return nil, nil
}

func (wget *WgetWithoutFlags) Next() Handler {
	return wget.NextHandler
}

func (wget *WgetWithOutput) setNext(handler Handler) {
	wget.NextHandler = handler
}

func (wget *WgetWithOutput) Handle(arr []string, flag *flag.FlagSet) ([]string, error) {
	if flag.Lookup("O").Value.String() != "" {
		ans, err := getInfoFromWebPage(arr)
		if err != nil {
			return nil, fmt.Errorf("something went wrong when getting info from WebPage: %s", err)
		}
		out, err := os.Create(flag.Lookup("O").Usage)
		if err != nil {
			return nil, fmt.Errorf("cannot create file: %v", err)
		}
		defer out.Close()
		out.Write([]byte(strings.Join(ans, "\n")))
		return []string{fmt.Sprintf("Downloaded to %s", flag.Lookup("O").Value.String())}, nil
	}
	return arr, nil
}

func (wget *WgetWithOutput) Next() Handler {
	return wget.NextHandler
}

func (w *Wget) Execute(params []string, handler Handler, fs *flag.FlagSet) ([]string, error) {
	tmpData := params
	for handler != nil {
		if val, err := handler.Handle(tmpData, fs); err == nil {
			tmpData = val
			handler = handler.Next()
		} else {
			return nil, err
		}
	}
	return []string{strings.Join(tmpData, "\n") + "\n"}, nil
}
