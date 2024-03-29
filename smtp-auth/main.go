package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/mail"
	"net/smtp"
	"os"
)

func main() {

	from := mail.Address{"", os.Getenv("SMTP_USER")}
	to := mail.Address{"", os.Getenv("TO_EMAIL")}

	// Setup headers
	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = "This is the email subject"

	// Setup message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + "This is an example body.\n With two lines."

	servername := os.Getenv("SMTP_SERVER")
	hostname, _, _ := net.SplitHostPort(servername)

	auth := smtp.PlainAuth("", os.Getenv("SMTP_USER"), os.Getenv("SMTP_PASSWORD"), hostname)

	err := smtp_send(servername, auth, from, to, []byte(message))
	if err != nil {
		log.Printf("smtp_dial [%s] failed %s", servername, err)
	}
}

func smtp_send(address string, auth smtp.Auth, from, to mail.Address, message []byte) error {

	hostname, port, err := net.SplitHostPort(address)
	if err != nil {
		return fmt.Errorf(" net.SplitHostPort failed %s", err)
	}

	conn, err := net.Dial("tcp", address)
	if err != nil {
		return fmt.Errorf("dialer.Dial failed: %s", err)
	}

	tlsConfig := tls.Config{
		ServerName: hostname,
	}

	if port == "465" {
		conn = tls.Client(conn, &tlsConfig)
	}

	client, err := smtp.NewClient(conn, hostname)
	if err != nil {
		log.Panic(err)
	}

	if port != "465" {
		err = client.StartTLS(&tlsConfig)
		if err != nil {
			return fmt.Errorf("client.Auth fialed: %s", err)
		}
	}

	// Auth
	err = client.Auth(auth)
	if err != nil {
		return fmt.Errorf("smtp.Auth %s", err)
	}

	// To && From
	err = client.Mail(from.Address)
	if err != nil {
		return fmt.Errorf("smtp.Mail %s", err)
	}

	err = client.Rcpt(to.Address)
	if err != nil {
		return fmt.Errorf("smtp.Rcpt %s", err)
	}

	// Data
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("smtp.Data %s", err)
	}

	_, err = w.Write(message)
	if err != nil {
		return fmt.Errorf("data.Write %s", err)
	}

	err = w.Close()
	if err != nil {
		return fmt.Errorf("data.Close %s", err)
	}

	return client.Quit()
}
