package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

func main() {
	lAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:10001")
	if err != nil {
		log.Println("Failed to resolve local UDP addr")
	}

	conn, err := net.ListenUDP("udp", lAddr)
	if err != nil {
		log.Fatalf("Failed to listen on %v: %v", lAddr, err)
	}
	defer conn.Close()

	for i := 0; i < 8; i++ {
		log.Printf("Starting server %v", i)
		go func(id int) {
			for {
				err := readMsg(conn, id)
				if err != nil {
					log.Printf("Failed to read message: %v\n", err)
				}
			}
		}(i)
	}
	for {
	}
}

func readMsg(conn *net.UDPConn, id int) error {
	buf := make([]byte, 1024)
	n, rAddr, err := conn.ReadFromUDP(buf)
	if err != nil {
		return err
	}
	in := string(buf[0:n])
	in = strings.Trim(in, "\n")
	out, err := strconv.Atoi(in)
	if err != nil {
		_, wErr := conn.WriteToUDP([]byte("Bad input"), rAddr)
		if wErr != nil {
			return fmt.Errorf("%v occurred and response could not be sent: %v", err, wErr)
		}
		return err
	}

	_, err = conn.WriteToUDP([]byte(fmt.Sprintf("Your new value is %v", out*2)), rAddr)
	if err != nil {
		return err
	}

	log.Printf("Server %v computed 2x%v=%v", id, in, out)

	return nil
}
