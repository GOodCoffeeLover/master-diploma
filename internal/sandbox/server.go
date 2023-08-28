package sandbox

import (
	"errors"
	"fmt"
	"io"

	"github.com/GOodCoffeeLover/master-diploma/internal/iotools"
	"github.com/GOodCoffeeLover/master-diploma/internal/remote"
	pb "github.com/GOodCoffeeLover/master-diploma/pkg/sandbox/api"
	log "github.com/rs/zerolog/log"
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
func (s *SandboxServer) Execute(stream pb.Sandbox_ExecuteServer) error {
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

	errs := make(chan error, 2)
	in := iotools.NewBuffer(1)
	out := iotools.NewBuffer(1)

	go s.readInput(in, stream, errs)
	go s.writeOutput(out, stream, errs)

	executor.Exec(req.GetCommand(), in, out)

	err = <-errs
	if err != nil {
		return err
	}
	return nil
}

func (ss *SandboxServer) writeOutput(out io.Reader, stream pb.Sandbox_ExecuteServer, errs chan<- error) {
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
			Output: string(b),
		})
		if err != nil {
			errs <- fmt.Errorf("error while sending: %w", err)
			break
		}
	}
	log.
		Info().
		Msg("finish with writing output")
}

func (s *SandboxServer) readInput(in io.WriteCloser, stream pb.Sandbox_ExecuteServer, errs chan<- error) {
	for {
		r, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			errs <- nil
			break
		}
		if err != nil {
			errs <- err
			break
		}
		in.Write([]byte(r.GetInput()))
	}
	in.Close()
	log.
		Info().
		Msg("finish with reading input")

}
