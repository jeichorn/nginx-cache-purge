// https://github.com/nathany/looper/blob/cb716912ee8c45628336f559d3e2fc8d398ae61b/print.go
package nginxcp

import (
	"fmt"

	"github.com/koyachi/go-term-ansicolor/ansicolor"
)

var DebugLevel int = 0

func Header() {
	fmt.Println(ansicolor.Cyan("Nginx-Cache-Purge 0.3.2 is watching for cache file changes"))
}

func DebugEnabled() {
	PrintInfo("Debug mode enabled.\n")
}

func PrintDebug(format string, a ...interface{}) {
    if (DebugLevel == 0) {
        return
    }
	msg := fmt.Sprintf(format, a...)
	fmt.Println(ansicolor.IntenseBlack(msg))
}

func PrintInfo(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	fmt.Println(ansicolor.IntenseBlack(msg))
}

func PrintTrace1(format string, a ...interface{}) {
    if (DebugLevel < 2) {
        return
    }
	msg := fmt.Sprintf(format, a...)
	fmt.Println(ansicolor.Black(msg))
}

func PrintTrace2(format string, a ...interface{}) {
    if (DebugLevel < 3) {
        return
    }
	msg := fmt.Sprintf(format, a...)
	fmt.Println(ansicolor.Yellow(msg))
}




func PrintError(msg error) {
	fmt.Println(ansicolor.IntenseRed(msg.Error()))
}
