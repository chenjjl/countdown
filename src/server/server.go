package server

import (
	"bufio"
	"countdown/src/event"
	"countdown/src/logger"
	"countdown/src/timeWheel"
	"io"
	"net"
	"time"
)

var log = logger.GetLogger("server")

type Server struct {
	timeWheel *timeWheel.TimeWheel
}

func NewServer() *Server {
	return &Server{
		timeWheel: timeWheel.NewTimeWheel(time.Second, 8, time.Minute, 8),
	}
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
			s.handle(conn)
		}
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	for {
		reader := bufio.NewReader(conn)
		data, err := reader.ReadSlice('\n')
		if err != nil {
			if err != io.EOF {
				log.Errorf("failed to receive data from tcp connection")
				log.Error(err)
			} else {
				break
			}
		}
		_event := event.Decode(string(data))
		err = s.timeWheel.Add(_event)
		if err != nil {
			log.Error(err)
		}
	}
}
