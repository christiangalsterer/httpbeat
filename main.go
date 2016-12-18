package main

import (
	httpbeat "github.com/christiangalsterer/httpbeat/beater"
	"github.com/elastic/beats/libbeat/beat"
	"os"
)

var version = "3.0.0-SNAPSHOT"
var name = "httpbeat"

func main() {
	err := beat.Run(name, version, httpbeat.New)
	if err != nil {
		os.Exit(1)
	}
}
