// https://github.com/nathany/looper/blob/cb716912ee8c45628336f559d3e2fc8d398ae61b/print.go
package nginxcp

import (
	"fmt"
    "log"

	"github.com/koyachi/go-term-ansicolor/ansicolor"
)

var DebugLevel int = 0
var AnsiOutput bool =true 

func Header() {
    PrintInfo("Nginx-Cache-Purge 0.3.2 is watching for cache file changes")
}

func DebugEnabled() {
	PrintInfo("Debug mode enabled level: %d", DebugLevel)
}

func PrintDebug(format string, a ...interface{}) {
    if (DebugLevel == 0) {
        return
    }
	var msg string = fmt.Sprintf(format, a...)
    if (AnsiOutput) {
        msg = ansicolor.IntenseBlack(msg)
    }
	log.Println(msg)
}

func PrintInfo(format string, a ...interface{}) {
	var msg string = fmt.Sprintf(format, a...)
    if (AnsiOutput) {
        msg = ansicolor.Green(msg)
    }
	log.Println(msg)
}

func PrintTrace1(format string, a ...interface{}) {
    if (DebugLevel < 2) {
        return
    }
	var msg string = fmt.Sprintf(format, a...)
    if (AnsiOutput) {
        msg = ansicolor.Black(msg)
    }
	log.Println(msg)
}

func PrintTrace2(format string, a ...interface{}) {
    if (DebugLevel < 3) {
        return
    }
	var msg = fmt.Sprintf(format, a...)
    if (AnsiOutput) {
        msg = ansicolor.Yellow(msg)
    }
	log.Println(msg)
}




func PrintError(msg error) {
    var out string = msg.Error()
    if (AnsiOutput) {
        out = ansicolor.IntenseRed(out)
    }
	log.Println(out)
}
