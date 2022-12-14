package server

import (
	"fmt"
	"net"
)

const (
	SERVER_HOST = "66.97.44.3"
	SERVER_PORT = "9945"
	SERVER_TYPE = "tcp"
)

func Serve() {

	server, err := net.Listen(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)
	if err != nil {
		panic(err)
	}

	defer server.Close()

	fmt.Println("Listening on " + SERVER_HOST + ":" + SERVER_PORT)

	for {
		connection, err := server.Accept()

		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			continue
			// os.Exit(1)
		}

		go NewConnection(connection)

	}
}
