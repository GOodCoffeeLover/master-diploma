package main

import (
	"os"

	"github.com/GOodCoffeeLover/MasterDiploma/internal/exec"
)

func main() {
	// executor := exec.NewExecutor()
	// executor.Exec()
	err := exec.ExecCmdExample("test", "ls", os.Stdin, os.Stdout, os.Stderr)
	if err != nil {
		panic(err)
	}
}
