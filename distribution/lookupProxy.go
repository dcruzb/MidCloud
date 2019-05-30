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

	err = mapstructure.Decode(termination.Result, &cp)
	if err != nil {
		return cp, err
	}

	return cp, nil
}

func (lp LookupProxy) Close() error {
	return lp.requestor.Close()
}
