package rrserver

import (
	"fmt"
	"github.com/songtianyi/rrframework/logs"
	"io"
	"net"
	"strconv"
	"strings"
)

var (
	CustomHandleConn = func(c *TCPConnection, packet []byte) {
		logs.Warn("forget rrserver.CustomHandleConn = YourCustomHandleConn in init func?")
	}
)

type TCPServer struct {
	ls   net.Listener
	port int
}

func CreateTCPServer(inf string, port int) (error, *TCPServer) {
	err, ipaddr := getIpAddrByInterface(inf)
	if err != nil {
		return err, nil
	}
	listener, err := net.Listen("tcp", net.JoinHostPort(ipaddr, strconv.Itoa(port)))
	if err != nil {
		return err, nil
	}
	s := &TCPServer{
		ls:   listener,
		port: port,
	}
	return nil, s
}

func CreateTCPServerByAddr(addr string) (error, *TCPServer) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err, nil
	}
	parts := strings.Split(addr, ":")
	if len(parts) != 2 {
		return fmt.Errorf("addr invalid len=%d", len(parts)), nil
	}
	p, err := strconv.Atoi(parts[1])
	if err != nil {
		return err, nil
	}
	s := &TCPServer{
		ls:   listener,
		port: p,
	}
	return nil, s
}

func (s *TCPServer) Start() {
	logs.Info("Server listening in [%s]", s.ls.Addr())
	for {
		conn, err := s.ls.Accept()
		if err != nil {
			logs.Error("Server Accept() return error, %s", err)
			break
		}
		go s.handleConn(NewTCPConnection(conn))
	}
	return
}

func (s *TCPServer) handleConn(c *TCPConnection) {
	for {
		err, packet := c.Read()
		if err != nil {
			// end goroutine
			if err != io.EOF {
				logs.Error("Error occurred when read packet, %s", err)
			}
			return
		}
		// Maybe thousands of packet coming in the same time
		// So lock for connection is necessary
		go CustomHandleConn(c, packet)
	}
}
