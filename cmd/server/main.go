package main

import (
	"net"
	"path/filepath"

	"github.com/GOodCoffeeLover/master-diploma/internal/sandbox"
	pb "github.com/GOodCoffeeLover/master-diploma/pkg/sandbox/api"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	// zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	// log.Logger = log.Output(iotools.NewBackSpacer(os.Stderr))
	lis, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.
			Fatal().
			Err(err).
			Msgf("failed to listen: %v", err)
	}

	home := homedir.HomeDir()
	config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(home, ".kube", "config"))
	if err != nil {
		log.
			Fatal().
			Err(err).
			Msg("can't create kube config")
	}

	server := grpc.NewServer(grpc.EmptyServerOption{})
	pb.RegisterSandboxServer(server, sandbox.NewSandboxServer(config))
	log.
		Info().
		Msg("Started")
	if err := server.Serve(lis); err != nil {
		log.
			Fatal().
			Err(err).
			Msg("can't run server")
	}
}
