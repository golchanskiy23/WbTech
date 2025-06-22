package commands

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

type Telnet struct {
	handler Handler
}

type TelnetConnect struct {
	NextHandler Handler
}

func (t *TelnetConnect) setNext(handler Handler) {
	t.NextHandler = handler
}

func (t *TelnetConnect) Handle(params []string, fs *flag.FlagSet) ([]string, error) {
	if len(params) < 2 {
		return nil, fmt.Errorf("usage: telnet <host> <port>")
	}

	host := params[0]
	port := params[1]
	address := net.JoinHostPort(host, port)

	conn, err := net.DialTimeout("tcp", address, fs.Lookup("timeout").Value.(flag.Getter).Get().(time.Duration))
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %v", err)
	}
	defer conn.Close()
	fmt.Printf("telnet connect success to %s\n", address)

	_, err = io.Copy(conn, os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error writing to connection: %v\n", err)
	}

	var ans []string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		ans = append(ans, text)
		conn.Write([]byte(text + "\n"))
	}

	conn.Close()
	fmt.Println("\nConnection closed")
	return ans, nil
}

func (t *TelnetConnect) Next() Handler {
	return t.NextHandler
}

func (t *Telnet) Execute(params []string, handler Handler, fs *flag.FlagSet) ([]string, error) {
	return handler.Handle(params, fs)
}
