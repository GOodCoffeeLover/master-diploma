package main

import (
	"os"

	"github.com/GOodCoffeeLover/master-diploma/internal/iotools"
	"github.com/GOodCoffeeLover/master-diploma/internal/remote"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	// zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	log.Logger = log.Output(iotools.NewBackSpacer(os.Stderr))
	log.
		Info().
		Msg("Started")
	s := remote.NewSession(os.Stdin, os.Stdout)
	s.Run("default", "test-0", "bash")
}
