package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)


func GetName() string {
	var res string
	fmt.Print("Enter your name: ")
	_, err := fmt.Scan(&res)
	if err != nil {
		log.Fatal(err)
	}
	return res
}

func startClient(name string) {
	conn, err := net.Dial("tcp", "localhost:2222")
	if err != nil {
		fmt.Println("\nSorry, server is down now.")
		return
	}

	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(conn)

	if _, err := conn.Write([]byte(name)); err != nil {
		log.Fatal(err)
	}

	// го рутина для обработки входящих сообщений
	go func(conn net.Conn) {
		for {
			message := make([]byte, 4096)
			_, err := conn.Read(message)
			if err == io.EOF {
				fmt.Println("\nSorry, server is down now.")
				os.Exit(0)
			}
			if err != nil {
				log.Fatal(err)
			}

			fmt.Print(string(message))
		}
	}(conn)

	for {
		inputReader := bufio.NewReader(os.Stdin)
		var input string
		var err error
		if input, err = inputReader.ReadString('\n'); err != nil {
			log.Fatal(err)
		}

		// отправляем сообщение
		if _, err := conn.Write([]byte(input)); err != nil {
			log.Fatal(err)
		}

		if input == "\\stop\n" {
			break
		}
	}

}

func main() {
	startClient(GetName())
}
