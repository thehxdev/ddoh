package server

import (
	"context"
	"log"
	"net"
	"strconv"

	"github.com/thehxdev/ddoh/config"
	"github.com/thehxdev/ddoh/resolver"
)

type Server struct {
	bufPool
	*resolver.Resolver
	Conn *net.UDPConn
	Addr *net.UDPAddr
	Ctx  context.Context
}

// var (
// 	running bool = true
// )

func Init() *Server {
	s := &Server{
		Addr: &net.UDPAddr{
			IP:   net.ParseIP("127.0.0.1"),
			Port: 53,
		},
		Resolver: resolver.Init(),
		bufPool:  newPool(config.Global.UDPBuffSize),
	}
	conn, err := net.ListenUDP("udp", s.Addr)
	if err != nil {
		log.Fatal(err)
	}
	s.Conn = conn
	return s
}

func (s *Server) Start() error {
	addr := s.Addr
	log.Printf("starting server on %s\n", net.JoinHostPort(addr.IP.String(), strconv.Itoa(addr.Port)))

	for {
		buff := s.bufPool.Get()
		_, addr, err := s.Conn.ReadFromUDP(buff[:cap(buff)])
		if err != nil {
			return err
		}
		go func(buff []byte) {
			log.Printf("new query from %s\n", addr.String())
			if err := s.Resolver.Resolve(s.Ctx, s.Conn, addr, buff[:cap(buff)]); err != nil {
				log.Println(err)
			}
			s.bufPool.Put(buff)
		}(buff)
	}
}

func (s *Server) Shutdown() {
	// running = false
	s.Conn.Close()
}
