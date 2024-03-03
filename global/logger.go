package global

import (
	"os"

	"github.com/whoisnian/glb/ansi"
	"github.com/whoisnian/glb/logger"
)

var LOG *logger.Logger

func SetupLogger() {
	opts := logger.NewOptions(logger.LevelInfo, false, false)
	if CFG.Debug {
		opts = logger.NewOptions(logger.LevelDebug, ansi.IsSupported(os.Stderr.Fd()), true)
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
