package main

import (
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"github.com/VP"
)

//type Message struct {
//	Name string
//	Text string
//}

type Server struct {
	clients map[net.Conn]string
	rwMutex sync.RWMutex
}

func (server *Server) StartServer() {

	
	l, err := net.Listen("tcp", ":2222")
	if err != nil {
		log.Fatal(err)
	}

	defer func(l net.Listener) {
		err := l.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(l)

	server.clients = make(map[net.Conn]string)

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		log.Println("accept client " + conn.RemoteAddr().String())

		name := make([]byte, 1024)

		if _, err := conn.Read(name); err != nil { // must change
			log.Fatal(err)
		}

		server.rwMutex.Lock()
		server.clients[conn] = string(name)
		server.rwMutex.Unlock()

		go func(conn net.Conn, name string) {

			defer func(conn net.Conn) {
				err := conn.Close()
				if err != nil { // mb must change
					log.Fatal(err)
				}
			}(conn)

			for {
				message := make([]byte, 4096)

				length, err := conn.Read(message)
				if err != nil { // must change
					return
				}
				
				fmt.Printf("%q\n", strings.Trim(string(message[:length]), " \n\t"))
				if strings.Trim(string(message[:length]), " \n\t") == "\\stop" {
					server.rwMutex.Lock()
					delete(server.clients, conn)
					fmt.Println("delete client: " + conn.RemoteAddr().String())
					server.rwMutex.Unlock()
					return
				}

				server.rwMutex.Lock()
				for clientConn := range server.clients {
					if clientConn != conn {
						_, err := clientConn.Write([]byte(name + ": " + string(message)))
						if err != nil {
							log.Println("fuck you")
							delete(server.clients, clientConn)
						}
					}
				}
				server.rwMutex.Unlock()
			}
		}(conn, server.clients[conn])
	}
}



func main() {
	
	server := Server{}
	server.StartServer()
}
