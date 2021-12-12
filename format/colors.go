package format

import (
	"os"
)

var colors = []string{
	"32", //green
	"33", //yellow
	"34", //blue
	"35", //magenta
	"92", //light green
	"93", //light yellow
	"94", //light blue
	"95", //light magenta
	"96", //light cyan
}

//returns true if ran from a tty
func istty() bool {
	if fileInfo, _ := os.Stdout.Stat(); (fileInfo.Mode() & os.ModeCharDevice) != 0 {
		return true
	}
	return false
}

// return true if should be color
func checkColor() bool {
	//TODO add command line flag for colors
	if istty() {
		return true
	}
	return false
}

var currColor int

//NextColor returns a new color each call (based on tty or by flag)
func NextColor() string {
	if checkColor() {
		currColor = (currColor + 1) % len(colors)
		return "\033[" + colors[currColor] + "m"
	}
	return "" //return no formatting if not tty or color flag disabled
}

//ResetColor returns color reset code (based on tty or by flag)
func ResetColor() string {
	if checkColor() {
		return "\033[0m"
	}
	return ""
}

//ErrorColor returns an error color
func ErrorColor() string {
	if checkColor() {
		return "\033[31m" //red
	}
	return ""
}
