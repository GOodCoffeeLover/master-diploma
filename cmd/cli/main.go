package main

import (
	"os"

	"github.com/GOodCoffeeLover/master-diploma/internal/iotools"
	"github.com/GOodCoffeeLover/master-diploma/internal/remote"
	pb "github.com/GOodCoffeeLover/master-diploma/pkg/sandbox/api"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// zerolog.SetGlobalLevel(zerolog.TraceLevel)
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	// zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	log.Logger = log.Output(iotools.NewBackSpacer(os.Stderr))

	conn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal().Err(err).Msgf("fail to dial: %v", err)
	}
	defer conn.Close()
	log.
		Info().
		Msg("Started")

	sandbox := pb.NewSandboxClient(conn)
	s := remote.NewSession(os.Stdin, os.Stdout, sandbox)
	s.Run("default", "test-0", "bash")
}
