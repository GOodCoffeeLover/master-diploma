package iotools

import (
	"io"
	"strings"

	"github.com/rs/zerolog/log"
)

type BackspaceWriter struct {
	writer io.Writer
}

func NewBackspaceWriter(w io.Writer) BackspaceWriter {
	log.
		Trace().
		Str("component", "BackspaceWriter").
		Msg("init new BackspaceWriter")
	return BackspaceWriter{
		writer: w,
	}
}

func (bsw BackspaceWriter) Write(b []byte) (int, error) {
	output := make([]byte, len(b))
	copy(output, b)
	output = append(output, []byte(strings.Repeat("\b", len(b)))...)
	return bsw.writer.Write(output)
}
