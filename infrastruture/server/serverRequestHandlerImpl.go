package server

import (
	"github.com/dcbCIn/MidCloud/lib"
	"net"
	"strconv"
)

type ServerRequestHandlerImpl struct {
	Port       int
	listener   net.Listener
	connection net.Conn
}

func NewServerRequestHandlerImpl(port int) (srh *ServerRequestHandlerImpl, err error) {
	srh = &ServerRequestHandlerImpl{Port: port}

	srh.listener, err = net.Listen("tcp", ":"+strconv.Itoa(srh.Port))
	if err != nil {
		return nil, err
	}
	return srh, nil
}

func (s *ServerRequestHandlerImpl) Start() (err error) {
	lib.PrintlnInfo("ServerRequestHandler", "Aceitando conexões...")

	s.connection, err = s.listener.Accept()

	if err != nil {
		lib.PrintlnInfo("ServerRequestHandler", "Erro ao abrir conexão")
		return err
	}

	lib.PrintlnInfo("ServerRequestHandler", "Conexão aceita...")
	return nil
}

func (s *ServerRequestHandlerImpl) CloseConnection() (err error) {
	lib.PrintlnInfo("ServerRequestHandler", "ServerRequestHandler.Stop - Closing connection")
	err = s.connection.Close()
	if err != nil {
		return err
	}
	lib.PrintlnInfo("ServerRequestHandler", "ServerRequestHandler.Stop - Connection closed")
	return nil
}

func (s *ServerRequestHandlerImpl) StopServer() (err error) {
	err = s.listener.Close()
	if err != nil {
		return err
	}
	lib.PrintlnInfo("ServerRequestHandler", "ServerRequestHandler.Stop - Listener closed")
	return nil
}

func (s *ServerRequestHandlerImpl) Receive() (msg []byte, err error) {
	msg = make([]byte, 10240)
	n, err := s.connection.Read(msg)
	if err != nil {
		return nil, err
	}

	return msg[:n], nil
}

func (s *ServerRequestHandlerImpl) Send(msg []byte) (err error) {
	_, err = s.connection.Write(msg)
	if err != nil {
		return err
	}
	return nil
}
