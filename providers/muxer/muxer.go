package muxer

import (
	"context"
	"sync"

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

func (p *muxer) Start(ctx context.Context) error {
	var wg sync.WaitGroup
	for _, provider := range p.p {
		wg.Add(1)
		go func() {
			provider.Start(ctx)
			wg.Done()
		}()
	}
	wg.Wait()
	return nil
}

func (p *muxer) GetMedia() []providers.Media {
	result := make([]providers.Media, 0)
	for _, provider := range p.p {
		result = append(result, provider.GetMedia()...)
	}
	return result
}
