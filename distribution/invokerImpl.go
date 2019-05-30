package dist

import (
	"MidCloud/infrastruture/server"
	"MidCloud/lib"
	"reflect"
)

type InvokerImpl struct {
	remoteObjects map[int]interface{}
}

func (inv InvokerImpl) Register(objectId int, remoteObject interface{}) {
	inv.remoteObjects[objectId] = remoteObject
}

func (inv InvokerImpl) Invoke(port int) (err error) {
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

			reflectedObject := reflect.ValueOf(remoteObject)
			function := reflectedObject.MethodByName(msgReceived.Body.RequestHeader.Operation)
			functionType := function.Type()
			if functionType.NumIn() != len(msgReceived.Body.RequestBody.Parameters) {
				lib.PrintlnError("Quantidade de parâmetros inválida para objeto remoto (", msgReceived.Body.RequestHeader.ObjectKey, ")/ operação (", msgReceived.Body.RequestHeader.Operation, ")")
			}
			var args []reflect.Value
			for i, parameter := range msgReceived.Body.RequestBody.Parameters {
				args[i] = parameter.(reflect.Value) // Todo Adjust value type
			}

			msgReceived.Body.ReplyHeader = ReplyHeader{"", msgReceived.Body.RequestHeader.RequestId, 1}
			msgReceived.Body.ReplyBody = nil
			var returned []reflect.Value
			returned = function.Call(args)

			//for i := 0; i <= functionType.NumOut(); i++ {

			msgReceived.Body.ReplyBody = returned
			//}

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
