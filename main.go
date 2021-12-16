package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"time"
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

	// args := os.Args[1:]
	// argsCmd := strings.Fields(args[0])
	// //TODO find better way to handle this to handle "bash -c 'sleep 2; cat go.mod'"
	// //TODO also handle cases such as 								"bash -c \"sleep 2; cat go.mod\""
	// // cmd := exec.Command("bash", "-c", "sleep 2; cat go.mod")
	// // fmt.Printf("%#v", argsCmd)

	//test for shell var, else try other shells
	//TODO: flag for explicitly selecting shell
	shell, _ := exec.LookPath(os.Getenv("SHELL"))
	if shell == "" {
		shells := []string{"bash", "sh", "ash", "zsh", "fish"}
		for _, item := range shells {
			var err error
			shell, err = exec.LookPath(item)
			if err == nil {
				break
			}
		}
		log.Fatal("ERROR: no shell found.") //TODO: pretty colors here
	}

	// println(shell)
	// RunCmd(argsCmd[0], argsCmd[1:]...)
	// RunCmd(shell, "-c", "echo d$SHELL")
	// fmt.Println(os.Args[1])
	RunCmd(shell, "-c", os.Args[1])
	// `go run . 'bash -c "echo wow"'` works with this
}

// func main() {
// 	//regex process input args to commands

// 	testParse()
// 	// parseInput()
// 	//prepare commands to []exec.Cmd
// 	//prepare prefix decorators for commands
// 	//spawn goroutines for each instance of command
// 	//watch for SIGINT
// 	//send SIGINT to all []exec.cmd
// 	//wait for commands to exit and output status
// }

//THE PLAM

//formatter should be called "prefixer"

//flags:
//-color
//-nocolor
//-wrap wrap output instead of truncating to terminal width?
//-fullcmd show full command on output
//-namelen number of characters to show before truncating name of commands
//-shell pass specific shell to run
//-about print LICENSE (embedded)
//-help print this menu
