package server

import (
	"bufio"
	"github.com/dcbCIn/MidCloud/lib"
	"net"
	"os"
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
	srh.clients = make([]Client, initialConnections)

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

func (srh *ServerRequestHandlerImpl) WaitForConnection(cliIdx int) (cl *Client) { // TODO if cliIdx >= inicitalConnections => need to append to the slice
	conn, err := srh.listener.Accept()
	if err != nil {
		lib.PrintlnError("Error while waiting for connection", err)
	}

	cl = &srh.clients[cliIdx]

	cl.connection = conn

	return cl
}

func (cl *Client) CloseConnection() {
	err := cl.connection.Close()
	if err != nil {
		lib.PrintlnError(err)
	}
}

func (cl *Client) Read() (message string) {
	var err error
	// recebe solicitações do cliente
	message, err = bufio.NewReader(cl.connection).ReadString('\n')
	if err != nil {
		lib.PrintlnError("Error while reading message from socket TCP. Details:", err)
	}

	return message
}

func (cl *Client) Write(message string) {
	// envia resposta

	// Vários tipos diferentes de se escrever utilizando Writer, todos funcionam
	//_, err := fmt.Fprintf(conn, msgToServer+"\n")
	//_, err := conn.Write([]byte( msgToServer + "\n"))
	/*reader := bufio.NewWriter(conn)
	_, err := reader.WriteString( msgToServer + "\n")
	reader.Flush()*/
	/*reader := bufio.NewWriter(conn)
	_, err := io.WriteString(reader, msgToServer + "\n")
	reader.Flush()*/
	//_, err := io.WriteString(conn, msgToServer+"\n")

	_, err := cl.connection.Write([]byte(message + "\n"))
	if err != nil {
		lib.PrintlnError("Error while writing message to socket TCP. Details:", err)
		os.Exit(1)
	}
}
