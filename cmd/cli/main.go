package main

import (
	"os"

	"github.com/GOodCoffeeLover/master-diploma/internal/iotools"
	"github.com/GOodCoffeeLover/master-diploma/internal/remote"
	"github.com/rs/zerolog"
	log "github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Level(zerolog.TraceLevel)
	log.Logger = log.Output(iotools.NewBackSpacer(os.Stderr))
	log.
		Info().
		Msg("Started")
	s := remote.NewSession(os.Stdin, os.Stdout)
	s.Run("default", "test-0", "bash")

}
