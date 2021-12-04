package main

import (
	"crypto/sha512"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os/exec"

	"github.com/creack/pty"
	"github.com/gliderlabs/ssh"
	gossh "golang.org/x/crypto/ssh"
)

func main() {
	var listen, keyPath, fingerprint, password string

	flag.StringVar(&listen, "address", ":2222", "Listen address")
	flag.StringVar(&password, "password", "123456", "Listen address")
	flag.StringVar(&keyPath, "k", "id_rsa", "Path to private key for SSH server")
	flag.StringVar(&fingerprint, "f", "OpenSSH_8.2p1 Debian-4", "SSH Fingerprint, excluding the SSH-2.0- prefix")
	flag.Parse()

	log.SetPrefix("SSH - ")

	salt := "1c362db832f3f864c8c2fe05f2002a05"
	hash := hashPassword(password, salt)

	privKeyBytes, err := ioutil.ReadFile(keyPath)
	if err != nil {
		log.Panicln("Error reading privkey:\t", err.Error())
	}

	privateKey, err := gossh.ParsePrivateKey(privKeyBytes)
	if err != nil {
		log.Panicln("Error parsing privkey:\t", err.Error())
	}

	server := ssh.Server{
		Addr:    listen,
		Version: fingerprint,
		// Password verifier
		PasswordHandler: func(sesion ssh.Context, password string) bool {
			log.Printf("Auth request %s:%s", sesion.User(), password)
			return hashPassword(password, salt) == hash
		},
		// Session handler
		Handler: func(s ssh.Session) {
			cmd := exec.Command("/bin/bash", "-i")
			ptyReq, _, isPty := s.Pty()
			if isPty {
				cmd.Env = append(cmd.Env, fmt.Sprintf("TERM=%s", ptyReq.Term))
				f, err := pty.Start(cmd)
				if err != nil {
					panic(err)
				}
				go func() {
					io.Copy(f, s) // stdin
				}()
				io.Copy(s, f) // stdout
				cmd.Wait()
			} else {
				io.WriteString(s, "No PTY requested.\n")
				s.Exit(1)
			}
		},
	}

	server.AddHostKey(privateKey)

	log.Println("Started SSH server on", server.Addr)

	log.Fatal(server.ListenAndServe())
}

func hashPassword(password string, salt string) string {
	hash := sha512.Sum512([]byte(password + salt))
	return fmt.Sprintf("%x", hash)
}
