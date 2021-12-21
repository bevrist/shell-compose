package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/pflag"
)

// func main() {
// 	c := make(chan os.Signal, 1)
// 	signal.Notify(c, os.Interrupt)
// 	go func() {
// 		var terminating bool = false
// 		for range c {
// 			if !terminating {
// 				fmt.Println("Gracefully stopping... (press Ctrl+C again to force)")
// 				terminating = true
// 				//TODO send sigint to processes
// 				continue
// 			}
// 			fmt.Println("ERROR: Aborting.")
// 			os.Exit(255)
// 		}
// 	}()
// 	time.Sleep(1 * time.Second)

// 	// args := os.Args[1:]
// 	// argsCmd := strings.Fields(args[0])
// 	// //TODO find better way to handle this to handle "bash -c 'sleep 2; cat go.mod'"
// 	// //TODO also handle cases such as 								"bash -c \"sleep 2; cat go.mod\""
// 	// // cmd := exec.Command("bash", "-c", "sleep 2; cat go.mod")
// 	// // fmt.Printf("%#v", argsCmd)

// //test for shell var, else try other shells
// //TODO: flag for explicitly selecting shell
// shell, _ := exec.LookPath(os.Getenv("SHELL"))
// if shell == "" {
// 	shells := []string{"bash", "sh", "ash", "zsh", "fish"}
// 	for _, item := range shells {
// 		var err error
// 		shell, err = exec.LookPath(item)
// 		if err == nil {
// 			break
// 		}
// 	}
// 	log.Fatal("ERROR: no shell found.") //TODO: pretty colors here
// }

// 	// println(shell)
// 	// RunCmd(argsCmd[0], argsCmd[1:]...)
// 	// RunCmd(shell, "-c", "echo d$SHELL")
// 	// fmt.Println(os.Args[1])
// 	RunCmd(shell, "-c", os.Args[1])
// 	// `go run . 'bash -c "echo wow"'` works with this
// }

func proc(cmd *exec.Cmd) {
	var out, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &stderr
	var err error
	fmt.Println("starting...")
	cmd.Start()
	go func() {
		for {
			if out.Len() != 0 {
				outs := strings.Split(string(out.Next(9999)), "\n") //FIXME this could probably be a .String()
				for _, line := range outs {
					if line == "" {
						continue
					}
					fmt.Println(PrintCmdName("test") + line)
				}
			}
		}
	}()
	err = cmd.Wait()
	if err != nil {
		println("err: " + stderr.String())
	}
	println(cmd.ProcessState.ExitCode())
}

func RunCmd(command string, args ...string) {
	cmd := exec.Command(command, args[0:]...)
	go proc(cmd)
	time.Sleep(3 * time.Second)
	cmd.Process.Signal(syscall.SIGINT)
	fmt.Println("SIGINT Sent to process...")
	time.Sleep(2 * time.Second)
}

// func getCommands() []string {
// 	for _, arg := range os.Args {

// 	}
// }

func main() {
	//input flags
	// fColor := pflag.BoolP("color", "c", false, "enable color for output")
	// fNoColor := pflag.BoolP("nocolor", "n", false, "disable color for output")
	fShell := pflag.StringP("shell", "s", "", "shell binary to launch commands with")
	pflag.Parse()

	//test for shell var, else try other shells
	var shell string
	if *fShell != "" {
		shell, _ = exec.LookPath(*fShell)
	} else {
		shell, _ = exec.LookPath(os.Getenv("SHELL"))
	}
	if shell == "" {
		shells := []string{"bash", "sh", "zsh", "ash", "fish"}
		for _, item := range shells {
			var err error
			shell, err = exec.LookPath(item)
			if err == nil {
				break
			}
		}
		log.Fatal("ERROR: no shell found.") //TODO: pretty colors here
	}

	// extract commands
	// testParse()
	// parseInput()
	//prepare commands to []exec.Cmd
	//prepare prefix decorators for commands
	//spawn goroutines for each instance of command
	//watch for SIGINT
	//send SIGINT to all []exec.cmd
	//wait for commands to exit and output status
}

//THE PLAM

//formatter should be called "prefixer"

//flags:
//-wrap wrap output instead of truncating to terminal width?
//-fullcmd show full command on output
//-namelen number of characters to show before truncating name of commands
//-about print LICENSE (embedded)
//-help print this menu
