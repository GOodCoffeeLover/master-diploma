package buffer

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/GOodCoffeeLover/master-diploma/internal/remoteExecuctor"
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

	b, ok := <-brw.ch
	buf[0] = b
	remoteExecuctor.PrintlnRaw(os.Stderr, fmt.Sprintf("read from chan %v (%v), %v", b, string(b), buf))
	// if b == 4 {
	// 	return 1, io.EOF
	// }
	if !ok {
		return 0, io.EOF
	}

	return 1, nil
}

func (brw BufferReadWriteCloser) Write(buf []byte) (int, error) {
	for _, b := range buf {
		brw.ch <- b

		remoteExecuctor.PrintlnRaw(os.Stderr, fmt.Sprintf("write to chan %v (%v)", b, string(b)))
	}
	return len(buf), nil
}

func (brw BufferReadWriteCloser) Close() error {
	close(brw.ch)
	remoteExecuctor.PrintlnRaw(os.Stderr, "close chan")
	return nil
}

func FromReaderToChan(ctx context.Context, in io.Reader, out chan<- byte) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			b := make([]byte, 1)
			n, err := in.Read(b)
			remoteExecuctor.PrintlnRaw(os.Stderr, fmt.Sprintf("read from in reader %v bytes", n))

			remoteExecuctor.PrintlnRaw(os.Stderr, fmt.Sprintf("frorm reader %v (%v)", b, string(b)))
			out <- b[0]
			if errors.Is(err, io.EOF) {
				remoteExecuctor.PrintlnRaw(os.Stderr, "closed reader")
				close(out)
				return
			}
			must(err, "error while reading")
		}

	}
}

func FromChanToWriter(ctx context.Context, in <-chan byte, out io.WriteCloser) {
	for {
		select {
		case <-ctx.Done():
			return
		case b, ok := <-in:
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
}

func must(err error, msg string) {
	if err != nil {
		panic(fmt.Errorf("%v: %w", msg, err))
	}
}
