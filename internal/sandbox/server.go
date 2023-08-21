package sandbox

import (
	"errors"
	"fmt"
	"io"

	"github.com/GOodCoffeeLover/master-diploma/internal/iotools"
	"github.com/GOodCoffeeLover/master-diploma/internal/remote"
	pb "github.com/GOodCoffeeLover/master-diploma/pkg/sandbox/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/client-go/rest"
)

type SandboxServer struct {
	pb.UnimplementedSandboxServer
	config *rest.Config
}

func NewSandboxServer(config *rest.Config) SandboxServer {
	return SandboxServer{
		config: config,
	}

}
func (s SandboxServer) Execute(stream pb.Sandbox_ExecuteServer) error {
	req, err := stream.Recv()
	if err != nil {
		return status.Error(codes.Internal, "execute: error while reading first msg")
	}
	executor, err := remote.NewExecutor(s.config, req.GetNamespace(), req.GetPod())
	if err != nil {
		if errors.Is(err, remote.ErrorInvalidExecutionTarget) {
			return status.Error(codes.InvalidArgument, err.Error())
		} else {
			return status.Error(codes.Internal, fmt.Sprintf("failed to create remote executor: %v", err))
		}
	}
	finish := make(chan struct{})
	errs := make(chan error, 2)
	in := iotools.NewBuffer(1)
	go func() {
	loop:
		for {
			select {
			case <-finish:
				{
					errs <- nil
					break loop
				}
			default:
				{
					r, err := stream.Recv()
					if errors.Is(err, io.EOF) {
						errs <- nil
						break loop
					}
					if err != nil {
						errs <- err
						break loop
					}
					in.Write([]byte(r.GetText()))
				}
			}
		}
		in.Close()
	}()
	out := iotools.NewBuffer(1)
	go func() {
		for {
			b := make([]byte, 1)
			_, err := out.Read(b)
			if errors.Is(err, io.EOF) {
				errs <- nil
				break
			}
			if err != nil {
				errs <- fmt.Errorf("error while reading: %w", err)
				break
			}
			err = stream.Send(&pb.ExecuteResponse{
				Text: string(b),
			})
			if err != nil {
				errs <- fmt.Errorf("error while sending: %w", err)
				break
			}
		}
		finish <- struct{}{}
	}()
	executor.Exec(req.GetCommand(), in, out)
	for i := 0; i < 2; i++ {
		err = <-errs
		if err != nil {
			return err
		}
	}
	return nil
}
