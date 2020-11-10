package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"runtime"
)

const (
	connHost = "localhost"
	connPort = "1337"
	connType = "tcp"
)

func main() {
	address := connHost + ":" + connPort

	fmt.Println("Starting", connType, "server on", address)

	l, err := net.Listen(connType, address)

	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}
	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println("Error:", err.Error())
			os.Exit(1)
		}
		fmt.Println("Client", c.RemoteAddr().String(), "connected.")

		go handleConnection(c)
	}
}

func handleConnection(conn net.Conn) {
	buffer, err := bufio.NewReader(conn).ReadBytes('\n')

	if err != nil {
		fmt.Println("Connection closed.")
		conn.Close()
		return
	}

	log.Println("Client message:", string(buffer[:len(buffer)-2]))

	switch string(buffer[:len(buffer)-2]) {
	case "shell_exec":
		conn.Write([]byte(shellExec(conn)))
	case "whichos":
		conn.Write([]byte(checkOS()))
	default:
		conn.Write([]byte("Command not found.\n"))
	}

	handleConnection(conn)
}

func checkOS() string {
	return "Running on " + runtime.GOOS + "\n"
}

func shellExec(conn net.Conn) string {
	buffer, err := bufio.NewReader(conn).ReadBytes('\n')

	if err != nil {
		conn.Close()
		return "Connection closed."
	}

	cmd := exec.Command(string(buffer[:len(buffer)-2]))
	stdout, err := cmd.Output()

	if err != nil {
		return err.Error() + "\n"
	}

	return string(stdout)
}
