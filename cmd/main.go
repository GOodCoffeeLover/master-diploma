package main

import (
	"fmt"
	"os"

	"github.com/GOodCoffeeLover/MasterDiploma/internal/buffer"
	remoteExecuctor "github.com/GOodCoffeeLover/MasterDiploma/internal/remoteExecuctor"
	"github.com/u-root/u-root/pkg/termios"
	"golang.org/x/term"
)

func main() {
	remoteExecuctor.PrintlnRaw(os.Stderr, "Started")
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	must(err, "Can't make raw term")
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	t, err := termios.GTTY(int(os.Stdin.Fd()))
	must(err, "Can't get termious terminal")

	must(t.SetOpts([]string{"~echo"}), "Can't turn off print to term")

	defer func() {
		must(t.SetOpts([]string{"echo"}), "Can't turn on print to term")
	}()

	out := os.Stdout
	in := os.Stdin
	inBuf := buffer.NewBufferReadWriteCloser(10)
	outBuf := buffer.NewBufferReadWriteCloser(10)
	inCh := make(chan byte, 10)
	outCh := make(chan byte, 10)

	go buffer.FromReaderToChan(in, inCh)
	go buffer.FromChanToWriter(inCh, inBuf)

	go buffer.FromReaderToChan(outBuf, outCh)
	go buffer.FromChanToWriter(outCh, out)

	must(remoteExecuctor.ExecCmdExample("test", "default", "bash", inBuf, outBuf, outBuf), "Error while remoteExecuctor to pod")

}

func must(err error, msg string) {
	if err != nil {
		panic(fmt.Errorf("%v: %w", msg, err))
	}
}
