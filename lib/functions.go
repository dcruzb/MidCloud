package lib

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
)

var SHOW_MESSAGES = []DebugLevel{ERROR} //, INFO, MESSAGE}

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
func Println(program string, messageLevel DebugLevel, message ...interface{}) {
	if len(SHOW_MESSAGES) > 0 {
		if inArrayDL(messageLevel, SHOW_MESSAGES) {
			switch messageLevel {
			case INFO:
				var logs []interface{}
				logs = append(logs, program, "- INFO -")
				logs = append(logs, message...)
				log.Println(logs...)
			case MESSAGE:
				fmt.Println(message...)
			case ERROR:
				_, file, line, ok := runtime.Caller(2)
				if !ok {
					file = "???"
					line = 0
				}

				log.Println(program, "\n          ***** ERROR *****",
					"\n          File:", file,
					"\n          Line:", strconv.Itoa(line),
					"\n          Message:\n               ", message)
			}
		}
	}
}

func PrintlnInfo(program string, message ...interface{}) {
	Println(program, INFO, message...)
}

func PrintlnMessage(program string, message ...interface{}) {
	Println(program, MESSAGE, message...)
}

func PrintlnError(program string, message ...interface{}) {
	Println(program, ERROR, message...)
}

func FailOnError(program string, err error, msg string) {
	if err != nil {
		Println(program, ERROR, msg, ":", err)
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
