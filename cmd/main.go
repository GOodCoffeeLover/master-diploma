package main

import (
	"os"

	"github.com/GOodCoffeeLover/MasterDiploma/internal/exec"
	// dockerterm "github.com/moby/term"
)

func main() {
	// executor := exec.NewExecutor()
	// executor.Exec()
	// in, out, _ := dockerterm.StdStreams()
	err := exec.ExecCmdExampleV2("test", "bash", os.Stdin, os.Stdout, os.Stdout)
	if err != nil {
		panic(err)
	}
}
