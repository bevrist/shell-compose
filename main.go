package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/bevrist/shell-compose/proc"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	var terminating bool = false
	go func() {
		for range c {
			if !terminating {
				fmt.Println("Gracefully stopping... (press Ctrl+C again to force)")
				terminating = true
				continue
			}
			fmt.Println("ERROR: Aborting.")
			os.Exit(255)
		}
	}()
	time.Sleep(1 * time.Second)

	args := os.Args[1:]
	argsCmd := strings.Fields(args[0])
	//TODO find better way to handle this to handle "bash -c 'sleep 2; cat go.mod'"
	//TODO also handle cases such as 								"bash -c \"sleep 2; cat go.mod\""
	// cmd := exec.Command("bash", "-c", "sleep 2; cat go.mod")
	fmt.Printf("%#v", argsCmd)

	// proc.RunProcess(argsCmd[0], argsCmd[1:]...)
	proc.Run(argsCmd[0], argsCmd[1:]...)

	// fmt.Println(formatter.NextColor() + "asdasd" + formatter.ResetColor() + "3123" + formatter.NextColor() + "asdasd" + formatter.ResetColor())
}
