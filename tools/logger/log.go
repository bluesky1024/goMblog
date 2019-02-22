package logger

import (
	"fmt"
)

//	"fmt"
//	"os"

func Info(log_type string, content string) {
	fmt.Println(log_type, "[Info] " + content)
}

func Err(log_type string, content string) {
	fmt.Println(log_type, "[Err]  " + content)
}

func Log(file_name string, log_type string, content string) {

}
