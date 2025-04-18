package main

import (
	"bufio"
	"catch-me-if-you-can/node"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
)

const (
	EXIT       = "exit"
	WINDOWS_OS = "windows"
)

func main() {
	inputCh := make(chan string)
	go getInput(inputCh)

	for {
		fmt.Print("Input code: ")
		input := <-inputCh

		if runtime.GOOS == WINDOWS_OS {
			input = strings.Replace(input, "\r\n", "", -1)
		} else {
			input = strings.Replace(input, "\n", "", -1)
		}

		if input == EXIT {
			os.Exit(0)
		}

		nodeLink, error := node.InitNodeByInput(input)
		if error != nil {
			log.Printf("!!! ERROR: %v\n", error)
			continue
		}
		fmt.Printf("Decoded code: %s\n", nodeLink.Decoded())
	}
}

func getInput(ch chan<- string) {
	reader := bufio.NewReader(os.Stdin)
	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			if err.Error() == "EOF" {
				continue
			}
			log.Println("Error input:", err)
			continue
		}
		input = strings.TrimSpace(input)
		ch <- input
	}
}
