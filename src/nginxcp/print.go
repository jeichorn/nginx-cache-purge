// https://github.com/nathany/looper/blob/cb716912ee8c45628336f559d3e2fc8d398ae61b/print.go
package nginxcp

import (
	"fmt"

	"github.com/koyachi/go-term-ansicolor/ansicolor"
)

func Header() {
	fmt.Println(ansicolor.Cyan("Nginx-Cache-Purge 0.3.2 is watching for cache file changes"))
}

func DebugEnabled() {
	DebugMessage("Debug mode enabled.\n")
}

func DebugMessage(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	fmt.Println(ansicolor.IntenseBlack(msg))
}

func DebugError(msg error) {
	fmt.Println(ansicolor.IntenseBlack(msg.Error()))
}

func DisplayHelp() {
}

func PrintWatching(folder string) {
	ClearPrompt()
	fmt.Println(ansicolor.Yellow("Watching path"), ansicolor.Yellow(folder))
}

func UnknownCommand(command string) {
	fmt.Println(ansicolor.Red("ERROR:")+" Unknown command", ansicolor.Magenta(command))
}

const CSI = "\x1b["

// remove from the screen anything that's been typed
// from github.com/kierdavis/ansi
func ClearPrompt() {
	fmt.Printf("%s2K", CSI)     // clear line
	fmt.Printf("%s%dG", CSI, 0) // go to column 0
}
