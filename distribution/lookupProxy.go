package dist

import (
	"github.com/dcbCIn/MidCloud/lib"
	"github.com/dcbCIn/MidCloud/services/common"
	"github.com/mitchellh/mapstructure"
)

type LookupProxy struct {
	host      string
	port      int
	requestor Requestor
}

func NewLookupProxy(host string, port int) *LookupProxy {
	return &LookupProxy{host, port, NewRequestorImpl(host, port)}
}

func (lp LookupProxy) Bind(sn string, cp common.ClientProxy) (err error) {
	inv := *NewInvocation(0, lp.host, lp.port, lib.FunctionName(), []interface{}{sn, cp})
	_, err = lp.requestor.Invoke(inv)
	if err != nil {
		return err
	}
	return nil
}

func (lp LookupProxy) Lookup(serviceName string) (cp common.ClientProxy, err error) {
	inv := *NewInvocation(0, lp.host, lp.port, lib.FunctionName(), []interface{}{serviceName})
	termination, err := lp.requestor.Invoke(inv)
	if err != nil {
		return cp, err
	}

	err = mapstructure.Decode(termination.Result.([]interface{})[0], &cp)
	if err != nil {
		return cp, err
	}

	return cp, nil
}

func (lp LookupProxy) List() (services []common.NamingRecord, err error) {
	inv := *NewInvocation(0, lp.host, lp.port, lib.FunctionName(), []interface{}{})
	termination, err := lp.requestor.Invoke(inv)
	if err != nil {
		return nil, err
	}

	err = mapstructure.Decode(termination.Result.([]interface{})[0], &services) // TODO change Termination.Result -> add 2 attributes ( 1st - Returns, 2nd Error )
	if err != nil {
		return services, err
	}

	return services, nil
}

func (lp LookupProxy) Close() error {
	return lp.requestor.Close()
}
