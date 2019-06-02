package common

import (
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

// A record for the name server, to
type NamingRecord struct {
	serviceName string
	clientProxy ClientProxy
}

// Lookup Remoting Pattern for location transparency
type ILookup interface {
	Bind(sn string, cp ClientProxy) (err error)
	Lookup(serviceName string) (cp ClientProxy, err error)
	List() []NamingRecord
}

type Lookup struct {
	services []NamingRecord
}

// Binds the name to the AOR.
// If already exists a service with the same name it is updated
func (l *Lookup) Bind(sn string, cp ClientProxy) (err error) {
	lib.PrintlnInfo("Lookup", "Service bind =", sn)
	for i, nr := range l.services {
		if nr.serviceName == sn {
			l.services[i] = NamingRecord{sn, cp}
			return nil
		}
	}

	l.services = append(l.services, NamingRecord{sn, cp})
	return nil
}

func (l Lookup) Lookup(serviceName string) (cp ClientProxy, err error) {
	lib.PrintlnInfo("Lookup", "Service lookup =", serviceName)
	for _, nr := range l.services {
		if nr.serviceName == serviceName {
			return nr.clientProxy, nil
		}
	}
	return ClientProxy{}, nil
}

func (l Lookup) List() []NamingRecord {
	lib.PrintlnInfo("Lookup", "Service list (", len(l.services), ")")
	return l.services
}
