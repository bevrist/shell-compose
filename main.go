package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

func main() {
	args := os.Args[1:]
	argsCmd := strings.Fields(args[0])
	//TODO find better way to handle this to handle "bash -c 'sleep 2; cat go.mod'"
	//TODO also handle cases such as 								"bash -c \"sleep 2; cat go.mod\""
	// cmd := exec.Command("bash", "-c", "sleep 2; cat go.mod")
	fmt.Printf("%#v", argsCmd)
	// cmd := exec.Command(argsCmd[0], argsCmd[1:]...)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		cmd := exec.CommandContext(ctx, argsCmd[0], argsCmd[1:]...)
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err := cmd.Run()
		println(cmd.ProcessState.ExitCode())
		if err != nil {
			println(stderr.String())
		}
		fmt.Println(string(out.String()))
	}()

	time.Sleep(2 * time.Second)
	cancel() //cancel async running command
	//INFO context not needed

	// fmt.Println(formatter.NextColor() + "asdasd" + formatter.ResetColor() + "3123" + formatter.NextColor() + "asdasd" + formatter.ResetColor())
}
