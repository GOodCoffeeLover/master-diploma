package main

import (
	"os"

	"github.com/GOodCoffeeLover/MasterDiploma/internal/exec"
)

func main() {
	// executor := exec.NewExecutor()
	// executor.Exec()
	// in, out, _ := dockerterm.StdStreams()
	in, out, _ := os.Stdin, os.Stdout, os.Stderr
	must(exec.ExecCmdExampleV2("test", "default", "ls", in, out, out))

}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
