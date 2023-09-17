package remote

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/GOodCoffeeLover/master-diploma/internal/iotools"
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

func (s *Session) setupTTY() (unsetup func(), err error) {

	oldState, err := term.MakeRaw(int(s.in.Fd()))
	if err != nil {
		err = fmt.Errorf("can't make raw term: %v", err)
		return
	}
	defer func() {
		if err != nil {
			term.Restore(int(s.in.Fd()), oldState)
		}
	}()
	t, err := termios.GTTY(int(s.in.Fd()))
	if err != nil {
		err = fmt.Errorf("can't get termious terminal: %v", err)
		return
	}
	err = t.SetOpts([]string{"~echo"})
	if err != nil {
		err = fmt.Errorf("failed turn off output: %w", err)
		return
	}

	unsetup = func() {
		term.Restore(int(os.Stdin.Fd()), oldState)
		t.SetOpts([]string{"echo"})
	}
	return
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
		inp := string(b)
		err = stream.Send(&pb.ExecuteRequest{
			Input: &inp,
		})
		if err != nil {
			return err
		}
	}
}

func (s *Session) writeOutput(stream pb.Sandbox_ExecuteClient) error {
	defer stream.CloseSend()
	for {
		resp, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			log.
				Debug().
				Str("component", "session").
				Msg("finished with stream to stdout")
			return nil
		}
		if err != nil {
			return err
		}
		n, err := s.out.Write([]byte(resp.GetOutput()))
		log.
			Trace().
			Err(err).
			Str("component", "FromChanToWriter").
			Msgf("write %v (%v) to out writer %v bytes", []byte(resp.GetOutput()), resp.GetOutput(), n)
		if err != nil {
			return err
		}

	}
	log.
		Debug().
		Str("component", "session").
		Msg("finished with chan to stdout")
	finish()
	iotools.NewBackspaceWriter(s.out).Write([]byte("Press ANY KEY to exit"))
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
		Namespace: &ns,
		Pod:       &pod,
		Command:   &cmd,
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
