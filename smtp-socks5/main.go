package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/smtp"
	"os"

	"golang.org/x/net/proxy"
)

func main() {
	auth := proxy.Auth{
		User:     os.Getenv("PROXY_USERNAME"),
		Password: os.Getenv("PROXY_PASSWORD"),
	}

	dialer, err := proxy.SOCKS5("tcp", os.Getenv("PROXY_ADDRESS"), &auth, proxy.Direct)
	if err != nil {
		log.Fatalf("can't connect to the proxy:", err)
	}

	err = smtp_dial(dialer, "smtp.gmail.com:587")
	if err != nil {
		log.Printf("smtp_dial [smtp.gmail.com:587] failed %s", err)
	}

	err = smtp_dial(dialer, "smtp.gmail.com:465")
	if err != nil {
		log.Printf("smtp_dial [smtp.gmail.com:465] failed %s", err)
	}
}

func smtp_dial(dialer proxy.Dialer, address string) error {
	hostname, port, err := net.SplitHostPort(address)
	if err != nil {
		return fmt.Errorf(" net.SplitHostPort failed %s", err)
	}

	conn, err := dialer.Dial("tcp", address)
	if err != nil {
		return fmt.Errorf("dialer.Dial failed: %s", err)
	}

	tlsConfig := tls.Config{
		ServerName: hostname,
	}

	if port == "465" {
		conn = tls.Client(conn, &tlsConfig)
	}

	client, err := smtp.NewClient(conn, address)
	if err != nil {
		return fmt.Errorf("smtp.NewClient failed: %s", err)
	}

	if port != "465" {
		err = client.StartTLS(&tlsConfig)
		if err != nil {
			return fmt.Errorf("client.Auth fialed: %s", err)
		}
	}

	smtpAuth := smtp.PlainAuth("", "mymailaddr@gmail.com", "", address)

	err = client.Auth(smtpAuth)
	if err != nil {
		return fmt.Errorf("client.Auth fialed: %s", err)
	}

	// TODO: The rest of the SMTP communication

	return nil
}
