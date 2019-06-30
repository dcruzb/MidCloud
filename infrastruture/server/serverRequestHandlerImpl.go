package server

import (
	"errors"
	"github.com/dcbCIn/MidCloud/lib"
	"net"
	"strconv"
)

type Client struct {
	connection net.Conn
}

type ServerRequestHandlerImpl struct {
	Port       int
	listener   net.Listener
	connection net.Conn
	clients    []Client
}

func NewServerRequestHandlerImpl(port int, initialConnections int) (srh *ServerRequestHandlerImpl, err error) {
	srh = &ServerRequestHandlerImpl{Port: port}
	srh.clients = make([]Client, initialConnections, initialConnections)

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

//------------------------------

func (srh *ServerRequestHandlerImpl) GetConnection(cliIdx int) (cl *Client, err error) {
	if cliIdx >= len(srh.clients) {
		return cl, errors.New("Invalid Client Index. Not enough clients (index=" + strconv.Itoa(cliIdx) + "/clients=" + strconv.Itoa(len(srh.clients)) + ")")
	}

	conn, err := srh.listener.Accept()
	if err != nil {
		lib.PrintlnError("Error while waiting for connection", err)
		return cl, err
	}

	cl = &srh.clients[cliIdx]

	cl.connection = conn

	return cl, nil
}

func (cl *Client) CloseConnection() {
	err := cl.connection.Close()
	if err != nil {
		lib.PrintlnError(err)
	}
}

func (cl *Client) Receive() (msg []byte, err error) {
	msg = make([]byte, 1024000) // TODO Verificar uma forma de obter informações da conexão sem precisar setar tamanho do array de bytes
	n, err := cl.connection.Read(msg)
	if err != nil {
		return nil, err
	}

	return msg[:n], nil
}

func (cl *Client) Send(msg []byte) (err error) {
	_, err = cl.connection.Write(msg)
	if err != nil {
		return err
	}
	return nil
}
