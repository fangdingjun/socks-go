package main

import (
	"bufio"
	"crypto/tls"
	"io"
	"log"
	"net"
	"net/http"
	"os"

	socks "github.com/fangdingjun/socks-go"
)

func main() {
	// connect to socks server
	c, err := net.Dial("tcp", "localhost:1080")
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	sc := &socks.Client{Conn: c, Username: "admin", Password: "passwd"}

	// connect to remote server
	if err := sc.Connect("httpbin.org", 443); err != nil {
		log.Fatal(err)
	}

	// tls
	conn := tls.Client(sc, &tls.Config{ServerName: "httpbin.org"})
	if err := conn.Handshake(); err != nil {
		log.Fatal(err)
	}

	// send http request
	req, err := http.NewRequest("GET", "https://httpbin.org/get", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Write(conn)

	bio := bufio.NewReader(conn)

	// read response
	res, err := http.ReadResponse(bio, req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	io.Copy(os.Stdout, res.Body)
}
