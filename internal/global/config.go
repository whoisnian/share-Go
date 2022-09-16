package global

type Config struct {
	Debug          bool   `flag:"d,false,Enable debug output"`
	ReadOnly       bool   `flag:"ro,false,ReadOnly mode"`
	HTTPListenAddr string `flag:"l,127.0.0.1:9000,Server listen addr"`
	RootPath       string `flag:"p,uploads,Storage root directory"`
}

var CFG Config
