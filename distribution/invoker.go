package dist

type Invoker interface {
	Invoke(port int) (err error)
	//StopServer()
}
