package lib

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var SHOW_MESSAGES = []DebugLevel{ERROR, INFO, MESSAGE, DEBUG}

type DebugLevel int

const (
	ERROR   DebugLevel = 0
	INFO    DebugLevel = 1
	MESSAGE DebugLevel = 2
	DEBUG   DebugLevel = 3
)

func (d DebugLevel) ToInt() int {
	return [...]int{0, 1, 2, 3}[d]
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
			case DEBUG:
				var logs []interface{}
				logs = append(logs, file, "- DEBUG -")
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

func PrintlnDebug(message ...interface{}) {
	Println(DEBUG, message...)
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
func Decode(mapValue map[string]interface{}, structValue interface{}) (decoded interface{}, err error) {
	//structValue = reflect.New( typeOfStruct.(reflect.Type))
	reflectedStructTemp := structValue.(*reflect.Value) //reflect.ValueOf(structValue) //.Elem()
	var reflectedStruct reflect.Value
	if reflectedStructTemp.Kind() == reflect.Interface {
		reflectedStruct = reflectedStructTemp.Elem()
	} else if reflectedStructTemp.Kind() == reflect.Ptr {
		reflectedStruct = reflectedStructTemp.Elem()
	} // else {
	//reflectedStruct = reflectedStructTemp
	//}

	for k, v := range mapValue {
		//fmt.Println("Decode -", "field name:", k, "value:", v)
		//fmt.Println("canAddr:", reflectedStruct.CanAddr())
		field := reflectedStruct.FieldByName(k)

		switch field.Kind() {
		case reflect.Struct:
			if field.Type() == reflect.TypeOf(time.Time{}) {
				dtString := v.(string)
				if len(dtString) == 33 { // Pattern returned by unmarshall => Example: 2019-06-30T16:23:34.5836417-03:00
					//dtString[28] = "Z"
					dtString = dtString[:27] + "00" + dtString[27:] // Need to include 2 extra zeros to match the layout
				}
				// time.RFC3339Nano => Layout: 2006-01-02T15:04:05.999999999Z07:00
				dateTime, err := time.Parse(time.RFC3339Nano, dtString)
				if err != nil {
					return decoded, err
				}
				field.Set(reflect.ValueOf(dateTime))
			} else {
				_, err = Decode(v.(map[string]interface{}), &field)
				if err != nil {
					return decoded, err
				}
			}
		case reflect.String:
			field.SetString(v.(string))
		case reflect.Int, reflect.Int32, reflect.Int64:
			//fmt.Println(int64(v.(float64)))
			field.SetInt(int64(v.(float64)))
		case reflect.Float64:
			field.SetFloat(v.(float64))
		}
	}
	//fmt.Println("Decode -", "structValue returned:", structValue)
	//fmt.Println("Decode -", "structValue returned:", &structValue)
	//fmt.Println("Decode -", "structValue returned:", structValue.(*reflect.Value))
	//fmt.Println("Decode -", "structValue returned:", structValue.(*reflect.Value).Elem())
	//
	//retorno := &structValue
	//fmt.Println("Decode -", "structValue returned:", retorno)
	//fmt.Println("Decode -", "structValue returned:", reflect.ValueOf(retorno).Elem())
	//fmt.Println("Decode -", "structValue returned:", reflect.ValueOf(retorno).Elem().Elem())
	//fmt.Println("Decode -", "structValue returned:", &(structValue.(reflect.Value)))
	return structValue, nil

	// Todo adicionar verificação abaixo (se um parametro é passível de conversão para outro tipo)
	//argType := argValue.Type()
	//if argType.ConvertibleTo(inType) {
	//	in[i] = argValue.Convert(inType)
	//} else {
	//	return reflect.ValueOf(nil), fmt.Errorf("Method %s. Param[%d] must be %s. Have %s", name, i, inType, argType)
	//}
}
