package common

import (
	"middleware/lib"
)

type ClientProxy struct {
	Ip   string
	Port int
	// Todo Add protocol to AOR if there is more than one
	//protocol string
	ObjectId int
}

type NamingRecord struct {
	serviceName string
	clientProxy ClientProxy
}

type ILookup interface {
	Bind(sn string, cp ClientProxy) (err error)
	Lookup(serviceName string) (cp ClientProxy, err error)
}

type Lookup struct {
	services []NamingRecord
}

func (l *Lookup) Bind(sn string, cp ClientProxy) (err error) {
	lib.PrintlnInfo("Lookup", "Service bind =", sn)
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
