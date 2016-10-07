package beater

import (
	"fmt"
	"github.com/christiangalsterer/httpbeat/config"
	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/cfgfile"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/publisher"
)

type Httpbeat struct {
	done     chan bool
	HbConfig config.ConfigSettings
	client   publisher.Client
}

func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	bt := &Httpbeat{
		done: make(chan bool),
	}

	err := cfgfile.Read(&bt.HbConfig, "")
	if err != nil {
		logp.Err("Error reading config file: %v", err)
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	return bt, nil
}

func (h *Httpbeat) Run(b *beat.Beat) error {
	var poller *Poller

	logp.Info("httpbeat is running! Hit CTRL-C to stop it.")
	h.client = b.Publisher.Connect()

	for i, urlConfig := range h.HbConfig.Httpbeat.Urls {
		logp.Debug("httpbeat", "Creating poller #%v with URL: %v", i, urlConfig.Url)
		poller = NewPooler(h, urlConfig)
		go poller.Run()
	}

	for {
		select {
		case <-h.done:
			return nil
		}
	}
}

func (h *Httpbeat) Cleanup(b *beat.Beat) error {
	return nil
}

func (h *Httpbeat) Stop() {
	h.client.Close()
	close(h.done)
}
