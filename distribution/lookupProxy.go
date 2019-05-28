package dist

import (
	"github.com/mitchellh/mapstructure"
	"middleware/lib"
	"middleware/lib/services/common"
)

type LookupProxy struct {
	Host      string
	Port      int
	requestor Requestor
}

func NewLookupProxy(host string, port int) *LookupProxy {
	return &LookupProxy{host, port, NewRequestorImpl(host, port)}
}

func (lp LookupProxy) Bind(sn string, cp common.ClientProxy) (err error) {
	inv := *NewInvocation(0, lp.Host, lp.Port, lib.FunctionName(), []interface{}{sn, cp})
	_, err = lp.requestor.Invoke(inv)
	if err != nil {
		return err
	}
	return nil
}

func (lp LookupProxy) Lookup(serviceName string) (cp common.ClientProxy, err error) {
	inv := *NewInvocation(0, lp.Host, lp.Port, lib.FunctionName(), []interface{}{serviceName})
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
