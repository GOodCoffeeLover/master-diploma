package remote

import (
	"context"
	"fmt"
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
	in        *os.File
	out       *os.File
	namespace string
	pod       string
}

func NewSession(in, out *os.File) Session {

	return Session{
		in:  in,
		out: out,
	}
}

func (s *Session) setupTTY() func() {
	oldState, err := term.MakeRaw(int(s.in.Fd()))
	must(err, "Can't make raw term")
	t, err := termios.GTTY(int(s.in.Fd()))
	must(err, "Can't get termious terminal")
	must(t.SetOpts([]string{"~echo"}), "Can't turn off print to term")

	return func() {
		term.Restore(int(os.Stdin.Fd()), oldState)
		must(t.SetOpts([]string{"echo"}), "Can't turn on print to term")
	}
}

func (s *Session) Run(cmd string) {
	home := homedir.HomeDir()
	config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(home, ".kube", "config"))
	must(err, "can't create kube config")

	unsetup := s.setupTTY()
	defer unsetup()

	inBuf := iotools.NewBuffer(10)
	outBuf := iotools.NewBuffer(10)
	inCh := make(chan byte, 1)
	outCh := make(chan byte, 1)
	ctx, finish := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}
	wg.Add(4)
	go func() {
		iotools.FromReaderToChan(ctx, s.in, inCh)
		log.
			Debug().
			Err(err).
			Str("component", "main").
			Msg("finished with stdin to chan")
		wg.Done()
	}()
	go func() {
		iotools.FromChanToWriter(ctx, inCh, inBuf)
		log.
			Debug().
			Err(err).
			Str("component", "main").
			Msg("finished with stdin chan to buf")
		wg.Done()
	}()

	go func() {
		iotools.FromReaderToChan(ctx, outBuf, outCh)
		log.
			Debug().
			Err(err).
			Str("component", "main").
			Msg("finished with outBuf to chan")
		wg.Done()
	}()
	go func() {
		iotools.FromChanToWriter(ctx, outCh, s.out)
		log.
			Debug().
			Err(err).
			Str("component", "main").
			Msg("finished with chan to stdout")
		wg.Done()
		finish()
	}()

	executor, err := NewExecutor(config, "default", "test")
	must(err, "can't create executor")
	must(executor.Exec(cmd, inBuf, outBuf), "Error while remoteExecuctor to pod")
	iotools.NewBackSpacer(s.out).Write([]byte("Press ANY KEY to exit"))
	wg.Wait()
}

func must(err error, msg string) {
	if err != nil {
		panic(fmt.Errorf("%v: %w", msg, err))
	}
}
