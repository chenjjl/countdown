package server

import (
	"bufio"
	"countdown/src/logger"
	"net"
)

var log = logger.GetLogger("server")

type Server struct {
}

func (s *Server) Start() {
	address := "127.0.0.1:9000"
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Errorf("failed to listen address %s", address)
		panic(err)
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			conn.Close()
			log.Errorf("failed to accept tcp connection, error is %+v", err)
		} else {
			handle(conn)
		}
	}
}

func handle(conn net.Conn) {
	defer conn.Close()
	for {
		reader := bufio.NewReader(conn)
		reader.ReadString
	}
}
