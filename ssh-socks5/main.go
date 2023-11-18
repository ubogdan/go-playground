package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/net/proxy"
)

func main() {
	sshHost := os.Getenv("SSH_HOST")
	if sshHost == "" {
		log.Fatal("ssh host address required")
	}

	sshPort := os.Getenv("SSH_PORT")
	if sshPort == "" {
		sshPort = "22"
	}

	sshUser := os.Getenv("SSH_USER")
	if sshUser == "" {
		log.Fatal("ssh user required")
	}

	sshPass := os.Getenv("SSH_PASSWORD")
	if sshPass == "" {
		log.Fatal("ssh password required")
	}

	auth := proxy.Auth{
		User:     os.Getenv("PROXY_USERNAME"),
		Password: os.Getenv("PROXY_PASSWORD"),
	}

	dialer, err := proxy.SOCKS5("tcp", os.Getenv("PROXY_ADDRESS"), &auth, proxy.Direct)
	if err != nil {
		log.Fatalf("can't connect to the proxy: %s", err)
	}

	conn, err := dialer.Dial("tcp", net.JoinHostPort(sshHost, sshPort))
	if err != nil {
		log.Fatalf("dial %s", err)
	}

	pConnect, pChans, pReqs, err := ssh.NewClientConn(conn, net.JoinHostPort(sshHost, sshPort), &ssh.ClientConfig{
		User:            sshUser,
		Auth:            []ssh.AuthMethod{ssh.Password(sshPass)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})
	if err != nil {
		fmt.Println(err)
	}

	client := ssh.NewClient(pConnect, pChans, pReqs)

	session, err := client.NewSession()
	defer session.Close()

	fd := int(os.Stdin.Fd())

	state, err := terminal.MakeRaw(fd)
	if err != nil {
		fmt.Println(err)
	}
	defer terminal.Restore(fd, state)

	w, h, err := terminal.GetSize(fd)
	if err != nil {
		fmt.Println(err)
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	err = session.RequestPty("xterm", h, w, modes)
	if err != nil {
		fmt.Println(err)
	}

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	err = session.Shell()
	if err != nil {
		fmt.Println(err)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGWINCH)
	go func() {
		for {
			switch <-sig {
			case syscall.SIGWINCH:
				fd := int(os.Stdout.Fd())
				w, h, _ = terminal.GetSize(fd)
				session.WindowChange(h, w)
			}
		}
	}()

	err = session.Wait()
	if err != nil {
		fmt.Println(err)
	}

}
