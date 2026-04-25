package main

import (
	"flag"
	"net"
	"bufio"
	"os"
	"fmt"
	
)

func main(){
	address := flag.String("address", "127.0.0.1:3223", "IP address client")
	flag.Parse()

	conn, err := net.Dial("tcp", *address)
	if err != nil {
		fmt.Printf("Error with connect server %s", err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	serverReader := bufio.NewScanner(conn)
	consoleReader := bufio.NewScanner(os.Stdin)

	for {
		fmt.Printf("> ")

		if !consoleReader.Scan() {
			break
		}

		text := consoleReader.Text()

		if text == "exit" {
			break
		}

		_, err := conn.Write([]byte(text + "\n"))
		if err != nil {
			break
		}

		if serverReader.Scan() {
			line := serverReader.Text()
			fmt.Println(line)
		} else {
			fmt.Println("Отключение от сервера")
		}

	}

}