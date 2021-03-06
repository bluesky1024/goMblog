package logger

import (
	"fmt"
	"runtime"
	"time"
)

func Info(log_type string, content interface{}) {
	pc, _, line, _ := runtime.Caller(1)
	fmt.Println("[Info]", log_type, time.Now(), "| File.Func:", runtime.FuncForPC(pc).Name(), "| file line:", line, "| content:", content)
}

func Err(log_type string, content interface{}) {
	pc, _, line, _ := runtime.Caller(1)
	fmt.Println("[Err]", log_type, time.Now(), "| File.Func:", runtime.FuncForPC(pc).Name(), "| file line:", line, "| content:", content)
}

func Log(file_name string, log_type string, content string) {

}
