package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

// Config ...
type Config struct {
	ListenAddr string
}

// Load ...
func Load(path string) *Config {
	fi, err := os.Open(path)
	if err != nil {
		log.Panicln(err)
	}
	defer fi.Close()

	content, err := ioutil.ReadAll(fi)
	if err != nil {
		log.Panicln(err)
	}

	config := new(Config)
	err = json.Unmarshal(content, config)
	if err != nil {
		log.Panicln(err)
	}
	return config
}

// Save ...
func Save(path string, config *Config) {
	fi, err := os.OpenFile(path, os.O_WRONLY, 0)
	if err != nil {
		log.Panicln(err)
	}
	defer fi.Close()

	content, err := json.MarshalIndent(*config, "", "  ")
	if err != nil {
		log.Panicln(err)
	}

	_, err = fi.Write(content)
	if err != nil {
		log.Panicln(err)
	}
}
