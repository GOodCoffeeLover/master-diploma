package iotools

import (
	"context"
	"errors"
	"io"

	"github.com/rs/zerolog/log"
)

type Buffer struct {
	ch chan byte
}

func NewBuffer(n uint) Buffer {
	ch := make(chan byte, n)
	return Buffer{
		ch: ch,
	}
}

func (brw Buffer) Read(buf []byte) (int, error) {

	b, ok := <-brw.ch

	buf[0] = b

	log.
		Debug().
		Str("component", "Buffer").
		Msgf("Read from chan %v (%v)", b, string(b))

	if !ok {
		log.
			Debug().
			Str("component", "Buffer").
			Msg("return EOF")
		return 0, io.EOF
	}

	return 1, nil
}

func (brw Buffer) Write(buf []byte) (int, error) {
	for _, b := range buf {
		brw.ch <- b
		log.
			Debug().
			Str("component", "Buffer").
			Msgf("write to chan %v (%v)", b, string(b))
	}
	return len(buf), nil
}

func (brw Buffer) Close() error {
	close(brw.ch)
	log.
		Debug().
		Str("component", "Buffer").
		Msg("close chan")
	return nil
}

func FromReaderToChan(ctx context.Context, in io.Reader, out chan<- byte) {
loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		default:
			b := make([]byte, 1)
			n, err := in.Read(b)
			log.
				Debug().
				Str("component", "FromReaderToChan").
				Err(err).
				Msgf("read from input %v bytes: %v (%v)", n, b, string(b))

			out <- b[0]
			if errors.Is(err, io.EOF) {
				break loop
			}

			if err != nil {
				panic(err)
			}
		}
	}
	close(out)
	log.
		Debug().
		Str("component", "FromReaderToChan").
		Msg("closed reader")
}

func FromChanToWriter(ctx context.Context, in <-chan byte, out io.WriteCloser) {
loop:
	for {
		select {
		case <-ctx.Done():
			break loop

		case b, ok := <-in:
			if !ok {
				break loop
			}
			log.
				Debug().
				Str("component", "FromChanToWriter").
				Msgf("from chan %v (%v)", b, string(b))
			n, err := out.Write([]byte{b})
			log.
				Debug().
				Str("component", "FromChanToWriter").
				Msgf("write to out writer %v bytes", n)
			if err != nil {
				panic(err)
			}

		}
	}
	err := out.Close()
	log.
		Debug().
		Err(err).
		Str("component", "FromChanToWriter").
		Msg("close writer")
}
