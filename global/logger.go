package global

import (
	"context"
	"os"

	"github.com/whoisnian/glb/ansi"
	"github.com/whoisnian/glb/logger"
)

var LOG *logger.Logger

func SetupLogger(_ context.Context) {
	opts := logger.Options{
		Level:     logger.LevelInfo,
		Colorful:  false,
		AddSource: false,
	}
	if CFG.Debug {
		opts.Level = logger.LevelDebug
		opts.Colorful = ansi.IsSupported(os.Stderr.Fd())
		opts.AddSource = true
	}

	switch CFG.LogFmt {
	case "nano":
		LOG = logger.New(logger.NewNanoHandler(os.Stderr, opts))
	case "text":
		LOG = logger.New(logger.NewTextHandler(os.Stderr, opts))
	case "json":
		LOG = logger.New(logger.NewJsonHandler(os.Stderr, opts))
	default:
		panic("unknown log format")
	}
}
