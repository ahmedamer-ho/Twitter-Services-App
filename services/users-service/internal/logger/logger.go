package logger

import (
	"os"

	"github.com/rs/zerolog"
)

var Log = zerolog.New(os.Stdout).
	With().
	Timestamp().
	Logger()
