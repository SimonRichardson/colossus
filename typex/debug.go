package typex

import (
	"fmt"
	"runtime"
)

func PrintStack(all bool) {
	buf := make([]byte, 1<<16)
	runtime.Stack(buf, all)
	fmt.Printf("\n%s\n", string(buf))
}
