package dist

import (
	"errors"
	"github.com/dcbCIn/MidCloud/infrastruture/client"
)

// Implements requestor
type RequestorImpl struct {
	crh client.ClientRequestHandlerImpl
}

func NewRequestorImpl(ipAddress string, portNumber int) *RequestorImpl {
	return &RequestorImpl{*client.NewClientRequestHandlerImpl(ipAddress, portNumber)}
}

func (r *RequestorImpl) Invoke(inv Invocation) (t Termination, err error) {
	requestHeader := RequestHeader{inv.IpAddress, inv.ObjectId, true, inv.ObjectId, inv.OperationName}
	requestBody := RequestBody{inv.Parameters}

	msg := Message{
		Header{"GIOP", 1, true, 0, 0},
		Body{requestHeader, requestBody, ReplyHeader{}, nil}}

	var bytes []byte
	bytes, err = Marshall(msg)
	if err != nil {
		return Termination{}, err
	}

	err = r.crh.Send(bytes)
	if err != nil {
		return Termination{}, err
	}

	var msgReturned Message
	msgToBeUnmarshalled, err := r.crh.Receive()
	if err != nil {
		return Termination{}, err
	}
	msgReturned, err = Unmarshall(msgToBeUnmarshalled)
	if err != nil {
		return Termination{}, err
	}

	if msgReturned.Body.ReplyHeader.ReplyStatus != 1 {
		// Todo identify errors by different ReplyStatus codes
		return Termination{}, errors.New("Server error while requesting remote operation. ")
	}

	//lib.PrintlnInfo("RequestorImpl", "RequestorImpl.Invoke - Reply recebido e unmarshalled")
	t = Termination{msgReturned.Body.ReplyBody}

	return t, err
}

func (r *RequestorImpl) Close() error {
	return r.crh.Close()
}
