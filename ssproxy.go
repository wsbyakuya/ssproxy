package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		return
	}

	setProxy()

	args := []string{}
	if len(os.Args) > 2 {
		args = os.Args[2:]
	}
	execute(os.Args[1], args)
}

func execute(name string, args []string) bool {
	cmd := exec.Command(name, args...)

	fmt.Println(">", name, strings.Join(args, " "))

	cmdStdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
		return false
	}

	cmdStderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Println(err)
		return false
	}
	cmd.Start()
	go io.Copy(os.Stdout, cmdStdout)
	go io.Copy(os.Stdout, cmdStderr)

	if err := cmd.Wait(); err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

const (
	socks5Proxy = "socks5://127.0.0.1:1080/"
)

func setProxy() {
	os.Setenv("http_proxy", socks5Proxy)
	os.Setenv("https_proxy", socks5Proxy)
}
