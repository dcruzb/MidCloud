package dist

import (
	"github.com/dcbCIn/MidCloud/infrastruture/client"
	"github.com/dcbCIn/MidCloud/lib"
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

	// Todo check if replyStatus of the message is valid

	lib.PrintlnInfo("RequestorImpl", "RequestorImpl.Invoke - Reply recebido e unmarshalled")
	t = Termination{msgReturned.Body.ReplyBody}

	return t, err
}

func (r *RequestorImpl) Close() error {
	return r.crh.Close()
}
