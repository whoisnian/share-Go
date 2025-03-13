package global

var (
	// build-time variables with `go build -ldflags="-X importpath.name=value"`
	// `-X` option from `go tool link --help`
	ModName   = "unknown"
	AppName   = "unknown"
	Version   = "unknown"
	BuildTime = "unknown"
)
