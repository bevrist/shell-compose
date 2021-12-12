package proc

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/bevrist/shell-compose/format"
)

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
					fmt.Println(format.PrintCmdName("test") + line)
				}
			}
		}
	}()
	err = cmd.Wait()
	if err != nil {
		// println("err: " + stderr.String())
	}
	println(cmd.ProcessState.ExitCode())
}

func Run(command string, args ...string) {
	cmd := exec.Command(command, args[0:]...)
	go proc(cmd)
	time.Sleep(2 * time.Second)
	cmd.Process.Signal(syscall.SIGINT)
	fmt.Println("SIGINT Sent to process...")
	time.Sleep(2 * time.Second)
}
