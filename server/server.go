package server

import (
	"log"
	"net"
	"strconv"

	"github.com/thehxdev/ddoh/config"
	"github.com/thehxdev/ddoh/resolver"
)

type Server struct {
	Conn     *net.UDPConn
	Addr     *net.UDPAddr
	Resolver *resolver.Resolver
}

var (
	running bool = true
)

func Init() *Server {
	s := &Server{
		Addr: &net.UDPAddr{
			IP:   net.ParseIP("127.0.0.1"),
			Port: 53,
		},
		Resolver: resolver.Init(),
	}
	conn, err := net.ListenUDP("udp", s.Addr)
	if err != nil {
		log.Fatal(err)
	}
	s.Conn = conn
	return s
}

func (s *Server) Start() {
	addr := s.Addr
	log.Printf("starting server on %s\n", net.JoinHostPort(addr.IP.String(), strconv.Itoa(addr.Port)))

	for running {
		buff := make([]byte, 0, config.Global.UDPBuffSize)
		_, addr, err := s.Conn.ReadFrom(buff[:cap(buff)])
		if err != nil {
			log.Println(err)
			break
		}
		go s.Resolver.Resolve(s.Conn, addr, buff[:cap(buff)])
	}
}

func (s *Server) Shutdown() {
	running = false
	s.Conn.Close()
}
