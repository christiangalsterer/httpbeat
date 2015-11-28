package main

import (
	"github.com/elastic/libbeat/beat"
	"github.com/elastic/libbeat/cfgfile"
	"github.com/elastic/libbeat/logp"
	"github.com/elastic/libbeat/publisher"
)

type Httpbeat struct {
	done                 chan struct{}
	HbConfig             ConfigSettings
	events               publisher.Client
}

func (h *Httpbeat) Config(b *beat.Beat) error {

	err := cfgfile.Read(&h.HbConfig, "")
	if err != nil {
		logp.Err("Error reading configuration file: %v", err)
		return err
	}

	logp.Info("httpbeat", "Init httpbeat")

	return nil
}

func (h *Httpbeat) Setup(b *beat.Beat) error {
	h.events = b.Events

	return nil
}

func (h *Httpbeat) Run(b *beat.Beat) error {
	var err error

	var poller *Poller
	for _, urlConfig := range h.HbConfig.Httpbeat.Urls {
		logp.Info("httpbeat", "Creating poller for URL: %v", urlConfig.Url)
		poller = NewSpooler(h, urlConfig)
		poller.Run()
	}

	return err
}

func (h *Httpbeat) Cleanup(b *beat.Beat) error {
	return nil
}

func (h *Httpbeat) Stop() {
	close(h.done)
}
