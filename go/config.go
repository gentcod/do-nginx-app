package main

import (
	// "time"

	"github.com/gentcod/environ"
)

type Config struct {
	HostAddr string
	ScriptPath string
	UpdateScriptPath string
	RemoteFilePath string
}

func LoadConfig(path string) (config Config, err error) {
	err = environ.Init(path, &config)
	if err != nil {
		return
	}

	return
}