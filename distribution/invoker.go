package dist

type Invoker interface {
	Register(objectId int, remoteObject interface{})
	Invoke(port int) (err error)
	//StopServer()
}
