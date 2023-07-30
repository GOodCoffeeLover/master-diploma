package main

import (
	"os"
	osExec "os/exec"

	"github.com/GOodCoffeeLover/MasterDiploma/internal/exec"
	"golang.org/x/term"
)

func main() {
	// executor := exec.NewExecutor()
	// executor.Exec()
	// in, out, _ := dockerterm.StdStreams()
	in, out, _ := os.Stdin, os.Stdout, os.Stderr

	oldState, err := term.MakeRaw(int(in.Fd()))
	must(err)
	defer term.Restore(int(in.Fd()), oldState)
	must(osExec.Command("stty", "-F", "/dev/tty", "-echo").Run())
	defer func() {
		must(osExec.Command("stty", "-F", "/dev/tty", "echo").Run())
	}()

	must(exec.ExecCmdExample("test", "default", "sh", in, out, out))

}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
