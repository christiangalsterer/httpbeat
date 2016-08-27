package main

import (
	httpbeat "github.com/christiangalsterer/httpbeat/beater"
	"github.com/elastic/beats/libbeat/beat"
)

var Version = "2.0.0-alpha.5"
var Name = "httpbeat"

func main() {
	beat.Run(Name, Version, httpbeat.New())
}
