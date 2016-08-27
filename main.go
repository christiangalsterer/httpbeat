package main

import (
	"github.com/elastic/beats/libbeat/beat"
	httpbeat "github.com/christiangalsterer/httpbeat/beater"
)

var Version = "2.0.0-alpha.5"
var Name = "httpbeat"

func main() {
	beat.Run(Name, Version, httpbeat.New())
}
