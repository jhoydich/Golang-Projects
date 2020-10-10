package main

import (
	"fmt"
	"net"
	"time"
	"log"
	"bufio"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		//Lets you have multiple connections
		go handleConn(conn)
	}
}

//
func handleConn(c net.Conn) {
	defer c.Close()
	for {

		data, err := bufio.NewReader(c).ReadLine('\n')
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(data)

		time.Sleep(time.Millisecond * 200)
	}
}