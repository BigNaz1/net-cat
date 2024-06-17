package main

import (
	"fmt"
	"os"
	"strconv"

	server "net-cat/internal"
)

var (
	usg1 string = "[USAGE]: ./TCPChat"
	usg2 string = "[USAGE]: ./TCPChat $port"
	usg3 string = "[USAGE]: ./TCPChat &host $port"
)

func main() {
	port := "8989"
	host := "localhost"
	if len(os.Args) == 2 {
		port = os.Args[1]
		_, err := strconv.Atoi(port)
		if err != nil {
			fmt.Println(usg2)
			fmt.Println("Please write correct port")
			return
		}
	} else if len(os.Args) == 3 {
		host = os.Args[1]
		port = os.Args[2]
		_, err := strconv.Atoi(port)
		if err != nil {
			fmt.Println(usg3)
			fmt.Println("Please write correct port")
			return
		}
	} else if len(os.Args) > 3 {
		fmt.Println(usg1)
		fmt.Println(usg2)
		fmt.Println(usg3)
		return
	}

	s := server.NewServer()
	err := s.Listen(host, port)
	if err != nil {
		fmt.Println("Something with listening the server")
		return
	}

	s.Run()
	defer s.Stop()
}
