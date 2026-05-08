package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
)

func main() {
	address := flag.String("address", "127.0.0.1:3223", "IP address client")
	flag.Parse()

	conn, err := net.Dial("tcp", *address)
	if err != nil {
		fmt.Printf("Error with connect server: %s\n", err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	go func() {
		serverReader := bufio.NewScanner(conn)
		for serverReader.Scan() {
			fmt.Printf("\n[Сервер]: %s\n> ", serverReader.Text())
		}

		fmt.Println("\nОтключение от сервера")
		os.Exit(0)
	}()

	consoleReader := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")

		if !consoleReader.Scan() {
			break
		}

		text := consoleReader.Text()

		if text == "exit" {
			break
		}

		_, err := conn.Write([]byte(text + "\n"))
		if err != nil {
			fmt.Println("Ошибка отправки данных")
			break
		}
	}
}
