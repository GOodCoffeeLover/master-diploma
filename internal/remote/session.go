package remote

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/GOodCoffeeLover/master-diploma/internal/iotools"
	"github.com/rs/zerolog/log"
	"github.com/u-root/u-root/pkg/termios"
	"golang.org/x/term"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type Session struct {
	in  *os.File
	out *os.File
}

func NewSession(in, out *os.File) Session {

	return Session{
		in:  in,
		out: out,
	}
}

func (s *Session) setupTTY() (func(), error) {

	oldState, err := term.MakeRaw(int(s.in.Fd()))
	if err != nil {
		return func() {}, fmt.Errorf("can't make raw term: %v", err)
	}

	t, err := termios.GTTY(int(s.in.Fd()))
	if err != nil {
		term.Restore(int(os.Stdin.Fd()), oldState)
		return func() {}, fmt.Errorf("can't get termious terminal: %v", err)
	}
	must(t.SetOpts([]string{"~echo"}), "Can't turn off print to term")

	return func() {
		term.Restore(int(os.Stdin.Fd()), oldState)
		t.SetOpts([]string{"echo"})
	}, nil
}
func (s *Session) readInput(ctx context.Context, ch chan<- byte) {
loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		default:
			b := make([]byte, 1)
			n, err := s.in.Read(b)
			log.
				Trace().
				Str("component", "session").
				Err(err).
				Msgf("read from input %v bytes: %v (%v)", n, b, string(b))

			ch <- b[0]
			if errors.Is(err, io.EOF) {
				break loop
			}

			if err != nil {
				panic(err)
			}
		}
	}
	close(ch)
	log.
		Debug().
		Str("component", "session").
		Msg("finished with stdin to chan")

}

func (s *Session) writeOutput(ctx context.Context, ch <-chan byte, finish context.CancelFunc) {
	run := true
	for run {
		select {
		case <-ctx.Done():
			run = false
			continue

		case b, ok := <-ch:
			if !ok {
				run = false
				continue
			}
			n, err := s.out.Write([]byte{b})
			log.
				Trace().
				Err(err).
				Str("component", "FromChanToWriter").
				Msgf("write %v (%v) to out writer %v bytes", b, string(b), n)
			if err != nil {
				panic(err)
			}

		}
	}
	log.
		Debug().
		Str("component", "session").
		Msg("finished with chan to stdout")
	finish()
	iotools.NewBackSpacer(s.out).Write([]byte("Press ANY KEY to exit"))
}

func (s *Session) Run(ns, pod, cmd string) error {
	unsetup, err := s.setupTTY()
	if err != nil {
		return err
	}
	defer unsetup()

	inCh := make(chan byte, 1)
	outCh := make(chan byte, 1)
	ctx, finish := context.WithCancel(context.Background())

	wg := sync.WaitGroup{}
	go func() {
		s.readInput(ctx, inCh)
	}()

	go func() {
		wg.Add(1)
		s.writeOutput(ctx, outCh, finish)
		wg.Done()
	}()

	go rpc(ctx, inCh, outCh, ns, pod, cmd)

	wg.Wait()
	return nil
}

func must(err error, msg string) {
	if err != nil {
		panic(fmt.Errorf("%v: %w", msg, err))
	}
}

func rpc(ctx context.Context, inCh, outCh chan byte, ns, pod, cmd string) {
	home := homedir.HomeDir()
	config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(home, ".kube", "config"))
	must(err, "can't create kube config")

	inBuf := iotools.NewBuffer(10)
	outBuf := iotools.NewBuffer(10)
	go func() {
		iotools.FromChanToWriter(ctx, inCh, inBuf)
		log.
			Debug().
			Str("component", "session").
			Msg("finished with stdin chan to buf")
	}()

	go func() {
		iotools.FromReaderToChan(ctx, outBuf, outCh)
		log.
			Debug().
			Str("component", "session").
			Msg("finished with outBuf to chan")
	}()

	executor, err := NewExecutor(config, ns, pod)
	must(err, "can't create executor")
	must(executor.Exec(cmd, inBuf, outBuf), "Error while remoteExecuctor to pod")

}
