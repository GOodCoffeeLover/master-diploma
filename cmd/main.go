package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	osExec "os/exec"

	"github.com/GOodCoffeeLover/MasterDiploma/internal/exec"
	"golang.org/x/term"
)

func main() {

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	must(err, "Can't make raw term")
	defer term.Restore(int(os.Stdin.Fd()), oldState)
	must(osExec.Command("stty", "-F", "/dev/tty", "-echo").Run(), "Can't turn off print to term")
	defer func() {
		must(osExec.Command("stty", "-F", "/dev/tty", "echo").Run(), "Can't turn on print to term")
	}()

	_, out := os.Stdin, os.Stdout
	in, _ := NewBufferReadWriteCloser(10), NewBufferReadWriteCloser(10)
	inCh := make(chan byte, 10)
	// outCh := make(chan byte, 10)
	go fromReaderToChan(bufio.NewReader(os.Stdin), inCh)
	go fromChanToWriter(inCh, in)

	// go fromReaderToChan(out, outCh)
	// go fromChanToWriter(outCh, os.Stdout)

	must(exec.ExecCmdExample("test", "default", "sh", in, out, out), "Error while exec to pod")

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

type BufferReadWriteCloser struct {
	ch chan byte
}

func NewBufferReadWriteCloser(n uint) BufferReadWriteCloser {
	return BufferReadWriteCloser{
		ch: make(chan byte, n),
	}
}

func (brw BufferReadWriteCloser) Read(buf []byte) (int, error) {
	var i int
	for i = 0; i < len(buf); i++ {
		b, ok := <-brw.ch
		if !ok && i == 0 {
			return 0, io.EOF
		}
		if !ok {
			return i, nil
		}
		buf[i] = b

		exec.PrintlnRaw(os.Stderr, fmt.Sprintf("read from chan %v (%v)", b, string(b)))
	}
	return i, nil
}

func (brw BufferReadWriteCloser) Write(buf []byte) (int, error) {
	for _, b := range buf {
		brw.ch <- b

		exec.PrintlnRaw(os.Stderr, fmt.Sprintf("write to chan %v (%v)", b, string(b)))
	}
	return len(buf), nil
}

func (brw BufferReadWriteCloser) Close() error {
	close(brw.ch)

	exec.PrintlnRaw(os.Stderr, fmt.Sprintf("Close chan"))
	return nil
}
