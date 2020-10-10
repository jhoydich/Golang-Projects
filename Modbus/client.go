package main

import(
	"net"
	"fmt"
	"bufio"
	"time"
)


func main() {

	//connecting to server
	c, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		fmt.Println(err)
		return
	}

	//for loop sending data every half a second
	for {
		data := "Data in here"
		fmt.Fprintf(c, data+"\n")
		time.Sleep(time.Millisecond * 500)

	}

}