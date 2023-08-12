package iotools

import (
	"io"
	"strings"

	"github.com/rs/zerolog/log"
)

type BackSpacer struct {
	writer io.Writer
}

func NewBackSpacer(w io.Writer) BackSpacer {
	log.
		Trace().
		Str("component", "backspacer").
		Msg("init new backspacer")
	return BackSpacer{
		writer: w,
	}
}

func (bsw BackSpacer) Write(b []byte) (int, error) {
	output := make([]byte, len(b))
	copy(output, b)
	output = append(output, []byte(strings.Repeat("\b", len(b)))...)
	return bsw.writer.Write(output)
}
