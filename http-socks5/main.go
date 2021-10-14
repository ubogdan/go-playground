package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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
		fmt.Fprintln(os.Stderr, "can't connect to the proxy:", err)
	}

	tr := http.Transport{Dial: dialer.Dial}
	myClient := &http.Client{
		Transport: &tr,
	}

	res, err := myClient.Get("https://ifconfig.co/ip")
	if err != nil {
		log.Fatalf("http.GET failed:%s", err)
	}
	defer res.Body.Close()

	ip, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("ioutil.Read failed %s", err)
	}

	log.Printf("IP %s", ip)
}
