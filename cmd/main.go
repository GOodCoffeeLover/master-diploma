package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/GOodCoffeeLover/MasterDiploma/internal/buffer"
	remoteExecuctor "github.com/GOodCoffeeLover/MasterDiploma/internal/remoteExecuctor"
	"github.com/u-root/u-root/pkg/termios"
	"golang.org/x/term"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	remoteExecuctor.PrintlnRaw(os.Stderr, "Started")
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	must(err, "Can't make raw term")
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	t, err := termios.GTTY(int(os.Stdin.Fd()))
	must(err, "Can't get termious terminal")

	must(t.SetOpts([]string{"~echo"}), "Can't turn off print to term")

	defer func() {
		must(t.SetOpts([]string{"echo"}), "Can't turn on print to term")
	}()

	out := os.Stdout
	in := os.Stdin
	inBuf := buffer.NewBufferReadWriteCloser(10)
	outBuf := buffer.NewBufferReadWriteCloser(10)
	inCh := make(chan byte, 10)
	outCh := make(chan byte, 10)
	wg := sync.WaitGroup{}
	wg.Add(4)
	go func() {
		buffer.FromReaderToChan(in, inCh)
		wg.Done()
	}()
	go func() {
		buffer.FromChanToWriter(inCh, inBuf)
		wg.Done()
	}()

	go func() {
		buffer.FromReaderToChan(outBuf, outCh)
		wg.Done()
	}()
	go func() {
		buffer.FromChanToWriter(outCh, out)
		wg.Done()
	}()

	home := homedir.HomeDir()
	config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(home, ".kube", "config"))
	must(err, "can't create kube config")
	executor, err := remoteExecuctor.NewRemoteExecutor(config, "default", "test")
	must(err, "can't create executor")
	must(executor.Exec("bash", inBuf, outBuf), "Error while remoteExecuctor to pod")

	wg.Wait()
}

func must(err error, msg string) {
	if err != nil {
		panic(fmt.Errorf("%v: %w", msg, err))
	}
}
