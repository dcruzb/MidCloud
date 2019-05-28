package dist

import (
	"encoding/json"
)

func Marshall(message Message) ([]byte, error) {
	return json.Marshal(message)
}

func Unmarshall(marshaledMessage []byte) (message Message, err error) {
	err = json.Unmarshal(marshaledMessage, &message)
	return message, err
}
