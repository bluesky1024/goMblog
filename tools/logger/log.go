package logger

import (
	"fmt"
	"runtime"
)

func Info(log_type string, content interface{}) {
	pc, _, _, _ := runtime.Caller(1)
	fmt.Println(log_type, "[Info]", runtime.FuncForPC(pc).Name(), content)
}

func Err(log_type string, content interface{}) {
	pc, _, _, _ := runtime.Caller(1)
	fmt.Println(log_type, "[Err]", runtime.FuncForPC(pc).Name(), content)
}

func Log(file_name string, log_type string, content string) {

}
