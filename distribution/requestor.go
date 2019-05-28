package dist

type Invocation struct {
	ObjectId      int
	IpAddress     string
	PortNumber    int
	OperationName string
	Parameters    []interface{}
}

func NewInvocation(objectId int, ipAddress string, portNumber int, operationName string, parameters []interface{}) *Invocation {
	return &Invocation{ObjectId: objectId, IpAddress: ipAddress, PortNumber: portNumber, OperationName: operationName, Parameters: parameters}
}

type Termination struct {
	Result interface{}
}

type Requestor interface {
	Invoke(inv Invocation) (t Termination, err error)
	Close() error
}
