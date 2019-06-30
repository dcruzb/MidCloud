package dist

import (
	"github.com/dcbCIn/MidCloud/infrastruture/server"
	"github.com/dcbCIn/MidCloud/lib"
	"reflect"
	"sync"
)

type InvokerImpl struct {
	remoteObjects map[int]interface{}
	srh           *server.ServerRequestHandlerImpl
}

func (inv *InvokerImpl) Register(objectId int, remoteObject interface{}) {
	if inv.remoteObjects == nil {
		inv.remoteObjects = make(map[int]interface{})
	}

	inv.remoteObjects[objectId] = remoteObject
}

func (inv *InvokerImpl) Invoke(port int, initialConnections int) (err error) {
	inv.srh, err = server.NewServerRequestHandlerImpl(port, initialConnections)
	if err != nil {
		return err
	}
	defer inv.srh.StopServer()
	lib.PrintlnInfo("InvokerImpl", "Invoker.invoke - Started to listen on port", port, "with", initialConnections, "connections")

	var wg = sync.WaitGroup{}

	for i := 0; i < initialConnections; i++ {

		wg.Add(1)
		go func(idx int) {
			inv.processConnection(idx)
			wg.Done()
		}(i)

	}

	wg.Wait()
	err = inv.srh.CloseConnection()
	if err != nil {
		return err
	}

	return nil
}

func (inv *InvokerImpl) processConnection(connectionIdx int) (err error) {

	for {
		cli, err := inv.srh.GetConnection(connectionIdx)
		if err != nil {
			lib.PrintlnError(err)
			return err
		}
		lib.PrintlnDebug("InvokerImpl", "Invoker.processConnection(", connectionIdx, ") - Connection established")

		for {
			msgToBeUnmarshalled, err := cli.Receive()
			if err != nil {
				if err.Error() == "EOF" {
					lib.PrintlnDebug("InvokerImpl", "Invoker.processConnection(", connectionIdx, ") - Connection has been closed!")
					break
				} else {
					lib.PrintlnDebug("InvokerImpl", "Invoker.processConnection(", connectionIdx, ") - Connection has been gracefully closed!")
					break
				}
			}
			lib.PrintlnDebug("InvokerImpl", "Invoker.processConnection(", connectionIdx, ") - Message received")

			msgReceived, err := Unmarshall(msgToBeUnmarshalled)
			if err != nil {
				return err
			}

			//			lib.PrintlnInfo("InvokerImpl", "Invoker.processConnection(",connectionIdx,") - Message unmarshalled")

			// Demultiplex
			remoteObject := inv.remoteObjects[msgReceived.Body.RequestHeader.ObjectKey]

			reflectedObject := reflect.ValueOf(remoteObject) //remoteObject.rcvr
			function := reflectedObject.MethodByName(msgReceived.Body.RequestHeader.Operation)
			functionType := function.Type()

			// Dispatch
			if functionType.NumIn() != len(msgReceived.Body.RequestBody.Parameters) {
				msgReceived.Body.ReplyHeader = ReplyHeader{"", msgReceived.Body.RequestHeader.RequestId, 0}
				msgReceived.Body.ReplyBody = nil
				lib.PrintlnError("Invoker.processConnection(", connectionIdx, ") Quantidade de parâmetros inválida para objeto remoto (", msgReceived.Body.RequestHeader.ObjectKey, ")/ operação (", msgReceived.Body.RequestHeader.Operation, ")")
			} else {
				args := make([]reflect.Value, functionType.NumIn())
				for i, parameter := range msgReceived.Body.RequestBody.Parameters {

					var arg reflect.Value
					switch reflect.TypeOf(parameter).Kind() {
					case reflect.Map:
						//reflectedArg := reflect.New( functionType.In(i) )
						//aux := &reflectedArg
						//
						//inter := aux.Elem() //.Interface() //reflectedArg.Elem().Interface() //.(reflectedArg.Type())
						inter := reflect.New(functionType.In(i))
						_, err := lib.Decode(parameter.(map[string]interface{}) /*reflect.TypeOf(common.ClientProxy{}) ,*/, &inter) //mapstructure.Decode(parameter, &inter)
						if err != nil {
							lib.PrintlnError("Invoker.processConnection(", connectionIdx, ") Erro ao realizar decode. Erro:", err)
						}

						//fmt.Println("Decode -", "structValue returned:", inter)
						//fmt.Println("Decode -", "structValue returned:", &inter)
						//retorno := &inter
						//fmt.Println("Decode -", "structValue returned:", retorno)
						//fmt.Println("Decode -", "structValue returned:", reflect.ValueOf(retorno).Elem())
						//fmt.Println("Decode -", "structValue returned:", reflect.ValueOf(retorno).Elem().Elem())
						//var retornoTipado common.ClientProxy
						//retornoTipado = reflect.ValueOf(retorno).Elem().Interface().(common.ClientProxy)
						//fmt.Println("Decode -", "structValue returned:",retornoTipado)
						//fmt.Println("Decode -", "structValue returned:", &retornoTipado)

						arg = inter.Elem() //inter.Elem()//reflect.ValueOf(inter) //inter.Addr().Elem()
						//lib.PrintlnInfo(arg) //par.(reflect.Value).Elem().Interface().(common.ClientProxy).Ip)
					default:
						arg = reflect.ValueOf(parameter)
					}

					args[i] = arg
				}

				msgReceived.Body.ReplyHeader = ReplyHeader{"", msgReceived.Body.RequestHeader.RequestId, 1}
				msgReceived.Body.ReplyBody = nil
				reflectedReturn := function.Call(args)

				returned := make([]interface{}, len(reflectedReturn))
				for i := 0; i < functionType.NumOut(); i++ {

					//reflectedArg := reflect.New(functionType.Out(i))
					//returned[i] = reflectedArg.Elem().Interface()
					returned[i] = reflectedReturn[i].Interface()

					/*var arg reflect.Value
					switch reflect.TypeOf(parameter).Kind() {
					case reflect.Map:
						reflectedArg := reflect.New(functionType.In(i))
						inter := reflectedArg.Elem().Interface()
						err = mapstructure.Decode(parameter, inter)
						arg = reflect.ValueOf(inter)
					default:
						arg = reflect.ValueOf(parameter)
					}

					args[i] = arg*/
				}
				//for i := 0; i <= functionType.NumOut(); i++ {
				msgReceived.Body.ReplyBody = returned
				msgReceived.Body.RequestBody.Parameters = []interface{}{}
				//}
			}

			var bytes []byte
			bytes, err = Marshall(msgReceived)
			if err != nil {
				return err
			}

			//			lib.PrintlnInfo("InvokerImpl", "Invoker.processConnection(",connectionIdx,") - Retorno marshalled")

			err = cli.Send(bytes)
			if err != nil {
				return err
			}

			lib.PrintlnDebug("InvokerImpl", "Invoker.processConnection(", connectionIdx, ") - Message sent")
		}
	}
}
