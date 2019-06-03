package lib

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"
)

var SHOW_MESSAGES = []DebugLevel{ERROR, INFO, MESSAGE}

type DebugLevel int

const (
	ERROR   DebugLevel = 0
	INFO    DebugLevel = 1
	MESSAGE DebugLevel = 2
)

func (d DebugLevel) ToInt() int {
	return [...]int{0, 1, 2}[d]
}

func FunctionName() string {
	pc, _, _, _ := runtime.Caller(1)

	name := strings.Split(runtime.FuncForPC(pc).Name(), ".")

	return name[len(name)-1]
}

func Println(messageLevel DebugLevel, message ...interface{}) {
	if len(SHOW_MESSAGES) > 0 {
		if inArrayDL(messageLevel, SHOW_MESSAGES) {
			_, file, line, ok := runtime.Caller(2)
			if !ok {
				file = "???"
				line = 0
			}

			switch messageLevel {
			case INFO:
				var logs []interface{}
				logs = append(logs, file, "- INFO -")
				logs = append(logs, message...)
				log.Println(logs...)
			case MESSAGE:
				fmt.Println(message...)
			case ERROR:
				log.Println(file, "\n          ***** ERROR *****",
					"\n          File:", file,
					"\n          Line:", strconv.Itoa(line),
					"\n          Message:\n               ", message)
			}
		}
	}
}

func PrintlnInfo(message ...interface{}) {
	Println(INFO, message...)
}

func PrintlnMessage(message ...interface{}) {
	Println(MESSAGE, message...)
}

func PrintlnError(message ...interface{}) {
	Println(ERROR, message...)
}

func FailOnError(err error, msg string) {
	if err != nil {
		Println(ERROR, msg, ":", err)
		os.Exit(1)
	}
}

func InArray(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func inArrayDL(a DebugLevel, list []DebugLevel) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// Decode expects that mapValue is actually the same Type as structValue
func Decode(mapValue map[string]interface{}, structValue interface{}) (err error) {
	reflectedStructTemp := reflect.ValueOf(structValue).Elem()
	var reflectedStruct reflect.Value
	if reflectedStructTemp.Kind() == reflect.Interface {
		reflectedStruct = reflectedStructTemp.Elem()
	} else {
		reflectedStruct = reflectedStructTemp
	}

	for k, v := range mapValue {
		fmt.Println("Decode -", "field name:", k, "value:", v)
		field := reflectedStruct.FieldByName(k)

		switch field.Kind() {
		case reflect.Struct:
			Decode(v.(map[string]interface{}), &field)
		case reflect.String:
			field.SetString(v.(string))
		case reflect.Int, reflect.Int32, reflect.Int64:
			fmt.Println(int64(v.(float64)))
			field.SetInt(int64(v.(float64)))
		case reflect.Float64:
			field.SetFloat(v.(float64))
		}
	}
	return nil
}
