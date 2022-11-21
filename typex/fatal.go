package typex

import (
	"fmt"
	"log"
	"runtime"
)

var (
	ShowStack = true
)

func Fatal(v ...interface{}) {
	if !ShowStack {
		log.Fatal(v...)
	}

	var (
		message = fmt.Sprint(v...)
		buf     = make([]byte, 1<<16)
	)

	runtime.Stack(buf, false)

	log.Fatalf("\n%s\n%s\n", message, string(buf))
}

func Fatalf(format string, v ...interface{}) {
	if !ShowStack {
		log.Fatalf(format, v...)
	}

	var (
		message = fmt.Sprintf(format, v...)
		buf     = make([]byte, 1<<16)
	)

	runtime.Stack(buf, false)

	log.Fatalf("\n%s\n%s\n", message, string(buf))
}
