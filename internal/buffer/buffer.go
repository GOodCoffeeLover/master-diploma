package buffer

import (
	"fmt"
	"io"
	"os"

	"github.com/GOodCoffeeLover/MasterDiploma/internal/remoteExecuctor"
	exec "github.com/GOodCoffeeLover/MasterDiploma/internal/remoteExecuctor"
)

type BufferReadWriteCloser struct {
	ch chan byte
}

func NewBufferReadWriteCloser(n uint) BufferReadWriteCloser {
	ch := make(chan byte, n)
	return BufferReadWriteCloser{
		ch: ch,
	}
}

func (brw BufferReadWriteCloser) Read(buf []byte) (int, error) {
	// var i int
	// for i = 0; i < len(buf); i++ {
	// 	select {
	// 	case b, ok := <-brw.ch:
	// 		{
	// 			if !ok && i == 0 {
	// 				return 0, io.EOF
	// 			}
	// 			if !ok {
	// 				return i, nil
	// 			}
	// 			buf[i] = b

	// 			exec.PrintlnRaw(os.Stderr, fmt.Sprintf("read from chan %v (%v), %v", b, string(b), buf))
	// 		}
	// 	case <-time.Tick(time.Second * 5):
	// 		{
	// 			exec.PrintlnRaw(os.Stderr, fmt.Sprintf("returns by default with %v", i))
	// 			return i, nil
	// 		}
	// 	}

	// }
	b, ok := <-brw.ch

	if !ok {
		return 0, io.EOF
	}
	exec.PrintlnRaw(os.Stderr, fmt.Sprintf("read from chan %v (%v), %v", b, string(b), buf))
	buf[0] = b
	return 1, nil
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

func FromReaderToChan(in io.Reader, out chan<- byte) {
	for {
		b := make([]byte, 1)
		n, err := in.Read(b)
		remoteExecuctor.PrintlnRaw(os.Stderr, fmt.Sprintf("read from in reader %v bytes", n))

		if err == io.EOF {
			remoteExecuctor.PrintlnRaw(os.Stderr, "closed reader")
			close(out)
			return
		}
		remoteExecuctor.PrintlnRaw(os.Stderr, fmt.Sprintf("frorm reader %v (%v)", b, string(b)))
		must(err, "error while reading")
		for i := 0; i < n; i++ {

			out <- b[i]
		}

	}
}

func FromChanToWriter(in <-chan byte, out io.WriteCloser) {
	for {
		b, ok := <-in
		if !ok {
			remoteExecuctor.PrintlnRaw(os.Stderr, "closed chan")
			must(out.Close(), "error while closing write")
			return
		}

		remoteExecuctor.PrintlnRaw(os.Stderr, fmt.Sprintf("from chan %v (%v)", b, string(b)))
		n, err := out.Write([]byte{b})
		remoteExecuctor.PrintlnRaw(os.Stderr, fmt.Sprintf("write to out writer %v bytes", n))
		must(err, "error while writing")
	}
}

func must(err error, msg string) {
	if err != nil {
		panic(fmt.Errorf("%v: %w", msg, err))
	}
}
