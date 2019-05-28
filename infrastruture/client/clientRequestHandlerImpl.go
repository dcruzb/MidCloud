package client

import (
	"fmt"
	"net"
	"os"
	"strconv"
)

type ClientRequestHandlerImpl struct {
	Host       string
	Port       int
	connection net.Conn
}

func NewClientRequestHandlerImpl(host string, port int) *ClientRequestHandlerImpl {
	address := host + ":" + strconv.Itoa(port)

	connection, err := net.Dial("tcp", address)

	if err != nil {
		fmt.Println(err)
		os.Exit(3)
	}

	return &ClientRequestHandlerImpl{Host: host, Port: port, connection: connection}
}

func (crh *ClientRequestHandlerImpl) Receive() (msg []byte, err error) {
	msg = make([]byte, 10240)
	n, err := crh.connection.Read(msg)
	if err != nil {
		return nil, err
	}

	return msg[:n], nil
}

func (crh *ClientRequestHandlerImpl) Send(msg []byte) (err error) {
	_, err = crh.connection.Write(msg)
	if err != nil {
		return err
	}
	return nil
}

func (crh *ClientRequestHandlerImpl) Close() error {
	return crh.connection.Close()
}
