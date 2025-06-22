package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"shell/internal"
	"strings"
)

// sort
// grep
// cut
// + cd
// + pwd
// + echo
// + kill
// + ps

// fork/exec commands
// pipeline
// wget
// telnet

// + quit (without flags and params - exit, else error)

// PATTERNS:
// Command - division shell commands into different struct
// with distributed interface
// Chain of responsibility - shell pipeline
// Strategy - commands behavior by flags

func main() {
	shell := &internal.Shell{CurrentCommand: nil}
	shell.InitShell()
	scanner := bufio.NewReader(os.Stdin)
	for {
		currDir, _ := os.Getwd()
		fmt.Printf("%s> ", currDir)
		words, _ := scanner.ReadString('\n')
		words = strings.TrimSpace(words)
		input := strings.Split(words, "|")

		if len(input) == 0 || input[0] == "" {
			continue
		} else {
			arr, flag_ := shell.CheckPipeline(input)
			if flag_ && strings.Contains(words, "|") {
				ans, err := shell.ExecutePipeline(arr)
				if err != nil {
					fmt.Printf("%s\n", err)
					continue
				}
				fmt.Printf("%s", strings.Join(ans, " "))
			} else {
				input = strings.Split(words, " ")
				if input[0] == "quit" {
					fmt.Println("Bye. Terminal is closing...")
					break
				}
				meta, ok := shell.Registry[input[0]]
				if !ok {
					fmt.Println("No  such command. Try again")
				} else {
					f := flag.NewFlagSet(input[0], flag.ContinueOnError)
					meta.Flags(f)

					if err := f.Parse(input[1:]); err != nil {
						fmt.Println("Error during parsing flagSet")
						continue
					}

					shell.SetCommand(input[0])
					handlerChain := meta.Chain()

					if val, err := shell.CurrentCommand.Execute(f.Args(), handlerChain, f); err != nil {
						fmt.Printf("Error during execution command: %v\n", err)
					} else {
						fmt.Printf("%s", strings.Join(val, " "))
					}
				}
			}
		}

	}
}
