package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
)

func main() {
	if len(os.Args) < 2 {
		return
	}

	if err:= setProxy(); err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	args := make([]string, 0)
	if len(os.Args) > 2 {
		args = os.Args[2:]
	}
	if err := execute(os.Args[1], args); err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
}

func execute(name string, args []string) error {
	cmd := exec.Command(name, args...)

	fmt.Println(">", name, strings.Join(args, " "))

	cmdStdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	cmdStderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	err = cmd.Start()
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		io.Copy(os.Stdout, cmdStdout)
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		io.Copy(os.Stderr, cmdStderr)
		wg.Done()
	}()
	wg.Wait()

	if err := cmd.Wait(); err != nil {
		return err
	}

	return nil
}

const (
	socks5Proxy = "socks5://127.0.0.1:1080/"
)

func setProxy() error {
	if err := os.Setenv("http_proxy", socks5Proxy); err != nil {
		return err
	}
	if err := os.Setenv("https_proxy", socks5Proxy); err != nil {
		return err
	}
	return nil
}
