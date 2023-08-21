package remote

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	pb "github.com/GOodCoffeeLover/master-diploma/pkg/sandbox/api"
	"google.golang.org/grpc"

	"github.com/rs/zerolog/log"
	"github.com/u-root/u-root/pkg/termios"
	"golang.org/x/term"
)

type Session struct {
	in      *os.File
	out     *os.File
	sandbox pb.SandboxClient
}

func NewSession(in, out *os.File, sandboxClient pb.SandboxClient) Session {

	return Session{
		sandbox: sandboxClient,
		in:      in,
		out:     out,
	}
}

func (s *Session) setupTTY() (func(), error) {

	oldState, err := term.MakeRaw(int(s.in.Fd()))
	if err != nil {
		return func() {}, fmt.Errorf("can't make raw term: %v", err)
	}

	t, err := termios.GTTY(int(s.in.Fd()))
	if err != nil {
		term.Restore(int(s.in.Fd()), oldState)
		return func() {}, fmt.Errorf("can't get termious terminal: %v", err)
	}
	err = t.SetOpts([]string{"~echo"})
	if err != nil {
		term.Restore(int(s.in.Fd()), oldState)
		return func() {}, fmt.Errorf("failed turn off output: %w", err)
	}

	return func() {
		term.Restore(int(os.Stdin.Fd()), oldState)
		t.SetOpts([]string{"echo"})
	}, nil
}
func (s *Session) readInput(stream pb.Sandbox_ExecuteClient) error {

	for {
		b := make([]byte, 1)
		n, err := s.in.Read(b)
		log.
			Trace().
			Str("component", "session").
			Err(err).
			Msgf("read from input %v bytes: %v (%v)", n, b, string(b))

		if errors.Is(err, io.EOF) {
			log.
				Debug().
				Str("component", "session").
				Msg("finished with stdin to chan")
			return nil
		}

		if err != nil {
			return err
		}
		err = stream.Send(&pb.ExecuteRequest{
			Text: string(b),
		})
		if err != nil {
			return err
		}
	}
}

func (s *Session) writeOutput(stream pb.Sandbox_ExecuteClient) error {
	for {
		resp, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			log.
				Debug().
				Str("component", "session").
				Msg("finished with chan to stdout")
			return nil
		}
		if err != nil {
			return err
		}
		n, err := s.out.Write([]byte(resp.GetText()))
		log.
			Trace().
			Err(err).
			Str("component", "FromChanToWriter").
			Msgf("write %v (%v) to out writer %v bytes", []byte(resp.GetText()), resp.GetText(), n)
		if err != nil {
			return err
		}

	}
}

func (s *Session) Run(ns, pod, cmd string) error {
	unsetup, err := s.setupTTY()
	if err != nil {
		return err
	}
	defer unsetup()
	stream, err := s.sandbox.Execute(context.Background(), grpc.EmptyCallOption{})
	if err != nil {
		return err
	}
	stream.Send(&pb.ExecuteRequest{
		Namespace: ns,
		Pod:       pod,
		Command:   cmd,
	})
	go func() {
		err := s.readInput(stream)
		if err != nil {
			panic(err)
		}
	}()

	err = s.writeOutput(stream)
	if err != nil {
		panic(err)
	}
	return nil
}
