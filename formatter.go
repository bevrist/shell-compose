package main

func PrintCmdName(command string, color string) string {
	return color + command + " | " + ResetColor()
}
