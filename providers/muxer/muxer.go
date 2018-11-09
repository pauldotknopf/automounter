package muxer

import (
	"github.com/pauldotknopf/automounter/providers"
)

type muxer struct {
	p []providers.MediaProvider
}

func Create(p []providers.MediaProvider) providers.MediaProvider {
	return &muxer{
		p,
	}
}

func (p *muxer) Name() string {
	return "muxer"
}

func (p *muxer) Start() error {
	var err error
	started := make([]providers.MediaProvider, 0)
	for _, provider := range p.p {
		if err != nil {
			err = provider.Start()
			if err == nil {
				started = append(started, provider)
			}
		}
	}
	if err != nil {
		// We failed to start a provider, let's stop any providers that did start
		for _, started := range started {
			started.Stop()
		}
	}
	return err
}

func (p *muxer) Stop() error {
	for _, provider := range p.p {
		err := provider.Stop()
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *muxer) GetMedia() []providers.Media {
	result := make([]providers.Media, 0)
	for _, provider := range p.p {
		result = append(result, provider.GetMedia()...)
	}
	return result
}
