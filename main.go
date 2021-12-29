package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"sync"
	"syscall"

	"github.com/spf13/pflag"
)

//proc handles output and lifecycle of commands
func proc(cmd *exec.Cmd, title string, color string) {
	//format title to be consistent length
	tlen := "5"
	tfmt := "%-" + tlen + "." + tlen + "s"
	title = fmt.Sprintf(tfmt, title) //TODO: make this more complicated and prettier
	//capture command output streams and start command
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	outReader := bufio.NewReader(stdout)
	errReader := bufio.NewReader(stderr)
	fmt.Println("starting..." + color + " \"" + cmd.String() + "\"" + ResetColor())
	cmd.Start()
	//process and print command output as it arrives
	go func() {
		for {
			//stdout
			outline, err := outReader.ReadString('\n')
			for err == nil {
				//skip printing empty lines
				if reEmpty.FindStringIndex(outline) != nil {
					continue
				}
				fmt.Print(PrintCmdName(title, color) + outline)
				outline, err = outReader.ReadString('\n')
			}

			//stderr
			errline, err := errReader.ReadString('\n')
			for err == nil {
				//skip printing empty lines
				if reEmpty.FindStringIndex(errline) != nil {
					continue
				}
				fmt.Print(PrintCmdName(title, color) + "stderr: " + errline) //TODO error color this
				errline, err = outReader.ReadString('\n')
			}
		}
	}()
	//keep goroutine running as long as command is running
	cmd.Wait()
	fmt.Println(PrintCmdName(title, color) + "Process Exited with code: " + fmt.Sprint(cmd.ProcessState.ExitCode()))
}

var (
	//arg flags
	fColor   = pflag.BoolP("color", "c", false, "enable color for output")
	fNoColor = pflag.BoolP("nocolor", "n", false, "disable color for output")
	fRestart = pflag.BoolP("restart", "r", false, "restart commands after failure (non zero exit code)")
	fShell   = pflag.StringP("shell", "s", "", "shell binary to launch commands with")
	//regex to capture all empty lines
	reEmpty = regexp.MustCompile(`^\s*$`)
)

func main() {
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

	//make list of command objects
	//use a subshell for each command for simplicity
	//execute and handle each command in own thread
	var wg sync.WaitGroup
	cmds := make([]*exec.Cmd, len(pflag.Args()))
	for i, cmdStr := range pflag.Args() {
		cmds[i] = exec.Command(shell, "-c", cmdStr)
		wg.Add(1)
		cmd := cmds[i] //copy to avoid data race
		title := "cmd " + fmt.Sprint(i+1)
		color := NextColor()
		go func() {
			defer wg.Done()
			proc(cmd, title, color)
		}()
	}

	//capture and handle interrupt signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		var terminating bool = false
		for range c {
			if !terminating {
				fmt.Println("Gracefully stopping... (press Ctrl+C again to force)") //TODO: pretty colors here
				terminating = true
				//send sigint to processes'
				for _, cmd := range cmds {
					cmd.Process.Signal(syscall.SIGINT)
				}
				continue
			}
			fmt.Println("ERROR: Aborting.") //TODO: pretty colors here
			os.Exit(255)
		}
	}()

	//wait for all commands to exit and output status
	wg.Wait()
	fmt.Println("Done.")
}

//THE PLAM

//formatter should be called "prefixer"

//flags:
//-wrap wrap output instead of truncating to terminal width?
//-fullcmd show full command on output
//-namelen number of characters to show before truncating name of commands
//-about print LICENSE (embedded)
//-help print this menu
//-restart restart commands after exiting
