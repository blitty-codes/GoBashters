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

		if input == "file\n" || input == "file\r\n" {
			fmt.Print("- name: ")
			title, _ := reader.ReadString('\n')

			fmt.Print("- body (use $ [enter] at the end): ")
			msg, _ := reader.ReadString('$')

			fmt.Println(msg)
			input = "file.000" + title[:len(title)-1] + ".000==." + string(msg[:len(msg)-1]) + "\n==."
			fmt.Println("input: " + input)
		}

		sendCommand(conn, input)

		if input == "shell_exec\n" || input == "shell_exec\r\n" {
			fmt.Print("> ")
			input, _ := reader.ReadString('\n')

			sendCommand(conn, input)
		}

		// TODO: Can we open a bash shell on another cmd?
		// to automize reverse shell
		// connect to it by nc -lnvp 4242"

		message, _ := bufio.NewReader(conn).ReadString('\n')

		log.Println("Server relay:", message)

		if input == "exit\n" || input == "exit\r\n" {
			conn.Close()
			break
		}
	}
}

func sendCommand(conn net.Conn, input string) {
	input = input + "\000"
	if runtime.GOOS == "windows" {
		conn.Write([]byte("W" + input))
	} else {
		conn.Write([]byte(input))
	}
}
