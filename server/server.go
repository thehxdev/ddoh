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

// NOTE: I know these three functions are wierd but anyway it's a solution
// to eliminate checking for verbose logging in the handler itself for each
// dns request.
// This way, I create tow handlers, one is quiet and one is verbose. The server
// can pick one of them before starting up and this causes better performance in
// Higher loads because there's no need to check verbosity settings.
func (s *Server) handler(buff []byte, addr *net.UDPAddr) {
	if err := s.Resolver.Resolve(s.Ctx, s.Conn, addr, buff[:cap(buff)]); err != nil {
		log.Println(err)
	}
	s.bufPool.Put(buff)
}

func (s *Server) quietHandler(buff []byte, addr *net.UDPAddr) {
	s.handler(buff, addr)
}

func (s *Server) verboseHandler(buff []byte, addr *net.UDPAddr) {
	log.Printf("new query from %s\n", addr.String())
	s.handler(buff, addr)
}

func (s *Server) Start() error {
	addr := s.Addr
	log.Printf("starting server on %s\n", net.JoinHostPort(addr.IP.String(), strconv.Itoa(addr.Port)))

	handler := s.quietHandler
	if config.Global.Verbose {
		handler = s.verboseHandler
	}

	// Maximum 20 concurrent handlers
	jobChan := make(chan struct{}, 20)
	defer close(jobChan)

	for {
		buff := s.bufPool.Get()
		_, addr, err := s.Conn.ReadFromUDP(buff[:cap(buff)])
		if err != nil {
			return err
		}
		go func(buff []byte, ch chan struct{}) {
			jobChan <- struct{}{}
			handler(buff, addr)
			<-jobChan
		}(buff, jobChan)
	}
}

func (s *Server) Shutdown() {
	s.Conn.Close()
}
