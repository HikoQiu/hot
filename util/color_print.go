package util

import "fmt"

const (
	COLOR_SUCC = "succ"
	COLOR_FAIL = "fail"
	COLOR_WARNING = "warning"
	COLOR_INFO = "info"
)

func Colorize(text string, status string) string {
	out := ""
	switch status {
	case COLOR_SUCC:
		out = "\033[32;1m"    // Green
	case COLOR_FAIL:
		out = "\033[31;1m"    // Red
	case COLOR_INFO:
		out = "\033[34;1m"    // Blue
	case COLOR_WARNING:
		out = "\033[33;1m"    // Yellow
	default:
		out = "\033[0m"    // Default
	}
	return out + text + "\033[0m"
}

func ColorPrintln(str string, color string) {
	fmt.Println(Colorize(str, color))
}
