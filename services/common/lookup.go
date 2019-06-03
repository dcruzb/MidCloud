package common

import (
	"fmt"
	"github.com/dcbCIn/MidCloud/lib"
)

// Absolute Object Reference (AOR) implementation.
type ClientProxy struct {
	Ip   string
	Port int
	// Todo Add protocol to AOR if there is more than one
	//protocol string
	ObjectId int
}

// A record for the name server
type NamingRecord struct {
	ServiceName string
	ClientProxy ClientProxy
}

// Lookup Remoting Pattern for location transparency
type ILookup interface {
	Bind(sn string, cp ClientProxy) (err error)
	Lookup(serviceName string) (cp ClientProxy, err error)
	List() (services []NamingRecord, err error)
	Close() (err error)
}

type Lookup struct {
	services []NamingRecord
}

// Binds the name to the AOR.
// If already exists a service with the same name it is updated
func (l *Lookup) Bind(sn string, cp ClientProxy) (err error) {
	lib.PrintlnInfo("Lookup", "Service bind =", sn)
	for i, nr := range l.services {
		if nr.ServiceName == sn {
			l.services[i] = NamingRecord{sn, cp}
			return nil
		}
	}

	lib.PrintlnInfo("Lookup", "Service bind. sn:", sn, "CP:", cp.Ip, cp.Port, cp.ObjectId)
	l.services = append(l.services, NamingRecord{sn, cp})
	return nil
}

func (l Lookup) Lookup(serviceName string) (cp ClientProxy, err error) {
	lib.PrintlnInfo("Lookup", "Service lookup =", serviceName)
	for _, nr := range l.services {
		if nr.ServiceName == serviceName {
			lib.PrintlnInfo("Lookup", "Service found = ", serviceName, "(", nr.ClientProxy.ObjectId, ")")
			return nr.ClientProxy, nil
		}
	}
	lib.PrintlnInfo("Lookup", "Service not found = ", serviceName)
	return ClientProxy{}, nil
}

func (l Lookup) List() (services []NamingRecord, err error) {
	lib.PrintlnInfo("Lookup", "Service list (", len(l.services), ")")
	for _, nr := range l.services {
		fmt.Println(nr, nr.ServiceName, nr.ClientProxy)
	}
	return l.services, nil
}

func (l Lookup) Close() (err error) {
	return nil
}
