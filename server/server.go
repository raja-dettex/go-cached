package server

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/raja-dettex/go-cache/cache"
	"github.com/raja-dettex/go-cache/commands"
)

type ServerOpts struct {
	ListenToAddr string
	IsLeader     bool
}
type Server struct {
	Opts  ServerOpts
	Cache cache.Cacher
}

func NewServer(opts ServerOpts, Cache cache.Cacher) *Server {
	return &Server{
		Opts:  opts,
		Cache: Cache,
	}
}

func (s *Server) Start() {
	ln, err := net.Listen("tcp", s.Opts.ListenToAddr)
	if err != nil {
		panic(err)
	}
	log.Printf("server listening on port %v", s.Opts.ListenToAddr)
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("error %v", err)
		}
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()
	buff := make([]byte, 2048)
	for {
		n, err := conn.Read(buff)
		if err == io.EOF {
			fmt.Println("closing connection...")
			conn.Write([]byte(fmt.Sprintf("end of line reached %v", err)))
		}
		if err != nil {
			log.Printf("conn read error %v", err)
			break
		}
		msg := buff[:n]
		s.handleCommand(conn, msg)
	}
}

func (s *Server) handleCommand(conn net.Conn, rawCommand []byte) {
	msg, err := commands.ParseMessage(rawCommand)
	if err != nil {
		// respond
	}
	switch msg.Cmd {
	case commands.CMDSet:
		s.handleSetCommand(conn, msg)
	case commands.CMDGet:
		s.handleGetCommand(conn, msg)
	case commands.CMDHas:
		s.handleHasCommand(conn, msg)
	case commands.CMDDelete:
		s.handleDeleteCommand(conn, msg)
	default:
		fmt.Println("this is a default case")
	}
}

func (s *Server) handleSetCommand(conn net.Conn, msg *commands.Msg) {
	err := s.Cache.Set(msg.Key, msg.Value, msg.TTL)
	if err != nil {
		conn.Write([]byte(fmt.Sprintf("set key error %v", err)))
	}
	conn.Write([]byte("set key successfull"))
}
func (s *Server) handleGetCommand(conn net.Conn, msg *commands.Msg) {
	res, err := s.Cache.Get(msg.Key)
	if err != nil {
		conn.Write([]byte(fmt.Sprintf("get key error %v", err)))
	}
	conn.Write(res)

}
func (s *Server) handleHasCommand(conn net.Conn, msg *commands.Msg) {
	ok := s.Cache.Has(msg.Key)
	if !ok {
		conn.Write([]byte("key does not exist"))
	}
	conn.Write([]byte("key exists"))
}
func (s *Server) handleDeleteCommand(conn net.Conn, msg *commands.Msg) {
	err := s.Cache.Delete(msg.Key)
	if err != nil {
		conn.Write([]byte(fmt.Sprintf("delete key error %v", err)))
	}
	conn.Write([]byte("delete successful"))
}
