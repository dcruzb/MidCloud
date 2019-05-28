package dist

import (
	"github.com/mitchellh/mapstructure"
	"middleware/app/server/remoteObjects"
	"middleware/lib"
	"middleware/lib/infra/server"
	"middleware/lib/services/common"
)

type InvokerImpl struct {
}

func (inv InvokerImpl) Invoke(port int) (err error) {
	srh, err := server.NewServerRequestHandlerImpl(port)
	if err != nil {
		return err
	}
	defer srh.StopServer()
	lib.PrintlnInfo("InvokerImpl", "Invoker.invoke - conexão aberta")

	var lookup = common.Lookup{}
	var jankenpo = remoteObjects.Jankenpo{}

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

			switch msgReceived.Body.RequestHeader.Operation { // Todo add the objectId to demultiplex
			case "Play":
				player1Move := msgReceived.Body.RequestBody.Parameters[0].(string)
				player2Move := msgReceived.Body.RequestBody.Parameters[1].(string)
				msgReceived.Body.ReplyHeader = ReplyHeader{"", msgReceived.Body.RequestHeader.RequestId, 1}
				msgReceived.Body.ReplyBody, _ = jankenpo.Play(player1Move, player2Move)
			case "Bind":
				serviceName := msgReceived.Body.RequestBody.Parameters[0].(string)
				var clientProxy common.ClientProxy
				err := mapstructure.Decode(msgReceived.Body.RequestBody.Parameters[1], &clientProxy)
				if err != nil {
					lib.PrintlnError("InvokerImpl", err)
				}
				//clientProxyMap := msgReceived.Body.RequestBody.Parameters[1].(map[string]interface{})
				//clientProxy := common.ClientProxy{clientProxyMap["Ip"].(string), int(clientProxyMap["Port"].(float64)), int(clientProxyMap["ObjectId"].(float64))}
				msgReceived.Body.ReplyHeader = ReplyHeader{"", msgReceived.Body.RequestHeader.RequestId, 1}
				msgReceived.Body.ReplyBody = lookup.Bind(serviceName, clientProxy)
			case "Lookup":
				serviceName := msgReceived.Body.RequestBody.Parameters[0].(string)
				msgReceived.Body.ReplyHeader = ReplyHeader{"", msgReceived.Body.RequestHeader.RequestId, 1}
				msgReceived.Body.ReplyBody, _ = lookup.Lookup(serviceName)
			default:
				msgReceived.Body.ReplyHeader = ReplyHeader{"", msgReceived.Body.RequestHeader.RequestId, 0}
				msgReceived.Body.ReplyBody = nil
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
