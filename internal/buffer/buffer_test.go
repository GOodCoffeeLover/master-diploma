package buffer

import (
	"strings"
	"testing"
)

func TestBasic(t *testing.T) {
	buf := NewBufferReadWriteCloser(10, "")
	input := []byte("string")
	output := make([]byte, len(input))

	n, err := buf.Write(input)
	if err != nil {
		t.Errorf("failed to write: %v", err)
	}
	if n != len(input) {
		t.Errorf("trying to write %v bytes, but write %v", len(input), n)
	}
	n, err = buf.Read(output)

	if err != nil {
		t.Errorf("failed to read: %v", err)
	}
	if n != len(output) {
		t.Errorf("trying to read %v bytes, but read %v", len(output), n)
	}
}

func TestFromReaderToChan(t *testing.T) {
	input := "someString\nstrstrs"
	reader := strings.NewReader(input)
	outChan := make(chan byte)

	go FromReaderToChan(reader, outChan)

	output := []byte{}
	for b := range outChan {
		output = append(output, b)
	}
	if string(output) != input {
		t.Errorf("input( %v ) not equal to output( %v )", input, string(output))
	}
}

type testBuffer struct {
}

var buf []byte

func (b testBuffer) Write(bytes []byte) (int, error) {
	buf = append(buf, bytes...)
	return len(bytes), nil
}

func (buf testBuffer) Close() error { return nil }

func TestFromChanToWriter(t *testing.T) {
	input := "someString\nstrstrs"
	ch := make(chan byte)
	outputBuf := testBuffer{}

	go func() {
		for _, b := range []byte(input) {
			ch <- b
		}
		close(ch)
	}()
	FromChanToWriter(ch, outputBuf)
	output := buf
	if string(output) != input {
		t.Errorf("input( %v ) not equal to output( %v )", input, string(output))
	}
}
