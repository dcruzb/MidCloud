package dist

import (
	"fmt"
	"github.com/dcbCIn/MidCloud/infrastruture/server"
	"github.com/dcbCIn/MidCloud/lib"
	"github.com/mitchellh/mapstructure"
	"reflect"
)

type InvokerImpl struct {
	remoteObjects map[int]interface{}
	//RemoteObject interface{}
}

func (inv *InvokerImpl) Register(objectId int, remoteObject interface{}) {
	if inv.remoteObjects == nil {
		inv.remoteObjects = make(map[int]interface{})
	}
	inv.remoteObjects[objectId] = remoteObject
	//inv.RemoteObject = remoteObject.(common.Lookup)

	fmt.Println(reflect.ValueOf(&remoteObject))
	fmt.Println(reflect.ValueOf(&remoteObject).Kind())
	fmt.Println(reflect.ValueOf(&remoteObject).Type())
	fmt.Println(reflect.ValueOf(&remoteObject).Elem())
	fmt.Println(reflect.ValueOf(&remoteObject).Elem().Kind())
	fmt.Println(reflect.ValueOf(&remoteObject).Elem().Type())
	fmt.Println(reflect.ValueOf(&remoteObject).Elem().Addr())
	fmt.Println(reflect.ValueOf(&remoteObject).Elem().Addr().Interface())
	fmt.Println("--------------")
	fmt.Println(reflect.ValueOf(remoteObject))
	fmt.Println(reflect.ValueOf(remoteObject).Kind())
	fmt.Println(reflect.ValueOf(remoteObject).Type())
	/*fmt.Println(reflect.ValueOf(remoteObject).Elem())
	fmt.Println(reflect.ValueOf(remoteObject).Elem().Kind())
	fmt.Println(reflect.ValueOf(remoteObject).Elem().Type())
	fmt.Println(reflect.ValueOf(remoteObject).Elem().Addr())
	fmt.Println(reflect.ValueOf(remoteObject).Elem().Addr().Interface())
	fmt.Println("--------------")

	fmt.Println(reflect.ValueOf(inv.RemoteObject))
	fmt.Println(reflect.ValueOf(inv.RemoteObject).Kind())
	fmt.Println(reflect.ValueOf(inv.RemoteObject).Type())
	fmt.Println(reflect.ValueOf(inv.RemoteObject).Elem())
	fmt.Println(reflect.ValueOf(inv.RemoteObject).Elem().Kind())*/
}

func (inv *InvokerImpl) Invoke(port int) (err error) {
	srh, err := server.NewServerRequestHandlerImpl(port)
	if err != nil {
		return err
	}
	defer srh.StopServer()
	lib.PrintlnInfo("InvokerImpl", "Invoker.invoke - conexão aberta")

	for {
		err = srh.Start()
		if err != nil {
			return err
		}

		for {
			lib.PrintlnInfo("InvokerImpl", "Invoker.invoke - Aguardando mensagem")

			msgToBeUnmarshalled, err := srh.Receive()
			if err != nil {
				if err.Error() == "EOF" {
					break
				} else {
					return err
				}
			}
			lib.PrintlnInfo("InvokerImpl", "Invoker.invoke - Mensagem recebida")

			msgReceived, err := Unmarshall(msgToBeUnmarshalled)

			if err != nil {
				return err
			}

			lib.PrintlnInfo("InvokerImpl", "Invoker.invoke - Mensagem unmarshalled")

			remoteObject := inv.remoteObjects[msgReceived.Body.RequestHeader.ObjectKey]

			//reflect.ValueOf(&remoteObject)
			//reflectedObject := reflect.ValueOf(&remoteObject)
			//functionType := reflectedObject.Type()
			//function, found := functionType.MethodByName(msgReceived.Body.RequestHeader.Operation)
			//function := reflectedObject.MethodByName()

			// Use Elem() only in pointers

			fmt.Println(reflect.ValueOf(remoteObject))
			fmt.Println(reflect.ValueOf(remoteObject).Kind())
			fmt.Println(reflect.ValueOf(remoteObject).Type())
			/*			fmt.Println(reflect.ValueOf(remoteObject).Elem())
						fmt.Println(reflect.ValueOf(remoteObject).Elem().Kind())
						fmt.Println(reflect.ValueOf(remoteObject).Elem().Type())
						fmt.Println(reflect.ValueOf(remoteObject).Elem().Addr())
						fmt.Println(reflect.ValueOf(remoteObject).Elem().Addr().Interface())
			*/
			/*fmt.Println(reflect.ValueOf(inv.RemoteObject))
			fmt.Println(reflect.ValueOf(inv.RemoteObject).Kind())
			fmt.Println(reflect.ValueOf(inv.RemoteObject).Type())
			*/
			//ro := reflect.New(reflect.TypeOf(remoteObject))
			reflectedObject := reflect.New(reflect.TypeOf(remoteObject)) //reflect.ValueOf( &ro ) //inv.RemoteObject)
			function := reflectedObject.MethodByName(msgReceived.Body.RequestHeader.Operation)
			functionType := function.Type()

			//found := true // Todo tirar
			validMessage := true
			/*if !found {
				msgReceived.Body.ReplyHeader = ReplyHeader{"", msgReceived.Body.RequestHeader.RequestId, 0}
				msgReceived.Body.ReplyBody = nil
				lib.PrintlnError("Operação inválida para objeto remoto (", msgReceived.Body.RequestHeader.ObjectKey, ")/ operação (", msgReceived.Body.RequestHeader.Operation, ")")
				validMessage = false
			} else */
			if functionType.NumIn() != len(msgReceived.Body.RequestBody.Parameters) {
				msgReceived.Body.ReplyHeader = ReplyHeader{"", msgReceived.Body.RequestHeader.RequestId, 0}
				msgReceived.Body.ReplyBody = nil
				lib.PrintlnError("Quantidade de parâmetros inválida para objeto remoto (", msgReceived.Body.RequestHeader.ObjectKey, ")/ operação (", msgReceived.Body.RequestHeader.Operation, ")")
				validMessage = false
			}

			if validMessage {
				args := make([]reflect.Value, functionType.NumIn())
				for i, parameter := range msgReceived.Body.RequestBody.Parameters {

					var arg reflect.Value
					switch reflect.TypeOf(parameter).Kind() {
					case reflect.String:
						arg = reflect.ValueOf(parameter)
					default:
						reflectedArg := reflect.New(functionType.In(i))
						inter := reflectedArg.Elem().Interface()
						err = mapstructure.Decode(parameter, inter)
						arg = reflect.ValueOf(inter)
					}

					args[i] = arg

					//args[i] =  // Todo Adjust value type
				}

				msgReceived.Body.ReplyHeader = ReplyHeader{"", msgReceived.Body.RequestHeader.RequestId, 1}
				msgReceived.Body.ReplyBody = nil
				var returned []reflect.Value
				returned = function.Call(args)

				//for i := 0; i <= functionType.NumOut(); i++ {

				msgReceived.Body.ReplyBody = returned
				//}
			}

			var bytes []byte
			bytes, err = Marshall(msgReceived)
			if err != nil {
				return err
			}

			lib.PrintlnInfo("InvokerImpl", "Invoker.invoke - Retorno marshalled")

			err = srh.Send(bytes)
			if err != nil {
				return err
			}

			lib.PrintlnInfo("InvokerImpl", "Invoker.invoke - Mensagem enviada")
		}
	}

	err = srh.CloseConnection()
	if err != nil {
		return err
	}

	return nil
}
