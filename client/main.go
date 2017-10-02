package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	rAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:10001")
	if err != nil {
		log.Printf("Failed to resolve local addr: %v", err)
		time.Sleep(1)
		os.Exit(1)
	}

	lAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:10002")
	if err != nil {
		log.Printf("Failed to resolve local addr: %v", err)
		time.Sleep(1)
		os.Exit(1)
	}

	conn, err := net.DialUDP("udp", lAddr, rAddr)
	if err != nil {
		log.Printf("Failed to Dial: %v", err)
		time.Sleep(1)
		os.Exit(1)
	}
	defer conn.Close()

	s := bufio.NewScanner(io.TeeReader(os.Stdin, conn))
	for s.Scan() {
		err := readAck(conn)
		if err != nil {
			log.Printf("error: %v", err)
		}
	}
	log.Printf("Err: %v", s.Err())
}

func readAck(conn *net.UDPConn) error {
	resp := make([]byte, 2048)
	_, err := bufio.NewReader(conn).Read(resp)
	if err != nil {
		return err
	}
	log.Printf("Server: %s\n", resp)
	return nil
}
