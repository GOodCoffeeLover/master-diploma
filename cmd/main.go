package main

import (
	"fmt"
	"os"
	osExec "os/exec"

	"github.com/GOodCoffeeLover/MasterDiploma/internal/exec"
	"golang.org/x/term"
)

func main() {
	in, out := os.Stdin, os.Stdout

	oldState, err := term.MakeRaw(int(in.Fd()))
	must(err, "Can't make raw term")
	defer term.Restore(int(in.Fd()), oldState)
	must(osExec.Command("stty", "-F", "/dev/tty", "-echo").Run(), "Can't turn off print to term")
	defer func() {
		must(osExec.Command("stty", "-F", "/dev/tty", "echo").Run(), "Can't turn on print to term")
	}()

	must(exec.ExecCmdExample("test", "default", "sh", in, out, out), "Error while exec to pod")

}

func must(err error, msg string) {
	if err != nil {
		panic(fmt.Errorf("%v %w", msg, err))
	}
}
