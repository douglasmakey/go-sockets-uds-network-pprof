package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
)

const socketPath = "/tmp/echo.sock"

func main() {
	// Create a Unix domain socket and listen for incoming connections.
	socket, err := net.Listen("unix", socketPath)
	if err != nil {
		panic(err)
	}

	// Cleanup the sockfile.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		os.Remove(socketPath)
		os.Exit(1)
	}()

	for {
		// Accept an incoming connection.
		conn, err := socket.Accept()
		if err != nil {
			panic(err)
		}

		fmt.Println("new connection")
		// Handle the connection in a separate goroutine.
		go func(conn net.Conn) {
			defer conn.Close()
			// Create a buffer for incoming data.
			buf := make([]byte, 4096)

			// Read data from the connection.
			n, err := conn.Read(buf)
			if err != nil {
				panic(err)
			}

			fmt.Println("got a message")
			// Echo the data back to the connection.
			_, err = conn.Write(buf[:n])
			if err != nil {
				panic(err)
			}
		}(conn)
	}
}
