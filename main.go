package main

import (
	"os"

	"github.com/elastic/libbeat/beat"
	"github.com/elastic/libbeat/logp"
)

// You can overwrite these, e.g.: go build -ldflags "-X main.Version 1.0.0-beta3"
var Version = "1.0.0-beta1"
var Name = "httpbeat"

func main() {

	hb := &Httpbeat{}

	b := beat.NewBeat(Name, Version, hb)

	b.CommandLineSetup()

	b.LoadConfig()

	err := hb.Config(b)
	if err != nil {
		logp.Critical("Config error: %v", err)
		os.Exit(1)
	}

	b.Run()
}
