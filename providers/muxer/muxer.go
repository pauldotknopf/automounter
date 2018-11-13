package muxer

import (
	"context"

	"golang.org/x/sync/errgroup"

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
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var eg errgroup.Group
	for _, provider := range p.p {
		eg.Go(func() error {
			err := provider.Start(ctx)
			if err != nil {
				cancel()
				return err
			}
			return nil
		})
	}

	return eg.Wait()
}

func (p *muxer) GetMedia() []providers.Media {
	result := make([]providers.Media, 0)
	for _, provider := range p.p {
		result = append(result, provider.GetMedia()...)
	}
	return result
}
