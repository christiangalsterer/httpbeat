package main

import (
	httpbeat "github.com/christiangalsterer/httpbeat/beater"
	"github.com/elastic/beats/libbeat/beat"
	"os"
)

var Version = "2.0.0-beta.1"
var Name = "httpbeat"

func main() {
	err := beat.Run(Name, Version, httpbeat.New)
	if err != nil {
		os.Exit(1)
	}
}
