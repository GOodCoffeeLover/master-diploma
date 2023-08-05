package main

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/GOodCoffeeLover/MasterDiploma/internal/buffer"
	exec "github.com/GOodCoffeeLover/MasterDiploma/internal/remoteExecuctor"
)

func main() {

	// oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	// must(err, "Can't make raw term")
	// defer term.Restore(int(os.Stdin.Fd()), oldState)
	// must(osExec.Command("stty", "-F", "/dev/tty", "-echo").Run(), "Can't turn off print to term")
	// defer func() {
	// 	must(osExec.Command("stty", "-F", "/dev/tty", "echo").Run(), "Can't turn on print to term")
	// }()

	out := os.Stdout
	in := os.Stdin
	// in := buffer.NewBufferReadWriteCloser(10)
	_ = buffer.NewBufferReadWriteCloser(10)
	inCh := make(chan byte, 10)
	// outCh := make(chan byte, 10)
	go fromReaderToChan(in, inCh)
	go fromChanToWriter(inCh, out)

	// go fromReaderToChan(out, outCh)
	// go fromChanToWriter(outCh, os.Stdout)

	// must(exec.ExecCmdExample("test", "default", "sh", in, out, out), "Error while exec to pod")
	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}

func must(err error, msg string) {
	if err != nil {
		panic(fmt.Errorf("%v %w", msg, err))
	}
}

func fromReaderToChan(in io.Reader, out chan<- byte) {
	for {
		b := make([]byte, 1)
		_, err := in.Read(b)

		if err == io.EOF {
			close(out)
			return
		}
		exec.PrintlnRaw(os.Stderr, fmt.Sprintf("readed %v (%v)", b, string(b)))
		must(err, "error while reading")
		out <- b[0]

	}
}

func fromChanToWriter(in <-chan byte, out io.WriteCloser) {
	for {
		b, ok := <-in
		if !ok {
			exec.PrintlnRaw(os.Stderr, fmt.Sprintf("closed"))
			must(out.Close(), "error while closing write")
			return
		}

		exec.PrintlnRaw(os.Stderr, fmt.Sprintf("get to write %v (%v)", b, string(b)))
		_, err := out.Write([]byte{b})
		must(err, "error while writing")
	}
}
