package main

import (
	"github.com/elastic/beats/libbeat/beat"
	httpbeat "github.com/christiangalsterer/httpbeat/beat"
)

var Version = "1.0.0"
var Name = "httpbeat"

func main() {
	beat.Run(Name, Version, httpbeat.New())
}
