package global

import "github.com/whoisnian/glb/config"

var CFG Config

type Config struct {
	Debug   bool   `flag:"d,false,Enable debug output"`
	LogFmt  string `flag:"log,nano,Log output format, one of nano, text and json"`
	Version bool   `flag:"v,false,Show version and quit"`

	ReadOnly   bool   `flag:"ro,false,ReadOnly mode"`
	RootPath   string `flag:"p,uploads,Storage root directory"`
	ListenAddr string `flag:"l,127.0.0.1:9000,Server listen addr"`
}

func SetupConfig() {
	err := config.FromCommandLine(&CFG)
	if err != nil {
		panic(err)
	}
}
