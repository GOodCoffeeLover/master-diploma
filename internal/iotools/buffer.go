package iotools

import (
	"io"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Buffer struct {
	ch  chan byte
	log *zerolog.Logger
}

func NewBuffer(n uint) Buffer {
	ch := make(chan byte, n)
	logger := log.With().Str("component", "Buffer").Logger()
	return Buffer{
		ch:  ch,
		log: &logger,
	}
}

func (brw Buffer) Read(buf []byte) (int, error) {

	b, ok := <-brw.ch

	buf[0] = b

	brw.log.
		Debug().
		Msgf("Read from chan %v (%v)", b, string(b))

	if !ok {
		brw.log.
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
		brw.log.
			Debug().
			Str("component", "Buffer").
			Msgf("write to chan %v (%v)", b, string(b))
	}
	return len(buf), nil
}

func (brw Buffer) Close() error {
	close(brw.ch)
	brw.log.
		Debug().
		Str("component", "Buffer").
		Msg("close chan")
	return nil
}
