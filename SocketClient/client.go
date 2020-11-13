package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"runtime"
)

const (
	connHost = "localhost"
	connPort = "1337"
	connType = "tcp"
)

func main() {
	address := connHost + ":" + connPort

	fmt.Println("Connecting with", connType, "on", address)

	conn, err := net.Dial(connType, address)

	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}

	reader := bufio.NewReader(os.Stdin)

	for {

		fmt.Print("$: ")

		input, _ := reader.ReadString('\n')

		fmt.Println(input)
		if runtime.GOOS == "windows" {
			conn.Write([]byte("W" + input))
		} else {
			conn.Write([]byte(input))
		}

		if input == "shell_exec\n" || input == "shell_exec\r\n" {
			fmt.Print("> ")
			input, _ := reader.ReadString('\n')

			if runtime.GOOS == "windows" {
				conn.Write([]byte("W" + input))
			} else {
				conn.Write([]byte(input))
			}
		}

		message, _ := bufio.NewReader(conn).ReadString('\n')

		log.Println("Server relay:", message)

		if input == "exit\n" || input == "exit\r\n" {
			conn.Close()
			break
		}
	}
}
