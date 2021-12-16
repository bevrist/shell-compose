package main

func PrintCmdName(command string) string {
	return NextColor() + command + " | " + ResetColor()
}
