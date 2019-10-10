package dist

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestMarshall(t *testing.T) {
	type args struct {
		message Message
	}
	message1 := Message{
		Header{"GIOP", 1, true, 0, 0},
		Body{RequestHeader{}, RequestBody{}, ReplyHeader{}, nil}}

	message2 := Message{
		Header{"GIOP", 2, true, 0, 0},
		Body{RequestHeader{}, RequestBody{}, ReplyHeader{}, nil}}

	jsonBytes :=
		func(message Message) []byte {
			bytes, _ := json.Marshal(message)
			return bytes
		}

	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{"Marshaller Test 1",
			args{message1},
			jsonBytes(message1),
			false,
		},

		{"Marshaller Test 2",
			args{message2},
			jsonBytes(message2),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Marshall(tt.args.message)
			if (err != nil) != tt.wantErr {
				t.Errorf("Marshall() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Marshall() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnmarshall(t *testing.T) {
	type args struct {
		marshaledMessage []byte
	}

	marshaledMessage1, _ := json.Marshal(Message{
		Header{"GIOP", 1, true, 0, 0},
		Body{RequestHeader{}, RequestBody{}, ReplyHeader{}, nil}})

	marshaledMessage2, _ := json.Marshal(Message{})

	message :=
		func(marshaledMessage []byte) Message {
		var message Message
			_ = json.Unmarshal(marshaledMessage, &message)
			return message
		}

	tests := []struct {
		name        string
		args        args
		wantMessage Message
		wantErr     bool
	}{
		{"UnMarshaller Test 1",
			args{marshaledMessage1},
			message(marshaledMessage1),
			false,
		},

		{"UnMarshaller Test 2",
			args{[]byte {'2','3'}},
			message(marshaledMessage2),
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMessage, err := Unmarshall(tt.args.marshaledMessage)
			if (err != nil) != tt.wantErr {
				t.Errorf("Unmarshall() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotMessage, tt.wantMessage) {
				t.Errorf("Unmarshall() = %v, want %v", gotMessage, tt.wantMessage)
			}
		})
	}
}
