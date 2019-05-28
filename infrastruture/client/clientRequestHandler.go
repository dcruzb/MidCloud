package client

type ClientRequestHandler interface {
	Send(msg []byte) (err error)
	Receive() (msg []byte, err error)
	Close() error
}
