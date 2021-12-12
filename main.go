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
	go func() {
		var terminating bool = false
		for range c {
			if !terminating {
				fmt.Println("Gracefully stopping... (press Ctrl+C again to force)")
				terminating = true
				//TODO send sigint to processes
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
}

//THE PLAM
//main
//regex process input args to commands
//prepare commands to []exec.Cmd
//prepare prefix decorators for commands
//spawn goroutines for each instance of command
//watch for SIGINT
//send SIGINT to all []exec.cmd

//formatter should be called "prefixer"

//flags:
//-color
//-nocolor
//-wrap wrap output instead of truncating to terminal width?
//-fullcmd show full command on output
//-namelen number of characters to show before truncating name of commands

















