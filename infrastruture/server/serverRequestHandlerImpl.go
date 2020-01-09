package server

import (
	"errors"
	"github.com/dcbCIn/MidCloud/lib"
	"net"
	"strconv"
)

// Client identifies the Client connected to the ServerRequestHandler
type Client struct {
	connection net.Conn
}

// ServerRequestHandlerImpl implementation of ServerRequestHandler.
// Start a server and wait for connections
type ServerRequestHandlerImpl struct {
	Port       int
	listener   net.Listener
	connection net.Conn
	clients    []Client
}

// NewServerRequestHandlerImpl creates a ServerRequestHandlerImpl
func NewServerRequestHandlerImpl(port int, initialConnections int) (srh *ServerRequestHandlerImpl, err error) {
	srh = &ServerRequestHandlerImpl{Port: port}
	srh.clients = make([]Client, initialConnections, initialConnections)

	srh.listener, err = net.Listen("tcp", ":"+strconv.Itoa(srh.Port))
	if err != nil {
		return nil, err
	}
	return srh, nil
}

// Start wait for new connection (one connection based implementation)
func (srh *ServerRequestHandlerImpl) Start() (err error) {
	lib.PrintlnInfo("ServerRequestHandler", "Aceitando conexões...")

	srh.connection, err = srh.listener.Accept()
	if err != nil {
		lib.PrintlnInfo("ServerRequestHandler", "Erro ao abrir conexão")
		return err
	}

	lib.PrintlnInfo("ServerRequestHandler", "Conexão aceita...")
	return nil
}

// CloseConnection closes the server connection (one connection based implementation)
func (srh *ServerRequestHandlerImpl) CloseConnection() (err error) {
	lib.PrintlnInfo("ServerRequestHandler", "ServerRequestHandler.Stop - Closing connection")
	err = srh.connection.Close()
	if err != nil {
		return err
	}
	lib.PrintlnInfo("ServerRequestHandler", "ServerRequestHandler.Stop - Connection closed")
	return nil
}

// StopServer close the server listener (one connection based implementation)
func (srh *ServerRequestHandlerImpl) StopServer() (err error) {
	err = srh.listener.Close()
	if err != nil {
		return err
	}
	lib.PrintlnInfo("ServerRequestHandler", "ServerRequestHandler.Stop - Listener closed")
	return nil
}

// Receive receives a message in sent direct to the server (one connection based implementation)
func (srh *ServerRequestHandlerImpl) Receive() (msg []byte, err error) {
	msg = make([]byte, 10240)
	n, err := srh.connection.Read(msg)
	if err != nil {
		return nil, err
	}

	return msg[:n], nil
}

// Send sends a message from the server (one connection based implementation)
func (srh *ServerRequestHandlerImpl) Send(msg []byte) (err error) {
	_, err = srh.connection.Write(msg)
	if err != nil {
		return err
	}
	return nil
}

//------------------------------

// GetConnection wait for new connection in the client cliIdx (multiple connections based implementation)
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

// CloseConnection closes the client connection (multiple connections based implementation)
func (cl *Client) CloseConnection() {
	err := cl.connection.Close()
	if err != nil {
		lib.PrintlnError(err)
	}
}

// Receive receives a message from a client (multiple connections based implementation)
func (cl *Client) Receive() (msg []byte, err error) {
	msg = make([]byte, 1024000) // TODO Verificar uma forma de obter informações da conexão sem precisar setar tamanho do array de bytes
	n, err := cl.connection.Read(msg)
	if err != nil {
		return nil, err
	}

	return msg[:n], nil
}

// Send send a message to a client (multiple connections based implementation)
func (cl *Client) Send(msg []byte) (err error) {
	_, err = cl.connection.Write(msg)
	if err != nil {
		return err
	}
	return nil
}
