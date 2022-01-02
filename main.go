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
	"time"

	_ "embed"

	"github.com/spf13/pflag"
)

var (
	//arg flags
	fColor   = pflag.Bool("color", false, "force color output")
	fHelp    = pflag.Bool("help", false, "show this help menu and exit")
	fLicense = pflag.Bool("license", false, "print the license")
	fNameLen = pflag.IntP("name-length", "n", 20, "max number of characters to show before truncating name of commands")
	fNoColor = pflag.Bool("nocolor", false, "disable color output")
	fRestart = pflag.BoolP("restart", "r", false, "restart commands after failure (non zero exit code)")
	fShell   = pflag.StringP("shell", "s", "", "shell to launch commands with")
	fVersion = pflag.BoolP("version", "v", false, "show version information and exit")

	//regex to capture all empty strings
	reEmpty = regexp.MustCompile(`^\s*$`)

	//Semantic Version Info
	Version, GitCommit, BuildDate string = "development", "null", "null"

	//go:embed LICENSE
	license string
)

//printCmdName outputs truncated command name in color
func printCmdName(commandTitle string, color string) string {
	return color + commandTitle + " | " + ResetColor()
}

//formatTitle creates string with padding for specific length and truncates if necessary
func formatTitle(title string, titleLen int) string {
	//truncate if title is too long
	if len(title) > titleLen {
		title = title[:titleLen-1]
		title += "…"
	}

	//format output with padding if necessary
	tlen := fmt.Sprint(titleLen)
	tfmt := "%-" + tlen + "." + tlen + "s"
	return fmt.Sprintf(tfmt, title)
}

//runProcess handles output and lifecycle of commands
func runProcess(cmd *exec.Cmd, title string, color string) {
	//capture command output streams
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	outReader := bufio.NewReader(stdout)
	errReader := bufio.NewReader(stderr)
	//loop when restart flag is enabled
	for {
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
					fmt.Print(printCmdName(title, color) + outline)
					outline, err = outReader.ReadString('\n')
				}

				//stderr
				errline, err := errReader.ReadString('\n')
				for err == nil {
					//skip printing empty lines
					if reEmpty.FindStringIndex(errline) != nil {
						continue
					}
					fmt.Print(printCmdName(title, color) + ErrorColor() + "STDERR: " + ResetColor() + errline)
					errline, err = outReader.ReadString('\n')
				}
			}
		}()
		//keep goroutine running as long as command is running
		cmd.Wait()
		exitCode := cmd.ProcessState.ExitCode()
		fmt.Println(printCmdName(title, color) + "Process Exited with code: " + fmt.Sprint(exitCode))
		// delay before restarting command
		time.Sleep(time.Second)
		if !*fRestart || exitCode == 0 {
			return
		}
	}
}

func init() {
	//update flags and help menu
	pflag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage of "+os.Args[0])
		fmt.Fprintln(os.Stderr, " shell-compose: run and view output of multiple commands at once")
		fmt.Fprintf(os.Stderr, "\n")
		pflag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintln(os.Stderr, "Written by @bevrist")
		fmt.Fprintln(os.Stderr, "https://github.com/bevrist")
		os.Exit(0)
	}
	pflag.Parse()

	//print license
	if *fLicense {
		fmt.Println(license)
		os.Exit(0)
	}

	// print version
	if *fVersion {
		fmt.Printf("Version: %s \nGit Commit: %s \nBuild Date: %s\n", Version, GitCommit, BuildDate)
		os.Exit(0)
	}

	//print help
	if *fHelp || len(pflag.Args()) < 1 {
		pflag.Usage()
		os.Exit(0)
	}
}

func main() {
	//get shell to launch commands with
	var shell string

	//if shell provided with flag, verify binary can be found
	if *fShell != "" {
		if shell, _ = exec.LookPath(*fShell); shell == "" {
			log.Fatal(ErrorColor() + "ERROR:" + ResetColor() + " '" + *fShell + "' shell binary not found.")
		}
	} else if shell, _ = exec.LookPath(os.Getenv("SHELL")); shell == "" {
		//test for shell var, else try other potential shells
		shells := []string{"bash", "sh", "zsh", "ash", "fish"}
		for _, item := range shells {
			var err error
			shell, err = exec.LookPath(item)
			if err == nil {
				break
			}
		}
		//err out if no shell is found after all checks
		if shell == "" {
			log.Fatal(ErrorColor() + "ERROR:" + ResetColor() + " no shell found.")
		}
	}

	//calculate max title length
	var nameLen int
	for _, cmd := range pflag.Args() {
		if len(cmd) > nameLen {
			nameLen = len(cmd)
		}
	}
	//do not allow titles longer than max length
	if nameLen > *fNameLen {
		nameLen = *fNameLen
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
		//format title to be consistent length with others
		title := formatTitle(cmdStr, nameLen)
		color := NextColor()
		go func() {
			defer wg.Done()
			runProcess(cmd, title, color)
		}()
	}

	//capture and handle interrupt signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		var terminating bool = false
		for range c {
			if !terminating {
				fmt.Println(SuccessColor() + "Gracefully stopping..." + ResetColor() + " (press Ctrl+C again to force)")
				terminating = true
				*fRestart = false //stop restarting processes
				//send sigint to processes'
				for _, cmd := range cmds {
					cmd.Process.Signal(syscall.SIGINT)
				}
				continue
			}
			fmt.Println(ErrorColor() + "ERROR:" + ResetColor() + " Aborting.")
			os.Exit(255)
		}
	}()

	//wait for all commands to exit and output status
	wg.Wait()
	fmt.Println("Done.")
}
