package server

type ServerRequestHandler interface {
	Start() (err error)
	Stop() (err error)
	Receive() (msg []byte, err error)
	Send(msg []byte) (err error)
}
